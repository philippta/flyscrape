// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package flyscrape_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"flyscrape"

	"github.com/stretchr/testify/require"
)

func TestWatch(t *testing.T) {
	f := tmpfile(t)
	defer os.Remove(f.Name())
	write(f, "test 1")

	calls := 0
	done := make(chan struct{})

	go func() {
		err := flyscrape.Watch(f.Name(), func(s string) error {
			calls++
			if calls == 1 {
				require.Equal(t, "test 1", s)
				return nil
			}
			if calls == 2 {
				require.Equal(t, "test 2", s)
				return flyscrape.StopWatch
			}
			return nil
		})
		require.NoError(t, err)
		close(done)
	}()

	write(f, "test 2")
	<-done
}

func TestWatchError(t *testing.T) {
	f := tmpfile(t)
	defer os.Remove(f.Name())

	done := make(chan struct{})
	go func() {
		err := flyscrape.Watch(f.Name(), func(s string) error {
			return errors.New("test")
		})
		require.Error(t, err)
		close(done)
	}()

	write(f, "test 2")
	<-done
}

func tmpfile(t *testing.T) *os.File {
	f, err := os.CreateTemp("", "scrape.js")
	require.NoError(t, err)
	return f
}

func write(f *os.File, s string) {
	time.Sleep(10 * time.Millisecond)
	f.Seek(0, 0)
	f.Truncate(0)
	f.WriteString(s)
	f.Sync()
}
