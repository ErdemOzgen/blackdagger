package executor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
)

type webhook struct {
	stdout    io.Writer
	stderr    io.Writer
	req       *resty.Request
	reqCancel context.CancelFunc
	cfg       *webhookConfig
}

type webhookConfig struct {
	URL                string            `json:"url" mapstructure:"url"`
	Method             string            `json:"method" mapstructure:"method"`
	Timeout            int               `json:"timeout" mapstructure:"timeout"`
	Headers            map[string]string `json:"headers" mapstructure:"headers"`
	Query              map[string]string `json:"query" mapstructure:"query"`
	Body               string            `json:"body" mapstructure:"body"`
	Silent             bool              `json:"silent" mapstructure:"silent"`
	SuccessStatusCodes []int             `json:"successStatusCodes" mapstructure:"successStatusCodes"`
}

var (
	errWebhookURLRequired = errors.New("webhook url is required")
	errWebhookStatusCode  = errors.New("webhook status code not successful")
)

func newWebhook(ctx context.Context, step dag.Step) (Executor, error) {
	cfg := &webhookConfig{}

	if len(step.Script) > 0 {
		if err := decodeWebhookConfigFromString(step.Script, cfg); err != nil {
			return nil, err
		}
	}

	if step.ExecutorConfig.Config != nil {
		if err := decodeWebhookConfig(step.ExecutorConfig.Config, cfg); err != nil {
			return nil, err
		}
	}

	if cfg.URL == "" && len(step.Args) > 0 {
		cfg.URL = step.Args[0]
	}

	cfg.URL = os.ExpandEnv(cfg.URL)
	cfg.Method = strings.ToUpper(os.ExpandEnv(cfg.Method))
	if cfg.Method == "" {
		cfg.Method = "POST"
	}
	cfg.Body = os.ExpandEnv(cfg.Body)
	for k, v := range cfg.Headers {
		cfg.Headers[k] = os.ExpandEnv(v)
	}
	for k, v := range cfg.Query {
		cfg.Query[k] = os.ExpandEnv(v)
	}

	if cfg.URL == "" {
		return nil, errWebhookURLRequired
	}

	ctx, cancel := context.WithCancel(ctx)
	client := resty.New()
	if cfg.Timeout > 0 {
		client.SetTimeout(time.Second * time.Duration(cfg.Timeout))
	}

	req := client.R().SetContext(ctx)
	if len(cfg.Headers) > 0 {
		req.SetHeaders(cfg.Headers)
	}
	if len(cfg.Query) > 0 {
		req.SetQueryParams(cfg.Query)
	}
	req.SetBody([]byte(cfg.Body))

	return &webhook{
		stdout:    os.Stdout,
		stderr:    os.Stderr,
		req:       req,
		reqCancel: cancel,
		cfg:       cfg,
	}, nil
}

func (e *webhook) SetStdout(out io.Writer) {
	e.stdout = out
}

func (e *webhook) SetStderr(out io.Writer) {
	e.stderr = out
}

func (e *webhook) Kill(_ os.Signal) error {
	e.reqCancel()
	return nil
}

func (e *webhook) Run() error {
	rsp, err := e.req.Execute(e.cfg.Method, e.cfg.URL)
	if err != nil {
		return err
	}

	success := isWebhookSuccess(rsp.StatusCode(), e.cfg.SuccessStatusCodes)

	if !success || !e.cfg.Silent {
		if _, err := e.stdout.Write([]byte(rsp.Status() + "\n")); err != nil {
			return err
		}
		if err := rsp.Header().Write(e.stdout); err != nil {
			return err
		}
	}

	if _, err := e.stdout.Write(rsp.Body()); err != nil {
		return err
	}

	if !success {
		return fmt.Errorf("%w: %d", errWebhookStatusCode, rsp.StatusCode())
	}

	return nil
}

func isWebhookSuccess(statusCode int, successCodes []int) bool {
	if len(successCodes) > 0 {
		return slices.Contains(successCodes, statusCode)
	}
	return statusCode >= 200 && statusCode < 300
}

func decodeWebhookConfig(dat map[string]any, cfg *webhookConfig) error {
	md, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		ErrorUnused: false,
		Result:      cfg,
	})
	return md.Decode(dat)
}

func decodeWebhookConfigFromString(s string, cfg *webhookConfig) error {
	if len(s) == 0 {
		return nil
	}
	ss := os.ExpandEnv(s)
	return json.Unmarshal([]byte(ss), cfg)
}

func init() {
	Register("webhook", newWebhook)
}
