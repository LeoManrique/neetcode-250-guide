package generator

import (
	"fmt"
	"strings"
)

// GenerateMarkdown creates the markdown content for the study plan
func GenerateMarkdown(plan []DayPlan) string {
	var sb strings.Builder

	sb.WriteString("# NeetCode 250 - Complete 125-Day Study Plan (All 250 Problems)\n\n")
	sb.WriteString("**Enhanced Study Strategy:**\n")
	sb.WriteString("- 2 problems per day for 125 days (both from the same category when possible)\n")
	sb.WriteString("- ALL 250 problems included with no gaps\n")
	sb.WriteString("- Spaced repetition: Categories cycle with intelligent spacing to optimize retention\n")
	sb.WriteString("- Progressive difficulty: Early days focus on Easy, gradually increasing complexity\n")
	sb.WriteString("- Category focus: Daily concentration on single topics for deeper pattern recognition\n")
	sb.WriteString("- Intelligent timing: Easier categories appear earlier, advanced topics later in the plan\n\n")
	sb.WriteString("---\n\n")

	difficultyEmoji := map[string]string{
		"Easy":   "🟢",
		"Medium": "🟡",
		"Hard":   "🔴",
	}

	for _, dayPlan := range plan {
		sb.WriteString(fmt.Sprintf("## Day %d - %s\n", dayPlan.Day, dayPlan.Date))
		sb.WriteString(fmt.Sprintf("**Topic:** %s\n\n", dayPlan.Category))
		sb.WriteString("**Problems:**\n")

		for _, problem := range dayPlan.Problems {
			emoji := difficultyEmoji[problem.Difficulty]
			sb.WriteString(fmt.Sprintf("- [ ] %s [%s](%s) - *%s*\n",
				emoji, problem.Name, problem.LeetcodeURL, problem.Category))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
