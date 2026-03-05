package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
)

const (
	targetHours   = 7
	targetMinutes = 5
	deadlineHour  = 16 // 4 pm
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")).Italic(true).Underline(true)
	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Bold(true)
	successStyle = lipgloss.NewStyle().Padding(1, 2).Foreground(lipgloss.Color("82")).Bold(false).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("70"))

	errorStyle   = lipgloss.NewStyle().Padding(1, 2).Foreground(lipgloss.Color("202")).Bold(false).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("202"))
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
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

	if projectedTotal >= targetDuration {
		surplus := projectedTotal - targetDuration
		msg := "Have some coffee my friend! You got extra: " + titleStyle.Render(formatDuration(surplus.Round(time.Minute)))
		fmt.Println(successStyle.Render(msg))

	} else {
		shortfall := targetDuration - projectedTotal
		msg := "Oops you are short by : " + titleStyle.Render(formatDuration(shortfall.Round(time.Minute)))
		fmt.Println(errorStyle.Render(msg))
	}
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
	showSpinner()
	printResults(storedDuration)
}

func PrintError(msg string) {
	fmt.Println(errorStyle.Render(msg))
}

func PrintSuccess(msg string) {
	fmt.Println(successStyle.Render(msg))
}

func showSpinner() {
	frames := []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}
	for i := 0; i < 12; i++ {
		frame := frames[i%len(frames)]
		fmt.Printf("\r%s %s...", spinnerStyle.Render(frame), valueStyle.Render("Calculating"))
		time.Sleep(30 * time.Millisecond)
	}
	fmt.Print("\r\033[K") // clear the spinner line
}
