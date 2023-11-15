// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cache

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteStore(file string) *SQLiteStore {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_timeout=5000&_journal=WAL", file))
	if err != nil {
		log.Printf("cache: failed to create database file %q: %v\n", file, err)
		os.Exit(1)
	}

	c := &SQLiteStore{db: db}
	c.migrate()

	return c
}

type SQLiteStore struct {
	db *sql.DB
}

func (s *SQLiteStore) Get(key string) ([]byte, bool) {
	var value []byte
	if err := s.db.QueryRow(`SELECT value FROM cache WHERE key = ? LIMIT 1`, key).Scan(&value); err != nil {
		return nil, false
	}
	return value, true
}

func (s *SQLiteStore) Set(key string, value []byte) {
	if _, err := s.db.Exec(`INSERT INTO cache (key, value) VALUES (?, ?)`, key, value); err != nil {
		log.Printf("cache: failed to insert cache key %q: %v\n", key, err)
	}
}

func (s *SQLiteStore) Close() {
	s.db.Close()
}

func (s *SQLiteStore) migrate() {
	if _, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS cache (key TEXT, value BLOB)`); err != nil {
		log.Printf("cache: failed to create cache table: %v\n", err)
		os.Exit(1)
	}
	if _, err := s.db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS cache_key_idx ON cache(key)`); err != nil {
		log.Printf("cache: failed to create cache index: %v\n", err)
		os.Exit(1)
	}
}
