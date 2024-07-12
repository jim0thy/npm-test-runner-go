package main

import (
	"bufio"
	"fmt"
	"fyne.io/fyne/v2"
	"log"
	"os"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2/app"
)

// Configurable parameters
const (
	MaxConsecutiveGreens = 3
	MaxConsecutiveFails  = 3
	MaxRuns              = 20
	TestCommand          = "npm test" // Command to run Jest tests
)

var (
	totalRuns         int
	failedTests       int
	successfulTests   int
	consecutiveGreens int
	consecutiveFails  int
	fyneApp           fyne.App
)

func init() {
	// Initialize the fyne app for notifications
	fyneApp = app.New()
}

func sendNotification(testNumber int, result string, consecutiveGreens int, consecutiveFails int) {
	notification := fyne.NewNotification(
		"Test Notification",
		fmt.Sprintf("Test %d %s. Consecutive greens: %d. Consecutive fails: %d", testNumber, result, consecutiveGreens, consecutiveFails),
	)
	fyneApp.SendNotification(notification)
}

func displayStats() {
	log.Printf("Total tests run: %d\n", totalRuns)
	log.Printf("Failed tests: %d\n", failedTests)
	log.Printf("Successful tests: %d\n", successfulTests)
	log.Printf("Consecutive green runs: %d\n", consecutiveGreens)
	log.Printf("Consecutive failed runs: %d\n", consecutiveFails)
}

func runTests(testNumber int) string {
	log.Printf("Running tests %d...\n", testNumber)

	cmd := exec.Command("sh", "-c", TestCommand)
	outputBytes, err := cmd.CombinedOutput()
	output := string(outputBytes)

	if err != nil {
		log.Printf("Error running tests: %v", err)
		return "FAILED"
	}
	log.Println(output)

	// Check for Jest's failure message in the output
	if strings.Contains(output, "FAIL") {
		return "FAILED"
	}
	return "PASSED"
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for testNumber := 1; testNumber <= MaxRuns; testNumber++ {
		result := runTests(testNumber)
		totalRuns++
		if result == "FAILED" {
			failedTests++
			consecutiveFails++
			consecutiveGreens = 0
		} else {
			successfulTests++
			consecutiveGreens++
			consecutiveFails = 0
		}

		displayStats()
		sendNotification(testNumber, result, consecutiveGreens, consecutiveFails)

		if consecutiveGreens >= MaxConsecutiveGreens {
			log.Printf("Achieved %d consecutive green runs. Exiting.\n", MaxConsecutiveGreens)
			break
		}

		if consecutiveFails >= MaxConsecutiveFails {
			log.Printf("Encountered %d consecutive failed runs. Exiting.\n", MaxConsecutiveFails)
			break
		}

		fmt.Printf("Press Enter to run test %d or Ctrl+C to stop...", testNumber+1)
		_, _ = reader.ReadString('\n')
	}

	if totalRuns >= MaxRuns {
		log.Printf("Reached maximum number of test runs (%d). Exiting.\n", MaxRuns)
	}
}
