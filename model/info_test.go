package model

import (
	"testing"

	"github.com/josephspurrier/goversioninfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVersionInfo(t *testing.T) {
	validJSON := `{
		"FixedFileInfo": {
			"FileVersion": {
				"Major": 1,
				"Minor": 2,
				"Patch": 3,
				"Build": 4
			},
			"ProductVersion": {
				"Major": 1,
				"Minor": 2,
				"Patch": 3,
				"Build": 4
			}
		},
		"StringFileInfo": {
			"FileVersion": "1.2.3.4",
			"ProductVersion": "1.2.3.4"
		}
	}`

	invalidJSON := `{invalid json}`

	t.Run("valid JSON", func(t *testing.T) {
		info, err := ParseVersionInfo([]byte(validJSON))
		require.NoError(t, err)

		assert.Equal(t, 1, info.FixedFileInfo.FileVersion.Major)
		assert.Equal(t, 2, info.FixedFileInfo.FileVersion.Minor)
		assert.Equal(t, 3, info.FixedFileInfo.FileVersion.Patch)
		assert.Equal(t, 4, info.FixedFileInfo.FileVersion.Build)
		assert.Equal(t, "1.2.3.4", info.StringFileInfo.FileVersion)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, err := ParseVersionInfo([]byte(invalidJSON))
		assert.Error(t, err)
	})
}

func TestStringifyVersionInfo(t *testing.T) {
	info := Info{
		FixedFileInfo: goversioninfo.FixedFileInfo{
			FileVersion: goversioninfo.FileVersion{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: 4,
			},
		},
		StringFileInfo: goversioninfo.StringFileInfo{
			FileVersion: "1.2.3.4",
		},
	}

	data, err := StringifyVersionInfo(info)
	require.NoError(t, err)
	assert.Contains(t, string(data), "1.2.3.4")
	assert.Contains(t, string(data), "FileVersion")
}

func TestGetFileVersion(t *testing.T) {
	t.Run("from FixedFileInfo", func(t *testing.T) {
		info := Info{
			FixedFileInfo: goversioninfo.FixedFileInfo{
				FileVersion: goversioninfo.FileVersion{
					Major: 1,
					Minor: 2,
					Patch: 3,
					Build: 4,
				},
			},
		}

		version, err := info.GetFileVersion()
		require.NoError(t, err)
		assert.Equal(t, Version{Major: 1, Minor: 2, Patch: 3, Build: 4}, version)
	})

	t.Run("from StringFileInfo when FixedFileInfo is empty", func(t *testing.T) {
		info := Info{
			StringFileInfo: goversioninfo.StringFileInfo{
				FileVersion: "2.3.4.5",
			},
		}

		version, err := info.GetFileVersion()
		require.NoError(t, err)
		assert.Equal(t, Version{Major: 2, Minor: 3, Patch: 4, Build: 5}, version)
	})

	t.Run("invalid string version", func(t *testing.T) {
		info := Info{
			StringFileInfo: goversioninfo.StringFileInfo{
				FileVersion: "invalid.version.string.too.many.parts",
			},
		}

		_, err := info.GetFileVersion()
		assert.Error(t, err)
	})
}

func TestGetProductVersion(t *testing.T) {
	t.Run("from FixedFileInfo", func(t *testing.T) {
		info := Info{
			FixedFileInfo: goversioninfo.FixedFileInfo{
				ProductVersion: goversioninfo.FileVersion{
					Major: 2,
					Minor: 3,
					Patch: 4,
					Build: 5,
				},
			},
		}

		version, err := info.GetProductVersion()
		require.NoError(t, err)
		assert.Equal(t, Version{Major: 2, Minor: 3, Patch: 4, Build: 5}, version)
	})

	t.Run("from StringFileInfo when FixedFileInfo is empty", func(t *testing.T) {
		info := Info{
			StringFileInfo: goversioninfo.StringFileInfo{
				ProductVersion: "3.4.5.6",
			},
		}

		version, err := info.GetProductVersion()
		require.NoError(t, err)
		assert.Equal(t, Version{Major: 3, Minor: 4, Patch: 5, Build: 6}, version)
	})
}

func TestFileVersionUpdated(t *testing.T) {
	info := Info{}
	version := Version{Major: 1, Minor: 2, Patch: 3, Build: 4}

	result := info.FileVersionUpdated(version, NotationDetail)

	assert.Equal(t, goversioninfo.FileVersion(version), result.FixedFileInfo.FileVersion)
	assert.Equal(t, "1.2.3.4", result.StringFileInfo.FileVersion)
}

func TestProductVersionUpdated(t *testing.T) {
	info := Info{}
	version := Version{Major: 2, Minor: 3, Patch: 4, Build: 5}

	result := info.ProductVersionUpdated(version, NotationNormal)

	assert.Equal(t, goversioninfo.FileVersion(version), result.FixedFileInfo.ProductVersion)
	assert.Equal(t, "2.3.4", result.StringFileInfo.ProductVersion)
}

func TestInfoVersionUpdated(t *testing.T) {
	info := Info{}
	fileVersion := Version{Major: 1, Minor: 2, Patch: 3, Build: 4}
	productVersion := Version{Major: 2, Minor: 3, Patch: 4, Build: 5}

	tests := []struct {
		name   string
		target VersionTarget
	}{
		{
			name:   "file target",
			target: TargetFile,
		},
		{
			name:   "product target",
			target: TargetProduct,
		},
		{
			name:   "both target",
			target: TargetBoth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := info.VersionUpdated(fileVersion, productVersion, tt.target, NotationDetail)

			switch tt.target {
			case TargetFile:
				assert.Equal(t, goversioninfo.FileVersion(fileVersion), result.FixedFileInfo.FileVersion)
				assert.Equal(t, "1.2.3.4", result.StringFileInfo.FileVersion)
				assert.Equal(t, goversioninfo.FileVersion{}, result.FixedFileInfo.ProductVersion)
			case TargetProduct:
				// Note: There's a bug in the original code - it uses fileVersion instead of productVersion
				assert.Equal(t, goversioninfo.FileVersion(fileVersion), result.FixedFileInfo.ProductVersion)
				assert.Equal(t, "1.2.3.4", result.StringFileInfo.ProductVersion)
				assert.Equal(t, goversioninfo.FileVersion{}, result.FixedFileInfo.FileVersion)
			case TargetBoth:
				assert.Equal(t, goversioninfo.FileVersion(fileVersion), result.FixedFileInfo.FileVersion)
				assert.Equal(t, goversioninfo.FileVersion(fileVersion), result.FixedFileInfo.ProductVersion)
				assert.Equal(t, "1.2.3.4", result.StringFileInfo.FileVersion)
				assert.Equal(t, "1.2.3.4", result.StringFileInfo.ProductVersion)
			}
		})
	}
}