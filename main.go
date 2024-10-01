package main

import (
	"fmt"
	"VC-Analyzer/pkg/analyzer"
)


func main(){
	fmt.Println("-----VC-Analyzer-----")

	repoPath := "path/to/repo"

	analyzer.AnalyzeCommitHistory(repoPath)

	analyzer.DetectBottlenecks(repoPath)

	analyzer.DetectAntiPatterns(repoPath)

	
}
