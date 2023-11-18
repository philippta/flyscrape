// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var arrayFields = []string{
	"urls",
	"follow",
	"allowedDomains",
	"blockedDomains",
	"allowedURLs",
	"blockedURLs",
	"proxies",
}

func parseConfigArgs(args []string) (map[string]any, error) {
	updates := map[string]any{}

	flag := ""
	for _, arg := range normalizeArgs(args) {
		if flag == "" && !isFlag(arg) {
			return nil, fmt.Errorf("expected flag, got %q instead", arg)
		}

		if flag != "" && isFlag(arg) {
			updates[flag[2:]] = true
			flag = ""
			continue
		}

		if flag != "" {
			if v, ok := updates[flag[2:]]; ok {
				if vv, ok := v.([]any); ok {
					updates[flag[2:]] = append(vv, parseArg(arg))
				} else {
					updates[flag[2:]] = []any{v, parseArg(arg)}
				}
			} else {
				if slices.Contains(arrayFields, flag[2:]) {
					updates[flag[2:]] = []any{parseArg(arg)}
				} else {
					updates[flag[2:]] = parseArg(arg)
				}
			}
			flag = ""
			continue
		}

		flag = arg
	}

	if flag != "" {
		updates[flag[2:]] = true
		flag = ""
	}

	return updates, nil
}

func normalizeArgs(args []string) []string {
	var norm []string

	for _, arg := range args {
		if !strings.HasPrefix(arg, "--") {
			norm = append(norm, arg)
		} else {
			norm = append(norm, strings.SplitN(arg, "=", 2)...)
		}
	}

	return norm
}

func parseArg(arg string) any {
	if arg == "true" {
		return true
	}
	if arg == "false" {
		return false
	}
	if num, err := strconv.Atoi(arg); err == nil {
		return num
	}
	return arg
}

func isFlag(arg string) bool {
	return strings.HasPrefix(arg, "--")
}
