package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Version
		hasError bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: Version{},
			hasError: false,
		},
		{
			name:     "major only",
			input:    "1",
			expected: Version{Major: 1, Minor: 0, Patch: 0, Build: 0},
			hasError: false,
		},
		{
			name:     "major.minor",
			input:    "1.2",
			expected: Version{Major: 1, Minor: 2, Patch: 0, Build: 0},
			hasError: false,
		},
		{
			name:     "major.minor.patch",
			input:    "1.2.3",
			expected: Version{Major: 1, Minor: 2, Patch: 3, Build: 0},
			hasError: false,
		},
		{
			name:     "major.minor.patch.build",
			input:    "1.2.3.4",
			expected: Version{Major: 1, Minor: 2, Patch: 3, Build: 4},
			hasError: false,
		},
		{
			name:     "too many parts",
			input:    "1.2.3.4.5",
			expected: Version{},
			hasError: true,
		},
		{
			name:     "invalid number",
			input:    "1.2.abc",
			expected: Version{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseVersion(tt.input)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestVersionIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		version  Version
		expected bool
	}{
		{
			name:     "empty version",
			version:  Version{},
			expected: true,
		},
		{
			name:     "non-empty version",
			version:  Version{Major: 1},
			expected: false,
		},
		{
			name:     "only build set",
			version:  Version{Build: 1},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.version.isEmpty()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVersionUpdated(t *testing.T) {
	baseVersion := Version{Major: 1, Minor: 2, Patch: 3, Build: 4}

	tests := []struct {
		name     string
		level    VersionLevel
		expected Version
	}{
		{
			name:     "major update",
			level:    LevelMajor,
			expected: Version{Major: 2, Minor: 0, Patch: 0, Build: 0},
		},
		{
			name:     "minor update",
			level:    LevelMinor,
			expected: Version{Major: 1, Minor: 3, Patch: 0, Build: 0},
		},
		{
			name:     "patch update",
			level:    LevelPatch,
			expected: Version{Major: 1, Minor: 2, Patch: 4, Build: 0},
		},
		{
			name:     "build update",
			level:    LevelBuild,
			expected: Version{Major: 1, Minor: 2, Patch: 3, Build: 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := baseVersion.Updated(tt.level)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVersionString(t *testing.T) {
	version := Version{Major: 1, Minor: 2, Patch: 3, Build: 4}

	tests := []struct {
		name     string
		notation VersionNotation
		expected string
	}{
		{
			name:     "simple notation",
			notation: NotationSimple,
			expected: "1.2",
		},
		{
			name:     "normal notation",
			notation: NotationNormal,
			expected: "1.2.3",
		},
		{
			name:     "detail notation",
			notation: NotationDetail,
			expected: "1.2.3.4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := version.String(tt.notation)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVersionLevelConstants(t *testing.T) {
	assert.Equal(t, VersionLevel("major"), LevelMajor)
	assert.Equal(t, VersionLevel("minor"), LevelMinor)
	assert.Equal(t, VersionLevel("patch"), LevelPatch)
	assert.Equal(t, VersionLevel("build"), LevelBuild)
}

func TestVersionNotationConstants(t *testing.T) {
	assert.Equal(t, VersionNotation("simple"), NotationSimple)
	assert.Equal(t, VersionNotation("normal"), NotationNormal)
	assert.Equal(t, VersionNotation("detail"), NotationDetail)
}

func TestVersionTargetConstants(t *testing.T) {
	assert.Equal(t, VersionTarget("both"), TargetBoth)
	assert.Equal(t, VersionTarget("file"), TargetFile)
	assert.Equal(t, VersionTarget("product"), TargetProduct)
}
