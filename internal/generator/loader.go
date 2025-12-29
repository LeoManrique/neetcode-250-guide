package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// findFile locates a file relative to the executable or working directory
func findFile(filename string) (string, error) {
	paths := []string{
		filepath.Join("data", filename),
		filepath.Join("..", "data", filename),
		filepath.Join("..", "..", "data", filename),
	}

	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		paths = append(paths,
			filepath.Join(execDir, "data", filename),
			filepath.Join(execDir, "..", "data", filename),
			filepath.Join(execDir, "..", "..", "data", filename),
		)
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("could not find %s in any expected location", filename)
}

// FindDataFile locates the problems JSON file
func FindDataFile() (string, error) {
	return findFile("exercises.json")
}

// FindCategoriesFile locates the categories JSON file
func FindCategoriesFile() (string, error) {
	return findFile("categories.json")
}

// LoadProblems loads problems from the JSON file
func LoadProblems(dataPath string) ([]Problem, error) {
	data, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	var problemsData ProblemsData
	if err := json.Unmarshal(data, &problemsData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return problemsData.Problems, nil
}

// LoadCategories loads categories from the JSON file and populates the global Categories variable
func LoadCategories(dataPath string) error {
	data, err := os.ReadFile(dataPath)
	if err != nil {
		return fmt.Errorf("failed to read categories file: %w", err)
	}

	var categoriesData CategoriesData
	if err := json.Unmarshal(data, &categoriesData); err != nil {
		return fmt.Errorf("failed to parse categories JSON: %w", err)
	}

	Categories = categoriesData.Categories
	return nil
}
