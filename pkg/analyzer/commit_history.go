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

type authorCommit struct {
	Author string
	Count  int
}

var (
	InActive string = "Inactive"
	Active   string = "Active"
)

// AnalyzeCommitHistory analyzes and prints commit history of the given repository
func AnalyzeCommitHistory(repoPath string) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatalf("Error opening repository: %v", err)
	}

	commitHistory(repo)
}

func commitHistory(repo *git.Repository) {
	// Get all authors and their commit count
	commitCounts, commitCount, err := getCommitCounts(repo)
	if err != nil {
		log.Fatalf("Error getting commit counts: %v", err)
	}
	authorCommits := getSortedAuthorCommits(commitCounts)
	printCommitHistoryAnalysis(commitCount, authorCommits)
}

func getSortedAuthorCommits(commitCounts map[string]int) []authorCommit {
	// Create a slice of author commits
	// for sorting by the number of commits in descending order
	var authorCommits []authorCommit
	for author, count := range commitCounts {
		authorCommits = append(authorCommits, authorCommit{Author: author, Count: count})
	}

	// Sort the slice by commit count in descending order
	sort.Slice(authorCommits, func(i, j int) bool {
		return authorCommits[i].Count > authorCommits[j].Count
	})

	return authorCommits
}

func getCommitCounts(repo *git.Repository) (map[string]int, int, error) {
	// Get the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		return nil, 0, fmt.Errorf("Error getting HEAD reference: %w", err)
	}

	// Iterate over the commit history starting from HEAD
	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, 0, fmt.Errorf("Error getting commit log: %w", err)
	}

	// Map to track the number of commits by each author
	commitCounts := make(map[string]int)
	commitCount := 0

	err = commitIter.ForEach(func(c *object.Commit) error {

		// Increment commit count for the author
		commitCounts[c.Author.Name]++
		commitCount++
		return nil
	})

	if err != nil {
		return nil, 0, fmt.Errorf("Error iterating over commits: %w", err)
	}

	return commitCounts, commitCount, nil
}

func printCommitHistoryAnalysis(commitCount int, authorCommits []authorCommit) {
	fmt.Println("Commit history analysis:")

	// Print total number of commits
	fmt.Printf("\nTotal number of commits: %d\n", commitCount)

	// Print the sorted list of authors and their commit counts
	fmt.Println("\nNumber of commits by each author (in decreasing order):")
	for _, ac := range authorCommits {
		fmt.Printf("%s: %d commits\n", ac.Author, ac.Count)
	}
}

func AnalyzeCommitSize(repoPath string) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatalf("Error opening repository: %v", err)
	}

	commitSize(repo)
}

func commitSize(repo *git.Repository) {
	commitCount, totalSize, err := getCommitStats(repo)
	if err != nil {
		log.Fatalf("Error getting commit stats: %v", err)
	}

	// Print commit size statistics
	printCommitSize(commitCount, totalSize)
}

func getCommitStats(repo *git.Repository) (int, int, error) {
	// Get the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		return 0, 0, fmt.Errorf("Error getting HEAD reference: %w", err)
	}

	// Iterate over the commit history starting from HEAD
	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return 0, 0, fmt.Errorf("Error getting commit log: %w", err)
	}

	totalSize := 0
	commitCount := 0

	err = commitIter.ForEach(func(c *object.Commit) error {
		totalSize += len(c.Message) // Approximate commit size as the message length
		commitCount++
		return nil
	})

	if err != nil {
		return 0, 0, fmt.Errorf("Error iterating over commits: %w", err)
	}
	return commitCount, totalSize, nil
}

func printCommitSize(commitCount int, totalSize int) {
	fmt.Printf("\nTotal number of commits: %d\n", commitCount)
	fmt.Printf("Total commit message size: %d bytes\n", totalSize)
	fmt.Printf("Average commit size: %.2f bytes\n", float64(totalSize)/float64(commitCount))
}

// AnalyzeBranchStats analyzes and prints branch stats of the given repository
func AnalyzeBranchStats(repoPath string) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatalf("Error opening repository: %v", err)
	}

	branchStats(repo, false)
}

func branchStats(repo *git.Repository, noColor bool) {
	// List all branches and their activity status
	branchesMap, activeBranchCount, inactiveBranchCount, err := getBranchCounts(repo)
	if err != nil {
		log.Fatalf("Error getting branch history: %v", err)
	}
	printBranchHistory(branchesMap, activeBranchCount, inactiveBranchCount, noColor)
}

func getBranchCounts(repo *git.Repository) (map[string]string, int, int, error) {
	branches, err := repo.Branches()
	if err != nil {
		return nil, 0, 0, fmt.Errorf("Error getting branches: %w", err)
	}

	branchesMap := make(map[string]string)
	activeBranchCount := 0
	inactiveBranchCount := 0

	err = branches.ForEach(func(ref *plumbing.Reference) error {
		commit, err := repo.CommitObject(ref.Hash())
		if err != nil {
			return err
		}

		branchName := ref.Name().Short()

		// Determine if the branch is active based on last commit date
		isActive := time.Since(commit.Committer.When) < 90*24*time.Hour

		branchStatus := InActive
		if isActive {
			branchStatus = Active
			activeBranchCount++
		} else {
			inactiveBranchCount++
		}

		branchesMap[branchName] = branchStatus
		return nil
	})

	if err != nil {
		return nil, 0, 0, fmt.Errorf("Error getting commit object: %w", err)
	}

	return branchesMap, activeBranchCount, inactiveBranchCount, nil
}

func printBranchHistory(branchesMap map[string]string, activeBranchCount int, inactiveBranchCount int, noColor bool) {
	fmt.Println("Branch analysis:")
	fmt.Println("\nBranches:")
	for branchName, branchStatus := range branchesMap {
		fmt.Printf("%s: ", branchName)
		statusColor := color.New(color.FgRed)
		if branchStatus == Active {
			statusColor = color.New(color.FgGreen)
		}
		if noColor {
			fmt.Printf("%s\n", branchStatus)
		} else {
			statusColor.Printf("%s\n", branchStatus)
		}
	}
	fmt.Printf("\nActive branches: %d\n", activeBranchCount)
	fmt.Printf("Inactive branches: %d\n", inactiveBranchCount)
}
