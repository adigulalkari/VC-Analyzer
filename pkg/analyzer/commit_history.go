package analyzer

import (
    "fmt"
    "log"

    "github.com/go-git/go-git/v5"              // Core Go-git library
    "github.com/go-git/go-git/v5/plumbing/object" // Used for commit objects
)

func AnalyzeCommitHistory(repoPath string) {
    repo, err := git.PlainOpen(repoPath)
    if err != nil {
        log.Fatalf("Error opening repository: %v", err)
    }

    // Get the HEAD reference
    ref, err := repo.Head()
    if err != nil {
        log.Fatalf("Error getting HEAD reference: %v", err)
    }

    // Iterate over the commit history starting from HEAD
    commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
    if err != nil {
        log.Fatalf("Error getting commit log: %v", err)
    }

    fmt.Println("Commit history analysis:")
    commitCount := 0
    err = commitIter.ForEach(func(c *object.Commit) error { 
        fmt.Printf("Commit: %s by %s\n", c.Hash, c.Author.Name)
        commitCount++
        return nil
    })

    if err != nil {
        log.Fatalf("Error iterating over commits: %v", err)
    }

    fmt.Printf("Total number of commits: %d\n", commitCount)
}
