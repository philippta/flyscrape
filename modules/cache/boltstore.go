// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cache

import (
	"errors"
	"log"
	"os"

	"go.etcd.io/bbolt"
)

var cache = []byte("cache")

func NewBoltStore(file string) *BoltStore {
	db, err := bbolt.Open(file, 0644, nil)
	if err != nil {
		log.Printf("cache: failed to create database file %q: %v\n", file, err)
		os.Exit(1)
	}

	c := &BoltStore{db: db}

	return c
}

type BoltStore struct {
	db *bbolt.DB
}

func (s *BoltStore) Get(key string) ([]byte, bool) {
	var value []byte

	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(cache)
		if bucket == nil {
			return errors.New("bucket not found")
		}

		v := bucket.Get([]byte(key))
		if v == nil {
			return errors.New("key not found")
		}

		value = make([]byte, len(v))
		copy(value, v)

		return nil
	})
	if err != nil {
		return nil, false
	}
	return value, true
}

func (s *BoltStore) Set(key string, value []byte) {
	err := s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(cache)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(key), value)
	})
	if err != nil {
		log.Printf("cache: failed to insert cache key %q: %v\n", key, err)
	}
}

func (s *BoltStore) Close() {
	s.db.Close()
}
