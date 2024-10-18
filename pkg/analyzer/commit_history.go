package analyzer

import (
	"fmt"
	"log"
	"sort"

	"time"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5" // Core Go-git library
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object" // Used for commit objects
)

// AnalyzeCommitHistory analyzes and prints commit history of the given repository
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

	// Map to track the number of commits by each author
	commitCounts := make(map[string]int)
	commitCount := 0

	fmt.Println("Commit history analysis:")
	err = commitIter.ForEach(func(c *object.Commit) error {

		// Increment commit count for the author
		commitCounts[c.Author.Name]++
		commitCount++
		return nil
	})

	if err != nil {
		log.Fatalf("Error iterating over commits: %v", err)
	}

	// Print total number of commits
	fmt.Printf("\nTotal number of commits: %d\n", commitCount)

	// Sort authors by the number of commits in descending order
	type authorCommit struct {
		Author string
		Count  int
	}

	// Create a slice of author commits for sorting
	var authorCommits []authorCommit
	for author, count := range commitCounts {
		authorCommits = append(authorCommits, authorCommit{Author: author, Count: count})
	}

	// Sort the slice by commit count in descending order
	sort.Slice(authorCommits, func(i, j int) bool {
		return authorCommits[i].Count > authorCommits[j].Count
	})

	// Print the sorted list of authors and their commit counts
	fmt.Println("\nNumber of commits by each author (in decreasing order):")
	for _, ac := range authorCommits {
		fmt.Printf("%s: %d commits\n", ac.Author, ac.Count)
	}

	// List all branches and their activity status
	fmt.Println("\nBranches:")
	branches, err := repo.Branches()
	if err != nil {
		log.Fatalf("Error getting branches: %v", err)
	}

	activeBranchCount := 0
	inactiveBranchCount := 0

	branches.ForEach(func(ref *plumbing.Reference) error {
		commit, err := repo.CommitObject(ref.Hash())
		if err != nil {
			log.Fatalf("Error getting commit object: %v", err)
		}

		// Determine if the branch is active based on last commit date
		isActive := time.Since(commit.Committer.When) < 90*24*time.Hour

		branchName := ref.Name().Short()
		branchStatus := "Inactive"
		statusColor := color.New(color.FgRed)

		if isActive {
			branchStatus = "Active"
			statusColor = color.New(color.FgGreen)
			activeBranchCount++
		} else {
			inactiveBranchCount++
		}

		fmt.Printf("%s: ", branchName)
		statusColor.Printf("%s\n", branchStatus)

		return nil
	})

	fmt.Printf("\nActive branches: %d\n", activeBranchCount)
	fmt.Printf("Inactive branches: %d\n", inactiveBranchCount)
}

func AnalyzeCommitSize(repoPath string) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatalf("Error opening repository: %v", err)
	}

	//Get the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		log.Fatalf("Error getting HEAD reference: %v", err)
	}

	//Iterate over the commit history starting from HEAD
	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Fatalf("Error getting commit log: %v", err)
	}

	totalSize := 0
	commitCount := 0

	err = commitIter.ForEach(func(c *object.Commit) error {
		totalSize += len(c.Message) // Approximate commit size as the message length
		commitCount++
		return nil
	})

	if err != nil {
		log.Fatalf("Error iterating over commits: %v", err)
	}

	// Print commit size statistics
	fmt.Printf("\nTotal number of commits: %d\n", commitCount)
	fmt.Printf("Total commit message size: %d bytes\n", totalSize)
	fmt.Printf("Average commit size: %.2f bytes\n", float64(totalSize)/float64(commitCount))
}
