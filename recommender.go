package main

import (
	"math"
	"sort"
)

type Entity struct {
	id   int64
	name string
}

type User struct {
	id      int64
	ratings map[int64]float64 // uses the id of the Entity!
}

func (u User) getSimilarity(u2 User) float64 {

	if u.id == u2.id {
		return 1.0
	}

	// Everyone is equally similar if they have no ratings
	if len(u.ratings) == 0 || len(u2.ratings) == 0 {
		return 0.5
	}

	// If the users don't share ratings, don't calculate anything; they're not similar
	sharedRatings := make([]int64, 0, len(u.ratings))

	for movieId, _ := range u.ratings {
		if u2.ratings[movieId] == 0 {
			continue
		}

		sharedRatings = append(sharedRatings, movieId)
		break
	}

	if len(sharedRatings) == 0 {
		return 0.01
	}

	var dotProduct = 0.0
	var sumRatingsUser1 = 0.0
	var sumRatingsUser2 = 0.0

	for id := range sharedRatings {
		dotProduct += u.ratings[int64(id)] * u2.ratings[int64(id)]

		sumRatingsUser1 += u.ratings[int64(id)] * u.ratings[int64(id)]
		sumRatingsUser2 += u2.ratings[int64(id)] * u2.ratings[int64(id)]
	}

	magnitude := math.Sqrt(sumRatingsUser1) * math.Sqrt(sumRatingsUser2)

	similarityScore := dotProduct / magnitude

	// todo: figure out how this becomes NaN
	if math.IsNaN(similarityScore) {
		return 0.01
	}

	return similarityScore
}

func (u User) getPredictedScore(entity Entity, users []User, amountOfNeighbours int, minimumSimilarityThreshold float64) float64 {

	if amountOfNeighbours == 0 {
		return 0.01
	}

	// Consider each Entity the same if the user hasn't rated before
	if len(u.ratings) == 0 {
		return 2.5
	}

	// Consider each Entity the same if there's no other users to compare to
	if len(users) == 0 {
		return 2.5
	}

	// If the user has already rated this, return their rating
	if u.ratings[entity.id] != 0 {
		return u.ratings[entity.id]
	}

	// Filter out any users that haven't rated this Entity and filter out the current user
	usersWithRating := make([]User, 0, len(users))

	for _, otherUser := range users {

		// Filter out current user
		if otherUser.id == u.id {
			continue
		}

		// Filter out user without rating for Entity
		if otherUser.ratings[entity.id] == 0 {
			continue
		}

		usersWithRating = append(usersWithRating, otherUser)
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
	minimumSimilarityUsers := 0

	for _, userWithRating := range usersWithRating {
		similarity := u.getSimilarity(userWithRating)

		similarityScoresWithUsers = append(similarityScoresWithUsers, similarityWithUser{
			user:       &userWithRating,
			similarity: similarity,
		})

		// Count the amount of users that reached minimum similarity threshold
		if minimumSimilarityThreshold >= similarity {
			minimumSimilarityUsers += 1
		}

		// Once enough users did, stop comparing for more
		if amountOfNeighbours == minimumSimilarityUsers {
			break
		}
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

		recommendedRating += neighbour.user.ratings[entity.id] * neighbour.similarity
	}

	recommendedRating = recommendedRating / sumOfSimilarityScores

	return recommendedRating
}
