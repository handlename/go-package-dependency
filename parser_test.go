package main

import (
	"os"
	"testing"
)

func TestParseDependencyContent(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		expectedLayers []Layer
		expectedDeps   map[LayerName][]LayerName
		expectError    bool
	}{
		{
			name: "valid dependency file",
			content: `# Dependencies

## Dependencies
- Infra layer -> Presentation layer -> Application layer -> Domain layer
- Another layer -> Domain layer

## Layers
- Domain layer
  - domain/entity
- Application layer
  - application/usecase
- Presentation layer
  - presentation/handler
- Infra layer
  - infra/repository
- Another layer
  - another/path
`,
			expectedLayers: []Layer{
				{Name: LayerName("Domain layer"), Path: LayerPath("domain/entity")},
				{Name: LayerName("Application layer"), Path: LayerPath("application/usecase")},
				{Name: LayerName("Presentation layer"), Path: LayerPath("presentation/handler")},
				{Name: LayerName("Infra layer"), Path: LayerPath("infra/repository")},
				{Name: LayerName("Another layer"), Path: LayerPath("another/path")},
			},
			expectedDeps: map[LayerName][]LayerName{
				LayerName("Infra layer"):        {LayerName("Presentation layer")},
				LayerName("Presentation layer"): {LayerName("Application layer")},
				LayerName("Application layer"):  {LayerName("Domain layer")},
				LayerName("Another layer"):      {LayerName("Domain layer")},
			},
			expectError: false,
		},
		{
			name:           "empty content",
			content:        ``,
			expectedLayers: []Layer{},
			expectedDeps:   map[LayerName][]LayerName{},
			expectError:    false,
		},
		{
			name: "only dependencies section",
			content: `## Dependencies
- Layer A -> Layer B`,
			expectedLayers: []Layer{},
			expectedDeps: map[LayerName][]LayerName{
				LayerName("Layer A"): {LayerName("Layer B")},
			},
			expectError: false,
		},
		{
			name: "only layers section",
			content: `## Layers
- Test layer
  - test/path`,
			expectedLayers: []Layer{
				{Name: LayerName("Test layer"), Path: LayerPath("test/path")},
			},
			expectedDeps: map[LayerName][]LayerName{},
			expectError:  false,
		},
		{
			name: "invalid layer path with ..",
			content: `## Layers
- Test layer
  - ../invalid/path`,
			expectedLayers: []Layer{},
			expectedDeps:   map[LayerName][]LayerName{},
			expectError:    true,
		},
		{
			name: "line that doesn't match layer format",
			content: `## Layers
-
  - test/path`,
			expectedLayers: []Layer{},
			expectedDeps:   map[LayerName][]LayerName{},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			config, err := parser.ParseDependencyContent(tt.content)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check layers
			if len(config.Layers) != len(tt.expectedLayers) {
				t.Errorf("Expected %d layers, got %d", len(tt.expectedLayers), len(config.Layers))
			}

			for i, expected := range tt.expectedLayers {
				if i >= len(config.Layers) {
					t.Errorf("Missing layer at index %d", i)
					continue
				}
				actual := config.Layers[i]
				if actual.Name != expected.Name || actual.Path != expected.Path {
					t.Errorf("Layer %d: expected {%s, %s}, got {%s, %s}",
						i, expected.Name, expected.Path, actual.Name, actual.Path)
				}
			}

			// Check dependencies
			if len(config.Dependencies) != len(tt.expectedDeps) {
				t.Errorf("Expected %d dependencies, got %d", len(tt.expectedDeps), len(config.Dependencies))
			}

			for layer, expectedDeps := range tt.expectedDeps {
				actualDeps, exists := config.Dependencies[layer]
				if !exists {
					t.Errorf("Missing dependencies for layer %s", layer)
					continue
				}

				if len(actualDeps) != len(expectedDeps) {
					t.Errorf("Layer %s: expected %d dependencies, got %d", layer, len(expectedDeps), len(actualDeps))
					continue
				}

				for i, expectedDep := range expectedDeps {
					if i >= len(actualDeps) || actualDeps[i] != expectedDep {
						t.Errorf("Layer %s dependency %d: expected %s, got %s", layer, i, expectedDep, actualDeps[i])
					}
				}
			}
		})
	}
}

func TestParseDependencyLine(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		expected    map[LayerName][]LayerName
		expectError bool
	}{
		{
			name: "simple dependency chain",
			line: "- Layer A -> Layer B -> Layer C",
			expected: map[LayerName][]LayerName{
				LayerName("Layer A"): {LayerName("Layer B")},
				LayerName("Layer B"): {LayerName("Layer C")},
			},
			expectError: false,
		},
		{
			name: "single dependency",
			line: "- Frontend -> Backend",
			expected: map[LayerName][]LayerName{
				LayerName("Frontend"): {LayerName("Backend")},
			},
			expectError: false,
		},
		{
			name: "long dependency chain",
			line: "- UI -> Service -> Repository -> Database",
			expected: map[LayerName][]LayerName{
				LayerName("UI"):         {LayerName("Service")},
				LayerName("Service"):    {LayerName("Repository")},
				LayerName("Repository"): {LayerName("Database")},
			},
			expectError: false,
		},
		{
			name: "dependency with extra spaces",
			line: "-   Layer A   ->   Layer B   ->   Layer C   ",
			expected: map[LayerName][]LayerName{
				LayerName("Layer A"): {LayerName("Layer B")},
				LayerName("Layer B"): {LayerName("Layer C")},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &DependencyConfig{
				Layers:       make([]Layer, 0),
				Dependencies: make(map[LayerName][]LayerName),
			}

			parser := NewParser()
			err := parser.ParseDependencyLine(tt.line, config)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check that all expected dependencies are present
			for layer, expectedDeps := range tt.expected {
				actualDeps, exists := config.Dependencies[layer]
				if !exists {
					t.Errorf("Missing dependencies for layer %s", layer)
					continue
				}

				if len(actualDeps) != len(expectedDeps) {
					t.Errorf("Layer %s: expected %d dependencies, got %d", layer, len(expectedDeps), len(actualDeps))
					continue
				}

				for i, expectedDep := range expectedDeps {
					if actualDeps[i] != expectedDep {
						t.Errorf("Layer %s dependency %d: expected %s, got %s", layer, i, expectedDep, actualDeps[i])
					}
				}
			}
		})
	}
}

func TestParseLayerLine(t *testing.T) {
	tests := []struct {
		name           string
		lines          []string
		expectedLayers []Layer
		expectError    bool
	}{
		{
			name: "layer with path",
			lines: []string{
				"- Domain layer",
				"  - domain/entity",
			},
			expectedLayers: []Layer{
				{Name: "Domain layer", Path: "domain/entity"},
			},
			expectError: false,
		},
		{
			name: "multiple layers with paths",
			lines: []string{
				"- Domain layer",
				"  - domain/entity",
				"- Application layer",
				"  - application/usecase",
			},
			expectedLayers: []Layer{
				{Name: "Domain layer", Path: "domain/entity"},
				{Name: "Application layer", Path: "application/usecase"},
			},
			expectError: false,
		},
		{
			name: "layer without path",
			lines: []string{
				"- Standalone layer",
			},
			expectedLayers: []Layer{
				{Name: "Standalone layer", Path: ""},
			},
			expectError: false,
		},
		{
			name: "invalid path with ..",
			lines: []string{
				"- Test layer",
				"  - ../invalid/path",
			},
			expectedLayers: []Layer{},
			expectError:    true,
		},
		{
			name: "line that doesn't match format",
			lines: []string{
				"- ",
				"  - valid/path",
			},
			expectedLayers: []Layer{},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &DependencyConfig{
				Layers: make([]Layer, 0),
			}

			parser := NewParser()
			var lastErr error
			for _, line := range tt.lines {
				err := parser.ParseLayerLine(line, config)
				if err != nil {
					lastErr = err
					break
				}
			}

			if tt.expectError {
				if lastErr == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if lastErr != nil {
				t.Errorf("Unexpected error: %v", lastErr)
				return
			}

			if len(config.Layers) != len(tt.expectedLayers) {
				t.Errorf("Expected %d layers, got %d", len(tt.expectedLayers), len(config.Layers))
			}

			for i, expected := range tt.expectedLayers {
				if i >= len(config.Layers) {
					t.Errorf("Missing layer at index %d", i)
					continue
				}
				actual := config.Layers[i]
				if actual.Name != expected.Name || actual.Path != expected.Path {
					t.Errorf("Layer %d: expected {%s, %s}, got {%s, %s}",
						i, expected.Name, expected.Path, actual.Name, actual.Path)
				}
			}
		})
	}
}

func TestGetModuleNameFromContent(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expected    ModuleName
		expectError bool
	}{
		{
			name: "valid module declaration",
			content: `module github.com/example/project

go 1.21
`,
			expected:    ModuleName("github.com/example/project"),
			expectError: false,
		},
		{
			name: "module with extra spaces",
			content: `module   github.com/example/project

go 1.21
`,
			expected:    ModuleName("github.com/example/project"),
			expectError: false,
		},
		{
			name: "module declaration in middle of file",
			content: `// This is a go.mod file
module github.com/test/module

require (
    github.com/dependency v1.0.0
)
`,
			expected:    ModuleName("github.com/test/module"),
			expectError: false,
		},
		{
			name: "no module declaration",
			content: `go 1.21

require (
    github.com/dependency v1.0.0
)
`,
			expected:    ModuleName(""),
			expectError: true,
		},
		{
			name:        "empty content",
			content:     ``,
			expected:    ModuleName(""),
			expectError: true,
		},
		{
			name: "invalid module name with spaces",
			content: `module invalid module name

go 1.21
`,
			expected:    ModuleName(""),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			result, err := parser.GetModuleNameFromContent(tt.content, "test.mod")

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestParseDependencyFile(t *testing.T) {
	// Create a temporary file
	content := `# Test Dependencies

## Dependencies
- Presentation layer -> Application layer -> Domain layer

## Layers
- Domain layer
  - domain/entity
- Application layer
  - application/usecase
- Presentation layer
  - presentation/handler
`

	tmpFile, err := os.CreateTemp("", "dependency-*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	parser := NewParser()
	config, err := parser.ParseDependencyFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to parse dependency file: %v", err)
	}

	// Verify the parsed content
	expectedLayers := []Layer{
		{Name: LayerName("Domain layer"), Path: LayerPath("domain/entity")},
		{Name: LayerName("Application layer"), Path: LayerPath("application/usecase")},
		{Name: LayerName("Presentation layer"), Path: LayerPath("presentation/handler")},
	}

	if len(config.Layers) != len(expectedLayers) {
		t.Errorf("Expected %d layers, got %d", len(expectedLayers), len(config.Layers))
	}

	for i, expected := range expectedLayers {
		if i >= len(config.Layers) {
			t.Errorf("Missing layer at index %d", i)
			continue
		}
		actual := config.Layers[i]
		if actual.Name != expected.Name || actual.Path != expected.Path {
			t.Errorf("Layer %d: expected {%s, %s}, got {%s, %s}",
				i, expected.Name, expected.Path, actual.Name, actual.Path)
		}
	}

	expectedDeps := map[LayerName][]LayerName{
		LayerName("Presentation layer"): {LayerName("Application layer")},
		LayerName("Application layer"):  {LayerName("Domain layer")},
	}

	if len(config.Dependencies) != len(expectedDeps) {
		t.Errorf("Expected %d dependencies, got %d", len(expectedDeps), len(config.Dependencies))
	}

	for layer, expectedDeps := range expectedDeps {
		actualDeps, exists := config.Dependencies[layer]
		if !exists {
			t.Errorf("Missing dependencies for layer %s", layer)
			continue
		}

		if len(actualDeps) != len(expectedDeps) {
			t.Errorf("Layer %s: expected %d dependencies, got %d", layer, len(expectedDeps), len(actualDeps))
			continue
		}

		for i, expectedDep := range expectedDeps {
			if actualDeps[i] != expectedDep {
				t.Errorf("Layer %s dependency %d: expected %s, got %s", layer, i, expectedDep, actualDeps[i])
			}
		}
	}
}

func TestGetModuleName(t *testing.T) {
	// Create a temporary go.mod file
	content := `module github.com/test/project

go 1.21

require (
    github.com/dependency v1.0.0
)
`

	tmpFile, err := os.CreateTemp("", "go.mod")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	parser := NewParser()
	result, err := parser.GetModuleName(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to get module name: %v", err)
	}

	expected := ModuleName("github.com/test/project")
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
