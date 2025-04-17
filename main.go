package main

import (
	"fmt"
	"math"
	"sort"
)

type Entity struct {
	id   int
	name string
}

type Rating struct {
	score  int
	entity *Entity
}

type User struct {
	id      int
	name    string
	ratings []Rating
}

func (u User) getSimilarity(u2 User) float64 {

	if u.id == u2.id {
		return 0.0
	}

	// Everyone is equally similar if they have no ratings
	if len(u.ratings) == 0 || len(u2.ratings) == 0 {
		return 0.5
	}

	var dotProduct = 0.0
	var sumRatingsUser1 = 0.0
	var sumRatingsUser2 = 0.0

	for _, ratingsUser1 := range u.ratings {

		for _, ratingsUser2 := range u2.ratings {

			// If they're not the same Entity rating, move on
			if ratingsUser1.entity.id != ratingsUser2.entity.id {
				continue
			}

			dotProduct += float64(ratingsUser1.score) * float64(ratingsUser2.score)

			sumRatingsUser1 += float64(ratingsUser1.score) * float64(ratingsUser1.score)
			sumRatingsUser2 += float64(ratingsUser2.score) * float64(ratingsUser2.score)
		}
	}

	magnitude := math.Sqrt(sumRatingsUser1) * math.Sqrt(sumRatingsUser2)

	similarityScore := dotProduct / magnitude

	return similarityScore
}

func (u User) getPredictedScore(entity Entity, users []User, amountOfNeighbours int) float64 {

	// Consider each Entity the same if the user hasn't rated before
	if len(u.ratings) == 0 {
		return 0.5
	}

	// Consider each Entity the same if there's no other users to compare to
	if len(users) == 0 {
		return 0.5
	}

	// If the user has already rated this, return their rating
	for _, rating := range u.ratings {
		if rating.entity.id == entity.id {
			return float64(rating.score)
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

	if len(usersWithRating) == 0 {
		return 0.5
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

			recommendedRating += float64(rating.score) * neighbour.similarity
			break
		}
	}

	recommendedRating = recommendedRating / sumOfSimilarityScores

	return recommendedRating
}

func main() {
	amountOfNeighbours := 1

	entities := []Entity{
		{1, "Bombastic side-eye"},
		{2, "Banana bus"},
		{3, "Farting4Fortnite"},
	}

	users := make([]User, 0, 10)

	users = append(users, User{
		id:   1,
		name: "Jennifer",
		ratings: []Rating{
			{
				score:  1,
				entity: &entities[0],
			},
			{
				score:  5,
				entity: &entities[1],
			},
			{
				score:  2,
				entity: &entities[2],
			},
		},
	})

	users = append(users, User{
		id:   2,
		name: "Maghba",
		ratings: []Rating{
			{
				score:  5,
				entity: &entities[0],
			},
			{
				score:  1,
				entity: &entities[1],
			},
			{
				score:  4,
				entity: &entities[2],
			},
		},
	})

	users = append(users, User{
		id:   3,
		name: "Bennifer",
		ratings: []Rating{
			{
				score:  2,
				entity: &entities[0],
			},
			{
				score:  5,
				entity: &entities[1],
			},
		},
	})

	recommendedRating := users[2].getPredictedScore(entities[2], users, amountOfNeighbours)

	fmt.Println(recommendedRating)
}
