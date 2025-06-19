package store

import (
	"net/http"
	"strconv"
	"strings"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
}

func (feedQuery *PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	queryString := r.URL.Query()

	limit := queryString.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return *feedQuery, nil
		}

		feedQuery.Limit = l
	}

	offset := queryString.Get("offset")
	if offset != "" {
		l, err := strconv.Atoi(offset)
		if err != nil {
			return *feedQuery, nil
		}

		feedQuery.Offset = l
	}

	sort := queryString.Get("sort")
	if sort != "" {
		feedQuery.Sort = sort
	}

	tags := queryString.Get("tags")
	if tags != "" {
		feedQuery.Tags = strings.Split(tags, ",")
	}

	search := queryString.Get("search")
	if search != "" {
		feedQuery.Search = search
	}

	return *feedQuery, nil

}
