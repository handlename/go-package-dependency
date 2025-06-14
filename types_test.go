package main

import (
	"errors"
	"testing"
)

// Test String methods for custom types
func TestLayerName_String(t *testing.T) {
	ln := LayerName("test-layer")
	expected := "test-layer"
	if ln.String() != expected {
		t.Errorf("LayerName.String() = %q, want %q", ln.String(), expected)
	}
}

func TestLayerPath_String(t *testing.T) {
	lp := LayerPath("path/to/layer")
	expected := "path/to/layer"
	if lp.String() != expected {
		t.Errorf("LayerPath.String() = %q, want %q", lp.String(), expected)
	}
}

func TestModuleName_String(t *testing.T) {
	mn := ModuleName("github.com/example/module")
	expected := "github.com/example/module"
	if mn.String() != expected {
		t.Errorf("ModuleName.String() = %q, want %q", mn.String(), expected)
	}
}

func TestPackageName_String(t *testing.T) {
	pn := PackageName("main")
	expected := "main"
	if pn.String() != expected {
		t.Errorf("PackageName.String() = %q, want %q", pn.String(), expected)
	}
}

func TestFilePath_String(t *testing.T) {
	fp := FilePath("path/to/file.go")
	expected := "path/to/file.go"
	if fp.String() != expected {
		t.Errorf("FilePath.String() = %q, want %q", fp.String(), expected)
	}
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
			if got := tt.input.IsValid(); got != tt.expected {
				t.Errorf("LayerName.IsValid() = %v, want %v", got, tt.expected)
			}
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
			if (err != nil) != tt.expectErr {
				t.Errorf("LayerName.Validate() error = %v, expectErr %v", err, tt.expectErr)
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
			if got := tt.input.IsValid(); got != tt.expected {
				t.Errorf("LayerPath.IsValid() = %v, want %v", got, tt.expected)
			}
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
			if (err != nil) != tt.expectErr {
				t.Errorf("LayerPath.Validate() error = %v, expectErr %v", err, tt.expectErr)
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
			if got := tt.input.IsValid(); got != tt.expected {
				t.Errorf("ModuleName.IsValid() = %v, want %v", got, tt.expected)
			}
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
			if (err != nil) != tt.expectErr {
				t.Errorf("ModuleName.Validate() error = %v, expectErr %v", err, tt.expectErr)
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
			if got := tt.input.IsValid(); got != tt.expected {
				t.Errorf("PackageName.IsValid() = %v, want %v", got, tt.expected)
			}
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
			if (err != nil) != tt.expectErr {
				t.Errorf("PackageName.Validate() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// Test custom error types
func TestUnsupportedReaderError(t *testing.T) {
	err := UnsupportedReaderError{ReaderType: "unknown"}
	expected := "unsupported reader type: unknown"
	if err.Error() != expected {
		t.Errorf("UnsupportedReaderError.Error() = %q, want %q", err.Error(), expected)
	}
}

func TestDirectoryCreationError(t *testing.T) {
	innerErr := errors.New("permission denied")
	err := DirectoryCreationError{Path: "/tmp/test", Err: innerErr}
	expected := "failed to create directory /tmp/test: permission denied"
	if err.Error() != expected {
		t.Errorf("DirectoryCreationError.Error() = %q, want %q", err.Error(), expected)
	}
}

func TestFileWriteError(t *testing.T) {
	innerErr := errors.New("disk full")
	err := FileWriteError{Path: "/tmp/file.txt", Err: innerErr}
	expected := "failed to write /tmp/file.txt: disk full"
	if err.Error() != expected {
		t.Errorf("FileWriteError.Error() = %q, want %q", err.Error(), expected)
	}
}

func TestFileFormatError(t *testing.T) {
	innerErr := errors.New("invalid syntax")
	err := FileFormatError{Path: "/tmp/code.go", Err: innerErr}
	expected := "failed to format /tmp/code.go: invalid syntax"
	if err.Error() != expected {
		t.Errorf("FileFormatError.Error() = %q, want %q", err.Error(), expected)
	}
}

func TestModuleNotFoundError(t *testing.T) {
	err := ModuleNotFoundError{Source: "go.mod"}
	expected := "module declaration not found in go.mod"
	if err.Error() != expected {
		t.Errorf("ModuleNotFoundError.Error() = %q, want %q", err.Error(), expected)
	}
}

// Test structs
func TestLayer(t *testing.T) {
	layer := Layer{
		Name: LayerName("test-layer"),
		Path: LayerPath("path/to/layer"),
	}

	if layer.Name.String() != "test-layer" {
		t.Errorf("Layer.Name = %q, want %q", layer.Name.String(), "test-layer")
	}

	if layer.Path.String() != "path/to/layer" {
		t.Errorf("Layer.Path = %q, want %q", layer.Path.String(), "path/to/layer")
	}
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

	if len(config.Layers) != 2 {
		t.Errorf("DependencyConfig.Layers length = %d, want %d", len(config.Layers), 2)
	}

	if len(config.Dependencies) != 1 {
		t.Errorf("DependencyConfig.Dependencies length = %d, want %d", len(config.Dependencies), 1)
	}

	deps, exists := config.Dependencies[LayerName("layer1")]
	if !exists {
		t.Error("Expected dependency for layer1 to exist")
	}

	if len(deps) != 1 || deps[0] != LayerName("layer2") {
		t.Errorf("Expected layer1 to depend on layer2, got %v", deps)
	}
}
