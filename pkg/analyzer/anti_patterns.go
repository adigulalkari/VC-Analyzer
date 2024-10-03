package analyzer

import (
    "fmt"
    "log"
    "time"

    "github.com/go-git/go-git/v5"
    "github.com/go-git/go-git/v5/plumbing/object"
)

// DetectAntiPatterns looks for common version control anti-patterns.
func DetectAntiPatterns(repoPath string) {
    // Open the Git repository
    repo, err := git.PlainOpen(repoPath)
    if err != nil {
        log.Fatalf("Error opening repository: %v", err)
    }

    // Get the HEAD reference
    ref, err := repo.Head()
    if err != nil {
        log.Fatalf("Error getting repository HEAD: %v", err)
    }

    // Get the commit history
    commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
    if err != nil {
        log.Fatalf("Error getting commit history: %v", err)
    }

    fmt.Println("Detecting anti-patterns...")

    // Initialize anti-pattern counters
    var largeCommitsCount int
    var forcePushesDetected bool
    var infrequentCommits bool

    // Iterate through the commits
    err = commitIter.ForEach(func(c *object.Commit) error {
        // Detect large commits
        if len(c.Message) > 1000 {
            largeCommitsCount++
        }

        // Simulate detecting force pushes (for demo purposes)
        // You can improve this logic by tracking ref changes
        if c.NumParents() > 1 {
            forcePushesDetected = true
        }

        return nil
    })
    if err != nil {
        log.Fatalf("Error iterating through commits: %v", err)
    }

    // Simulate checking for infrequent commits
    // Check the time between the most recent and older commits
    recentCommit, _ := repo.CommitObject(ref.Hash())
    recentTime := recentCommit.Committer.When

    err = commitIter.ForEach(func(c *object.Commit) error {
        if recentTime.Sub(c.Committer.When) > (7 * 24 * time.Hour) {
            infrequentCommits = true
        }
        return nil
    })
    if err != nil {
        log.Fatalf("Error iterating through commits: %v", err)
    }

    // Display anti-pattern results
    if largeCommitsCount > 0 {
        fmt.Printf("Detected %d large commit(s).\n", largeCommitsCount)
    } else {
        fmt.Println("No large commits detected.")
    }

    if forcePushesDetected {
        fmt.Println("Force pushes detected.")
    } else {
        fmt.Println("No force pushes detected.")
    }

    if infrequentCommits {
        fmt.Println("Detected infrequent commits (more than 7 days between commits).")
    } else {
        fmt.Println("No infrequent commit patterns detected.")
    }

    fmt.Println("Anti-pattern detection complete.")
}
