package generator

// AnalyzePlanDistribution analyzes the distribution of categories and difficulties
func AnalyzePlanDistribution(plan []DayPlan) (map[string]map[string]int, map[string]map[string]int) {
	categoryByPhase := map[string]map[string]int{
		"Early (1-30)":   make(map[string]int),
		"Middle (31-80)": make(map[string]int),
		"Late (81-125)":  make(map[string]int),
	}
	difficultyByPhase := map[string]map[string]int{
		"Early (1-30)":   make(map[string]int),
		"Middle (31-80)": make(map[string]int),
		"Late (81-125)":  make(map[string]int),
	}

	for _, dayPlan := range plan {
		phase := getPhase(dayPlan.Day)
		categoryByPhase[phase][dayPlan.Category]++
		for _, problem := range dayPlan.Problems {
			difficultyByPhase[phase][problem.Difficulty]++
		}
	}

	return categoryByPhase, difficultyByPhase
}

func getPhase(day int) string {
	switch {
	case day <= 30:
		return "Early (1-30)"
	case day <= 80:
		return "Middle (31-80)"
	default:
		return "Late (81-125)"
	}
}

// CountProblems returns the total number of problems in the plan
func CountProblems(plan []DayPlan) int {
	total := 0
	for _, day := range plan {
		total += len(day.Problems)
	}
	return total
}

// GetUsedProblemNames returns a set of all problem names in the plan
func GetUsedProblemNames(plan []DayPlan) map[string]bool {
	used := make(map[string]bool)
	for _, day := range plan {
		for _, problem := range day.Problems {
			used[problem.Name] = true
		}
	}
	return used
}

// FindMissingProblems returns problems not included in the plan
func FindMissingProblems(allProblems []Problem, plan []DayPlan) []string {
	usedNames := GetUsedProblemNames(plan)
	var missing []string
	for _, p := range allProblems {
		if !usedNames[p.Name] {
			missing = append(missing, p.Name)
		}
	}
	return missing
}
