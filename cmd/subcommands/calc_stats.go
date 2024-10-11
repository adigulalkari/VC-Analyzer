package subcommands

import (
	"errors"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/adigulalkari/VC-Analyzer/pkg/analyzer"
)

var (
	authorStats bool
	commitSize  bool
)

var CalcStatsCmd = &cobra.Command{
	Use:   "calc-stats <path/to/repo>",
	Short: "Calculate statistics for the local repo",
	Long:  `This command allows you to calculate various statistics for a local Git repository, including author statistics and commit size statistics.`,
	Example: heredoc.Doc(`
        # Calculate author statistics
        $ vc-analyze calc-stats --author-stats path/to/local/repo

        # Calculate commit size statistics
        $ vc-analyze calc-stats --commit-size path/to/local/repo
    `),
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a path to the repository")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		repoPath := args[0] // Get the repository path from the arguments

		// Check if the repository exists (optional)
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			return fmt.Errorf("repository path does not exist: %s", repoPath)
		}

		// Check which flag is set and call the appropriate function
		if authorStats {
			// Call a function to calculate author statistics
			analyzer.AnalyzeCommitHistory(repoPath)
		} else if commitSize {
			// Call a function to calculate commit size statistics
			analyzer.AnalyzeCommitSize(repoPath)
		} else {
			return errors.New("no valid flag provided, use --author-stats or --commit-size")
		}

		return nil
	},
}

func init() {
    CalcStatsCmd.Flags().BoolVar(&authorStats, "author-stats", false, "Calculate statistics for each author")
    CalcStatsCmd.Flags().BoolVar(&commitSize, "commit-size", false, "Calculate the size of commits")
}