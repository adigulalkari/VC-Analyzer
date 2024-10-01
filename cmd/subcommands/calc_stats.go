package subcommands

import (
    "fmt"
    "github.com/spf13/cobra"
)

// CalcStatsCmd represents the calc-stats command
var CalcStatsCmd = &cobra.Command{
    Use:   "calc-stats",
    Short: "Calculate statistics for a repository",
    Long:  "The calc-stats command allows you to calculate various statistics for a specified Git repository in the owner/repo format.",
    Run: func(cmd *cobra.Command, args []string) {
        repository, _ := cmd.Flags().GetString("repository")
        fmt.Printf("Calculating stats for repository: %s\n", repository)

        
    },
}

func init() {
    // Add the --repository flag
    CalcStatsCmd.Flags().StringP("repository", "r", "", "Repository in owner/repo format")
    CalcStatsCmd.MarkFlagRequired("repository")
}
