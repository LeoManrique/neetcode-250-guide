# NeetCode 250 Study Plan Generator

A CLI tool that generates a 125-day study plan covering all 250 NeetCode problems.

## Usage

```bash
# Build
go build -o bin/neetcode-plan ./cmd/neetcode-plan

# Run with start date
./bin/neetcode-plan --start 2025-01-06
./bin/neetcode-plan --start today
./bin/neetcode-plan --start monday

# Run interactively (prompts for date)
./bin/neetcode-plan
```

Or run directly:
```bash
go run ./cmd/neetcode-plan --start today
```

Output is saved to `output/`.

## Features

- 2 problems per day, same category when possible
- Progressive difficulty (Easy → Medium → Hard)
- Spaced repetition for category cycling
- Checkboxes and LeetCode links for tracking

## Project Structure

```
├── cmd/neetcode-plan/    # CLI entry point
├── internal/generator/   # Core logic
├── data/                 # Problem & category data
└── output/               # Generated plans
```

## License

MIT
