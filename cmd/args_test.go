// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseConfigUpdates(t *testing.T) {
	tests := []struct {
		flags   string
		err     bool
		updates map[string]any
	}{
		{
			flags:   `--foo bar`,
			updates: map[string]any{"foo": "bar"},
		},
		{
			flags:   `--foo=bar`,
			updates: map[string]any{"foo": "bar"},
		},
		{
			flags:   `--foo`,
			updates: map[string]any{"foo": true},
		},
		{
			flags:   `--foo false`,
			updates: map[string]any{"foo": false},
		},
		{
			flags:   `--foo a --foo b`,
			updates: map[string]any{"foo": []any{"a", "b"}},
		},
		{
			flags:   `--foo a --foo=b`,
			updates: map[string]any{"foo": []any{"a", "b"}},
		},
		{
			flags:   `--foo 69`,
			updates: map[string]any{"foo": 69},
		},
		{
			flags:   `--foo.bar a`,
			updates: map[string]any{"foo.bar": "a"},
		},
		{
			flags: `foo`,
			err:   true,
		},
		{
			flags: `--foo a b`,
			err:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.flags, func(t *testing.T) {
			args, err := parseConfigArgs(strings.Fields(test.flags))

			if test.err {
				require.Error(t, err)
				require.Empty(t, args)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.updates, args)
		})
	}
}
