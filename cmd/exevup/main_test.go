package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/simp7/goversioninfo-toolkit/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVersionInfoFromFile(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_versioninfo.json")

	validContent := `{
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

	t.Run("valid file", func(t *testing.T) {
		err := os.WriteFile(testFile, []byte(validContent), 0644)
		require.NoError(t, err)

		info, err := parseVersionInfoFromFile(testFile)
		require.NoError(t, err)
		assert.Equal(t, 1, info.FixedFileInfo.FileVersion.Major)
		assert.Equal(t, "1.2.3.4", info.StringFileInfo.FileVersion)
	})

	t.Run("non-existent file creates empty file", func(t *testing.T) {
		nonExistentFile := filepath.Join(tempDir, "non_existent.json")

		// Write empty JSON object to make parsing successful
		err := os.WriteFile(nonExistentFile, []byte("{}"), 0644)
		require.NoError(t, err)

		info, err := parseVersionInfoFromFile(nonExistentFile)
		require.NoError(t, err)

		// File should be created but empty, so version info should be empty
		assert.Equal(t, 0, info.FixedFileInfo.FileVersion.Major)

		// Clean up
		os.Remove(nonExistentFile)
	})

	t.Run("invalid JSON file", func(t *testing.T) {
		invalidFile := filepath.Join(tempDir, "invalid.json")
		err := os.WriteFile(invalidFile, []byte("{invalid json"), 0644)
		require.NoError(t, err)

		_, err = parseVersionInfoFromFile(invalidFile)
		assert.Error(t, err)
	})
}

func TestOverwriteVersionInfoToFile(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "output_test.json")

	// Create a minimal valid info object
	info := model.Info{}

	t.Run("write to new file", func(t *testing.T) {
		err := overwriteVersionInfoToFile(testFile, info)
		assert.NoError(t, err)

		// Verify file was created
		_, err = os.Stat(testFile)
		assert.NoError(t, err)
	})

	t.Run("overwrite existing file", func(t *testing.T) {
		// Write some content first
		err := os.WriteFile(testFile, []byte("existing content"), 0644)
		require.NoError(t, err)

		err = overwriteVersionInfoToFile(testFile, info)
		assert.NoError(t, err)

		// Read back and verify it was overwritten
		content, err := os.ReadFile(testFile)
		require.NoError(t, err)
		assert.NotContains(t, string(content), "existing content")
	})
}

func TestIntegrationVersionUpdate(t *testing.T) {
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "integration_test.json")
	outputFile := filepath.Join(tempDir, "integration_output.json")

	// Create initial version info file
	initialContent := `{
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
			"FileVersion": "1.2.3",
			"ProductVersion": "1.2.3"
		}
	}`

	err := os.WriteFile(inputFile, []byte(initialContent), 0644)
	require.NoError(t, err)

	t.Run("complete workflow", func(t *testing.T) {
		// Parse input file
		info, err := parseVersionInfoFromFile(inputFile)
		require.NoError(t, err)

		// Get versions
		fileVersion, err := info.GetFileVersion()
		require.NoError(t, err)
		productVersion, err := info.GetProductVersion()
		require.NoError(t, err)

		// Update versions (patch level)
		updatedFileVersion := fileVersion.Updated("patch")
		updatedProductVersion := productVersion.Updated("patch")

		// Apply updates
		updatedInfo := info.VersionUpdated(updatedFileVersion, updatedProductVersion, "both", "detail")

		// Write to output file
		err = overwriteVersionInfoToFile(outputFile, updatedInfo)
		require.NoError(t, err)

		// Verify the update
		resultInfo, err := parseVersionInfoFromFile(outputFile)
		require.NoError(t, err)

		resultFileVersion, err := resultInfo.GetFileVersion()
		require.NoError(t, err)
		resultProductVersion, err := resultInfo.GetProductVersion()
		require.NoError(t, err)

		// Check that patch was incremented and build was reset
		assert.Equal(t, 1, resultFileVersion.Major)
		assert.Equal(t, 2, resultFileVersion.Minor)
		assert.Equal(t, 4, resultFileVersion.Patch) // 3 + 1
		assert.Equal(t, 0, resultFileVersion.Build) // Reset to 0

		assert.Equal(t, 1, resultProductVersion.Major)
		assert.Equal(t, 2, resultProductVersion.Minor)
		assert.Equal(t, 4, resultProductVersion.Patch) // 3 + 1
		assert.Equal(t, 0, resultProductVersion.Build) // Reset to 0
	})
}
