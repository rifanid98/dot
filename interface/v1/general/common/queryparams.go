package common

import (
	"net/url"
	"strconv"
)

type CleanQueryParams struct {
	Limit       int
	Offset      int
	QueryParams url.Values
}

func QueryParams(queryParams url.Values) CleanQueryParams {
	limit, err := strconv.Atoi(queryParams.Get("per_page"))
	if err != nil {
		limit = 10
	}
	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil {
		page = 1
	}
	offset := (page - 1) * limit
	queryParams.Del("per_page")
	queryParams.Del("page")
	return CleanQueryParams{
		Limit:       limit,
		Offset:      offset,
		QueryParams: queryParams,
	}
}
