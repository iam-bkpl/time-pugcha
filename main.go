package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	timeFlag := flag.String("t", "4:30", "Current stored time (format HH:MM)")
	versionFlag := flag.Bool("v", false, "Print version information")

	flag.Parse()

	fmt.Println("Time Flag: ", *timeFlag)
	fmt.Println("Version Flag: ", *versionFlag)

	var input string

	// 1. Handle Input: Check Args first, otherwise ask the user
	if len(os.Args) > 1 {
		input = strings.TrimSpace(os.Args[1])
	} else {
		fmt.Print("Enter your current stored time (e.g., 4:30): ")
		fmt.Scanln(&input)
	}

	// 2. Parse the input string into a Duration
	// Expecting format "HH:MM"
	storedDuration, err := parseDuration(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid format: %v. Please use HH:MM (e.g. 4:30)\n", err)
		os.Exit(1)
	}

	// 3. Define the Constants
	now := time.Now()
	targetTotal := 7*time.Hour + 5*time.Minute
	// Deadline is 4:00 PM today
	deadline := time.Date(now.Year(), now.Month(), now.Day(), 16, 0, 0, 0, now.Location())

	// 4. Calculate Remaining Time until 4 PM
	remainingWindow := time.Until(deadline)
	if remainingWindow < 0 {
		remainingWindow = 0
		fmt.Println("⚠️ Note: The 4:00 PM deadline has already passed.")
	}

	// 5. Logic: Stored + Remaining vs Target
	projectedTotal := storedDuration + remainingWindow

	fmt.Println("--- Results ---")
	fmt.Printf("Current Time:     %s\n", now.Format("03:04 PM"))
	fmt.Printf("Stored Progress:  %v\n", storedDuration)
	fmt.Printf("Time until 4PM:   %v\n", remainingWindow.Round(time.Minute))
	fmt.Printf("Projected Total:  %v\n", projectedTotal.Round(time.Minute))
	fmt.Println("----------------")

	if projectedTotal >= targetTotal {
		surplus := projectedTotal - targetTotal
		fmt.Printf("🥳 You're on track! Surplus: %v\n", surplus.Round(time.Minute))
	} else {
		shortfall := targetTotal - projectedTotal
		fmt.Printf("😢 You'll be short by: %v\n", shortfall.Round(time.Minute))
	}
}

// parseDuration converts a "4:30" string into a time.Duration
func parseDuration(input string) (time.Duration, error) {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("must be in HH:MM format")
	}

	var h, m int
	_, err := fmt.Sscanf(input, "%d:%d", &h, &m)
	if err != nil {
		return 0, err
	}

	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute, nil
}
