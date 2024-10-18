package analyzer

import (
	"fmt"
	"reflect"
	"testing"

	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
)

func Example_printCommitHistoryAnalysis() {
	commitCount := 5
	authorCommits := []authorCommit{
		{"Alice", 3},
		{"Bob", 2},
	}

	printCommitHistoryAnalysis(commitCount, authorCommits)

	// Output:
	// Commit history analysis:
	//
	// Total number of commits: 5
	//
	// Number of commits by each author (in decreasing order):
	// Alice: 3 commits
	// Bob: 2 commits
}

func Example_printBranchHistory() {
	branchesMap := map[string]string{
		"main": Active,
	}
	activeBranchCount := 1
	inactiveBranchCount := 0

	printBranchHistory(branchesMap, activeBranchCount, inactiveBranchCount, true)

	// Output:
	//
	// Branch analysis:
	//
	// Branches:
	// main: Active
	//
	// Active branches: 1
	// Inactive branches: 0
}

func Example_commitHistory() {
	// Create a new in-memory repository
	fs := memfs.New()
	repo, err := git.Init(memory.NewStorage(), fs)
	if err != nil {
		fmt.Printf("Failed to initialize in-memory repository: %v\n", err)
		return
	}
	_, err = createCommit(repo, time.Now())
	if err != nil {
		fmt.Printf("Failed to create commit for tests: %v\n", err)
		return
	}

	commitHistory(repo)

	// Output:
	// Commit history analysis:
	//
	// Total number of commits: 1
	//
	// Number of commits by each author (in decreasing order):
	// Test Author: 1 commits
}

func TestListBranches(t *testing.T) {
	// Create a new in-memory repository
	fs := memfs.New()
	repo, err := git.Init(memory.NewStorage(), fs)
	if err != nil {
		t.Fatalf("Failed to initialize in-memory repository: %v", err)
	}
	expectedBranches := map[string]string{
		"master": "Active",
		"foo":    "Inactive",
	}
	// Create a new commit in master branch
	lastCommit, err := createCommit(repo, time.Now())
	if err != nil {
		t.Fatalf("Failed to create commit for tests: %v", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		t.Fatalf("Failed to get worktree: %v", err)
	}

	targetBranch := plumbing.NewBranchReferenceName("foo")
	err = wt.Checkout(&git.CheckoutOptions{
		Hash:   lastCommit,
		Create: true,
		Branch: targetBranch,
	})
	if err != nil {
		t.Fatalf("Failed to create and checkout branch: %v", err)
	}

	// Create a new commit in foo branch
	_, err = createCommit(repo, time.Now().Add(-91*24*time.Hour))
	if err != nil {
		t.Fatalf("Failed to create commit for tests: %v", err)
	}

	branchesMap, activeBranchCount, inactiveBranchCount, err := getBranchCounts(repo)
	if err != nil {
		t.Fatalf("Failed to get branches for tests: %v", err)
	}

	// Check local branches
	if len(branchesMap) != 2 {
		t.Errorf("Expected 2 local branch , got %v", branchesMap)
	}
	if activeBranchCount != 1 {
		t.Errorf("Expected 1 active branch , got %v", activeBranchCount)
	}
	if inactiveBranchCount != 1 {
		t.Errorf("Expected 1 inactive branch , got %v", inactiveBranchCount)
	}
	if !reflect.DeepEqual(branchesMap, expectedBranches) {
		t.Errorf("Expected  branch %v, got %v", expectedBranches, branchesMap)
	}
}

func createCommit(repo *git.Repository, when time.Time) (plumbing.Hash, error) {
	wt, err := repo.Worktree()
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("Failed to get worktree: %w", err)
	}
	util.WriteFile(wt.Filesystem, "foo", []byte("foo"), 0644)
	_, err = wt.Add("foo")
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("Failed to add file: %w", err)
	}

	h, err := wt.Commit("Initial commit ", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test Author",
			Email: "test@example.com",
			When:  when,
		},
	})
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("Failed to create commit: %w", err)
	}
	return h, nil
}

func TestCountCommits(t *testing.T) {
	// Create a new in-memory repository
	fs := memfs.New()
	repo, err := git.Init(memory.NewStorage(), fs)
	if err != nil {
		t.Fatalf("Failed to initialize in-memory repository: %v", err)
	}

	_, err = createCommit(repo, time.Now())
	if err != nil {
		t.Fatalf("Failed to create commit for tests: %v", err)
	}

	// Count commits
	commitsPerAuthor, totalCommits, err := getCommitCounts(repo)
	if err != nil {
		t.Fatalf("Failed to count commits: %v", err)
	}

	// Check total commits
	if totalCommits != 1 {
		t.Errorf("Expected 1 commit, got %d", totalCommits)
	}

	if len(commitsPerAuthor) != 1 {
		t.Errorf("Expected 1 author, got %v", len(commitsPerAuthor))
	}
	// Check commits per author
	if commitsPerAuthor["Test Author"] != 1 {
		t.Errorf("Expected 1 commit for 'Test Author', got %v", commitsPerAuthor)
	}
}

func TestGetSortedAuthorCommits(t *testing.T) {
	tests := []struct {
		name         string
		commitCounts map[string]int
		expected     []authorCommit
	}{
		{
			name: "basic test",
			commitCounts: map[string]int{
				"Alice": 5,
				"Bob":   3,
				"Carol": 8,
			},
			expected: []authorCommit{
				{Author: "Carol", Count: 8},
				{Author: "Alice", Count: 5},
				{Author: "Bob", Count: 3},
			},
		},
		{
			name:         "empty map",
			commitCounts: map[string]int{},
			expected:     nil,
		},
		{
			name: "single author",
			commitCounts: map[string]int{
				"Alice": 5,
			},
			expected: []authorCommit{
				{Author: "Alice", Count: 5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getSortedAuthorCommits(tt.commitCounts)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func Example_printCommitSize() {
	printCommitSize(5, 250)
	// Output:
	// Total number of commits: 5
	// Total commit message size: 250 bytes
	// Average commit size: 50.00 bytes
}

func Example_commitSize() {
	// Create a new in-memory repository
	fs := memfs.New()
	repo, err := git.Init(memory.NewStorage(), fs)
	if err != nil {
		fmt.Printf("Failed to initialize in-memory repository: %v\n", err)
		return
	}

	_, err = createCommit(repo, time.Now())
	if err != nil {
		fmt.Printf("Failed to create commit for tests: %v\n", err)
		return
	}
	commitSize(repo)

	// Output:
	// Total number of commits: 1
	// Total commit message size: 15 bytes
	// Average commit size: 15.00 bytes
}

func TestCommitStats(t *testing.T) {
	// Create a new in-memory repository
	fs := memfs.New()
	repo, err := git.Init(memory.NewStorage(), fs)
	if err != nil {
		t.Fatalf("Failed to initialize in-memory repository: %v", err)
	}

	_, err = createCommit(repo, time.Now())
	if err != nil {
		t.Fatalf("Failed to create commit for tests: %v", err)
	}

	_, err = createCommit(repo, time.Now())
	if err != nil {
		t.Fatalf("Failed to create commit for tests: %v", err)
	}

	commitCount, totalSize, err := getCommitStats(repo)
	if err != nil {
		t.Errorf("Expected nil Error, got: %v", err)
	}
	// Check total commits
	if commitCount != 2 {
		t.Errorf("Expected 2 commits, got %d", commitCount)
	}

	// Check total commits size
	if totalSize != 30 {
		t.Errorf("Expected size 30, got %d", totalSize)
	}
}

func Example_branchStats() {
	// Create a new in-memory repository
	fs := memfs.New()
	repo, err := git.Init(memory.NewStorage(), fs)
	if err != nil {
		fmt.Printf("Failed to initialize in-memory repository: %v\n", err)
		return
	}

	_, err = createCommit(repo, time.Now())
	if err != nil {
		fmt.Printf("Failed to create commit for tests: %v\n", err)
		return
	}
	branchStats(repo, true)

	// Output:
	// Branch analysis:
	//
	// Branches:
	// master: Active
	//
	// Active branches: 1
	// Inactive branches: 0
}
