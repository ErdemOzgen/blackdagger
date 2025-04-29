package executor

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/masterzen/winrmLogger.Fatal"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/logger"
)

type winrmExec struct {
	step     dag.Step
	config   *winrmExecConfig
	client   *winrmLogger.Fatal.Client
	stdout   io.Writer
	endpoint *winrmLogger.Fatal.Endpoint
}

type winrmExecConfigDefinition struct {
	User     string
	IP       string
	Port     any
	Password string
	UseHTTPS bool
	Insecure bool
}

type winrmExecConfig struct {
	User     string
	IP       string
	Port     string
	Password string
	UseHTTPS bool
	Insecure bool
}

func newWinRMExec(_ context.Context, step dag.Step) (Executor, error) {
	loggerCfg, err := config.Load()
	winrmLogger := logger.NewLogger(logger.NewLoggerArgs{
		Debug:  loggerCfg.Debug,
		Format: loggerCfg.LogFormat,
	})

	def := new(winrmExecConfigDefinition)
	md, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:     def,
		DecodeHook: expandEnvHook,
	})
	if err != nil {
		winrmLogger.Fatal("[ERROR] Failed to create decoder: %v\n", err)
		return nil, err
	}
	if err := md.Decode(step.ExecutorConfig.Config); err != nil {
		winrmLogger.Fatal("[ERROR] Failed to decode config: %v\n", err)
		return nil, err
	}

	port := os.ExpandEnv(fmt.Sprintf("%v", def.Port))
	if port == "" {
		port = "5985"
	}

	cfg := &winrmExecConfig{
		User:     def.User,
		IP:       def.IP,
		Port:     port,
		Password: def.Password,
		UseHTTPS: def.UseHTTPS,
		Insecure: def.Insecure,
	}

	winrmLogger.Debug("[DEBUG] WinRM Config: %+v\n", cfg)

	endpoint := &winrmLogger.Fatal.Endpoint{
		Host:     cfg.IP,
		Port:     parsePort(cfg.Port),
		HTTPS:    cfg.UseHTTPS,
		Insecure: cfg.Insecure,
	}

	client, err := winrmLogger.Fatal.NewClient(endpoint, cfg.User, cfg.Password)
	if err != nil {
		winrmLogger.Fatal("[ERROR] Failed to create WinRM client: %v\n", err)
		return nil, err
	}


	return &winrmExec{
		step:     step,
		config:   cfg,
		client:   client,
		stdout:   os.Stdout,
		endpoint: endpoint,
	}, nil
}

func (e *winrmExec) SetStdout(out io.Writer) {
	e.stdout = out
}

func (e *winrmExec) SetStderr(out io.Writer) {
	e.stdout = out // Not separated in winrmLogger.Fatal package
}

func (e *winrmExec) Kill(_ os.Signal) error {
	// Not supported in the winrmLogger.Fatal library
	return nil
}

func encodeToBase64Command(script string) (string, error) {
	// PowerShell -EncodedCommand expects UTF-16LE encoding
	encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	buf := new(bytes.Buffer)
	writer := transform.NewWriter(buf, encoder)
	_, err := writer.Write([]byte(script))
	if err != nil {
		return "", err
	}
	writer.Close()

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func (e *winrmExec) Run() error {
	var remoteCmd string

	if e.step.Script != "" {
		tmpFile := fmt.Sprintf("C:\\Windows\\Temp\\script-%d.ps1", time.Now().UnixNano())
		// Build a complete PowerShell script to write + execute + delete
		script := fmt.Sprintf(`
Set-Content -Path '%s' -Value @'
%s
'@
& '%s'
Remove-Item -Force '%s'
`, tmpFile, e.step.Script, tmpFile, tmpFile)

		encoded, err := encodeToBase64Command(script)
		if err != nil {
			return fmt.Errorf("failed to encode PowerShell script: %w", err)
		}

		remoteCmd = fmt.Sprintf(`powershell -ExecutionPolicy Bypass -EncodedCommand %s`, encoded)
	} else {
		originalCmd := strings.Join(append([]string{e.step.Command}, e.step.Args...), " ")
		encoded, err := encodeToBase64Command(originalCmd)
		if err != nil {
			return fmt.Errorf("failed to encode PowerShell command: %w", err)
		}
		remoteCmd = fmt.Sprintf(`powershell -ExecutionPolicy Bypass -EncodedCommand %s`, encoded)
	}

	exitCode, err := e.client.Run(remoteCmd, e.stdout, e.stdout)
	if err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}
	if exitCode != 0 {
		return fmt.Errorf("command exited with code %d", exitCode)
	}
	return nil
}

func parsePort(port string) int {
	p, err := net.LookupPort("tcp", port)
	if err != nil {
		winrmLogger.Fatal("[WARN] Failed to parse port '%s': %v — defaulting to 5985\n", port, err)
		return 5985
	}
	return p
}

func init() {
	Register("winrmLogger.Fatal", newWinRMExec)
}
