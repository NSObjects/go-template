// Package pagination normalizes bounded offset pagination.
package pagination

import (
	"errors"
	"math"
)

const defaultPage = 1

// ErrInvalid reports invalid or overflowing pagination input.
var ErrInvalid = errors.New("invalid pagination")

// Options defines the pagination policy for one list endpoint.
type Options struct {
	DefaultPageSize int
	MaxPageSize     int
}

// Window is a validated page request translated into store-facing bounds.
type Window struct {
	Page     int
	PageSize int
	Offset   int
	Limit    int
}

// Normalize applies defaults, validates bounds, and rejects offset overflow.
func Normalize(page, pageSize int, options Options) (Window, error) {
	if page == 0 {
		page = defaultPage
	}
	if pageSize == 0 {
		pageSize = options.DefaultPageSize
	}
	if page < 1 ||
		pageSize < 1 ||
		options.DefaultPageSize < 1 ||
		options.MaxPageSize < options.DefaultPageSize ||
		pageSize > options.MaxPageSize {
		return Window{}, ErrInvalid
	}

	pageIndex := page - 1
	if pageIndex > math.MaxInt/pageSize {
		return Window{}, ErrInvalid
	}

	return Window{
		Page:     page,
		PageSize: pageSize,
		Offset:   pageIndex * pageSize,
		Limit:    pageSize,
	}, nil
}

// Bounds converts validated offset pagination into safe slice indexes.
func Bounds(total, offset, limit int) (start, end int, ok bool, err error) {
	if total < 0 || offset < 0 || limit < 1 {
		return 0, 0, false, ErrInvalid
	}
	if offset >= total {
		return 0, 0, false, nil
	}
	if limit > math.MaxInt-offset {
		return offset, total, true, nil
	}
	end = offset + limit
	if end > total {
		end = total
	}
	return offset, end, true, nil
}
