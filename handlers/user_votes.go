package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Struct to match the API response
type ApiResponse struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
	Data       []struct {
		City          string `json:"city"`
		Name          string `json:"name"`
		EstimatedCost int32  `json:"estimated_cost"`
		UserRating    struct {
			AverageRating float64 `json:"average_rating"`
			Votes         int32   `json:"votes"`
		} `json:"user_rating"`
	} `json:"data"`
}

func getVoteCount(cityName string, estimatedCost int32) int32 {
	// Construct the API URL
	url := fmt.Sprintf("https://jsonmock.hackerrank.com/api/food_outlets?city=%s&estimated_cost=%d", cityName, estimatedCost)

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return 0
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return 0
	}

	// Parse the JSON response
	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return 0
	}

	// Check if any data matches the criteria
	if len(apiResponse.Data) == 0 {
		fmt.Println("No matching restaurant found.")
		return 0
	}

	var voteCount int32
	for _, data := range apiResponse.Data {
		voteCount += data.UserRating.Votes
	}

	// Return the vote count of the first matching restaurant
	return voteCount
}

type VoteRequest struct {
	City          string `json:"city"`
	EstimatedCost int32  `json:"estimatedCost"`
}

// Struct to match the response body
type VoteResponse struct {
	City      string `json:"city"`
	UserCount int32  `json:"usercount"`
}

// Handler function for /user-vote endpoint
func GetVoteCountHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var request VoteRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the GetVoteCount function from the handler package
	voteCount := getVoteCount(request.City, request.EstimatedCost)

	// Prepare the response
	response := VoteResponse{
		City:      request.City,
		UserCount: voteCount,
	}

	// Send the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
