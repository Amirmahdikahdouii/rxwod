package http

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	defaultPage  = 1
	defaultLimit = 20
	maxLimit     = 100
)

func parsePagination(c echo.Context) (page int, limit int, err error) {
	page = defaultPage
	if raw := c.QueryParam("page"); raw != "" {
		page, err = strconv.Atoi(raw)
		if err != nil || page < 1 {
			return 0, 0, errors.New("page must be a positive integer")
		}
	}

	limit = defaultLimit
	if raw := c.QueryParam("limit"); raw != "" {
		limit, err = strconv.Atoi(raw)
		if err != nil || limit < 1 || limit > maxLimit {
			return 0, 0, fmt.Errorf("limit must be an integer between 1 and %d", maxLimit)
		}
	}

	return page, limit, nil
}
