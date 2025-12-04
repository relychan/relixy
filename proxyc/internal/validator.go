package internal

import (
	"net/http"

	orderedmap "github.com/pb33f/ordered-map/v2"
	"github.com/relychan/relyx/proxyc/handler/graphqlhandler"
	"github.com/relychan/relyx/schema"
)

// ValidateOperations validate operations of a route.
func ValidateOperations(pair *orderedmap.Pair[string, *schema.RelyProxyPathItem]) error {
	operation := pair.Value

	err := validateOperation(operation.Get)
	if err != nil {
		return MetadataValidationError{
			Method:  http.MethodGet,
			Path:    pair.Key,
			Message: err.Error(),
		}
	}

	err = validateOperation(operation.Post)
	if err != nil {
		return MetadataValidationError{
			Method:  http.MethodPost,
			Path:    pair.Key,
			Message: err.Error(),
		}
	}

	err = validateOperation(operation.Patch)
	if err != nil {
		return MetadataValidationError{
			Method:  http.MethodPatch,
			Path:    pair.Key,
			Message: err.Error(),
		}
	}

	err = validateOperation(operation.Put)
	if err != nil {
		return MetadataValidationError{
			Method:  http.MethodPut,
			Path:    pair.Key,
			Message: err.Error(),
		}
	}

	err = validateOperation(operation.Delete)
	if err != nil {
		return MetadataValidationError{
			Method:  http.MethodDelete,
			Path:    pair.Key,
			Message: err.Error(),
			Details: map[string]any{
				"error": err,
			},
		}
	}

	return nil
}

func validateOperation(operation *schema.RelyProxyOperation) error {
	if operation == nil {
		return nil
	}

	switch operation.Proxy.Type {
	case schema.ProxyTypeGraphQL:
		_, err := graphqlhandler.ValidateGraphQLString(operation.Proxy.Request.Query)
		if err != nil {
			return err
		}
	default:
	}

	return nil
}
