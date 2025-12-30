package generator

import (
	"sort"
)

var difficulties = []string{"Easy", "Medium", "Hard"}

func organizeProblemsByCategoryAndDifficulty(problems []Problem) CategoryProblems {
	categoryProblems := make(CategoryProblems)

	for _, category := range CategoryOrder() {
		categoryProblems[category] = map[string][]Problem{
			"Easy":   {},
			"Medium": {},
			"Hard":   {},
		}
	}

	for _, problem := range problems {
		if _, exists := categoryProblems[problem.Category]; exists {
			categoryProblems[problem.Category][problem.Difficulty] = append(
				categoryProblems[problem.Category][problem.Difficulty],
				problem,
			)
		}
	}

	return categoryProblems
}

// SpacingState tracks spaced repetition state
type SpacingState struct {
	CategoryGaps      map[string]int // How many picks to skip before this category
	CategoryCounters  map[string]int // Picks since last use of this category
	CurrentDifficulty string         // Current difficulty tier
}

// NewSpacingState creates a new spacing state
func NewSpacingState() *SpacingState {
	gaps := make(map[string]int)
	counters := make(map[string]int)
	for _, cat := range CategoryOrder() {
		gaps[cat] = 0
		counters[cat] = 0
	}
	return &SpacingState{
		CategoryGaps:      gaps,
		CategoryCounters:  counters,
		CurrentDifficulty: "",
	}
}

// Reset resets all gaps and counters (called on difficulty change)
func (s *SpacingState) Reset() {
	for cat := range s.CategoryGaps {
		s.CategoryGaps[cat] = 0
		s.CategoryCounters[cat] = 0
	}
}

// selectNextProblem picks the next problem using strict priority:
// 1. Difficulty tier (Easy > Medium > Hard)
// 2. Category with spaced repetition (increasing gaps)
// 3. LeetCode number (lowest first)
func selectNextProblem(categoryProblems CategoryProblems, usedProblems map[string]bool, state *SpacingState, usedToday map[string]bool, usedYesterday map[string]bool) (*Problem, string) {
	// Try each difficulty tier in order
	for _, difficulty := range difficulties {
		// Check if this difficulty has any problems before potentially resetting
		if !hasProblemsInDifficulty(difficulty, categoryProblems, usedProblems) {
			continue
		}

		// Reset spacing only when actually switching to a new difficulty with problems
		if state.CurrentDifficulty != difficulty {
			state.Reset()
			state.CurrentDifficulty = difficulty
		}

		problem, category := selectFromDifficultyWithSpacing(difficulty, categoryProblems, usedProblems, state, usedToday, usedYesterday)
		if problem != nil {
			return problem, category
		}
	}
	return nil, ""
}

// hasProblemsInDifficulty checks if any category has unused problems in this difficulty
func hasProblemsInDifficulty(difficulty string, categoryProblems CategoryProblems, usedProblems map[string]bool) bool {
	for _, category := range CategoryOrder() {
		for _, problem := range categoryProblems[category][difficulty] {
			if !usedProblems[problem.Name] {
				return true
			}
		}
	}
	return false
}

// selectFromDifficultyWithSpacing finds the next problem with increasing gaps (capped)
func selectFromDifficultyWithSpacing(difficulty string, categoryProblems CategoryProblems, usedProblems map[string]bool, state *SpacingState, usedToday map[string]bool, usedYesterday map[string]bool) (*Problem, string) {
	categories := CategoryOrder() // Ordered by difficulty
	const maxGap = 6              // Cap gaps to keep spacing uniform

	// First pass: skip categories used today, yesterday, and those not ready
	for _, category := range categories {
		if usedToday[category] || usedYesterday[category] {
			continue
		}
		if state.CategoryCounters[category] < state.CategoryGaps[category] {
			continue
		}
		problem := pickFromCategory(category, difficulty, categoryProblems, usedProblems, state, categories, maxGap)
		if problem != nil {
			return problem, category
		}
	}

	// Second pass: allow categories used yesterday but not today
	for _, category := range categories {
		if usedToday[category] {
			continue
		}
		if state.CategoryCounters[category] < state.CategoryGaps[category] {
			continue
		}
		problem := pickFromCategory(category, difficulty, categoryProblems, usedProblems, state, categories, maxGap)
		if problem != nil {
			return problem, category
		}
	}

	// Third pass: allow categories used today (fallback)
	for _, category := range categories {
		if state.CategoryCounters[category] < state.CategoryGaps[category] {
			continue
		}
		problem := pickFromCategory(category, difficulty, categoryProblems, usedProblems, state, categories, maxGap)
		if problem != nil {
			return problem, category
		}
	}

	// Fourth pass: ignore gaps (when all ready categories exhausted)
	for _, category := range categories {
		problems := categoryProblems[category][difficulty]
		problem := findLowestUnused(problems, usedProblems)
		if problem != nil {
			state.CategoryGaps[category] = 1 // Reset gap
			state.CategoryCounters[category] = 0
			for _, other := range categories {
				if other != category {
					state.CategoryCounters[other]++
				}
			}
			return problem, category
		}
	}

	return nil, ""
}

func pickFromCategory(category, difficulty string, categoryProblems CategoryProblems, usedProblems map[string]bool, state *SpacingState, categories []string, maxGap int) *Problem {
	problems := categoryProblems[category][difficulty]
	problem := findLowestUnused(problems, usedProblems)
	if problem != nil {
		if state.CategoryGaps[category] < maxGap {
			state.CategoryGaps[category]++
		}
		state.CategoryCounters[category] = 0
		for _, other := range categories {
			if other != category {
				state.CategoryCounters[other]++
			}
		}
	}
	return problem
}

// findLowestUnused returns the unused problem with the lowest LeetCode number
func findLowestUnused(problems []Problem, usedProblems map[string]bool) *Problem {
	var candidates []*Problem
	for i := range problems {
		if !usedProblems[problems[i].Name] {
			candidates = append(candidates, &problems[i])
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	// Sort by LeetCode number and return lowest
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].LeetCodeNumber < candidates[j].LeetCodeNumber
	})

	return candidates[0]
}
