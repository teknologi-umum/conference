package nocodb

import (
	"io"
	"net/http"
)

type Client struct {
	apiToken   string
	baseUrl    string
	httpClient *http.Client
	logger     io.Writer
}

type ClientOptions struct {
	ApiToken   string
	BaseUrl    string
	HttpClient *http.Client
	Logger     io.Writer
}

func NewClient(options ClientOptions) (*Client, error) {
	if options.HttpClient == nil {
		options.HttpClient = http.DefaultClient
	}

	return &Client{
		apiToken:   options.ApiToken,
		baseUrl:    options.BaseUrl,
		httpClient: options.HttpClient,
		logger:     options.Logger,
	}, nil
}

// PageInfo attributes are particularly valuable when dealing with large datasets that are divided into multiple pages.
// They enable you to determine whether additional pages of records are available for retrieval or if you've reached
// the end of the dataset.
type PageInfo struct {
	// TotalRows indicates the total number of rows available for the specified conditions (if any).
	TotalRows int64 `json:"totalRows"`
	// Page specifies the current page number.
	Page int64 `json:"page"`
	// PageSize defaults to 25 and defines the number of records on each page.
	PageSize int64 `json:"pageSize"`
	// IsFirstPage is a boolean value that indicates whether the current page is the first page of records in the dataset.
	IsFirstPage bool `json:"isFirstPage"`
	// IsLastPage is a boolean value that indicates whether the current page is the last page of records in the dataset.
	IsLastPage bool `json:"isLastPage"`
}

type Sort interface {
	// parse is a private function to make the Sort interface can't be implemented outside of this package.
	parse() string
}

// sortImpl is an implementation of Sort interface.
type sortImpl struct {
	value string
}

func (s *sortImpl) parse() string {
	return s.value
}

func SortAscending(fieldName string) Sort {
	return &sortImpl{value: fieldName}
}

func SortDescending(fieldName string) Sort {
	return &sortImpl{value: "-" + fieldName}
}
