# Recommender

A collaborative filtering movie recommender written in Go. It predicts how a user would rate unseen movies by finding similar users (nearest neighbours) and combining their ratings weighted by similarity.

## How It Works

1. **Data loading** — Reads movie metadata from `data/movies.csv` and user ratings from `data/ratings.csv` ([MovieLens](https://grouplens.org/datasets/movielens/) CSV format).
2. **User similarity** — Computes cosine similarity between users based on their shared ratings.
3. **Rating prediction** — For a given user and movie, selects the _k_ most similar neighbours who rated that movie and produces a weighted average as the predicted score.
4. **Recommendations** — Iterates over all movies and surfaces the top predictions above a configurable threshold.

## Data Format

Place two CSV files under a `data/` directory:

| File | Columns |
|---|---|
| `movies.csv` | `movieId, title, genres` |
| `ratings.csv` | `userId, movieId, rating, timestamp` |

The ratings CSV must be sorted by `userId`.

## Configuration

Tuneable parameters are set at the top of `main()`:

| Parameter | Default | Description |
|---|---|---|
| `amountOfNeighbours` | 100 000 | Max neighbours to consider per prediction |
| `minimumSimilarityThreshold` | 0.02 | Minimum cosine similarity to qualify as a neighbour |
| `amountOfUsers` | 100 000 | Max users to load from the dataset |
| `minRating` / `maxRating` | 1.0 / 5.0 | Valid rating range (rows outside this are skipped) |
| `amountOfRecommendedEntities` | 50 | Number of recommendations to generate |
| `minimumRatingForRecommendation` | 4.5 | Predicted score threshold for a recommendation |

## Running

```sh
go run .
```

Requires Go 1.24+ and the CSV data files described above.

## Project Structure

```
main.go          – Data loading (movies & ratings) and main recommendation loop
recommender.go   – Core types (Entity, User) and recommendation logic
go.mod           – Module definition
```
