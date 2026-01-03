package graphqlhandler

import (
	"errors"
	"fmt"

	"github.com/hasura/goenvconf"
	orderedmap "github.com/pb33f/ordered-map/v2"
	"github.com/relychan/relixy/schema/base_schema"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

var (
	ErrProxyActionInvalid           = errors.New("proxy action must exist with the graphql type")
	ErrGraphQLQueryEmpty            = errors.New("query is required for graphql proxy")
	ErrGraphQLUnsupportedQueryBatch = errors.New("graphql query batch is not supported")
)

// ValidateGraphQLString parses and validates the GraphQL query string.
func ValidateGraphQLString(query string) (*GraphQLHandler, error) {
	if query == "" {
		return nil, ErrGraphQLQueryEmpty
	}

	doc, err := parser.ParseQuery(&ast.Source{
		Input: query,
	})
	if err != nil {
		return nil, err
	}

	switch len(doc.Operations) {
	case 0:
		return nil, ErrGraphQLQueryEmpty
	case 1:
		graphqlOperation := doc.Operations[0]

		handler := &GraphQLHandler{
			query:               query,
			variableDefinitions: graphqlOperation.VariableDefinitions,
			operationName:       graphqlOperation.Name,
			operation:           graphqlOperation.Operation,
		}

		return handler, nil
	default:
		return nil, ErrGraphQLUnsupportedQueryBatch
	}
}

func validateGraphQLVariables(
	inputs *orderedmap.OrderedMap[string, *base_schema.GraphQLVariableDefinition],
	getEnvFunc goenvconf.GetEnvFunc,
) (map[string]graphqlVariable, error) {
	result := make(map[string]graphqlVariable)

	if inputs == nil {
		return result, nil
	}

	for iter := inputs.Oldest(); iter != nil; iter = iter.Next() {
		variable := graphqlVariable{
			Path: iter.Value.Path,
		}

		if iter.Value.Default != nil {
			defaultValue, err := iter.Value.Default.GetCustom(getEnvFunc)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to evaluate default value of variable %s: %w",
					iter.Key,
					err,
				)
			}

			variable.Default = defaultValue
		}

		result[iter.Key] = variable
	}

	return result, nil
}
