package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func clearScreen() {
	fmt.Print("\033[2J")      // Clear screen
	fmt.Print("\033[H")       // Move cursor to top-left
	fmt.Print("\033[?25l")    // Hide cursor
}
func showWelcome() {
	clearScreen()
    fmt.Println("\033[1;36m") // Cyan bold
	welcome := `
╔══════════════════════════════════════╗
║                                      ║
║       🚀 CDK Deploy Wizard 🚀        ║
║                                      ║
╚══════════════════════════════════════╝
`
	fmt.Println(welcome)
	fmt.Println("Starting Please Wait...")
    fmt.Println("\033[0m")    // Reset
}

func main() {
    // Run `cdk ls` and capture output
    showWelcome()
    cmd := exec.Command("bash", "-c", "cdk ls 2>/dev/null")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        fmt.Println("Error running 'cdk ls':", err)
        return
    }

    rawStacks := out.String()
    stacks := parseStacks(rawStacks)
    if len(stacks) == 0 {
        fmt.Println("No stacks found.")
        return
    }

    // Add an "All stacks" option
    options := append(stacks, "All stacks")

    // Prompt user to select stacks (multi-select)
    var selected []string
    prompt := &survey.MultiSelect{
        Message: "Select stack(s) to operate on:",
        Options: options,
    }
    err = survey.AskOne(prompt, &selected)
    if err != nil {
        fmt.Println("Prompt failed:", err)
        return
    }

    if len(selected) == 0 {
        fmt.Println("No stacks selected, exiting.")
        return
    }

    // If "All stacks" selected, override selected with all stacks
    for _, s := range selected {
        if s == "All stacks" {
            selected = stacks
            break
        }
    }

    // Prompt user for action
    var action string
    actionPrompt := &survey.Select{
        Message: "Choose action:",
        Options: []string{"Deploy", "Destroy"},
    }
    err = survey.AskOne(actionPrompt, &action)
    if err != nil {
        fmt.Println("Prompt failed:", err)
        return
    }

  
    finalCmd:=fmt.Sprintf("cdk %s %s\n", strings.ToLower(action), strings.Join(selected, " "))
    fmt.Println("Executing >> ",finalCmd)
    cdkcmd := exec.Command("bash", "-c", finalCmd)
    cdkcmd.Stdout = os.Stdout
    cdkcmd.Stderr = os.Stderr
    cdkcmd.Stdin = os.Stdin
    err = cdkcmd.Run()
    if err != nil {
        fmt.Println("Error running ",finalCmd, err)
        return
    }
}

func parseStacks(raw string) []string {
    lines := strings.Split(raw, "\n")
    var stacks []string
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line != "" {
            stacks = append(stacks, line)
        }
    }
    return stacks
}
