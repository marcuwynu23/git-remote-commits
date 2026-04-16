package main

import (
	"fmt"
	"git-remote-commits/git"
	"git-remote-commits/model"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	remoteName, exitCode, shouldExit := parseArgs(os.Args[1:])
	if shouldExit {
		os.Exit(exitCode)
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: unable to read current directory")
		os.Exit(1)
	}
	if err := git.IsGitRepo(wd); err != nil {
		fmt.Fprintln(os.Stderr, "Error: current directory is not a git repository.")
		fmt.Fprintln(os.Stderr, "Run this command inside a git repository.")
		os.Exit(1)
	}

	p := tea.NewProgram(model.Initial(remoteName), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func parseArgs(args []string) (remoteName string, exitCode int, shouldExit bool) {
	remoteName = "origin"
	if len(args) == 0 {
		return remoteName, 0, false
	}
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "Error: too many arguments.")
		printUsage()
		return "", 1, true
	}

	arg := strings.TrimSpace(args[0])
	switch arg {
	case "", "origin":
		return "origin", 0, false
	case "-h", "--help":
		printUsage()
		return "", 0, true
	default:
		if strings.HasPrefix(arg, "-") {
			fmt.Fprintf(os.Stderr, "Error: unknown flag %q.\n", arg)
			printUsage()
			return "", 1, true
		}
		return arg, 0, false
	}
}

func printUsage() {
	fmt.Println("Usage: git-remote-commits [remote]")
	fmt.Println("")
	fmt.Println("Args:")
	fmt.Println("  remote    Remote name to compare against current branch (default: origin)")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -h, --help    Show this help message")
}
