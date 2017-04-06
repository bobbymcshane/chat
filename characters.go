// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build !windows

package main

const TOP_RIGHT = '┐'
const VERTICAL_LINE = '│'
const HORIZONTAL_LINE = '─'
const TOP_LEFT = '┌'
const BOTTOM_RIGHT = '┘'
const BOTTOM_LEFT = '└'
const VERTICAL_LEFT = '┤'
const VERTICAL_RIGHT = '├'
const HORIZONTAL_DOWN = '┬'
const HORIZONTAL_UP = '┴'
const QUOTA_LEFT = '«'
const QUOTA_RIGHT = '»'

var connectors = map[rune]bool{
	VERTICAL_LINE:   true,
	HORIZONTAL_LINE: true,
}

// [toInsert][above][below][left][right]
var characterConnectorMap = map[string]rune{
	// vertical insert
	"│││─ ": VERTICAL_LEFT,
	"│││ ─": VERTICAL_RIGHT,

	// horizontal insert
	"─│ ──": HORIZONTAL_UP,
	"─ │──": HORIZONTAL_DOWN,
}
