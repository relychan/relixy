package openapi

import (
	"slices"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v2high "github.com/pb33f/libopenapi/datamodel/high/v2"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/pb33f/libopenapi/utils"
	"go.yaml.in/yaml/v4"
)

func convertSwaggerToOpenAPIv3Document(src *v2high.Swagger) *v3high.Document {
	result := &v3high.Document{
		Version:      "3.0.0",
		Info:         src.Info,
		Tags:         src.Tags,
		ExternalDocs: src.ExternalDocs,
		Extensions:   src.Extensions,
		Security:     src.Security,
		Components:   &v3high.Components{},
	}

	if src.Host != "" {
		scheme := "https"

		if len(src.Schemes) > 0 && !slices.Contains(src.Schemes, scheme) {
			scheme = src.Schemes[0]
		}

		serverURL := scheme + "://" + src.Host + src.BasePath

		if src.BasePath != "" && src.BasePath != "/" {
			serverURL += src.BasePath
		}

		result.Servers = []*v3high.Server{
			{
				URL: serverURL,
			},
		}
	}

	result.Components.Parameters = convertParametersV2ToV3(src.Parameters)
	result.Components.SecuritySchemes = convertSecurityDefinitionsV2ToV3(src.SecurityDefinitions)

	if src.Responses != nil {
		result.Components.Responses = convertResponsesMapV2ToV3(src.Responses.Definitions)
	}

	if src.Definitions != nil && src.Definitions.Definitions != nil {
		result.Components.Schemas = src.Definitions.Definitions
	}

	result.Paths = convertPathsV2ToV3(src.Paths)

	return result
}

func convertPathsV2ToV3(paths *v2high.Paths) *v3high.Paths {
	if paths == nil {
		return nil
	}

	result := &v3high.Paths{
		Extensions: paths.Extensions,
	}

	if paths.PathItems == nil {
		return result
	}

	result.PathItems = orderedmap.New[string, *v3high.PathItem]()

	for iter := paths.PathItems.Oldest(); iter != nil; iter = iter.Next() {
		if iter.Value == nil {
			continue
		}

		item := convertPathItemV2ToV3(iter.Value)
		result.PathItems.Set(iter.Key, item)
	}

	return result
}

func convertPathItemV2ToV3(pathItem *v2high.PathItem) *v3high.PathItem {
	result := &v3high.PathItem{
		Extensions: pathItem.Extensions,
		Reference:  pathItem.Ref,
		Parameters: make([]*v3high.Parameter, len(pathItem.Parameters)),
	}

	for i, param := range pathItem.Parameters {
		item := convertParameterV2ToV3(param)
		result.Parameters[i] = item
	}

	result.Get = convertOperationV2ToV3(pathItem.Get)
	result.Post = convertOperationV2ToV3(pathItem.Post)
	result.Put = convertOperationV2ToV3(pathItem.Put)
	result.Delete = convertOperationV2ToV3(pathItem.Delete)
	result.Patch = convertOperationV2ToV3(pathItem.Patch)
	result.Head = convertOperationV2ToV3(pathItem.Head)
	result.Options = convertOperationV2ToV3(pathItem.Options)

	return result
}

func convertOperationV2ToV3(operation *v2high.Operation) *v3high.Operation {
	if operation == nil {
		return nil
	}

	// Consumes     []string
	// Produces     []string
	// Responses    *Responses

	result := &v3high.Operation{
		Tags:         operation.Tags,
		Extensions:   operation.Extensions,
		Summary:      operation.Summary,
		Description:  operation.Description,
		ExternalDocs: operation.ExternalDocs,
		OperationId:  operation.OperationId,
		Deprecated:   &operation.Deprecated,
		Security:     operation.Security,
	}

	for i, param := range operation.Parameters {
		item := convertParameterV2ToV3(param)
		result.Parameters[i] = item
	}

	if operation.Responses == nil {
		return result
	}

	result.Responses = &v3high.Responses{
		Extensions: operation.Responses.Extensions,
		Default:    convertResponseV2ToV3(operation.Responses.Default),
		Codes:      convertResponsesMapV2ToV3(operation.Responses.Codes),
	}

	return result
}

func convertResponsesMapV2ToV3(
	responses *orderedmap.Map[string, *v2high.Response],
) *orderedmap.Map[string, *v3high.Response] {
	if responses == nil || responses.Len() == 0 {
		return nil
	}

	result := orderedmap.New[string, *v3high.Response]()

	for iter := responses.Oldest(); iter != nil; iter = iter.Next() {
		resp := convertResponseV2ToV3(iter.Value)
		if resp != nil {
			result.Set(iter.Key, resp)
		}
	}

	return result
}

func convertResponseV2ToV3(value *v2high.Response) *v3high.Response {
	if value == nil {
		return nil
	}

	result := &v3high.Response{
		Description: value.Description,
		Extensions:  value.Extensions,
		Headers:     convertHeadersV2ToV3(value.Headers),
	}

	// TODO
	// Schema      *base.SchemaProxy

	return result
}

func convertHeadersV2ToV3(
	value *orderedmap.Map[string, *v2high.Header],
) *orderedmap.Map[string, *v3high.Header] {
	if value == nil || value.Len() == 0 {
		return nil
	}

	result := orderedmap.New[string, *v3high.Header]()

	for iter := value.Oldest(); iter != nil; iter = iter.Next() {
		if iter.Value == nil {
			continue
		}

		header := convertHeaderV2ToV3(iter.Value)
		result.Set(iter.Key, header)
	}

	return result
}

func convertHeaderV2ToV3(value *v2high.Header) *v3high.Header {
	baseSchema := &base.Schema{
		Format:      value.Format,
		Pattern:     value.Pattern,
		UniqueItems: &value.UniqueItems,
	}

	if value.Default != nil {
		baseSchema.Default = utils.CreateYamlNode(value.Default)
	}

	if value.Type != "" {
		baseSchema.Type = []string{value.Type}
	}

	if len(value.Enum) > 0 {
		baseSchema.Enum = make([]*yaml.Node, len(value.Enum))

		for i, enum := range value.Enum {
			baseSchema.Enum[i] = utils.CreateYamlNode(enum)
		}
	}

	maximum := float64(value.Maximum)
	baseSchema.Maximum = &maximum
	minimum := float64(value.Minimum)
	baseSchema.Minimum = &minimum
	multipleOf := float64(value.MultipleOf)
	baseSchema.MultipleOf = &multipleOf
	baseSchema.ExclusiveMinimum = &base.DynamicValue[bool, float64]{
		A: value.ExclusiveMinimum,
	}

	baseSchema.ExclusiveMaximum = &base.DynamicValue[bool, float64]{
		A: value.ExclusiveMaximum,
	}

	if value.MaxLength > 0 {
		maxLength := int64(value.MaxLength)
		baseSchema.MaxLength = &maxLength
	}

	if value.MinLength > 0 {
		minLength := int64(value.MinLength)
		baseSchema.MaxLength = &minLength
	}

	if value.MaxItems > 0 {
		maxItems := int64(value.MaxItems)
		baseSchema.MaxItems = &maxItems
	}

	if value.MinItems > 0 {
		minItems := int64(value.MinItems)
		baseSchema.MinItems = &minItems
	}

	if value.Items != nil {
		baseSchema.Items = convertItemsV2ToV3(value.Items)
	}

	result := &v3high.Header{
		Description: value.Description,
		Extensions:  value.Extensions,
		Schema:      base.CreateSchemaProxy(baseSchema),
	}

	return result
}

func convertSecurityDefinitionsV2ToV3(
	securityDefinitions *v2high.SecurityDefinitions,
) *orderedmap.Map[string, *v3high.SecurityScheme] {
	if securityDefinitions == nil || securityDefinitions.Definitions == nil ||
		securityDefinitions.Definitions.Len() == 0 {
		return nil
	}

	result := orderedmap.New[string, *v3high.SecurityScheme]()

	for iter := securityDefinitions.Definitions.Oldest(); iter != nil; iter = iter.Next() {
		if iter.Value == nil {
			continue
		}

		security := convertSecuritySchemeV2ToV3(iter.Value)
		result.Set(iter.Key, security)
	}

	return result
}

func convertSecuritySchemeV2ToV3(value *v2high.SecurityScheme) *v3high.SecurityScheme {
	result := &v3high.SecurityScheme{
		Description: value.Description,
		Name:        value.Name,
		Extensions:  value.Extensions,
		In:          value.In,
		Type:        value.Type,
	}

	switch value.Type {
	case string(BasicAuthScheme):
		result.Type = string(HTTPAuthScheme)
		result.Scheme = string(BasicAuthScheme)
	case string(OAuth2Scheme):
		flow := &v3high.OAuthFlow{
			AuthorizationUrl: value.AuthorizationUrl,
			TokenUrl:         value.TokenUrl,
			Scopes:           value.Scopes.Values,
		}

		if value.Scopes != nil && value.Scopes.Values != nil {
			flow.Scopes = value.Scopes.Values
		}

		result.Flows = &v3high.OAuthFlows{}

		switch value.Flow {
		case "accessCode":
			result.Flows.AuthorizationCode = flow
		case "implicit":
			result.Flows.Implicit = flow
		case "password":
			result.Flows.Password = flow
		case "application":
			result.Flows.ClientCredentials = flow
		default:
		}
	default:
	}

	return result
}

func convertParametersV2ToV3(
	parameters *v2high.ParameterDefinitions,
) *orderedmap.Map[string, *v3high.Parameter] {
	if parameters == nil || parameters.Definitions == nil && parameters.Definitions.Len() == 0 {
		return nil
	}

	result := orderedmap.New[string, *v3high.Parameter]()

	for iter := parameters.Definitions.Oldest(); iter != nil; iter = iter.Next() {
		if iter.Value == nil {
			continue
		}

		param := convertParameterV2ToV3(iter.Value)
		result.Set(iter.Key, param)
	}

	return result
}

func convertParameterV2ToV3(parameter *v2high.Parameter) *v3high.Parameter {
	param := &v3high.Parameter{
		Name:        parameter.Name,
		In:          parameter.In,
		Description: parameter.Description,
		Required:    parameter.Required,
		Schema:      parameter.Schema,
		Example:     parameter.Default,
		Extensions:  parameter.Extensions,
	}

	if parameter.AllowEmptyValue != nil {
		param.AllowEmptyValue = *parameter.AllowEmptyValue
	}

	if param.Schema != nil {
		return param
	}

	// CollectionFormat string
	baseSchema := &base.Schema{
		Format:      parameter.Format,
		Pattern:     parameter.Pattern,
		Default:     parameter.Default,
		UniqueItems: parameter.UniqueItems,
		Enum:        parameter.Enum,
	}

	if parameter.Type != "" {
		baseSchema.Type = []string{parameter.Type}
	}

	if parameter.Maximum != nil {
		maximum := float64(*parameter.Maximum)

		baseSchema.Maximum = &maximum
	}

	if parameter.Minimum != nil {
		minimum := float64(*parameter.Minimum)

		baseSchema.Minimum = &minimum
	}

	if parameter.MultipleOf != nil {
		multipleOf := float64(*parameter.MultipleOf)

		baseSchema.MultipleOf = &multipleOf
	}

	if parameter.ExclusiveMinimum != nil {
		baseSchema.ExclusiveMinimum = &base.DynamicValue[bool, float64]{
			A: *parameter.ExclusiveMinimum,
		}
	}

	if parameter.ExclusiveMaximum != nil {
		baseSchema.ExclusiveMaximum = &base.DynamicValue[bool, float64]{
			A: *parameter.ExclusiveMaximum,
		}
	}

	if parameter.MaxLength != nil && *parameter.MaxLength > 0 {
		maxLength := int64(*parameter.MaxLength)
		baseSchema.MaxLength = &maxLength
	}

	if parameter.MinLength != nil && *parameter.MinLength > 0 {
		minLength := int64(*parameter.MinLength)
		baseSchema.MaxLength = &minLength
	}

	if parameter.MaxItems != nil && *parameter.MaxItems > 0 {
		maxItems := int64(*parameter.MaxItems)
		baseSchema.MaxItems = &maxItems
	}

	if parameter.MinItems != nil && *parameter.MaxItems > 0 {
		minItems := int64(*parameter.MinItems)
		baseSchema.MinItems = &minItems
	}

	if parameter.Items != nil {
		baseSchema.Items = convertItemsV2ToV3(parameter.Items)
	}

	param.Schema = base.CreateSchemaProxy(baseSchema)

	return param
}

func convertItemsV2ToV3(items *v2high.Items) *base.DynamicValue[*base.SchemaProxy, bool] {
	baseSchema := &base.Schema{
		Format:      items.Format,
		Pattern:     items.Pattern,
		Default:     items.Default,
		UniqueItems: &items.UniqueItems,
		Enum:        items.Enum,
	}

	if items.Type != "" {
		baseSchema.Type = []string{items.Type}
	}

	maximum := float64(items.Maximum)
	baseSchema.Maximum = &maximum
	minimum := float64(items.Minimum)
	baseSchema.Minimum = &minimum
	multipleOf := float64(items.MultipleOf)
	baseSchema.MultipleOf = &multipleOf
	baseSchema.ExclusiveMinimum = &base.DynamicValue[bool, float64]{
		A: items.ExclusiveMinimum,
	}

	baseSchema.ExclusiveMaximum = &base.DynamicValue[bool, float64]{
		A: items.ExclusiveMaximum,
	}

	if items.MaxLength > 0 {
		maxLength := int64(items.MaxLength)
		baseSchema.MaxLength = &maxLength
	}

	if items.MinLength > 0 {
		minLength := int64(items.MinLength)
		baseSchema.MaxLength = &minLength
	}

	if items.MaxItems > 0 {
		maxItems := int64(items.MaxItems)
		baseSchema.MaxItems = &maxItems
	}

	if items.MinItems > 0 {
		minItems := int64(items.MinItems)
		baseSchema.MinItems = &minItems
	}

	if items.Items != nil {
		baseSchema.Items = convertItemsV2ToV3(items.Items)
	}

	return &base.DynamicValue[*base.SchemaProxy, bool]{
		N: 0,
		A: base.CreateSchemaProxy(baseSchema),
	}
}
