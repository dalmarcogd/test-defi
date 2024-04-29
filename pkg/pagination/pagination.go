package pagination

const (
	DefaultPaginationLimit = 25
	MaxPaginationLimit     = 50
)

type (
	PageParams struct {
		Page int `query:"page"`
		Size int `query:"size"`
	}

	PageResponse[T any] struct {
		Items       []T `json:"items"`
		TotalItems  int `json:"total_items"`
		TotalPages  int `json:"total_pages"`
		TotalInPage int `json:"total_in_page"`
	}

	PaginatedList[T any] struct {
		Items []T
		Total int
	}
)

var Disabled = PageParams{Page: -1}

func ParseParams(params PageParams) PageParams {
	var page, size int
	if params.Page == 0 {
		page = 1
	} else {
		page = params.Page
	}

	switch {
	case params.Size > 0 && params.Size <= MaxPaginationLimit:
		size = params.Size
	case params.Size > MaxPaginationLimit:
		size = MaxPaginationLimit
	default:
		size = DefaultPaginationLimit
	}

	return PageParams{
		Page: page,
		Size: size,
	}
}

func ParseResponse[T any, F ~func(T) (O, error), O any](params PageParams, paginatedList PaginatedList[T], convertFunc F) (PageResponse[O], error) {
	totalPages := 1
	if params.Page > 0 {
		totalPages = ((paginatedList.Total - 1) / params.Size) + 1
	}

	listItems := make([]O, len(paginatedList.Items))
	for i, item := range paginatedList.Items {
		o, err := convertFunc(item)
		if err != nil {
			return PageResponse[O]{}, err
		}

		listItems[i] = o
	}

	return PageResponse[O]{
		Items:       listItems,
		TotalItems:  paginatedList.Total,
		TotalPages:  totalPages,
		TotalInPage: len(paginatedList.Items),
	}, nil
}
