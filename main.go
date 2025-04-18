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
	start := time.Now()

	movies := extractMovies()

	users := extractRatings(1000, movies, 1.0, 5.0)

	elapsed := time.Since(start)

	fmt.Printf("Function took %s \n", elapsed)
	fmt.Printf("Amount of users: %d \n", len(users))

	usersList := make([]User, 0, len(users))

	for _, user := range users {
		usersList = append(usersList, *user)
	}

	start = time.Now()

	predictionsForUser := make(map[int64]map[int64]float64)

	for _, user := range users {
		index := 0

		predictionsForUser[user.id] = make(map[int64]float64)

		for _, movie := range movies {
			predictionsForUser[user.id][movie.id] = user.getPredictedScore(movie, usersList, 10)

			index += 1

			if index == 150 {
				break
			}
		}
	}

	elapsed = time.Since(start)

	fmt.Printf("Function took %s \n", elapsed)
	fmt.Println("Predictions done for amount of users: ", len(predictionsForUser))

	for id, predictions := range predictionsForUser {
		fmt.Printf("Predictions for user with id %d are:\n", id)

		for movieId, prediction := range predictions {
			if prediction < 4.0 {
				continue
			}

			fmt.Printf("Movie name: %s\n", movies[movieId].name)
			fmt.Printf("Predicted score: %.2f\n", prediction)
		}

		break
	}
}
