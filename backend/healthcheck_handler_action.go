package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/urfave/cli/v2"
)

func HealthcheckHandlerAction(c *cli.Context) error {
	port := c.String("port")
	timeout := c.Duration("timeout")
	if timeout == 0 {
		timeout = time.Minute
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:"+port+"/api/public/ping", nil)
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode >= 200 && response.StatusCode < 400 {
		return nil
	}

	return fmt.Errorf("not healthy")
}
