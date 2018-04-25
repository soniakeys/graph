// Copyright 2018 Sonia Keys
// License MIT: https://opensource.org/licenses/MIT

package io

import "testing"

func TestParseNIs(t *testing.T) {
	defer func() {
		if x := recover(); x == nil {
			t.Fatal("expected panic")
		}
	}()
	parseNIs([]string{" 3"}, 10)
}
