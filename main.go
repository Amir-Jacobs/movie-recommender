package main

import (
	"fmt"
	"errors"
	"math"
)

func cosineSimilarity(user1 []int, user2 []int) (float64, error) {
	if len(user1) != len(user2) {
		return 0.0, errors.New("vectors must be of same length")
	}

	var dotProduct float64 = 0.0
	var sum1 float64 = 0.0
	var sum2 float64 = 0.0

	for i := range len(user1) {
		
		if user1[i] == -1 {
			continue
		}

		if user2[i] == -1 {
			continue
		}
		
		dotProduct += float64(user1[i]) * float64(user2[i])

		sum1 += float64(user1[i]) * float64(user1[i])
		sum2 += float64(user2[i]) * float64(user2[i])
	}

	magnitude := math.Sqrt(sum1) * math.Sqrt(sum2)

	cosineSimilarity := dotProduct / magnitude

	return cosineSimilarity, nil
}

type RecommendedRating struct {
	index int
	rating float64
}

func createRecommendation(user []int, otherUsers [][]int) {
	var similarityScores []float64 = make([]float64, 0, len(otherUsers))  

	for _, otherUser := range otherUsers {
		similarityScore, _ := cosineSimilarity(user, otherUser)

		similarityScores = append(similarityScores, similarityScore)
	}
	
	
	var recommendedRatings []RecommendedRating = make([]RecommendedRating, 0, 10)

	for i, rating := range(user) {
		
		if rating != -1 {
			continue
		}

		var recommendedRating float64 = 0.0
		var sumOfSimilarityScores float64 = 0.0


		for _, otherUser := range otherUsers {
			if otherUser[i] == -1 {
				continue
			}

			sumOfSimilarityScores += similarityScores[i]
			recommendedRating += float64(otherUser[i]) * similarityScores[i]
		}

		recommendedRatings = append(recommendedRatings, RecommendedRating{index: i, rating: recommendedRating / sumOfSimilarityScores})
	}

	fmt.Println(recommendedRatings)
}

func main() {
	ratings := [][]int{
		{5, 1, 4, -1},
		{4, 1, 5, 1},
		{1, 5, 1, 4},
	}

	createRecommendation(ratings[0], [][]int{ratings[1], ratings[2]})
} 

