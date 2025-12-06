package schema

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestNewSecurityRequirement(t *testing.T) {
	scopes := []string{"read", "write"}
	sr := NewSecurityRequirement("oauth2", scopes)

	assert.Equal(t, "oauth2", sr.Name())
	assert.DeepEqual(t, scopes, sr.Scopes())
}

func TestSecurityRequirement_Name(t *testing.T) {
	testCases := []struct {
		name     string
		sr       SecurityRequirement
		expected string
	}{
		{
			name:     "with name",
			sr:       SecurityRequirement{"apiKey": []string{}},
			expected: "apiKey",
		},
		{
			name:     "empty requirement",
			sr:       SecurityRequirement{},
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.sr.Name()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSecurityRequirement_Scopes(t *testing.T) {
	testCases := []struct {
		name     string
		sr       SecurityRequirement
		expected []string
	}{
		{
			name:     "with scopes",
			sr:       SecurityRequirement{"oauth2": []string{"read", "write"}},
			expected: []string{"read", "write"},
		},
		{
			name:     "empty scopes",
			sr:       SecurityRequirement{"apiKey": []string{}},
			expected: []string{},
		},
		{
			name:     "empty requirement",
			sr:       SecurityRequirement{},
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.sr.Scopes()
			assert.DeepEqual(t, tc.expected, result)
		})
	}
}

func TestSecurityRequirement_IsOptional(t *testing.T) {
	testCases := []struct {
		name     string
		sr       SecurityRequirement
		expected bool
	}{
		{
			name:     "not optional",
			sr:       SecurityRequirement{"apiKey": []string{}},
			expected: false,
		},
		{
			name:     "optional (empty)",
			sr:       SecurityRequirement{},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.sr.IsOptional()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSecurityRequirements_IsEmpty(t *testing.T) {
	testCases := []struct {
		name     string
		srs      SecurityRequirements
		expected bool
	}{
		{
			name:     "empty",
			srs:      SecurityRequirements{},
			expected: true,
		},
		{
			name: "not empty",
			srs: SecurityRequirements{
				{"apiKey": []string{}},
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.srs.IsEmpty()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSecurityRequirements_IsOptional(t *testing.T) {
	testCases := []struct {
		name     string
		srs      SecurityRequirements
		expected bool
	}{
		{
			name:     "empty is optional",
			srs:      SecurityRequirements{},
			expected: true,
		},
		{
			name: "contains optional requirement",
			srs: SecurityRequirements{
				{"apiKey": []string{}},
				{},
			},
			expected: true,
		},
		{
			name: "all required",
			srs: SecurityRequirements{
				{"apiKey": []string{}},
				{"oauth2": []string{"read"}},
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.srs.IsOptional()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSecurityRequirements_Add(t *testing.T) {
	srs := SecurityRequirements{}
	sr1 := SecurityRequirement{"apiKey": []string{}}
	sr2 := SecurityRequirement{"oauth2": []string{"read", "write"}}

	srs.Add(sr1)
	assert.Equal(t, 1, len(srs))
	assert.DeepEqual(t, sr1, srs[0])

	srs.Add(sr2)
	assert.Equal(t, 2, len(srs))
	assert.DeepEqual(t, sr2, srs[1])
}

func TestSecurityRequirements_Get(t *testing.T) {
	srs := SecurityRequirements{
		{"apiKey": []string{}},
		{"oauth2": []string{"read", "write"}},
		{"basic": []string{}},
	}

	testCases := []struct {
		name     string
		key      string
		expected SecurityRequirement
	}{
		{
			name:     "found apiKey",
			key:      "apiKey",
			expected: SecurityRequirement{"apiKey": []string{}},
		},
		{
			name:     "found oauth2",
			key:      "oauth2",
			expected: SecurityRequirement{"oauth2": []string{"read", "write"}},
		},
		{
			name:     "not found",
			key:      "notfound",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := srs.Get(tc.key)
			assert.DeepEqual(t, tc.expected, result)
		})
	}
}

func TestSecurityRequirements_First(t *testing.T) {
	testCases := []struct {
		name     string
		srs      SecurityRequirements
		expected SecurityRequirement
	}{
		{
			name: "returns first",
			srs: SecurityRequirements{
				{"apiKey": []string{}},
				{"oauth2": []string{"read"}},
			},
			expected: SecurityRequirement{"apiKey": []string{}},
		},
		{
			name:     "empty returns nil",
			srs:      SecurityRequirements{},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.srs.First()
			assert.DeepEqual(t, tc.expected, result)
		})
	}
}

