package executor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"syscall"

	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/util"
	"github.com/mitchellh/mapstructure"
)

type ansibleExecutor struct {
	ctx    context.Context
	step   dag.Step
	cfg    *ansibleConfig
	stdout io.Writer
	stderr io.Writer
	lock   sync.Mutex
	cmd    *exec.Cmd
	args   []string
}

type ansibleConfig struct {
	Binary            string            `mapstructure:"binary"`
	Playbook          string            `mapstructure:"playbook"`
	Inventory         string            `mapstructure:"inventory"`
	ExtraVars         map[string]string `mapstructure:"extraVars"`
	Tags              []string          `mapstructure:"tags"`
	SkipTags          []string          `mapstructure:"skipTags"`
	Limit             string            `mapstructure:"limit"`
	Check             bool              `mapstructure:"check"`
	Diff              bool              `mapstructure:"diff"`
	Forks             int               `mapstructure:"forks"`
	Become            bool              `mapstructure:"become"`
	User              string            `mapstructure:"user"`
	VaultPasswordFile string            `mapstructure:"vaultPasswordFile"`
	Env               map[string]string `mapstructure:"env"`
}

var (
	errAnsiblePlaybookRequired = errors.New("ansible playbook is required")
)

func newAnsible(ctx context.Context, step dag.Step) (Executor, error) {
	cfg := &ansibleConfig{}
	md, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:     cfg,
		DecodeHook: expandEnvHook,
	})
	if err != nil {
		return nil, err
	}
	if err := md.Decode(step.ExecutorConfig.Config); err != nil {
		return nil, err
	}

	if cfg.Binary == "" {
		cfg.Binary = "ansible-playbook"
	}

	playbook := strings.TrimSpace(cfg.Playbook)
	if playbook == "" {
		playbook = strings.TrimSpace(step.Command)
	}
	if playbook == "" {
		return nil, errAnsiblePlaybookRequired
	}

	args := buildAnsibleArgs(playbook, step.Args, cfg)

	return &ansibleExecutor{
		ctx:    ctx,
		step:   step,
		cfg:    cfg,
		stdout: os.Stdout,
		stderr: os.Stderr,
		args:   args,
	}, nil
}

func buildAnsibleArgs(
	playbook string, stepArgs []string, cfg *ansibleConfig,
) []string {
	args := []string{playbook}

	if cfg.Inventory != "" {
		args = append(args, "-i", cfg.Inventory)
	}
	if len(cfg.Tags) > 0 {
		args = append(args, "--tags", strings.Join(cfg.Tags, ","))
	}
	if len(cfg.SkipTags) > 0 {
		args = append(args, "--skip-tags", strings.Join(cfg.SkipTags, ","))
	}
	if cfg.Limit != "" {
		args = append(args, "--limit", cfg.Limit)
	}
	if cfg.Check {
		args = append(args, "--check")
	}
	if cfg.Diff {
		args = append(args, "--diff")
	}
	if cfg.Forks > 0 {
		args = append(args, "--forks", fmt.Sprintf("%d", cfg.Forks))
	}
	if cfg.Become {
		args = append(args, "--become")
	}
	if cfg.User != "" {
		args = append(args, "--user", cfg.User)
	}
	if cfg.VaultPasswordFile != "" {
		args = append(args, "--vault-password-file", cfg.VaultPasswordFile)
	}
	if len(cfg.ExtraVars) > 0 {
		keys := make([]string, 0, len(cfg.ExtraVars))
		for k := range cfg.ExtraVars {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			args = append(args, "-e", fmt.Sprintf("%s=%s", k, cfg.ExtraVars[k]))
		}
	}

	args = append(args, stepArgs...)
	return args
}

func (e *ansibleExecutor) SetStdout(out io.Writer) {
	e.stdout = out
}

func (e *ansibleExecutor) SetStderr(out io.Writer) {
	e.stderr = out
}

func (e *ansibleExecutor) Kill(sig os.Signal) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.cmd == nil || e.cmd.Process == nil {
		return nil
	}
	return syscall.Kill(-e.cmd.Process.Pid, sig.(syscall.Signal))
}

func (e *ansibleExecutor) Run() error {
	cmd := exec.CommandContext(e.ctx, e.cfg.Binary, e.args...)
	if e.step.Dir != "" && !util.FileExists(e.step.Dir) {
		return fmt.Errorf("directory %q does not exist", e.step.Dir)
	}

	dagCtx, err := dag.GetContext(e.ctx)
	if err != nil {
		return err
	}

	cmd.Dir = e.step.Dir
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, e.step.Variables...)
	cmd.Env = append(cmd.Env, dagCtx.Envs.All()...)
	for k, v := range e.cfg.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	e.step.OutputVariables.Range(func(_, value any) bool {
		cmd.Env = append(cmd.Env, value.(string))
		return true
	})
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
	cmd.Stdout = e.stdout
	cmd.Stderr = e.stderr

	e.lock.Lock()
	e.cmd = cmd
	e.lock.Unlock()

	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

func init() {
	Register("ansible", newAnsible)
}
