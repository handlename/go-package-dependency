package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseDependencyFile(filePath string) (*DependencyConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return p.ParseDependencyContent(file)
}

func (p *Parser) ParseDependencyContent(reader any) (*DependencyConfig, error) {
	var scanner *bufio.Scanner

	switch r := reader.(type) {
	case *os.File:
		scanner = bufio.NewScanner(r)
	case *strings.Reader:
		scanner = bufio.NewScanner(r)
	case string:
		scanner = bufio.NewScanner(strings.NewReader(r))
	default:
		return nil, UnsupportedReaderError{ReaderType: fmt.Sprintf("%T", reader)}
	}

	config := &DependencyConfig{
		Layers:       make([]Layer, 0),
		Dependencies: make(map[LayerName][]LayerName),
	}

	inDependenciesSection := false
	inLayersSection := false

	for scanner.Scan() {
		rawLine := scanner.Text()
		line := strings.TrimSpace(rawLine)

		// Check for section headers
		if strings.HasPrefix(line, "## Dependencies") {
			inDependenciesSection = true
			inLayersSection = false
			continue
		}
		if strings.HasPrefix(line, "## Layers") {
			inLayersSection = true
			inDependenciesSection = false
			continue
		}

		// Parse dependencies section
		if inDependenciesSection && strings.Contains(line, "->") {
			err := p.ParseDependencyLine(line, config)
			if err != nil {
				return nil, err
			}
		}

		// Parse layers section
		if inLayersSection && (strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "  - ")) {
			err := p.ParseLayerLine(rawLine, config)
			if err != nil {
				return nil, err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}

func (p *Parser) ParseDependencyLine(line string, config *DependencyConfig) error {
	// Parse "- Infra layer -> Presentation layer -> Application layer -> Domain layer"
	line = strings.TrimPrefix(line, "- ")
	parts := strings.Split(line, " -> ")

	for i := range len(parts) - 1 {
		dependent := LayerName(strings.TrimSpace(parts[i]))
		dependency := LayerName(strings.TrimSpace(parts[i+1]))

		if config.Dependencies[dependent] == nil {
			config.Dependencies[dependent] = make([]LayerName, 0)
		}
		config.Dependencies[dependent] = append(config.Dependencies[dependent], dependency)
	}

	return nil
}

func (p *Parser) ParseLayerLine(line string, config *DependencyConfig) error {
	// Parse layer definitions like "- Domain layer" and "  - domain/entity"
	trimmed := strings.TrimSpace(line)

	if strings.HasPrefix(trimmed, "- ") && line == trimmed {
		// This is a layer name (not indented)
		layerName := LayerName(strings.TrimPrefix(trimmed, "- "))
		if err := layerName.Validate(); err != nil {
			return fmt.Errorf("invalid layer name: %v", err)
		}
		config.Layers = append(config.Layers, Layer{Name: layerName, Path: ""})
	} else if strings.HasPrefix(trimmed, "- ") && line != trimmed {
		// This is an indented line with "- " (a layer path)
		layerPath := LayerPath(strings.TrimPrefix(trimmed, "- "))
		if err := layerPath.Validate(); err != nil {
			return fmt.Errorf("invalid layer path: %v", err)
		}
		if len(config.Layers) > 0 {
			config.Layers[len(config.Layers)-1].Path = layerPath
		}
	}

	return nil
}

func (p *Parser) GetModuleName(goModPath string) (ModuleName, error) {
	file, err := os.Open(goModPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return p.GetModuleNameFromContent(file, goModPath)
}

func (p *Parser) GetModuleNameFromContent(reader any, sourceName string) (ModuleName, error) {
	var scanner *bufio.Scanner

	switch r := reader.(type) {
	case *os.File:
		scanner = bufio.NewScanner(r)
	case *strings.Reader:
		scanner = bufio.NewScanner(r)
	case string:
		scanner = bufio.NewScanner(strings.NewReader(r))
	default:
		return "", UnsupportedReaderError{ReaderType: fmt.Sprintf("%T", reader)}
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			moduleName := ModuleName(strings.TrimSpace(strings.TrimPrefix(line, "module ")))
			if err := moduleName.Validate(); err != nil {
				return "", fmt.Errorf("invalid module name: %v", err)
			}
			return moduleName, nil
		}
	}

	return "", ModuleNotFoundError{Source: sourceName}
}
