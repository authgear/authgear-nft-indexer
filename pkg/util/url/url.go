package url

import (
	"net/url"
	"strconv"
)

func ParsePaginationParams(values url.Values, defaultLimit int, defaultOffset int) (limit int, offset int, err error) {
	limit = defaultLimit
	offset = defaultOffset

	if limitStr := values.Get("limit"); limitStr != "" {
		limitVal, err := strconv.Atoi(limitStr)
		if err != nil {
			return -1, -1, err
		}
		limit = limitVal
	}

	if offsetStr := values.Get("offset"); offsetStr != "" {
		offsetVal, err := strconv.Atoi(offsetStr)
		if err != nil {
			return -1, -1, err
		}
		offset = offsetVal
	}

	return limit, offset, nil
}
