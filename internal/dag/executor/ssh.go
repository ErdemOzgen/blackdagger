package executor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/ssh"

	"github.com/ErdemOzgen/blackdagger/internal/dag"
)

type sshExec struct {
	step      dag.Step
	config    *sshExecConfig
	sshConfig *ssh.ClientConfig
	stdout    io.Writer
	session   *ssh.Session
}

type sshExecConfigDefinition struct {
	User                  string
	IP                    string
	Port                  any
	Key                   string
	Password              string
	StrictHostKeyChecking bool
}

type sshExecConfig struct {
	User     string
	IP       string
	Port     string
	Key      string
	Password string
}

// selectSSHAuthMethod selects the authentication method based on the configuration.
// If the key is provided, it will use the public key authentication method.
// Otherwise, it will use the password authentication method.
func selectSSHAuthMethod(cfg *sshExecConfig) (ssh.AuthMethod, error) {
	var (
		signer ssh.Signer
		err    error
	)

	if len(cfg.Key) != 0 {
		// Create the Signer for this private key.
		if signer, err = getPublicKeySigner(cfg.Key); err != nil {
			return nil, err
		}
		return ssh.PublicKeys(signer), nil
	}

	return ssh.Password(cfg.Password), nil
}

// expandEnvHook is a mapstructure decode hook that expands environment variables in string fields.
func expandEnvHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if f.Kind() != reflect.String || t.Kind() != reflect.String {
		return data, nil
	}
	return os.ExpandEnv(data.(string)), nil
}

func newSSHExec(_ context.Context, step dag.Step) (Executor, error) {
	def := new(sshExecConfigDefinition)
	md, err := mapstructure.NewDecoder(
		&mapstructure.DecoderConfig{
			Result:     def,
			DecodeHook: expandEnvHook,
		},
	)
	if err != nil {
		return nil, err
	}

	if err := md.Decode(step.ExecutorConfig.Config); err != nil {
		return nil, err
	}

	cfg := &sshExecConfig{
		User:     def.User,
		IP:       def.IP,
		Key:      def.Key,
		Password: def.Password,
	}

	// Handle Port as either string or int.
	port := os.ExpandEnv(fmt.Sprintf("%v", def.Port))
	if port == "" {
		port = "22"
	}
	cfg.Port = port

	// StrictHostKeyChecking is not supported yet.
	if def.StrictHostKeyChecking {
		return nil, errStrictHostKey
	}

	// Select the authentication method.
	authMethod, err := selectSSHAuthMethod(cfg)
	if err != nil {
		return nil, err
	}

	sshConfig := &ssh.ClientConfig{
		User: cfg.User,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		// nolint: gosec
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return &sshExec{
		step:      step,
		config:    cfg,
		sshConfig: sshConfig,
		stdout:    os.Stdout,
	}, nil
}

var errStrictHostKey = errors.New("StrictHostKeyChecking is not supported yet")

func (e *sshExec) SetStdout(out io.Writer) {
	e.stdout = out
}

func (e *sshExec) SetStderr(out io.Writer) {
	e.stdout = out
}

func (e *sshExec) Kill(_ os.Signal) error {
	if e.session != nil {
		return e.session.Close()
	}
	return nil
}

func (e *sshExec) Run() error {
	addr := net.JoinHostPort(e.config.IP, e.config.Port)
	conn, err := ssh.Dial("tcp", addr, e.sshConfig)
	if err != nil {
		return err
	}
	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	e.session = session
	defer session.Close()

	session.Stdout = e.stdout
	session.Stderr = e.stdout

	var remoteCmd string

	// If a script is provided, use it; otherwise, use the command.
	if e.step.Script != "" {
		// Create a unique temporary file path on the remote host.
		tmpFile := fmt.Sprintf("/tmp/script-%d.sh", time.Now().UnixNano())
		// Construct a heredoc command to write the script to the temporary file,
		// execute it with bash, and remove the temporary file afterwards.
		remoteCmd = fmt.Sprintf("cat > %s <<'EOF'\n%s\nEOF\nbash %s; rm -f %s", tmpFile, e.step.Script, tmpFile, tmpFile)
	} else {
		originalCmd := strings.Join(append([]string{e.step.Command}, e.step.Args...), " ")
		// Wrap the command in a shell to ensure proper parsing of shell operators.
		remoteCmd = fmt.Sprintf("sh -c %q", originalCmd)
	}

	return session.Run(remoteCmd)
}

// referenced code:
//
//	https://go.googlesource.com/crypto/+/master/ssh/example_test.go
//	https://gist.github.com/boyzhujian/73b5ecd37efd6f8dd38f56e7588f1b58
func getPublicKeySigner(path string) (ssh.Signer, error) {
	// A public key may be used to authenticate against the remote host
	// by using an unencrypted PEM-encoded private key file.
	key, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	return signer, nil
}

func init() {
	Register("ssh", newSSHExec)
}
