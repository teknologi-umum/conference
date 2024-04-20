package nocodb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// CreateTableRecords allows the creation of new records within a specified table. Records to be inserted are input as
// an array of key-value pair objects, where each key corresponds to a field name. Ensure that all the required fields
// are included in the payload, with exceptions for fields designated as auto-increment or those having default values.
//
// When dealing with 'Links' or 'Link To Another Record' field types, you should utilize the 'Create Link' API to insert
// relevant data.
//
// Certain read-only field types will be disregarded if included in the request. These field types include 'Look Up,'
// 'Roll Up,' 'Formula,' 'Auto Number,' 'Created By,' 'Updated By,' 'Created At,' 'Updated At,' 'Barcode,' and 'QR Code.'
func (c *Client) CreateTableRecords(ctx context.Context, tableId string, records []any) error {
	requestUrl, err := url.Parse(c.baseUrl + "/api/v2/tables/" + tableId + "/records")
	if err != nil {
		return fmt.Errorf("parsing url: %w", err)
	}

	requestBody, err := json.Marshal(records)
	if err != nil {
		return fmt.Errorf("marshaling records: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl.String(), bytes.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	request.Header.Add("xc-auth", c.apiToken)
	request.Header.Add("Content-Type", "application/json")

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

	return nil
}
