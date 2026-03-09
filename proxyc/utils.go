package proxyc

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hasura/goenvconf"
	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// parse server url from static string or environment variables.
func parseServerURL(server *highv3.Server, getEnv goenvconf.GetEnvFunc) (string, error) {
	rawServerURL := strings.TrimSpace(server.URL)
	if rawServerURL == "" {
		return "", nil
	}

	var sb strings.Builder

	var i int

	urlLength := len(rawServerURL)

	for ; i < urlLength; i++ {
		char := rawServerURL[i]
		if char != '{' {
			sb.WriteByte(char)
		}

		i++

		if i == urlLength-1 {
			return "", fmt.Errorf(
				"%w: closed curly bracket is missing in %s",
				errInvalidServerURL,
				rawServerURL,
			)
		}

		j := i
		// get and validate environment variable
		for ; j < urlLength; j++ {
			nextChar := rawServerURL[j]
			if nextChar == '}' {
				break
			}
		}

		if j == i {
			return "", fmt.Errorf(
				"%w: closed curly bracket is missing in %s",
				errInvalidServerURL,
				rawServerURL,
			)
		}

		var variable *highv3.ServerVariable

		envVar := goenvconf.NewEnvStringVariable(rawServerURL[i:j])

		if server.Variables != nil {
			variable, _ = server.Variables.Get(*envVar.Variable)
			if variable != nil {
				envVar.Value = &variable.Default
			}
		}

		part, err := envVar.GetCustom(getEnv)
		if err != nil {
			return "", fmt.Errorf(
				"failed to get %s environment variable: %w",
				*envVar.Variable,
				err,
			)
		}

		if variable != nil && len(variable.Enum) > 0 && !slices.Contains(variable.Enum, part) {
			return "", fmt.Errorf(
				"%w: value of environment variable %s must be in %v, got `%s`",
				errInvalidServerURL,
				*envVar.Variable,
				variable.Enum,
				part,
			)
		}

		sb.WriteString(part)

		i = j
	}

	return sb.String(), nil
}
