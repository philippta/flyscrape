// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cache_test

import (
	"os"
	"testing"

	"github.com/philippta/flyscrape/modules/cache"
	"github.com/stretchr/testify/require"
)

func TestBoltStore(t *testing.T) {
	dir, err := os.MkdirTemp("", "boltstore")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	store := cache.NewBoltStore(dir + "/test.db")

	v, ok := store.Get("foo")
	require.Nil(t, v)
	require.False(t, ok)

	store.Set("foo", []byte("bar"))

	v, ok = store.Get("foo")
	require.NotNil(t, v)
	require.True(t, ok)
	require.Equal(t, []byte("bar"), v)
}
