package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	targetHours   = 7
	targetMinutes = 5
	deadlineHour  = 16 // 4 pm
)

var targetDuration = time.Duration(targetHours)*time.Hour + time.Duration(targetMinutes)*time.Minute

func resolveInput(timeFlag string, args []string) (string, error) {
	hasFlagInput := timeFlag != ""
	hasArgsInput := len(args) >= 1

	if hasFlagInput && hasArgsInput {
		return "", fmt.Errorf("Ambiguous Input: Please provide time via -t flag, OR positional argument, not both")
	}
	if hasFlagInput {
		return timeFlag, nil
	}
	if hasArgsInput {
		return args[0], nil
	}

	// Fallback to interactive stdin
	fmt.Print("Enter your current stored time (e.g., 4:30): ")
	var input string
	if _, err := fmt.Scanln(&input); err != nil || strings.TrimSpace(input) == "" {
		return "", fmt.Errorf("No input provided")
	}
	return strings.TrimSpace(input), nil
}

func parseDuration(input string) (time.Duration, error) {
	input = strings.TrimSpace(input)
	parts := strings.Split(input, ":")

	if len(parts) != 2 {
		return 0, fmt.Errorf("Invalid Format %q — expected HH:MM (e.g. 4:30)", input)
	}

	h, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("Invalid Hours %q — must be a whole number", parts[0])
	}

	m, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("Invalid Minutes %q — must be a whole number", parts[1])
	}

	if h < 0 || h > 12 {
		return 0, fmt.Errorf("Hours must be between 1 and 12 (got %d)", h)
	}

	if m < 0 || m > 59 {
		return 0, fmt.Errorf("Minutes must be between 0 and 59 (got %d)", m)
	}

	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute, nil
}

func printResults(storedDuration time.Duration) {
	now := time.Now()
	deadline := time.Date(now.Year(), now.Month(), now.Day(), deadlineHour, 0, 0, 0, now.Location())

	remainingWindow := time.Until(deadline)
	deadlinePassed := remainingWindow < 0
	if deadlinePassed {
		remainingWindow = 0
	}

	projectedTotal := storedDuration + remainingWindow

	fmt.Println()
	fmt.Println("--- Results ---")
	fmt.Printf("Current Time:     %s\n", now.Format("15:04"))
	fmt.Printf("Stored Progress:  %s\n", formatDuration(storedDuration))
	fmt.Printf("Target:           %s\n", formatDuration(targetDuration))

	if deadlinePassed {
		fmt.Println("⚠️  Deadline:        4:00 PM has already passed")
	} else {
		fmt.Printf("Time until 4PM:   %s\n", formatDuration(remainingWindow.Round(time.Minute)))

	}

	fmt.Printf("Projected Total:  %s\n", formatDuration(projectedTotal.Round(time.Minute)))
	fmt.Println("----------------")

	if projectedTotal >= targetDuration {
		surplus := projectedTotal - targetDuration
		fmt.Printf("✅ You're on track! Surplus: %s\n", formatDuration(surplus.Round(time.Minute)))

	} else {
		shortfall := targetDuration - projectedTotal
		fmt.Printf("⚠️  You'll be short by: %s\n", formatDuration(shortfall.Round(time.Minute)))
	}
	fmt.Println()
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h == 0 {
		return fmt.Sprintf("%dm", m)
	}
	return fmt.Sprintf("%dh %dm", h, m)
}

func main() {

	// Accept using flag
	timeFlag := flag.String("t", "", "Current stored time (format HH:MM)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-t HH:MM] [time]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  -t HH:MM    Stored time as a flag\n")
		fmt.Fprintf(os.Stderr, "  time        Stored time as a positional argument\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  %s -t 4:30\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s 4:30\n", os.Args[0])
	}
	flag.Parse()

	// Take user input
	input, err := resolveInput(*timeFlag, flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}
	storedDuration, err := parseDuration(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	printResults(storedDuration)
}
