package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func extractMovies() map[int64]Entity {
	f, err := os.Open("data/movies.csv")

	if err != nil {
		log.Fatal("Unable to read input file ", err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	movies := make(map[int64]Entity)

	for {
		record, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Unable to read line")
			continue
		}

		movieId, err := strconv.ParseInt(record[0], 10, 64)

		if err != nil {
			continue
		}

		movies[movieId] = Entity{
			id:   movieId,
			name: record[1],
		}
	}

	return movies
}

func extractRatings(maxUsers int, uniqueMovies map[int64]Entity, minRating float64, maxRating float64) map[int64]*User {
	f, err := os.Open("data/ratings.csv")

	if err != nil {
		log.Fatal("Unable to read input file ", err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)

	uniqueUsersWithRatings := make(map[int64]*User)

	for {
		record, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Unable to read line")
			continue
		}

		userId, err := strconv.ParseInt(record[0], 10, 64)

		if err != nil {
			continue
		}

		movieId, err := strconv.ParseInt(record[1], 10, 64)

		if err != nil {
			continue
		}

		// If movie doesn't exist don't include the rating
		movie := uniqueMovies[movieId]

		if movie.id == 0 {
			continue
		}

		rating, err := strconv.ParseFloat(record[2], 64)

		if err != nil {
			continue
		}

		if rating < minRating || rating > maxRating {
			continue
		}

		// If max users is reached, don't add new ones
		if uniqueUsersWithRatings[userId] == nil && len(uniqueUsersWithRatings) >= maxUsers {

			// Am able to do this because the CSV is sorted by userId
			break
		}

		// If user didn't exist yet
		if uniqueUsersWithRatings[userId] == nil {
			ratings := make(map[int64]float64)

			ratings[movieId] = rating

			uniqueUsersWithRatings[userId] = &User{
				id:      userId,
				ratings: ratings,
			}

			continue
		}

		uniqueUsersWithRatings[userId].ratings[movieId] = rating
	}

	return uniqueUsersWithRatings
}

func main() {
	amountOfNeighbours := 100_000
	minimumSimilarityThreshold := 0.02
	amountOfUsers := 100_000

	minRating := 1.0
	maxRating := 5.0

	amountOfRecommendedEntities := 50
	minimumRatingForRecommendation := 4.5

	start := time.Now()

	movies := extractMovies()

	users := extractRatings(amountOfUsers, movies, minRating, maxRating)

	elapsed := time.Since(start)

	fmt.Printf("Function took %s \n", elapsed)
	fmt.Printf("Amount of users: %d \n", len(users))

	usersList := make([]User, 0, len(users))

	for _, user := range users {
		usersList = append(usersList, *user)
	}

	//	start = time.Now()
	//
	//	predictionsForUser := make(map[int64]map[int64]float64)
	//
	//recommend:
	//	for _, user := range users {
	//		predictionsForUser[user.id] = make(map[int64]float64)
	//
	//		for _, movie := range movies {
	//			predictedScore := user.getPredictedScore(movie, usersList, amountOfNeighbours, minimumSimilarityThreshold)
	//
	//			if predictedScore < minimumRatingForRecommendation {
	//				continue
	//			}
	//
	//			predictionsForUser[user.id][movie.id] = predictedScore
	//
	//			if len(predictionsForUser[user.id]) >= amountOfRecommendedEntities {
	//				continue recommend
	//			}
	//		}
	//	}
	//
	//	elapsed = time.Since(start)
	//
	//	fmt.Printf("Function took %s \n", elapsed)
	//	fmt.Println("Predictions done for amount of users: ", len(predictionsForUser))
	//
	//	for id, predictions := range predictionsForUser {
	//		fmt.Printf("Predictions for user with id %d are:\n", id)
	//
	//		for movieId, prediction := range predictions {
	//			fmt.Printf("Movie name: %s\n", movies[movieId].name)
	//			fmt.Printf("Predicted score: %.2f\n", prediction)
	//		}
	//
	//		fmt.Println("")
	//		fmt.Println("")
	//	}

	myRatings := make(map[int64]float64)

	myRatings[8368] = 5.0
	myRatings[40815] = 5.0
	myRatings[55768] = 5.0

	amir := User{
		id:      200000,
		ratings: myRatings,
	}

	predictionsForUser := make(map[int64]map[int64]float64)

	predictionsForUser[amir.id] = make(map[int64]float64)

	for _, movie := range movies {
		predictedScore := amir.getPredictedScore(movie, usersList, amountOfNeighbours, minimumSimilarityThreshold)

		if predictedScore < minimumRatingForRecommendation {
			continue
		}

		predictionsForUser[amir.id][movie.id] = predictedScore

		if len(predictionsForUser[amir.id]) >= amountOfRecommendedEntities {
			break
		}
	}

	for id, predictions := range predictionsForUser {
		fmt.Printf("Predictions for user with id %d are:\n", id)

		for movieId, prediction := range predictions {
			fmt.Printf("Movie name: %s\n", movies[movieId].name)
			fmt.Printf("Predicted score: %.2f\n", prediction)
		}

		fmt.Println("")
		fmt.Println("")
	}
}
