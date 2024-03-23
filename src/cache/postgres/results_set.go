package postgres

import (
	"fmt"

	"github.com/hearchco/hearchco/src/search/result"
)

func (db DB) SetResults(query string, results []result.Result) error {
	for _, r := range results {
		// create new result in DB and retrieves it's ID
		id, err := db.queries.AddResult(db.ctx, AddResultParams{
			Query:       query,
			Url:         r.URL,
			Rank:        int64(r.Rank),
			Score:       r.Score,
			Title:       r.Title,
			Description: r.Description,
		})

		if err != nil {
			return fmt.Errorf("failed to add result: %w", err)
		}

		// adds all engine ranks to the result's ID foreign key
		for _, er := range r.EngineRanks {
			err := db.queries.AddEngineRank(db.ctx, AddEngineRankParams{
				EngineName:       er.SearchEngine.String(),
				EngineRank:       int64(er.Rank),
				EnginePage:       int64(er.Page),
				EngineOnPageRank: int64(er.OnPageRank),
				ResultID:         id,
			})

			if err != nil {
				return fmt.Errorf("failed to add engine rank: %w", err)
			}
		}
	}

	return nil
}

func (db DB) SetImageResults(query string, results []result.Result) error {
	for _, r := range results {
		// create new result in DB and retrieves it's ID
		id, err := db.queries.AddResult(db.ctx, AddResultParams{
			Query:       query,
			Url:         r.URL,
			Rank:        int64(r.Rank),
			Score:       r.Score,
			Title:       r.Title,
			Description: r.Description,
		})

		if err != nil {
			return fmt.Errorf("failed to add result: %w", err)
		}

		err = db.queries.AddImageResult(db.ctx, AddImageResultParams{
			ImageOriginalHeight:  int64(r.ImageResult.Original.Height),
			ImageOriginalWidth:   int64(r.ImageResult.Original.Width),
			ImageThumbnailHeight: int64(r.ImageResult.Thumbnail.Height),
			ImageThumbnailWidth:  int64(r.ImageResult.Thumbnail.Width),
			ImageThumbnailUrl:    r.ImageResult.ThumbnailURL,
			ImageSource:          r.ImageResult.Source,
			ImageSourceUrl:       r.ImageResult.SourceURL,
			ResultID:             id,
		})

		if err != nil {
			return fmt.Errorf("failed to add image result: %w", err)
		}

		// adds all engine ranks to the result's ID foreign key
		for _, er := range r.EngineRanks {
			err := db.queries.AddEngineRank(db.ctx, AddEngineRankParams{
				EngineName:       er.SearchEngine.String(),
				EngineRank:       int64(er.Rank),
				EnginePage:       int64(er.Page),
				EngineOnPageRank: int64(er.OnPageRank),
				ResultID:         id,
			})

			if err != nil {
				return fmt.Errorf("failed to add engine rank: %w", err)
			}
		}
	}

	return nil
}
