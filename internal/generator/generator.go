package generator

import (
	"fmt"
	"time"
)

// GeneratePlan creates a 125-day study plan including all 250 problems
func GeneratePlan(problems []Problem, startDate time.Time) []DayPlan {
	categoryProblems := organizeProblemsByCategoryAndDifficulty(problems)
	usedProblems := make(map[string]bool)
	categoryProgress := initCategoryProgress()

	var plan []DayPlan
	var categoryUsageHistory []string

	day := 0
	const maxDays = 150

	for len(usedProblems) < 250 && day < maxDays {
		currentDate := startDate.AddDate(0, 0, day)
		selectedCategory := selectCategoryForDay(day, categoryProblems, categoryUsageHistory, categoryProgress, usedProblems)

		if selectedCategory == "" {
			fmt.Printf("Warning: No category available on day %d, %d problems remaining\n", day+1, 250-len(usedProblems))
			break
		}

		dayProblems := selectProblemsForDay(day, selectedCategory, categoryProblems, usedProblems)
		if len(dayProblems) == 0 {
			fmt.Printf("Warning: No problems found for %s on day %d\n", selectedCategory, day+1)
			day++
			continue
		}

		plan = append(plan, DayPlan{
			Date:     currentDate.Format("2006-01-02"),
			Day:      day + 1,
			Problems: dayProblems,
			Category: selectedCategory,
		})

		categoryUsageHistory = updateHistory(categoryUsageHistory, selectedCategory)
		day++
	}

	// Add any remaining problems
	if len(usedProblems) < 250 {
		plan = appendRemainingProblems(plan, categoryProblems, usedProblems, startDate, day, maxDays)
	}

	return plan
}

func initCategoryProgress() CategoryProgress {
	progress := make(CategoryProgress)
	for _, cat := range CategoryOrder() {
		progress[cat] = map[string]int{"Easy": 0, "Medium": 0, "Hard": 0}
	}
	return progress
}

func updateHistory(history []string, category string) []string {
	newHistory := []string{category}
	for _, cat := range history {
		if cat != category {
			newHistory = append(newHistory, cat)
		}
	}
	if len(newHistory) > 10 {
		newHistory = newHistory[:10]
	}
	return newHistory
}

func appendRemainingProblems(plan []DayPlan, categoryProblems CategoryProblems, usedProblems map[string]bool, startDate time.Time, day, maxDays int) []DayPlan {
	var remaining []Problem
	for _, category := range CategoryOrder() {
		for _, difficulty := range difficulties {
			for _, problem := range categoryProblems[category][difficulty] {
				if !usedProblems[problem.Name] {
					remaining = append(remaining, problem)
				}
			}
		}
	}

	for len(remaining) > 0 && day < maxDays {
		currentDate := startDate.AddDate(0, 0, day)
		var dayProblems []Problem

		for i := 0; i < 2 && len(remaining) > 0; i++ {
			problem := remaining[0]
			remaining = remaining[1:]
			dayProblems = append(dayProblems, problem)
			usedProblems[problem.Name] = true
		}

		if len(dayProblems) > 0 {
			plan = append(plan, DayPlan{
				Date:     currentDate.Format("2006-01-02"),
				Day:      day + 1,
				Problems: dayProblems,
				Category: determineCategoryLabel(dayProblems),
			})
			day++
		}
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
