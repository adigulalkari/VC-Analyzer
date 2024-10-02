package subcommands

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

type RepositoryInfo struct {
	FullName        string `json:"full_name"`
	Description     string `json:"description"`
	Stars           int    `json:"stargazers_count"`
	Forks           int    `json:"forks_count"`
	OpenIssuesCount int    `json:"open_issues_count"`
	Language        string `json:"language"`
}

// Fetch repository data from GitHub API
func fetchRepositoryData(repo string) (*RepositoryInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s", repo)
	log.Printf("Making request to URL: %s", url) // Log the URL

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Log the response status
	log.Printf("Response Status: %s", resp.Status)

	// Check if the request was successful
	if resp.StatusCode == http.StatusForbidden {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("rate-limited or forbidden: %s\nResponse: %s", resp.Status, string(body))
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch repository data: %s", resp.Status)
	}

	// Parse the JSON response
	var repoInfo RepositoryInfo
	if err := json.NewDecoder(resp.Body).Decode(&repoInfo); err != nil {
		return nil, fmt.Errorf("failed to parse repository data: %v", err)
	}

	return &repoInfo, nil
}

// calcStatsCmd represents the calc-stats command
var CalcStatsCmd = &cobra.Command{
	Use:   "calc-stats",
	Short: "Calculate statistics for a repository",
	Long:  "The calc-stats command fetches statistics for a specified Git repository in the owner/repo format.",
	Run: func(cmd *cobra.Command, args []string) {
		repository, _ := cmd.Flags().GetString("repository")

		if !strings.Contains(repository, "/") {
			log.Fatalf("Invalid repository format. Expected 'owner/repo'.")
		}

		// Fetch repository data
		repoData, err := fetchRepositoryData(repository)
		if err != nil {
			log.Fatalf("Error fetching repository data: %v", err)
		}

		// Display the repository statistics
		fmt.Printf("Repository: %s\nDescription: %s\nStars: %d\nForks: %d\nOpen Issues: %d\nLanguage: %s\n",
			repoData.FullName, repoData.Description, repoData.Stars, repoData.Forks, repoData.OpenIssuesCount, repoData.Language)
	},
}

func init() {
	CalcStatsCmd.Flags().StringP("repository", "r", "", "Repository in owner/repo format example: -r owner/reponame")
	CalcStatsCmd.MarkFlagRequired("repository")
}
