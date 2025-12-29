package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
	"neetcode-study-plan/internal/generator"
)

const outputDir = "output"

var startFlag = flag.String("start", "", "Start date (YYYY-MM-DD, 'today', or 'monday')")

func parseDate(input string) (time.Time, bool) {
	input = strings.TrimSpace(strings.ToLower(input))

	switch input {
	case "", "monday":
		return getNextMonday(), true
	case "today":
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), true
	default:
		parsed, err := time.Parse("2006-01-02", input)
		if err != nil {
			return time.Time{}, false
		}
		return parsed, true
	}
}

func promptStartDate() time.Time {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("📅 When would you like to start your 125-day study plan?")
	fmt.Println("Examples:")
	fmt.Println("  - 2025-01-01 (New Year)")
	fmt.Println("  - today (starts today)")
	fmt.Println("  - monday (starts next Monday)")
	fmt.Println("  - Press Enter for default (next Monday)")

	fmt.Print("\nEnter start date (YYYY-MM-DD, 'today', 'monday', or Enter for default): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("\n📅 Using default start date (next Monday)")
		return getNextMonday()
	}

	if date, ok := parseDate(input); ok {
		return date
	}

	fmt.Println("❌ Invalid date format. Using default (next Monday)")
	return getNextMonday()
}

func getNextMonday() time.Time {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	daysAhead := int(time.Monday - today.Weekday())
	if daysAhead <= 0 {
		daysAhead += 7
	}
	return today.AddDate(0, 0, daysAhead)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func main() {
	flag.Parse()

	fmt.Println("🔧 Generating 125-day plan with ALL 250 problems...")

	var startDate time.Time
	if *startFlag != "" {
		if date, ok := parseDate(*startFlag); ok {
			startDate = date
		} else {
			fmt.Printf("❌ Invalid date: %s\n", *startFlag)
			os.Exit(1)
		}
	} else {
		startDate = promptStartDate()
	}
	fmt.Printf("📅 Study plan will start on: %s\n", startDate.Format("Monday, January 02, 2006"))

	// Load categories
	categoriesPath, err := generator.FindCategoriesFile()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if err := generator.LoadCategories(categoriesPath); err != nil {
		fmt.Printf("Error loading categories: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("📂 Loaded %d categories\n", len(generator.Categories))

	// Find and load problems
	dataPath, err := generator.FindDataFile()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	allProblems, err := generator.LoadProblems(dataPath)
	if err != nil {
		fmt.Printf("Error loading problems: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("📚 Loaded %d total problems\n", len(allProblems))

	fmt.Println("\n🗓️ Generating complete plan...")
	plan := generator.GeneratePlan(allProblems, startDate)

	totalProblemsInPlan := generator.CountProblems(plan)
	fmt.Printf("\n✅ Plan includes %d problems out of 250 total\n", totalProblemsInPlan)

	if totalProblemsInPlan < 250 {
		fmt.Printf("❌ WARNING: Missing %d problems!\n", 250-totalProblemsInPlan)
		return
	}

	fmt.Println("\n📊 Analyzing plan distribution...")
	_, difficultyByPhase := generator.AnalyzePlanDistribution(plan)

	markdownContent := generator.GenerateMarkdown(plan)

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Find next available filename
	baseFilename := fmt.Sprintf("neetcode-250-study-plan-%s", startDate.Format("2006-01-02"))
	filename := filepath.Join(outputDir, baseFilename+".md")
	counter := 1

	for fileExists(filename) {
		filename = filepath.Join(outputDir, fmt.Sprintf("%s-%d.md", baseFilename, counter))
		counter++
	}

	if err := os.WriteFile(filename, []byte(markdownContent), 0644); err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Generated complete 125-day plan with %d days\n", len(plan))
	fmt.Printf("📄 Saved to: %s\n", filename)

	// Summary statistics
	fmt.Printf("\n📈 Plan Statistics:\n")
	fmt.Printf("  Total days: %d\n", len(plan))
	fmt.Printf("  Total problems: %d\n", totalProblemsInPlan)
	fmt.Printf("  Average problems per day: %.1f\n", float64(totalProblemsInPlan)/float64(len(plan)))

	// Category focus statistics
	sameCategoryDays := 0
	for _, day := range plan {
		if day.Category != "Mixed" {
			sameCategoryDays++
		}
	}
	mixedCategoryDays := len(plan) - sameCategoryDays

	fmt.Printf("\n🎯 Category Focus:\n")
	fmt.Printf("  Same category days: %d (%.1f%%)\n", sameCategoryDays, float64(sameCategoryDays)/float64(len(plan))*100)
	fmt.Printf("  Mixed category days: %d (%.1f%%)\n", mixedCategoryDays, float64(mixedCategoryDays)/float64(len(plan))*100)

	// Difficulty distribution by phase
	fmt.Printf("\n📊 Difficulty Distribution by Phase:\n")
	for _, phase := range []string{"Early (1-30)", "Middle (31-80)", "Late (81-125)"} {
		difficulties := difficultyByPhase[phase]
		totalPhase := 0
		for _, count := range difficulties {
			totalPhase += count
		}
		if totalPhase > 0 {
			fmt.Printf("  %s:\n", phase)
			for _, diff := range []string{"Easy", "Medium", "Hard"} {
				count := difficulties[diff]
				percentage := float64(count) / float64(totalPhase) * 100
				fmt.Printf("    %s: %d (%.1f%%)\n", diff, count, percentage)
			}
		}
	}

	// Verify all problems are included
	missingProblems := generator.FindMissingProblems(allProblems, plan)

	if len(missingProblems) > 0 {
		sort.Strings(missingProblems)
		fmt.Printf("\n❌ Missing Problems (%d):\n", len(missingProblems))
		for _, problem := range missingProblems {
			fmt.Printf("  - %s\n", problem)
		}
	} else {
		fmt.Println("\n✅ All 250 problems successfully included!")
	}
}
