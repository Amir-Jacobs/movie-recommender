package main

import (
	"log"
	"math"
	"sort"
)

type Entity struct {
	id   int64
	name string
}

type Rating struct {
	score  float64
	entity *Entity
}

type User struct {
	id      int64
	name    string
	ratings []Rating
}

func (u User) getSimilarity(u2 User) float64 {

	if u.id == u2.id {
		return 1.0
	}

	// Everyone is equally similar if they have no ratings
	if len(u.ratings) == 0 || len(u2.ratings) == 0 {
		return 0.5
	}

	var dotProduct = 0.0
	var sumRatingsUser1 = 0.0
	var sumRatingsUser2 = 0.0

	var usersShareRatings bool

	for _, ratingsUser1 := range u.ratings {

		for _, ratingsUser2 := range u2.ratings {

			// If they're not the same Entity rating, move on
			if ratingsUser1.entity.id != ratingsUser2.entity.id {
				continue
			}

			usersShareRatings = true

			dotProduct += ratingsUser1.score * ratingsUser2.score

			sumRatingsUser1 += ratingsUser1.score * ratingsUser1.score
			sumRatingsUser2 += ratingsUser2.score * ratingsUser2.score
		}
	}

	if !usersShareRatings {
		return 0.01
	}

	magnitude := math.Sqrt(sumRatingsUser1) * math.Sqrt(sumRatingsUser2)

	similarityScore := dotProduct / magnitude

	if math.IsNaN(similarityScore) {
		log.Fatal("\n", u.ratings, "\n", u2.ratings)
	}

	return similarityScore
}

func (u User) getPredictedScore(entity Entity, users []User, amountOfNeighbours int) float64 {

	// Consider each Entity the same if the user hasn't rated before
	if len(u.ratings) == 0 {
		return 2.5
	}

	// Consider each Entity the same if there's no other users to compare to
	if len(users) == 0 {
		return 2.5
	}

	// If the user has already rated this, return their rating
	for _, rating := range u.ratings {
		if rating.entity.id == entity.id {
			return rating.score
		}
	}

	// Filter out any users that haven't rated this Entity and filter out the current user
	usersWithRating := make([]User, 0, len(users))

	for _, otherUser := range users {

		if otherUser.id == u.id {
			continue
		}

		for _, rating := range otherUser.ratings {
			if rating.entity.id != entity.id {
				continue
			}

			usersWithRating = append(usersWithRating, otherUser)
			break
		}
	}

	// If there's no Users that ranked this Entity
	if len(usersWithRating) == 0 {
		return 2.5
	}

	// Get the similarity score for each user
	type similarityWithUser struct {
		user       *User
		similarity float64
	}

	similarityScoresWithUsers := make([]similarityWithUser, 0, len(usersWithRating))

	for _, userWithRating := range usersWithRating {
		similarityScoresWithUsers = append(similarityScoresWithUsers, similarityWithUser{
			user:       &userWithRating,
			similarity: u.getSimilarity(userWithRating),
		})
	}

	// Sort the similarity scores to be descending, and only use the closest neighbours for rating
	sort.SliceStable(similarityScoresWithUsers, func(i, j int) bool {
		return similarityScoresWithUsers[i].similarity > similarityScoresWithUsers[j].similarity
	})

	if len(similarityScoresWithUsers) < amountOfNeighbours {
		amountOfNeighbours = len(similarityScoresWithUsers)
	}

	closestNeighbours := similarityScoresWithUsers[:amountOfNeighbours]

	recommendedRating := 0.0
	sumOfSimilarityScores := 0.0

	for _, neighbour := range closestNeighbours {
		sumOfSimilarityScores += neighbour.similarity

		for _, rating := range neighbour.user.ratings {
			if rating.entity.id != entity.id {
				continue
			}

			recommendedRating += rating.score * neighbour.similarity
			break
		}
	}

	recommendedRating = recommendedRating / sumOfSimilarityScores

	return recommendedRating
}
