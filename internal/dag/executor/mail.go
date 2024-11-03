package executor

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/mailer"
	"github.com/mitchellh/mapstructure"
)

type mail struct {
	stdout io.Writer
	stderr io.Writer
	mailer *mailer.Mailer
	cfg    *mailConfig
}

type mailConfig struct {
	From    string `mapstructure:"from"`
	To      string `mapstructure:"to"`
	Subject string `mapstructure:"subject"`
	Message string `mapstructure:"message"`
}

func newMail(ctx context.Context, step dag.Step) (Executor, error) {
	var cfg mailConfig
	if err := decodeMailConfig(step.ExecutorConfig.Config, &cfg); err != nil {
		return nil, err
	}

	cfg.From = os.ExpandEnv(cfg.From)
	cfg.To = os.ExpandEnv(cfg.To)
	cfg.Subject = os.ExpandEnv(cfg.Subject)
	cfg.Message = os.ExpandEnv(cfg.Message)

	exec := &mail{cfg: &cfg}

	dagCtx, err := dag.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	m := mailer.New(&mailer.NewMailerArgs{
		Host:     dagCtx.DAG.SMTP.Host,
		Port:     dagCtx.DAG.SMTP.Port,
		Username: dagCtx.DAG.SMTP.Username,
		Password: dagCtx.DAG.SMTP.Password,
	})
	exec.mailer = m

	return exec, nil
}

func (e *mail) SetStdout(out io.Writer) {
	e.stdout = out
}

func (e *mail) SetStderr(out io.Writer) {
	e.stderr = out
}

func (*mail) Kill(_ os.Signal) error {
	return nil
}

const mailLogTemplate = `sending email
-----
from: %s
to: %s
subject: %s
message: %s
-----
`

func (e *mail) Run() error {
	_, _ = e.stdout.Write(
		[]byte(fmt.Sprintf(
			mailLogTemplate,
			e.cfg.From,
			e.cfg.To,
			e.cfg.Subject,
			e.cfg.Message,
		)),
	)
	err := e.mailer.Send(
		e.cfg.From,
		[]string{e.cfg.To},
		e.cfg.Subject,
		e.cfg.Message,
		[]string{},
	)
	if err != nil {
		_, _ = e.stdout.Write([]byte("error occurred."))
	} else {
		_, _ = e.stdout.Write([]byte("sending email succeed."))
	}
	return err
}

func decodeMailConfig(dat map[string]any, cfg *mailConfig) error {
	md, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		ErrorUnused: false,
		Result:      cfg,
	})
	return md.Decode(dat)
}

func init() {
	Register("mail", newMail)
}
