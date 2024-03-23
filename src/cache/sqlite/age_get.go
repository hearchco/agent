package sqlite

import "time"

func (db DB) GetAge(query string) (time.Duration, error) {
	tstmp, err := db.queries.GetResultsTTLByQuery(db.ctx, query)
	if err != nil {
		return 0, err
	}

	return time.Since(tstmp), nil
}
