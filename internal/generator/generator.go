package generator

import (
	"time"
)

// GeneratePlan creates a 125-day study plan including all 250 problems
// Selection priority: Difficulty (E>M>H) > Category difficulty > LeetCode number
func GeneratePlan(problems []Problem, startDate time.Time) []DayPlan {
	categoryProblems := organizeProblemsByCategoryAndDifficulty(problems)
	usedProblems := make(map[string]bool)

	var plan []DayPlan
	spacingState := NewSpacingState()
	usedYesterday := make(map[string]bool)
	day := 0
	const maxDays = 150

	for len(usedProblems) < 250 && day < maxDays {
		currentDate := startDate.AddDate(0, 0, day)
		var dayProblems []Problem
		usedToday := make(map[string]bool)

		// Select 2 problems for each day
		for i := 0; i < 2 && len(usedProblems) < 250; i++ {
			problem, category := selectNextProblem(categoryProblems, usedProblems, spacingState, usedToday, usedYesterday)
			if problem != nil {
				dayProblems = append(dayProblems, *problem)
				usedProblems[problem.Name] = true
				usedToday[category] = true
			}
		}

		if len(dayProblems) == 0 {
			day++
			continue
		}

		plan = append(plan, DayPlan{
			Date:     currentDate.Format("2006-01-02"),
			Day:      day + 1,
			Problems: dayProblems,
			Category: determineCategoryLabel(dayProblems),
		})

		// Update yesterday's categories for next day
		usedYesterday = usedToday
		day++
	}

	return plan
}

func determineCategoryLabel(problems []Problem) string {
	categories := make(map[string]bool)
	for _, p := range problems {
		categories[p.Category] = true
	}
	if len(categories) == 1 {
		for c := range categories {
			return c
		}
	}
	return "Mixed"
}
