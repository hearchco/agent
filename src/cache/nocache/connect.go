package nocache

import (
	"time"

	"github.com/hearchco/hearchco/src/search/result"
)

type DB struct{}

func Connect() (DB, error) {
	return DB{}, nil
}

func (db DB) Close() error {
	return nil
}

func (db DB) GetResults(query string) ([]result.Result, error) {
	return nil, nil
}

func (db DB) GetImageResults(query string, salt string) ([]result.Result, error) {
	return nil, nil
}

func (db DB) SetResults(query string, results []result.Result) error {
	return nil
}

func (db DB) SetImageResults(query string, results []result.Result) error {
	return nil
}

func (db DB) GetAge(query string) (time.Duration, error) {
	return 0, nil
}
