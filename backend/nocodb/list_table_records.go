package nocodb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ListTableRecordOptions struct {
	// Fields llows you to specify the fields that you wish to include in your API response.
	// By default, all the fields are included in the response.
	Fields []string
	// Sort llows you to specify the fields by which you want to sort the records in your API response.
	// By default, sorting is done in ascending order for the designated fields
	Sort []Sort
	// Where enables you to define specific conditions for filtering records in your API response.
	// Multiple conditions can be combined using logical operators such as 'and' and 'or'.
	// Each condition consists of three parts: a field name, a comparison operator, and a value.
	//
	// Example: where=(field1,eq,value1)~and(field2,eq,value2) will filter records where 'field1' is equal to 'value1'
	// AND 'field2' is equal to 'value2'.
	//
	// You can also use other comparison operators like 'ne' (not equal), 'gt' (greater than), 'lt' (less than),
	// and more, to create complex filtering rules.
	//
	// If ViewId parameter is also included, then the filters included here will be applied over the filtering
	// configuration defined in the view.
	//
	// Please remember to maintain the specified format, and do not include spaces between the different
	// condition components
	//
	// SDK implementation note: I am too lazy to create a proper struct that can make the where clause
	// easy, so I'll leave this to the users.
	Where string
	// Offset enables you to control the pagination of your API response by specifying the number of records you
	// want to skip from the beginning of the result set. The default value for this parameter is set to 0, meaning
	// no records are skipped by default.
	//
	// Example: offset=25 will skip the first 25 records in your API response, allowing you to access records starting
	// from the 26th position.
	//
	// Please note that the 'offset' value represents the number of records to exclude, not an index value, so an
	// offset of 25 will skip the first 25 records.
	Offset int64
	// Limit enables you to set a limit on the number of records you want to retrieve in your API response.
	// By default, your response includes all the available records, but by using this parameter, you can control
	// the quantity you receive.
	Limit int64
	// ViewId View Identifier. Allows you to fetch records that are currently visible within a specific view.
	// API retrieves records in the order they are displayed if the SORT option is enabled within that view.
	//
	// Additionally, if you specify a sort query parameter, it will take precedence over any sorting configuration
	// defined in the view. If you specify a where query parameter, it will be applied over the filtering configuration
	// defined in the view.
	//
	// By default, all fields, including those that are disabled within the view, are included in the response.
	// To explicitly specify which fields to include or exclude, you can use the fields query parameter to customize
	// the output according to your requirements.
	ViewId string
}

type listTableRecordsResponse struct {
	List     any      `json:"list"`
	PageInfo PageInfo `json:"pageInfo"`
}

// ListTableRecords allows you to retrieve records from a specified table. You can customize the response by applying
// various query parameters for filtering, sorting, and formatting.
//
// Pagination: The response is paginated by default, with the first page being returned initially. The response includes
// the following additional information in the pageInfo JSON block.
//
// Note: `out` parameter MUST BE a pointer to a struct array.
func (c *Client) ListTableRecords(ctx context.Context, tableId string, out any, options ListTableRecordOptions) (PageInfo, error) {
	queryParams := &url.Values{}
	if len(options.Fields) > 0 {
		queryParams.Set("fields", strings.Join(options.Fields, ","))
	}

	if len(options.Sort) > 0 {
		var sortStrings []string
		for _, sort := range options.Sort {
			sortStrings = append(sortStrings, sort.parse())
		}
		queryParams.Set("sort", strings.Join(sortStrings, ","))
	}

	if options.Where != "" {
		queryParams.Set("where", options.Where)
	}

	if options.Offset > 0 {
		queryParams.Set("offset", strconv.FormatInt(options.Offset, 10))
	}

	if options.Limit > 0 {
		queryParams.Set("limit", strconv.FormatInt(options.Offset, 10))
	}

	if options.ViewId != "" {
		queryParams.Set("viewId", options.ViewId)
	}

	requestUrl, err := url.Parse(c.baseUrl + "/api/v2/tables/" + tableId + "/records")
	if err != nil {
		return PageInfo{}, fmt.Errorf("parsing url: %w", err)
	}

	requestUrl.RawQuery = queryParams.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return PageInfo{}, fmt.Errorf("creating request: %w", err)
	}

	request.Header.Add("xc-auth", c.apiToken)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return PageInfo{}, fmt.Errorf("executing http request: %w", err)
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
			return PageInfo{}, fmt.Errorf("unmarshaling bad request error: %w", err)
		}
		return PageInfo{}, badRequestError
	}

	var responseBody listTableRecordsResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return PageInfo{}, fmt.Errorf("deserializing json: %w", err)
	}

	// Re-marshal list response, unmarshal to request output
	marshalledList, err := json.Marshal(responseBody.List)
	if err != nil {
		return PageInfo{}, fmt.Errorf("re-marshalling list: %w", err)
	}

	err = json.Unmarshal(marshalledList, out)
	if err != nil {
		return PageInfo{}, fmt.Errorf("unmarshalling list: %w", err)
	}

	return responseBody.PageInfo, nil
}
