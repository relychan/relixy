package internal

import (
	"errors"
	"net/http"

	"github.com/relychan/goutils"
)

var (
	ErrWildcardMustBeLast = errors.New(
		"wildcard '*' must be the last value in a route. trim trailing text or use a '{param}' instead",
	)
	ErrMissingClosingBracket            = errors.New("route param closing delimiter '}' is missing")
	ErrDuplicatedParamKey               = errors.New("routing pattern contains duplicate param key")
	ErrInvalidRegexpPatternParamInRoute = errors.New("invalid regexp pattern in route param")
	ErrReplaceMissingChildNode          = errors.New("replacing missing child node")
)

func newInvalidOperationMetadataError(method string, pattern string, err error) error {
	return goutils.RFC9457Error{
		Type:     "about:blank",
		Title:    "Invalid Operation Metadata",
		Detail:   err.Error(),
		Status:   http.StatusBadRequest,
		Code:     "invalid-operation-metadata",
		Instance: method + " " + pattern,
	}
}
