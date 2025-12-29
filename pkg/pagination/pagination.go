package pagination

import (
	"encoding/json"
	"errors"
	"fmt"
	"paymentserviceklink/app/web"
	"paymentserviceklink/pkg/util"
	"strconv"
	"strings"
	"time"
)

var (
	// DefaultPageSize specifies the default page size
	DefaultPageSize = 100
	// MaxPageSize specifies the maximum page size
	MaxPageSize = 100
	// PageVar specifies the query parameter name for page number
	PageVar = "page"
	// PageSizeVar specifies the query parameter name for page size
	PageSizeVar = "per_page"
)

// Pages represents a paginated list of data items.
type Pages struct {
	Page         int         `json:"page"`
	PerPage      int         `json:"per_page"`
	PageCount    int         `json:"page_count"`
	TotalCount   int         `json:"total_count"`
	Items        interface{} `json:"items"`
	Sort         string      `json:"sort;default:asc"`
	Unscoped     bool        `json:"Unscoped"`
	Filters      []*Filter   `json:"filter"`
	JoinOperator string      `json:"joinOperator"`
}

type Filter struct {
	ID       string `json:"id"`
	Value    any    `json:"value"`
	Variant  string `json:"variant"`
	Operator string `json:"operator"`
	FilterID string `json:"filterId"`
	Table    string
}

func New(page, perPage string, total int, sortBy, unscoped, filter, joinOperator string) (*Pages, error) {

	pageInt := ParseIntFallback(page, 1)
	perPageInt := ParseIntFallback(perPage, DefaultPageSize)

	if pageInt <= 0 {
		pageInt = DefaultPageSize
	}
	if perPageInt > MaxPageSize {
		perPageInt = MaxPageSize
	}
	pageCount := 0
	if total > 0 {
		pageCount = -1
		pageCount = (total + perPageInt - 1) / perPageInt
		if pageInt > pageCount {
			pageInt = pageCount
		}
	}
	if pageInt < 1 {
		pageInt = 1
	}

	if sortBy != "asc" && sortBy != "desc" {
		sortBy = "desc"
	}

	if joinOperator != "and" && joinOperator != "or" {
		joinOperator = "and"
	}

	var filters []*Filter

	if filter != "" {
		err := json.Unmarshal([]byte(filter), &filters)
		if err != nil {
			return nil, errors.New("invalid format filter " + err.Error())
		}
	}

	pages := &Pages{
		Page:         pageInt,
		PerPage:      perPageInt,
		TotalCount:   total,
		PageCount:    pageCount,
		Sort:         sortBy,
		Unscoped:     util.StringToBoolean(unscoped),
		Filters:      filters,
		JoinOperator: joinOperator,
	}

	err := pages.ValidationFilterVariant()
	if err != nil {
		return nil, err
	}

	return pages, nil
}

// NewFromRequest creates a Pages object using the query parameters found in the given HTTP request.
// count stands for the total number of items. Use -1 if this is unknown.
func NewFromRequest(page, perPage string, count int, sortBy string, unscoped string, filter string, joinOperator string) (*Pages, error) {
	return New(page, perPage, count, sortBy, unscoped, filter, joinOperator)
}

// parseInt parses a string into an integer. If parsing is failed, defaultValue will be returned.
func ParseIntFallback(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if result, err := strconv.Atoi(value); err == nil {
		return result
	}
	return defaultValue
}

// Offset returns the OFFSET value that can be used in a SQL statement.
func (p *Pages) Offset() int {
	return (p.Page - 1) * p.PerPage
}

// Limit returns the LIMIT value that can be used in a SQL statement.
func (p *Pages) Limit() int {
	return p.PerPage
}

// BuildLinkHeader returns an HTTP header containing the links about the pagination.
func (p *Pages) BuildLinkHeader(baseURL string, defaultPerPage int) string {
	links := p.BuildLinks(baseURL, defaultPerPage)
	header := ""
	if links[0] != "" {
		header += fmt.Sprintf("<%v>; rel=\"first\", ", links[0])
		header += fmt.Sprintf("<%v>; rel=\"prev\"", links[1])
	}
	if links[2] != "" {
		if header != "" {
			header += ", "
		}
		header += fmt.Sprintf("<%v>; rel=\"next\"", links[2])
		if links[3] != "" {
			header += fmt.Sprintf(", <%v>; rel=\"last\"", links[3])
		}
	}
	return header
}

// BuildLinks returns the first, prev, next, and last links corresponding to the pagination.
// A link could be an empty string if it is not needed.
// For example, if the pagination is at the first page, then both first and prev links
// will be empty.
func (p *Pages) BuildLinks(baseURL string, defaultPerPage int) [4]string {
	var links [4]string
	pageCount := p.PageCount
	page := p.Page
	if pageCount >= 0 && page > pageCount {
		page = pageCount
	}
	if strings.Contains(baseURL, "?") {
		baseURL += "&"
	} else {
		baseURL += "?"
	}
	if page > 1 {
		links[0] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, 1)
		links[1] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, page-1)
	}
	if pageCount >= 0 && page < pageCount {
		links[2] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, page+1)
		links[3] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, pageCount)
	} else if pageCount < 0 {
		links[2] = fmt.Sprintf("%v%v=%v", baseURL, PageVar, page+1)
	}
	if perPage := p.PerPage; perPage != defaultPerPage {
		for i := 0; i < 4; i++ {
			if links[i] != "" {
				links[i] += fmt.Sprintf("&%v=%v", PageSizeVar, perPage)
			}
		}
	}

	return links
}

func (p *Pages) GetMetadata() *web.Metadata {
	pageCount := 0
	if p.TotalCount > 0 {
		pageCount = -1
		pageCount = (p.TotalCount + p.PerPage - 1) / p.PerPage
		if p.Page > pageCount {
			p.Page = pageCount
		}

		p.PageCount = pageCount
	}
	return &web.Metadata{
		PerPage:    p.PerPage,
		TotalCount: p.TotalCount,
		PageCount:  p.PageCount,
		OrderBy:    p.Sort,
	}
}

func (p *Pages) ValidationFilterVariant() error {
	for _, filter := range p.Filters {
		if filter.Variant == "boolean" {
			if filter.Value != "true" && filter.Value != "false" {
				return fmt.Errorf("invalid filter value %v is not suitable for variant %v", filter.Value, filter.Variant)
			}
		}

		if filter.Variant == "number" {
			_, err := strconv.Atoi(filter.Value.(string))
			if err != nil {
				return fmt.Errorf("invalid filter value %v is not suitable for variant %v", filter.Value, filter.Variant)
			}
		}

		if filter.Variant == "date" {
			_, err := time.Parse(time.DateOnly, filter.Value.(string))
			if err != nil {
				return fmt.Errorf("invalid filter value %v is not suitable for variant %v", filter.Value, filter.Variant)
			}
		}

		if filter.Variant == "time" {
			_, err := time.Parse(time.TimeOnly, filter.Value.(string))
			if err != nil {
				return fmt.Errorf("invalid filter value %v is not suitable for variant %v", filter.Value, filter.Variant)
			}
		}
	}

	return nil
}
