// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cache

import (
	"github.com/cornelk/hashmap"
)

func NewMemStore() *hashmap.Map[string, []byte] {
	return hashmap.New[string, []byte]()
}
