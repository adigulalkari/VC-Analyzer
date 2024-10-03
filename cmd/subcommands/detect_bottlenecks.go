package subcommands

import (
    "errors"
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/MakeNowJust/heredoc/v2"

    "github.com/adigulalkari/VC-Analyzer/pkg/analyzer" 
)

var DetectBottlenecksCmd = &cobra.Command{ 
    Use:     "detect-bottlenecks <detail>",
    Short:   "Find out the detect-bottlenecks present in your repository",
    Example: heredoc.Doc(`
        Find out the detect-bottlenecks present in your repository
        $ vc-analyze detect-bottlenecks path/to/local/repo
    `),
    Args: func(cmd *cobra.Command, args []string) error {
        if len(args) < 1 {
            return errors.New("requires a detail argument")
        }
        return nil
    },
    RunE: func(cmd *cobra.Command, args []string) error {
        repoPath := args[0] // Get the repository path from the arguments

        // Check if the repository exists (optional)
        if _, err := os.Stat(repoPath); os.IsNotExist(err) {
            return fmt.Errorf("repository path does not exist: %s", repoPath)
        }

        // Call the AnalyzeCommitHistory function
        analyzer.DetectBottlenecks(repoPath)

        return nil
    },
}
