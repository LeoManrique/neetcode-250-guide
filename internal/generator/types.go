package generator

// Problem represents a single coding problem
type Problem struct {
	Name           string `json:"name"`
	Difficulty     string `json:"difficulty"`
	Category       string `json:"category"`
	NeetcodeURL    string `json:"neetcode_url"`
	LeetcodeURL    string `json:"leetcode_url"`
	Slug           string `json:"slug"`
	LeetCodeNumber int    `json:"leetcode_number"`
}

// ProblemsData represents the JSON file structure
type ProblemsData struct {
	Problems []Problem `json:"problems"`
}

// DayPlan represents a single day's study plan
type DayPlan struct {
	Date     string
	Day      int
	Problems []Problem
	Category string
}

// CategoryProblems organizes problems by difficulty within a category
type CategoryProblems map[string]map[string][]Problem

// Category represents a problem category with its difficulty rating
type Category struct {
	Name       string `json:"name"`
	Difficulty int    `json:"difficulty"` // difficulty out of 20
}

// CategoriesData represents the JSON file structure for categories
type CategoriesData struct {
	Categories []Category `json:"categories"`
}

// Categories holds the loaded category data
var Categories []Category

// CategoryOrder returns category names in order (for backward compatibility)
func CategoryOrder() []string {
	names := make([]string, len(Categories))
	for i, cat := range Categories {
		names[i] = cat.Name
	}
	return names
}
