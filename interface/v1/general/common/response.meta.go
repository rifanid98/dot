package common

import "math"

func GetMeta(page, limit, total int64) Meta {
	lastPage := math.Ceil(float64(total) / float64(limit))
	if total < limit {
		lastPage = 1
	}
	from := page*limit - limit
	to := page * limit
	if to > total {
		to = total - 1
		if to < 0 {
			to = 0
		}
	}

	return Meta{
		CurrentPage: int(page),
		PerPage:     int(limit),
		From:        int(from),
		To:          int(to),
		Total:       int(total),
		LastPage:    int(lastPage),
	}
}
