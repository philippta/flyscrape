package flyscrape

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ScrapeResult struct {
	URL        string      `json:"url"`
	Body       []byte      `json:"body"`
	StatusCode int         `json:"status_code"`
	Header     http.Header `json:"header"`
	Timestamp  time.Time   `json:"timestamp"`
}

func (r ScrapeResult) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *ScrapeResult) Scan(v any) error {
	switch vt := v.(type) {
	case []byte:
		return json.Unmarshal(vt, r)
	case string:
		return json.Unmarshal([]byte(vt), r)
	default:
		return fmt.Errorf("unable to scan type: %T", v)
	}
}

func Fetch(client *http.Client, url string) (*ScrapeResult, error) {
	result := &ScrapeResult{
		URL:       url,
		Timestamp: time.Now().UTC(),
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching url %q: %w", url, err)
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.Header = resp.Header

	result.Body, err = io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("error fetching url %q: %w", url, err)
	}

	return result, nil
}

func NewScrapeStore(file string) (*ScrapeStore, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%v?_journal=WAL&_timeout=5000", file))
	if err != nil {
		return nil, fmt.Errorf("error opening db file: %w", err)
	}
	return &ScrapeStore{db: db}, nil
}

func NewScrapeStoreInMemory() (*ScrapeStore, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("error opening db file: %w", err)
	}
	return &ScrapeStore{db: db}, nil
}

type ScrapeStore struct {
	db *sql.DB
}

func (s *ScrapeStore) Migrate() error {
	migrations := []string{
		`create table scrape_results(value blob);`,
		`alter table scrape_results add column url text as (json_extract(value, '$.url'))`,
		`create index scrape_results_url_idx on scrape_results(url)`,
		`alter table scrape_results add column status_code text as (json_extract(value, '$.status_code'))`,
		`create index scrape_results_status_code_idx on scrape_results(status_code)`,
	}

	var version int
	if err := s.db.QueryRow(`PRAGMA user_version`).Scan(&version); err != nil {
		return fmt.Errorf("error reading user_version: %w", err)
	}

	for i, mig := range migrations {
		if i < version {
			continue
		}

		tx, err := s.db.Begin()
		if err != nil {
			return fmt.Errorf("error creating tx for migration: %w", err)
		}
		defer tx.Rollback()

		if _, err := tx.Exec(mig); err != nil {
			return fmt.Errorf("error running migration #%d: %w", i+1, err)
		}

		if _, err := tx.Exec(fmt.Sprintf("PRAGMA user_version = %d;", i+1)); err != nil {
			return fmt.Errorf("error running migration #%d: %w", i+1, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("error committing tx for migration: %w", err)
		}
	}

	return nil
}

func (s *ScrapeStore) InsertScrapeResult(result *ScrapeResult) error {
	if _, err := s.db.Exec(`insert into scrape_results values (?)`, result); err != nil {
		return fmt.Errorf("error inserting scrape result: %w", err)
	}
	return nil
}

func (s *ScrapeStore) Close() {
	s.db.Close()
}
