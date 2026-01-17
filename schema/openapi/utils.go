package openapi

import (
	"slices"

	highv3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

// ExtractCommonParametersOfOperation extracts common parameters from operation's parameters.
func ExtractCommonParametersOfOperation(
	pathParams []*highv3.Parameter,
	operation *highv3.Operation,
) []*highv3.Parameter {
	if operation == nil || len(operation.Parameters) == 0 {
		return pathParams
	}

	remainParams := make([]*highv3.Parameter, 0, len(operation.Parameters))

	for _, param := range operation.Parameters {
		if slices.ContainsFunc(pathParams, func(originalParam *highv3.Parameter) bool {
			return param.Name == originalParam.Name && param.In == originalParam.In
		}) {
			continue
		}

		if param.In == string(InPath) {
			pathParams = append(pathParams, param)
		} else {
			remainParams = append(remainParams, param)
		}
	}

	operation.Parameters = slices.Clip(remainParams)

	return pathParams
}

// MergeParameters merge parameter slices by unique name and location.
func MergeParameters(dest []*highv3.Parameter, src []*highv3.Parameter) []*highv3.Parameter {
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

func mergeOrderedMaps[K comparable, V any](dest, src *orderedmap.Map[K, V]) *orderedmap.Map[K, V] {
	if src == nil || src.Len() == 0 {
		return dest
	}

	if dest == nil {
		return src
	}

	for iter := src.Oldest(); iter != nil; iter = iter.Next() {
		dest.Set(iter.Key, iter.Value)
	}

	return dest
}
