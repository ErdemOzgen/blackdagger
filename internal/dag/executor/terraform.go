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

type terraformExecutor struct {
	ctx      context.Context
	step     dag.Step
	cfg      *terraformConfig
	stdout   io.Writer
	stderr   io.Writer
	lock     sync.Mutex
	cmd      *exec.Cmd
	mainArgs []string
	initArgs []string
}

type terraformConfig struct {
	Binary        string            `mapstructure:"binary"`
	WorkingDir    string            `mapstructure:"workingDir"`
	Subcommand    string            `mapstructure:"subcommand"`
	Init          bool              `mapstructure:"init"`
	InitArgs      []string          `mapstructure:"initArgs"`
	VarFiles      []string          `mapstructure:"varFiles"`
	Vars          map[string]string `mapstructure:"vars"`
	BackendConfig map[string]string `mapstructure:"backendConfig"`
	Targets       []string          `mapstructure:"targets"`
	PlanFile      string            `mapstructure:"planFile"`
	AutoApprove   bool              `mapstructure:"autoApprove"`
	Env           map[string]string `mapstructure:"env"`
}

var (
	errTerraformSubcommandRequired = errors.New("terraform subcommand is required")
)

func newTerraform(ctx context.Context, step dag.Step) (Executor, error) {
	cfg := &terraformConfig{}
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
		cfg.Binary = "terraform"
	}

	subcommand := strings.TrimSpace(cfg.Subcommand)
	if subcommand == "" {
		subcommand = strings.TrimSpace(step.Command)
	}
	if subcommand == "" {
		return nil, errTerraformSubcommandRequired
	}

	mainArgs := buildTerraformMainArgs(subcommand, step.Args, cfg)
	initArgs := buildTerraformInitArgs(cfg)

	return &terraformExecutor{
		ctx:      ctx,
		step:     step,
		cfg:      cfg,
		stdout:   os.Stdout,
		stderr:   os.Stderr,
		mainArgs: mainArgs,
		initArgs: initArgs,
	}, nil
}

func buildTerraformMainArgs(
	subcommand string, stepArgs []string, cfg *terraformConfig,
) []string {
	var args []string
	if cfg.WorkingDir != "" {
		args = append(args, fmt.Sprintf("-chdir=%s", cfg.WorkingDir))
	}

	subcommand = strings.ToLower(subcommand)
	args = append(args, subcommand)
	args = append(args, stepArgs...)

	if cfg.PlanFile != "" {
		switch subcommand {
		case "plan":
			args = append(args, fmt.Sprintf("-out=%s", cfg.PlanFile))
		case "apply":
			args = append(args, cfg.PlanFile)
		}
	}

	for _, v := range cfg.VarFiles {
		args = append(args, fmt.Sprintf("-var-file=%s", v))
	}

	if len(cfg.Vars) > 0 {
		keys := make([]string, 0, len(cfg.Vars))
		for k := range cfg.Vars {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			args = append(args, "-var", fmt.Sprintf("%s=%s", k, cfg.Vars[k]))
		}
	}

	if len(cfg.Targets) > 0 {
		sortedTargets := append([]string(nil), cfg.Targets...)
		sort.Strings(sortedTargets)
		for _, t := range sortedTargets {
			args = append(args, fmt.Sprintf("-target=%s", t))
		}
	}

	if cfg.AutoApprove && (subcommand == "apply" || subcommand == "destroy") {
		args = append(args, "-auto-approve")
	}

	return args
}

func buildTerraformInitArgs(cfg *terraformConfig) []string {
	if !cfg.Init {
		return nil
	}

	var args []string
	if cfg.WorkingDir != "" {
		args = append(args, fmt.Sprintf("-chdir=%s", cfg.WorkingDir))
	}
	args = append(args, "init")
	args = append(args, cfg.InitArgs...)

	if len(cfg.BackendConfig) > 0 {
		keys := make([]string, 0, len(cfg.BackendConfig))
		for k := range cfg.BackendConfig {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			args = append(args,
				fmt.Sprintf("-backend-config=%s=%s", k, cfg.BackendConfig[k]),
			)
		}
	}

	return args
}

func (e *terraformExecutor) SetStdout(out io.Writer) {
	e.stdout = out
}

func (e *terraformExecutor) SetStderr(out io.Writer) {
	e.stderr = out
}

func (e *terraformExecutor) Kill(sig os.Signal) error {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.cmd == nil || e.cmd.Process == nil {
		return nil
	}
	return syscall.Kill(-e.cmd.Process.Pid, sig.(syscall.Signal))
}

func (e *terraformExecutor) Run() error {
	if len(e.initArgs) > 0 {
		if err := e.runCommand(e.ctx, e.initArgs); err != nil {
			return err
		}
	}
	return e.runCommand(e.ctx, e.mainArgs)
}

func (e *terraformExecutor) runCommand(ctx context.Context, args []string) error {
	cmd := exec.CommandContext(ctx, e.cfg.Binary, args...)
	if e.step.Dir != "" && !util.FileExists(e.step.Dir) {
		return fmt.Errorf("directory %q does not exist", e.step.Dir)
	}

	dagCtx, err := dag.GetContext(ctx)
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
	Register("terraform", newTerraform)
}
