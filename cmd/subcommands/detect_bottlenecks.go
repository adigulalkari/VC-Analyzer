package subcommands

import (
    "errors"
    "fmt"
    "os"
    "os/exec"
    "strings"

    "github.com/spf13/cobra"
    "github.com/MakeNowJust/heredoc/v2"
)

// Struct to hold commit information
type CommitInfo struct {
    Hash     string
    Author   string
    Date     string
    Message  string
    Files    map[string]int // File changes with number of lines added/removed
}

// Function to gather commit history data
func getCommitHistory(repoPath string) ([]CommitInfo, error) {
    // Navigate to the repository path
    cmd := exec.Command("git", "-C", repoPath, "log", "--pretty=format:%H,%an,%ad,%s", "--numstat")
    output, err := cmd.CombinedOutput()
    if err != nil {
        return nil, fmt.Errorf("failed to run git command: %v", err)
    }

    // Parse the output
    lines := strings.Split(string(output), "\n")
    var commits []CommitInfo
    var currentCommit *CommitInfo
    for _, line := range lines {
        if strings.Contains(line, ",") {
            // New commit entry
            fields := strings.SplitN(line, ",", 4)
            currentCommit = &CommitInfo{
                Hash:    fields[0],
                Author:  fields[1],
                Date:    fields[2],
                Message: fields[3],
                Files:   make(map[string]int),
            }
            commits = append(commits, *currentCommit)
        } else if currentCommit != nil {
            // File changes for the current commit
            fileFields := strings.Fields(line)
            if len(fileFields) == 3 {
                // Update file change details
                currentCommit.Files[fileFields[2]] = 1 // Mark file as changed
            }
        }
    }
    return commits, nil
}

// Function to identify bottlenecks based on commit history
func detectBottlenecks(commits []CommitInfo) {
    fileChangeCounts := make(map[string]int)
    for _, commit := range commits {
        for file := range commit.Files {
            fileChangeCounts[file]++
        }
    }

    // Display files with the highest change frequency
    fmt.Println("Potential bottleneck files (most frequently changed):")
    for file, count := range fileChangeCounts {
        if count > 2 { // Arbitrary threshold for example
            fmt.Printf("%s: %d changes\n", file, count)
        }
    }
}

var DetectBottlenecksCmd = &cobra.Command{
    Use:   "detect-bottlenecks <repository-path>",
    Short: "Find bottlenecks in the commit history of a local repository",
    Example: heredoc.Doc(`
        $ vc-analyze detect-bottlenecks path/to/local/repo
    `),
    Args: func(cmd *cobra.Command, args []string) error {
        if len(args) < 1 {
            return errors.New("requires a repository path argument")
        }
        return nil
    },
    RunE: func(cmd *cobra.Command, args []string) error {
        repoPath := args[0] // Get the repository path from the arguments

        // Check if the repository exists
        if _, err := os.Stat(repoPath); os.IsNotExist(err) {
            return fmt.Errorf("repository path does not exist: %s", repoPath)
        }

        // Get the commit history for the repository
        commits, err := getCommitHistory(repoPath)
        if err != nil {
            return fmt.Errorf("error fetching commit history: %v", err)
        }

        // Detect bottlenecks based on the commit history
        detectBottlenecks(commits)

        return nil
    },
}
