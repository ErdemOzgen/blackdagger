package logforward

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type HTTPSink struct {
	client  *resty.Client
	url     string
	headers map[string]string
}

func NewHTTPSink(url string, timeout time.Duration, headers map[string]string) *HTTPSink {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	client := resty.New().SetTimeout(timeout)
	return &HTTPSink{client: client, url: strings.TrimSpace(url), headers: headers}
}

func (s *HTTPSink) Forward(ctx context.Context, rec Record) error {
	if s == nil || s.url == "" {
		return nil
	}

	body, err := json.Marshal(rec)
	if err != nil {
		return err
	}

	req := s.client.R().SetContext(ctx).SetBody(body)
	if len(s.headers) > 0 {
		req.SetHeaders(s.headers)
	}
	req.SetHeader("Content-Type", "application/json")

	resp, err := req.Post(s.url)
	if err != nil {
		return err
	}
	if resp.StatusCode() < http.StatusOK || resp.StatusCode() >= http.StatusMultipleChoices {
		return fmt.Errorf("http sink returned status %d", resp.StatusCode())
	}
	return nil
}
