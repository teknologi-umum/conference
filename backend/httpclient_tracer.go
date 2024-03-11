package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/getsentry/sentry-go"
)

func NewSentryRoundTripper(originalRoundTripper http.RoundTripper, tracePropagationTargets []string) http.RoundTripper {
	if originalRoundTripper == nil {
		originalRoundTripper = http.DefaultTransport
	}

	return &SentryRoundTripper{
		originalRoundTripper:    originalRoundTripper,
		tracePropagationTargets: tracePropagationTargets,
	}
}

type SentryRoundTripper struct {
	originalRoundTripper    http.RoundTripper
	tracePropagationTargets []string
}

func (s *SentryRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	// Respect trace propagation targets
	if len(s.tracePropagationTargets) > 0 {
		requestUrlString := request.URL.String()
		for _, t := range s.tracePropagationTargets {
			if strings.Contains(requestUrlString, t) {
				continue
			}

			return s.originalRoundTripper.RoundTrip(request)
		}
	}

	// Start Sentry trace
	ctx := request.Context()
	cleanRequestURL := request.URL.Path

	span := sentry.StartSpan(ctx, "http.client", sentry.WithTransactionName(fmt.Sprintf("%s %s", request.Method, cleanRequestURL)))
	defer span.Finish()

	span.SetData("http.query", request.URL.Query().Encode())
	span.SetData("http.fragment", request.URL.Fragment)
	span.SetData("http.request.method", request.Method)

	request.Header.Add("Baggage", span.ToBaggage())
	request.Header.Add("Sentry-Trace", span.ToSentryTrace())

	response, err := s.originalRoundTripper.RoundTrip(request)

	if response != nil {
		span.Status = sentry.HTTPtoSpanStatus(response.StatusCode)
		span.SetData("http.response.status_code", response.Status)
		span.SetData("http.response_content_length", strconv.FormatInt(response.ContentLength, 10))
	}

	return response, err
}
