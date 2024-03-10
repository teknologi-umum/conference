package nocodb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type ReadTableRecordsOptions struct {
	// Fields a llows you to specify the fields that you wish to include in your API response. By default, all the fields are included in the response.
	Fields []string
}

// ReadTableRecords allows you to retrieve a single record identified by Record-ID, serving as unique identifier for
// the record from a specified table.
//
// Note: `out` parameter MUST BE a pointer to a struct.
func (c *Client) ReadTableRecords(ctx context.Context, tableId string, recordId string, out any, options ReadTableRecordsOptions) error {
	queryParams := &url.Values{}
	if len(options.Fields) > 0 {
		queryParams.Set("fields", strings.Join(options.Fields, ","))
	}

	requestUrl, err := url.Parse(c.baseUrl + "/api/v2/tables/" + tableId + "/records/" + recordId)
	if err != nil {
		return fmt.Errorf("parsing url: %w", err)
	}

	requestUrl.RawQuery = queryParams.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	request.Header.Add("xc-auth", c.apiToken)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("executing http request: %w", err)
	}
	defer func() {
		if response.Body != nil {
			err := response.Body.Close()
			if err != nil {
				if c.logger != nil {
					_, _ = c.logger.Write([]byte("Closing response body: " + err.Error()))
				}
			}
		}
	}()

	if response.StatusCode == 400 {
		var badRequestError BadRequestError
		err = json.NewDecoder(response.Body).Decode(&badRequestError)
		if err != nil {
			return fmt.Errorf("unmarshaling bad request error: %w", err)
		}
		return badRequestError
	}

	err = json.NewDecoder(response.Body).Decode(&out)
	if err != nil {
		return fmt.Errorf("decoding response body: %w", err)
	}

	return nil
}
