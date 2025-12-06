package schema

import "slices"

// ExtractCommonParametersOfOperation extracts common parameters from operation's parameters.
func ExtractCommonParametersOfOperation(
	pathParams []Parameter,
	operation *RelyProxyOperation,
) []Parameter {
	if operation == nil || len(operation.Parameters) == 0 {
		return pathParams
	}

	remainParams := make([]Parameter, 0, len(operation.Parameters))

	for _, param := range operation.Parameters {
		if slices.ContainsFunc(pathParams, func(originalParam Parameter) bool {
			return param.Name == originalParam.Name && param.In == originalParam.In
		}) {
			continue
		}

		if param.In == InPath {
			pathParams = append(pathParams, param)
		} else {
			remainParams = append(remainParams, param)
		}
	}

	operation.Parameters = slices.Clip(remainParams)

	return pathParams
}

// MergeParameters merge parameter slices by unique name and location.
func MergeParameters(dest []Parameter, src []Parameter) []Parameter {
L:
	for _, srcParam := range src {
		for j, destParam := range dest {
			if destParam.Name == srcParam.Name && destParam.In == srcParam.In {
				dest[j] = srcParam

				continue L
			}
		}

		dest = append(dest, srcParam)
	}

	return dest
}
