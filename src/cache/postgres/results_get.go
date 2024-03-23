package postgres

import (
	"fmt"

	"github.com/hearchco/hearchco/src/anonymize"
	"github.com/hearchco/hearchco/src/search/engines"
	"github.com/hearchco/hearchco/src/search/result"
)

func createResultFromRow(row GetResultsByQueryWithEngineRanksNotOlderThanTimestampRow) (result.Result, error) {
	engName, err := engines.NameString(row.EngineName)
	if err != nil {
		return result.Result{}, err
	}

	return result.Result{
		URL: row.Url,
		// URLHash not needed for general results
		Rank:        uint(row.Rank),
		Score:       row.Score,
		Title:       row.Title,
		Description: row.Description,
		ImageResult: result.ImageResult{}, // not needed for general results
		EngineRanks: []result.RetrievedRank{
			{
				SearchEngine: engName,
				Rank:         uint(row.EngineRank),
				Page:         uint(row.EnginePage),
				OnPageRank:   uint(row.EngineOnPageRank),
			},
		},
	}, nil
}

func (db DB) GetResults(query string) ([]result.Result, error) {
	rows, err := db.queries.GetResultsByQueryWithEngineRanksNotOlderThanTimestamp(db.ctx, GetResultsByQueryWithEngineRanksNotOlderThanTimestampParams{
		Query:     query,
		CreatedAt: db.Timestamp(),
	})
	if err != nil {
		return []result.Result{}, fmt.Errorf("failed to get results: %w", err)
	}

	started := false
	var currentID int64
	var newResult result.Result
	var results []result.Result

	for _, row := range rows {
		if !started {
			// if on first row, create the new result and set the currentID to the row's ID
			started = true
			currentID = row.ID

			newResult, err = createResultFromRow(row)
			if err != nil {
				return []result.Result{}, fmt.Errorf("failed to create result from row: %w", err)
			}
		} else if row.ID != currentID {
			// if not on the first row and the ID is different from the currentID
			// set the currentID to the new row's ID and append the new result to the results slice
			// since we are done with the previous result's engine ranks
			results = append(results, newResult)
			currentID = row.ID

			newResult, err = createResultFromRow(row)
			if err != nil {
				return []result.Result{}, fmt.Errorf("failed to create result from row: %w", err)
			}
		} else {
			// otherwise, append the engine rank to the new result
			engName, err := engines.NameString(row.EngineName)
			if err != nil {
				return []result.Result{}, fmt.Errorf("failed to get engine name: %w", err)
			}

			newResult.EngineRanks = append(newResult.EngineRanks, result.RetrievedRank{
				SearchEngine: engName,
				Rank:         uint(row.EngineRank),
				Page:         uint(row.EnginePage),
				OnPageRank:   uint(row.EngineOnPageRank),
			})
		}
	}

	return results, nil
}

func createImageResultFromRow(row GetImageResultsByQueryWithEngineRanksNotOlderThanTimestampRow, salt string) (result.Result, error) {
	engName, err := engines.NameString(row.EngineName)
	if err != nil {
		return result.Result{}, err
	}

	return result.Result{
		URL:         row.Url,
		URLHash:     anonymize.HashToSHA256B64Salted(row.Url, salt),
		Rank:        uint(row.Rank),
		Score:       row.Score,
		Title:       row.Title,
		Description: row.Description,
		ImageResult: result.ImageResult{
			Original: result.ImageFormat{
				Height: uint(row.ImageOriginalHeight),
				Width:  uint(row.ImageOriginalWidth),
			},
			Thumbnail: result.ImageFormat{
				Height: uint(row.ImageThumbnailHeight),
				Width:  uint(row.ImageThumbnailWidth),
			},
			ThumbnailURL:     row.ImageThumbnailUrl,
			ThumbnailURLHash: anonymize.HashToSHA256B64Salted(row.ImageThumbnailUrl, salt),
			Source:           row.ImageSource,
			SourceURL:        row.ImageSourceUrl,
		},
		EngineRanks: []result.RetrievedRank{
			{
				SearchEngine: engName,
				Rank:         uint(row.EngineRank),
				Page:         uint(row.EnginePage),
				OnPageRank:   uint(row.EngineOnPageRank),
			},
		},
	}, nil
}

func (db DB) GetImageResults(query string, salt string) ([]result.Result, error) {
	rows, err := db.queries.GetImageResultsByQueryWithEngineRanksNotOlderThanTimestamp(db.ctx, GetImageResultsByQueryWithEngineRanksNotOlderThanTimestampParams{
		Query:     query,
		CreatedAt: db.Timestamp(),
	})
	if err != nil {
		return []result.Result{}, fmt.Errorf("failed to get image results: %w", err)
	}

	started := false
	var currentID int64
	var newResult result.Result
	var results []result.Result

	for _, row := range rows {
		if !started {
			// if on first row, create the new result and set the currentID to the row's ID
			started = true
			currentID = row.ID

			newResult, err = createImageResultFromRow(row, salt)
			if err != nil {
				return []result.Result{}, fmt.Errorf("failed to create image result from row: %w", err)
			}

		} else if row.ID != currentID {
			// if not on the first row and the ID is different from the currentID
			// set the currentID to the new row's ID and append the new result to the results slice
			// since we are done with the previous result's engine ranks
			results = append(results, newResult)
			currentID = row.ID

			newResult, err = createImageResultFromRow(row, salt)
			if err != nil {
				return []result.Result{}, fmt.Errorf("failed to create image result from row: %w", err)
			}

		} else {
			// otherwise, append the engine rank to the new result
			engName, err := engines.NameString(row.EngineName)
			if err != nil {
				return []result.Result{}, fmt.Errorf("failed to get engine name: %w", err)
			}

			newResult.EngineRanks = append(newResult.EngineRanks, result.RetrievedRank{
				SearchEngine: engName,
				Rank:         uint(row.EngineRank),
				Page:         uint(row.EnginePage),
				OnPageRank:   uint(row.EngineOnPageRank),
			})
		}
	}

	return results, nil
}
