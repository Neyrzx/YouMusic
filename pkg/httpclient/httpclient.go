package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/neyrzx/youmusic/internal/config"
	"github.com/neyrzx/youmusic/pkg/logger"
	"github.com/rogpeppe/retry"
	"github.com/rs/zerolog"
)

const (
	packageKey  = "gateways"
	packageName = "musicInfo"
)

type HTTPClient struct {
	logger         *zerolog.Logger
	retryStrategy  retry.Strategy
	badStatusCodes []int
}

func NewHTTPClient(cfg config.GatewayHTTPClient) *HTTPClient {
	logger := logger.DefaultLogger().With().Str(packageKey, packageName).Logger()

	return &HTTPClient{
		logger: &logger,
		retryStrategy: retry.Strategy{
			Delay:       cfg.RetryStratagyDelay,
			MaxDelay:    cfg.RetryStrategyMaxDelay,
			MaxDuration: cfg.RetryStrategyMaxDuration,
			Factor:      cfg.RetryStrategyFactor,
		},
		badStatusCodes: []int{http.StatusBadRequest, http.StatusInternalServerError},
	}
}

func (c *HTTPClient) Get(ctx context.Context, url string) (io.ReadCloser, error) {
	var (
		body []byte
		err  error
	)

	for i := c.retryStrategy.Start(); ; {
		if body, err = c.get(ctx, url); err == nil {
			break
		}

		c.logger.Err(err).Msg("failed getting response")

		if !i.Next(nil) {
			return nil, fmt.Errorf("failed to getting response from %s after %d tries: %w", url, i.Count(), err)
		}
	}

	return io.NopCloser(bytes.NewBuffer(body)), nil
}

func (c *HTTPClient) get(ctx context.Context, url string) (body []byte, err error) {
	var (
		res *http.Response
		req *http.Request
	)

	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
		return nil, fmt.Errorf("failed to http.NewRequest: %w", err)
	}

	if res, err = http.DefaultClient.Do(req); err != nil {
		return nil, fmt.Errorf("failed to http.DefaultClient.Do: %w", err)
	}

	if slices.Contains(c.badStatusCodes, res.StatusCode) {
		return nil, fmt.Errorf("failed with statusCode: %d", res.StatusCode)
	}

	defer res.Body.Close()

	if body, err = io.ReadAll(res.Body); err != nil {
		return nil, fmt.Errorf("failed to io.ReadAll(response body): %w", err)
	}

	return body, nil
}
