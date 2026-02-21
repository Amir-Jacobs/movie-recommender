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

func (u User) averageRating() float64 {
	if len(u.ratings) == 0 {
		return 0
	}
	sum := 0.0
	for _, r := range u.ratings {
		sum += r
	}
	return sum / float64(len(u.ratings))
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

	for movieId := range u.ratings {
		if _, ok := u2.ratings[movieId]; !ok {
			continue
		}

		sharedRatings = append(sharedRatings, movieId)
	}

	if len(sharedRatings) == 0 {
		return 0.01
	}

	var dotProduct = 0.0
	var sumRatingsUser1 = 0.0
	var sumRatingsUser2 = 0.0

	for _, movieId := range sharedRatings {
		dotProduct += u.ratings[movieId] * u2.ratings[movieId]

		sumRatingsUser1 += u.ratings[movieId] * u.ratings[movieId]
		sumRatingsUser2 += u2.ratings[movieId] * u2.ratings[movieId]
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

	// If the user hasn't rated before, use the average rating from other users
	if len(u.ratings) == 0 {
		sum := 0.0
		count := 0
		for _, other := range users {
			if r, ok := other.ratings[entity.id]; ok {
				sum += r
				count++
			}
		}
		if count == 0 {
			return 0
		}
		return sum / float64(count)
	}

	// Consider each Entity the same if there's no other users to compare to
	if len(users) == 0 {
		return u.averageRating()
	}

	// If the user has already rated this, return their rating
	if rating, ok := u.ratings[entity.id]; ok {
		return rating
	}

	// Filter out any users that haven't rated this Entity and filter out the current user
	usersWithRating := make([]User, 0, len(users))

	for _, otherUser := range users {

		// Filter out current user
		if otherUser.id == u.id {
			continue
		}

		// Filter out user without rating for Entity
		if _, ok := otherUser.ratings[entity.id]; !ok {
			continue
		}

		usersWithRating = append(usersWithRating, otherUser)
	}

	// If there's no Users that ranked this Entity
	if len(usersWithRating) == 0 {
		return u.averageRating()
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
		if similarity >= minimumSimilarityThreshold {
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
