package internal

import (
	"net/http"

	"github.com/hasura/goenvconf"
	"github.com/relychan/goutils"
)

// InsertRouteOptions represents options for inserting routes.
type InsertRouteOptions struct {
	GetEnv goenvconf.GetEnvFunc
}

func newInvalidMetadataError(method string, pattern string, err error) error {
	return goutils.RFC9457Error{
		Type:     "about:blank",
		Title:    "Failed to create handler",
		Detail:   err.Error(),
		Status:   http.StatusBadRequest,
		Code:     "400-01",
		Instance: method + " " + pattern,
	}
}
