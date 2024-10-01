package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

type RepositoryResponse struct {
	StargazersCount int64 `json:"stargazers_count"`
}

func GetStarCount(repoName string)(*int64, error){
	response, err := http.Get("https://api.github.com/repos/" + repoName)

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	
	var responseStruct RepositoryResponse
	err = json.Unmarshal(body, &responseStruct)
	if err != nil {
		return nil, err
	}

	return &responseStruct.StargazersCount, nil
}