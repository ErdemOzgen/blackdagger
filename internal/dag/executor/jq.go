package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/itchyny/gojq"
	"github.com/mitchellh/mapstructure"
)

type jq struct {
	stdout io.Writer
	stderr io.Writer
	query  string
	input  map[string]any
	cfg    *jqConfig
}

type jqConfig struct {
	Raw bool `mapstructure:"raw"`
}

func newJQ(_ context.Context, step dag.Step) (Executor, error) {
	var jqCfg jqConfig
	if step.ExecutorConfig.Config != nil {
		if err := decodeJqConfig(
			step.ExecutorConfig.Config, &jqCfg,
		); err != nil {
			return nil, err
		}
	}
	s := os.ExpandEnv(step.Script)
	input := map[string]any{}
	if err := json.Unmarshal([]byte(s), &input); err != nil {
		return nil, err
	}
	return &jq{
		stdout: os.Stdout,
		input:  input,
		query:  step.CmdWithArgs,
		cfg:    &jqCfg,
	}, nil
}

func (e *jq) SetStdout(out io.Writer) {
	e.stdout = out
}

func (e *jq) SetStderr(out io.Writer) {
	e.stderr = out
}

func (*jq) Kill(_ os.Signal) error {
	return nil
}

func (e *jq) Run() error {
	query, err := gojq.Parse(e.query)
	if err != nil {
		return err
	}
	iter := query.Run(e.input)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			_, _ = e.stderr.Write([]byte(fmt.Sprintf("%#v", err)))
			continue
		}
		val, err := json.MarshalIndent(v, "", "    ")
		if err != nil {
			_, _ = e.stderr.Write([]byte(fmt.Sprintf("%#v", err)))
			continue
		}
		if e.cfg.Raw {
			s := string(val)
			_, _ = e.stdout.Write([]byte(strings.Trim(s, `"`)))
		} else {
			_, _ = e.stdout.Write(val)
		}
	}
	return nil
}

func decodeJqConfig(dat map[string]any, cfg *jqConfig) error {
	md, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		ErrorUnused: false,
		Result:      cfg,
	})
	return md.Decode(dat)
}

func init() {
	Register("jq", newJQ)
}
