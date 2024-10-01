package main

import (
    "os"
    "fmt"
    "github.com/spf13/cobra"
	// "github.com/MakeNowJust/heredoc/v2"
)

var rootCmd = &cobra.Command{
    Use:     "vc-analyze",
    Version: "0.1.0",
    Short:   "A command-line tool for analyzing version control activity in Git repositories.",
}




func main() {
    fmt.Println("<<-----VC-Analyzer----->>")

    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}
