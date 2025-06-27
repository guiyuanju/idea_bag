package tui

import (
	"fmt"
	"slices"
)

const (
	StyleReset             = "\x1b[0m"
	StyleBold              = "\x1b[1m"
	StyleBackgroundBlue    = "\x1b[44m"
	StyleBackgroundBlack   = "\x1b[40m"
	StyleBackgroundDefault = "\x1b[49m"
	StyleEsc               = "\x1b"
)

func HideCursor() {
	fmt.Printf("%s[?25l", StyleEsc)
}

func ShowCursor() {
	fmt.Printf("%s[?25h", StyleEsc)
}

type Text string

func style(s string) Text {
	return Text(s + StyleReset)
}

func (s Text) String() string {
	return string(s)
}

func (s Text) Bold() Text {
	return style(StyleBold + s.String())
}

func (s Text) BackgroundBlue() Text {
	return Text(StyleBackgroundBlue + s.String() + "\x1b[49m")
}

func (s Text) BackgroundBlack() Text {
	return Text(StyleBackgroundBlack + s.String() + "\x1b[49m")
}

func matchedRanges(s string, match []string) [][]int {
	if len(s) == 0 || len(match) == 0 {
		return [][]int{}
	}

	var matchRanges [][]int
	for _, m := range match {
		var i, j int
		for i < len(s) {
			if s[i] == m[j] {
				j++
			}
			i++
			if j == len(m) {
				matchRanges = append(matchRanges, []int{i - len(m), i})
				j = 0
			}
		}
	}

	// merge overlapping ranges
	var mergedRanges [][]int
	slices.SortFunc(matchRanges, func(a, b []int) int { return a[0] - b[0] })
	var i int
	for ; i < len(matchRanges)-1; i++ {
		r1 := matchRanges[i]
		r2 := matchRanges[i+1]
		if r1[1] >= r2[0] {
			r2[0] = r1[0]
			r2[1] = max(r2[1], r1[1])
		} else {
			mergedRanges = append(mergedRanges, r1)
		}
	}
	mergedRanges = append(mergedRanges, matchRanges[i])

	return mergedRanges
}

func styleMatched(s string, match []string, styleFunc func(s string) string) string {
	ranges := matchedRanges(s, match)

	var res []byte
	var prevEndIdx int
	for _, r := range ranges {
		res = append(res, s[prevEndIdx:r[0]]...)
		bolded := styleFunc(s[r[0]:r[1]])
		res = append(res, []byte(bolded)...)
		prevEndIdx = r[1]
	}
	res = append(res, s[prevEndIdx:]...)
	return string(res)
}

func BoldMatched(s string, match []string) string {
	return styleMatched(s, match, func(s string) string { return Text(s).Bold().String() })
}

func BgBlueMatched(s string, match []string) string {
	return styleMatched(s, match, func(s string) string { return Text(s).BackgroundBlue().String() })
}
