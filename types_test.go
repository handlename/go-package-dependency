package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test String methods for custom types
func TestLayerName_String(t *testing.T) {
	ln := LayerName("test-layer")
	expected := "test-layer"
	assert.Equal(t, expected, ln.String())
}

func TestLayerPath_String(t *testing.T) {
	lp := LayerPath("path/to/layer")
	expected := "path/to/layer"
	assert.Equal(t, expected, lp.String())
}

func TestModuleName_String(t *testing.T) {
	mn := ModuleName("github.com/example/module")
	expected := "github.com/example/module"
	assert.Equal(t, expected, mn.String())
}

func TestPackageName_String(t *testing.T) {
	pn := PackageName("main")
	expected := "main"
	assert.Equal(t, expected, pn.String())
}

func TestFilePath_String(t *testing.T) {
	fp := FilePath("path/to/file.go")
	expected := "path/to/file.go"
	assert.Equal(t, expected, fp.String())
}

// Test LayerName validation
func TestLayerName_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		input    LayerName
		expected bool
	}{
		{"valid layer name", LayerName("valid-layer"), true},
		{"empty string", LayerName(""), false},
		{"whitespace only", LayerName("   "), false},
		{"single character", LayerName("a"), true},
		{"with spaces", LayerName("layer with spaces"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsValid())
		})
	}
}

func TestLayerName_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     LayerName
		expectErr bool
	}{
		{"valid layer name", LayerName("valid-layer"), false},
		{"empty string", LayerName(""), true},
		{"whitespace only", LayerName("   "), true},
		{"single character", LayerName("a"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test LayerPath validation
func TestLayerPath_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		input    LayerPath
		expected bool
	}{
		{"valid path", LayerPath("path/to/layer"), true},
		{"empty string", LayerPath(""), false},
		{"whitespace only", LayerPath("   "), false},
		{"path with parent directory", LayerPath("path/../other"), false},
		{"path with double dots", LayerPath("path/.."), false},
		{"single dot", LayerPath("."), true},
		{"root path", LayerPath("/"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsValid())
		})
	}
}

func TestLayerPath_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     LayerPath
		expectErr bool
	}{
		{"valid path", LayerPath("path/to/layer"), false},
		{"empty string", LayerPath(""), true},
		{"whitespace only", LayerPath("   "), true},
		{"path with parent directory", LayerPath("path/../other"), true},
		{"path with double dots", LayerPath("path/.."), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test ModuleName validation
func TestModuleName_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		input    ModuleName
		expected bool
	}{
		{"valid module name", ModuleName("github.com/example/module"), true},
		{"empty string", ModuleName(""), false},
		{"whitespace only", ModuleName("   "), false},
		{"module with spaces", ModuleName("module with spaces"), false},
		{"single word", ModuleName("module"), true},
		{"with dashes", ModuleName("my-module"), true},
		{"with underscores", ModuleName("my_module"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsValid())
		})
	}
}

func TestModuleName_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     ModuleName
		expectErr bool
	}{
		{"valid module name", ModuleName("github.com/example/module"), false},
		{"empty string", ModuleName(""), true},
		{"whitespace only", ModuleName("   "), true},
		{"module with spaces", ModuleName("module with spaces"), true},
		{"single word", ModuleName("module"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test PackageName validation
func TestPackageName_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		input    PackageName
		expected bool
	}{
		{"valid package name", PackageName("main"), true},
		{"empty string", PackageName(""), false},
		{"whitespace only", PackageName("   "), false},
		{"package with slash", PackageName("package/name"), false},
		{"package with spaces", PackageName("package name"), true},
		{"single character", PackageName("p"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsValid())
		})
	}
}

func TestPackageName_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     PackageName
		expectErr bool
	}{
		{"valid package name", PackageName("main"), false},
		{"empty string", PackageName(""), true},
		{"whitespace only", PackageName("   "), true},
		{"package with slash", PackageName("package/name"), true},
		{"package with spaces", PackageName("package name"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test custom error types
func TestUnsupportedReaderError(t *testing.T) {
	err := UnsupportedReaderError{ReaderType: "unknown"}
	expected := "unsupported reader type: unknown"
	assert.Equal(t, expected, err.Error())
}

func TestDirectoryCreationError(t *testing.T) {
	innerErr := errors.New("permission denied")
	err := DirectoryCreationError{Path: "/tmp/test", Err: innerErr}
	expected := "failed to create directory /tmp/test: permission denied"
	assert.Equal(t, expected, err.Error())
}

func TestFileWriteError(t *testing.T) {
	innerErr := errors.New("disk full")
	err := FileWriteError{Path: "/tmp/file.txt", Err: innerErr}
	expected := "failed to write /tmp/file.txt: disk full"
	assert.Equal(t, expected, err.Error())
}

func TestFileFormatError(t *testing.T) {
	innerErr := errors.New("invalid syntax")
	err := FileFormatError{Path: "/tmp/code.go", Err: innerErr}
	expected := "failed to format /tmp/code.go: invalid syntax"
	assert.Equal(t, expected, err.Error())
}

func TestModuleNotFoundError(t *testing.T) {
	err := ModuleNotFoundError{Source: "go.mod"}
	expected := "module declaration not found in go.mod"
	assert.Equal(t, expected, err.Error())
}

// Test structs
func TestLayer(t *testing.T) {
	layer := Layer{
		Name: LayerName("test-layer"),
		Path: LayerPath("path/to/layer"),
	}

	assert.Equal(t, "test-layer", layer.Name.String())
	assert.Equal(t, "path/to/layer", layer.Path.String())
}

func TestDependencyConfig(t *testing.T) {
	config := DependencyConfig{
		Layers: []Layer{
			{Name: LayerName("layer1"), Path: LayerPath("path1")},
			{Name: LayerName("layer2"), Path: LayerPath("path2")},
		},
		Dependencies: map[LayerName][]LayerName{
			LayerName("layer1"): {LayerName("layer2")},
		},
	}

	assert.Len(t, config.Layers, 2)
	assert.Len(t, config.Dependencies, 1)

	deps, exists := config.Dependencies[LayerName("layer1")]
	assert.True(t, exists, "Expected dependency for layer1 to exist")
	assert.Len(t, deps, 1)
	assert.Equal(t, LayerName("layer2"), deps[0])
}
