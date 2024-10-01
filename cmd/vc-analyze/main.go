package main

import (
    "os"
    "fmt"
    "github.com/spf13/cobra"
	"github.com/MakeNowJust/heredoc/v2"
)

var rootCmd = &cobra.Command{
    Use:     "vc-analyze",
    Version: "0.1.0",
    Short:   "A command-line tool for analyzing version control activity in Git repositories.",
}


var getCmd = &cobra.Command{
    Use:   "get <details>",
    Short: "Display one or many repositories details",
    Example: heredoc.Doc(`
        Get stars count for a given repository
        $ vc-analyze get stars -r golang/go
    `),
    Args: func(cmd *cobra.Command, args []string) error {
        return nil
    },
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Println("Hello from get command")
        return nil
    },
}


func main() {
    fmt.Println("<<-----VC-Analyzer----->>")

	rootCmd.AddCommand(getCmd)
	
    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}
