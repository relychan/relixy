package openapi

// SecurityRequirement wraps the security requirement specification.
type SecurityRequirement map[string][]string

// NewSecurityRequirement creates a SecurityRequirement instance from name and scope.
func NewSecurityRequirement(name string, scopes []string) SecurityRequirement {
	return SecurityRequirement{
		name: scopes,
	}
}

// Name returns the name of security requirement.
func (as SecurityRequirement) Name() string {
	if len(as) > 0 {
		for k := range as {
			return k
		}
	}

	return ""
}

// Scopes returns scopes of security requirement.
func (as SecurityRequirement) Scopes() []string {
	if len(as) > 0 {
		for _, scopes := range as {
			return scopes
		}
	}

	return []string{}
}

// IsOptional checks if the security is optional.
func (as SecurityRequirement) IsOptional() bool {
	return len(as) == 0
}

// SecurityRequirements wraps list of security requirements with helpers.
type SecurityRequirements []SecurityRequirement //nolint:recvcheck

// IsEmpty checks if there is no security.
func (ass SecurityRequirements) IsEmpty() bool {
	return len(ass) == 0
}

// IsOptional checks if the security is optional.
func (ass SecurityRequirements) IsOptional() bool {
	if ass.IsEmpty() {
		return true
	}

	for _, as := range ass {
		if as.IsOptional() {
			return true
		}
	}

	return false
}

// Add adds a security with name and scope.
func (ass *SecurityRequirements) Add(item SecurityRequirement) {
	*ass = append(*ass, item)
}

// Get gets a security by name.
func (ass SecurityRequirements) Get(name string) SecurityRequirement {
	for _, as := range ass {
		if as.Name() == name {
			return as
		}
	}

	return nil
}

// First returns the first security.
func (ass SecurityRequirements) First() SecurityRequirement {
	for _, as := range ass {
		return as
	}

	return nil
}
