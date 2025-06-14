package main

import (
	"fmt"
	"strings"
)

// Custom types for better type safety
type LayerName string
type LayerPath string
type ModuleName string
type PackageName string
type FilePath string

func (ln LayerName) String() string   { return string(ln) }
func (lp LayerPath) String() string   { return string(lp) }
func (mn ModuleName) String() string  { return string(mn) }
func (pn PackageName) String() string { return string(pn) }
func (fp FilePath) String() string    { return string(fp) }

// Validation methods for custom types
func (ln LayerName) IsValid() bool {
	return ln != "" && strings.TrimSpace(string(ln)) != ""
}

func (ln LayerName) Validate() error {
	if !ln.IsValid() {
		return fmt.Errorf("layer name cannot be empty")
	}
	return nil
}

func (lp LayerPath) IsValid() bool {
	return lp != "" && strings.TrimSpace(string(lp)) != "" && !strings.Contains(string(lp), "..")
}

func (lp LayerPath) Validate() error {
	if lp == "" || strings.TrimSpace(string(lp)) == "" {
		return fmt.Errorf("layer path cannot be empty")
	}
	if strings.Contains(string(lp), "..") {
		return fmt.Errorf("layer path cannot contain '..' for security reasons")
	}
	return nil
}

func (mn ModuleName) IsValid() bool {
	return mn != "" && strings.TrimSpace(string(mn)) != "" && !strings.Contains(string(mn), " ")
}

func (mn ModuleName) Validate() error {
	if mn == "" || strings.TrimSpace(string(mn)) == "" {
		return fmt.Errorf("module name cannot be empty")
	}
	if strings.Contains(string(mn), " ") {
		return fmt.Errorf("module name cannot contain spaces")
	}
	return nil
}

func (pn PackageName) IsValid() bool {
	return pn != "" && strings.TrimSpace(string(pn)) != "" && !strings.Contains(string(pn), "/")
}

func (pn PackageName) Validate() error {
	if pn == "" || strings.TrimSpace(string(pn)) == "" {
		return fmt.Errorf("package name cannot be empty")
	}
	if strings.Contains(string(pn), "/") {
		return fmt.Errorf("package name cannot contain '/'")
	}
	return nil
}

// Custom error types for better error handling
type UnsupportedReaderError struct {
	ReaderType string
}

func (e UnsupportedReaderError) Error() string {
	return fmt.Sprintf("unsupported reader type: %s", e.ReaderType)
}

type DirectoryCreationError struct {
	Path string
	Err  error
}

func (e DirectoryCreationError) Error() string {
	return fmt.Sprintf("failed to create directory %s: %v", e.Path, e.Err)
}

type FileWriteError struct {
	Path string
	Err  error
}

func (e FileWriteError) Error() string {
	return fmt.Sprintf("failed to write %s: %v", e.Path, e.Err)
}

type FileFormatError struct {
	Path string
	Err  error
}

func (e FileFormatError) Error() string {
	return fmt.Sprintf("failed to format %s: %v", e.Path, e.Err)
}

type ModuleNotFoundError struct {
	Source string
}

func (e ModuleNotFoundError) Error() string {
	return fmt.Sprintf("module declaration not found in %s", e.Source)
}

// Structs
type Layer struct {
	Name LayerName
	Path LayerPath
}

type DependencyConfig struct {
	Layers       []Layer
	Dependencies map[LayerName][]LayerName
}
