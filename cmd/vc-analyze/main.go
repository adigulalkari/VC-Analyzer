package main

import (
    "os"

    "github.com/spf13/cobra"
    "github.com/adigulalkari/VC-Analyzer/cmd/subcommands"
    "github.com/common-nighthawk/go-figure"
)

var rootCmd = &cobra.Command{
    Use:     "vc-analyze",
    Version: "0.1.0",
    Short:   "A command-line tool for analyzing version control activity in Git repositories.",
    
    CompletionOptions: cobra.CompletionOptions{
        DisableDefaultCmd: true,  // This line disables the 'completion' command
    },
}

func init() {
    subcommands.GetCmd.Flags().StringVarP(&subcommands.Repository, "repository", "r", "", "The GitHub repository in the format 'owner/repo'")

    rootCmd.AddCommand(subcommands.GetCmd)
    rootCmd.AddCommand(subcommands.CalcStatsCmd)
    rootCmd.AddCommand(subcommands.AntiPatternsCmd)
    rootCmd.AddCommand(subcommands.DetectBottlenecksCmd)
    
}

func main() {
    myFigure := figure.NewFigure("VC-Analyze", "", true)
    myFigure.Print()

    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}
