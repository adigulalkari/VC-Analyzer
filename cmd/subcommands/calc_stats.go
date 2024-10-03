package subcommands

import (
    "errors"
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/MakeNowJust/heredoc/v2"

    "github.com/adigulalkari/VC-Analyzer/pkg/analyzer" 
)

var CalcStatsCmd = &cobra.Command{ 
    Use:     "calc-stats <detail>",
    Short:   "Calculate the number of commits for the local repo",
    Example: heredoc.Doc(`
        Calculate the number of commits for the local repo
        $ vc-analyze calc-stats path/to/local/repo
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
        analyzer.AnalyzeCommitHistory(repoPath)

        return nil
    },
}
