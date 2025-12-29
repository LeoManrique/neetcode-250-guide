package generator

import "math/rand"

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

func calculateCategoryWeights(day int, categoryProblems CategoryProblems, categoryUsageHistory []string, categoryProgress CategoryProgress) map[string]float64 {
	weights := make(map[string]float64)
	dayProgress := float64(day) / 125.0

	for _, category := range CategoryOrder() {
		// Use actual difficulty rating (1-20) normalized to 0-1
		categoryDifficulty := float64(GetCategoryDifficulty(category)) / 20.0
		baseWeight := calculateBaseWeight(dayProgress, categoryDifficulty)
		remainingProblems := countRemainingProblems(category, categoryProblems, categoryProgress)

		if remainingProblems > 4 {
			baseWeight *= 1.2
		} else if remainingProblems == 0 {
			baseWeight = 0
		}

		weights[category] = max(0.0001, baseWeight)
	}

	applySpacedRepetitionPenalty(weights, categoryUsageHistory)
	return weights
}

func calculateBaseWeight(dayProgress, categoryDifficulty float64) float64 {
	// dayProgress: 0/125 to 124/125 (each day has unique value)
	// categoryDifficulty: 1/20 to 20/20 (each category has unique value)

	// Target difficulty = day progress (day 1 targets easiest, day 125 targets hardest)
	diff := categoryDifficulty - dayProgress

	// Moderate penalty for distance from target
	weight := 1.0 - diff*diff*15

	return max(0.0001, weight)
}

func countRemainingProblems(category string, categoryProblems CategoryProblems, categoryProgress CategoryProgress) int {
	remaining := 0
	for _, difficulty := range difficulties {
		available := len(categoryProblems[category][difficulty])
		used := categoryProgress[category][difficulty]
		if available-used > 0 {
			remaining += available - used
		}
	}
	return remaining
}

func applySpacedRepetitionPenalty(weights map[string]float64, history []string) {
	const recentUsagePenalty = 0.7
	for i, category := range history {
		if i >= 6 {
			break
		}
		penalty := recentUsagePenalty * (1 - float64(i)/6)
		if _, exists := weights[category]; exists {
			weights[category] = weights[category] * (1 - penalty)
		}
	}
}

func selectCategoryForDay(day int, categoryProblems CategoryProblems, categoryUsageHistory []string, categoryProgress CategoryProgress, usedProblems map[string]bool) string {
	categoryWeights := calculateCategoryWeights(day, categoryProblems, categoryUsageHistory, categoryProgress)
	availableCategories := findAvailableCategories(categoryProblems, usedProblems)

	if len(availableCategories) == 0 {
		return ""
	}

	return weightedRandomSelect(availableCategories, categoryWeights)
}

func findAvailableCategories(categoryProblems CategoryProblems, usedProblems map[string]bool) []string {
	var categories2Plus, categories1Plus []string

	for _, category := range CategoryOrder() {
		available := countAvailableInCategory(category, categoryProblems, usedProblems)
		if available >= 2 {
			categories2Plus = append(categories2Plus, category)
		} else if available >= 1 {
			categories1Plus = append(categories1Plus, category)
		}
	}

	if len(categories2Plus) > 0 {
		return categories2Plus
	}
	return categories1Plus
}

func countAvailableInCategory(category string, categoryProblems CategoryProblems, usedProblems map[string]bool) int {
	count := 0
	for _, difficulty := range difficulties {
		for _, problem := range categoryProblems[category][difficulty] {
			if !usedProblems[problem.Name] {
				count++
			}
		}
	}
	return count
}

func weightedRandomSelect(categories []string, weights map[string]float64) string {
	var totalWeight float64
	categoryWeights := make([]float64, len(categories))

	for i, cat := range categories {
		w := weights[cat]
		if w == 0 {
			w = 0.0001
		}
		categoryWeights[i] = w
		totalWeight += w
	}

	if totalWeight > 0 {
		r := rand.Float64() * totalWeight
		cumulative := 0.0
		for i, w := range categoryWeights {
			cumulative += w
			if r <= cumulative {
				return categories[i]
			}
		}
	}

	return categories[rand.Intn(len(categories))]
}

func selectProblemsForDay(day int, selectedCategory string, categoryProblems CategoryProblems, usedProblems map[string]bool) []Problem {
	var dayProblems []Problem
	preferences := getDifficultyPreferences(day)

	for problemNum := 0; problemNum < 2; problemNum++ {
		difficultyOrder := difficulties
		if problemNum < len(preferences) {
			difficultyOrder = preferences[problemNum]
		}

		problem := findUnusedProblem(selectedCategory, difficultyOrder, categoryProblems, usedProblems)
		if problem == nil {
			problem = findUnusedProblem(selectedCategory, difficulties, categoryProblems, usedProblems)
		}

		if problem != nil {
			usedProblems[problem.Name] = true
			dayProblems = append(dayProblems, *problem)
		} else {
			break
		}
	}

	return dayProblems
}

func getDifficultyPreferences(day int) [][]string {
	switch {
	case day < 30:
		secondPref := []string{"Easy", "Easy", "Medium"}
		if day%3 == 0 {
			secondPref = []string{"Medium", "Easy", "Hard"}
		}
		return [][]string{{"Easy", "Medium", "Hard"}, secondPref}
	case day < 80:
		return [][]string{{"Easy", "Medium", "Hard"}, {"Medium", "Easy", "Hard"}}
	default:
		return [][]string{{"Medium", "Hard", "Easy"}, {"Medium", "Hard", "Easy"}}
	}
}

func findUnusedProblem(category string, difficultyOrder []string, categoryProblems CategoryProblems, usedProblems map[string]bool) *Problem {
	// Collect unused problems from preferred difficulties
	var candidates []*Problem
	for _, difficulty := range difficultyOrder {
		for i := range categoryProblems[category][difficulty] {
			problem := &categoryProblems[category][difficulty][i]
			if !usedProblems[problem.Name] {
				candidates = append(candidates, problem)
			}
		}
		// Stop at first difficulty that has candidates
		if len(candidates) > 0 {
			break
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	// Single candidate - return directly
	if len(candidates) == 1 {
		return candidates[0]
	}

	// Weighted random selection favoring lower LeetCode numbers
	return selectByLeetCodeWeight(candidates)
}

func selectByLeetCodeWeight(candidates []*Problem) *Problem {
	// Find max LeetCode number for normalization
	maxNum := 0
	for _, p := range candidates {
		if p.LeetCodeNumber > maxNum {
			maxNum = p.LeetCodeNumber
		}
	}

	// Calculate weights: lower numbers get higher weights
	weights := make([]float64, len(candidates))
	totalWeight := 0.0
	for i, p := range candidates {
		// Invert and normalize, then square for stronger preference
		normalized := float64(maxNum-p.LeetCodeNumber+1) / float64(maxNum)
		weight := normalized * normalized
		weights[i] = weight
		totalWeight += weight
	}

	// Weighted random selection
	r := rand.Float64() * totalWeight
	cumulative := 0.0
	for i, w := range weights {
		cumulative += w
		if r <= cumulative {
			return candidates[i]
		}
	}

	return candidates[len(candidates)-1]
}
