package main

import (
    "os"
    "fmt"

    "github.com/spf13/cobra"
    "github.com/adigulalkari/VC-Analyzer/cmd/subcommands"
)

var rootCmd = &cobra.Command{
    Use:     "vc-analyze",
    Version: "0.1.0",
    Short:   "A command-line tool for analyzing version control activity in Git repositories.",
}

func init() {
    subcommands.GetCmd.Flags().StringVarP(&subcommands.Repository, "repository", "r", "", "The GitHub repository in the format 'owner/repo'")
    rootCmd.AddCommand(subcommands.GetCmd)
}

func main() {
    fmt.Println("<<-----VC-Analyzer----->>")

    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}
