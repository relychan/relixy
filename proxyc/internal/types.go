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
	ErrParamKeyRequired                 = errors.New("param key must not be empty")
	ErrDuplicatedParamKey               = errors.New("routing pattern contains duplicate param key")
	ErrInvalidParamPattern              = errors.New("invalid param pattern")
	ErrDuplicatedRoutingPattern         = errors.New("routing pattern is duplicated")
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
