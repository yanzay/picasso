package templates

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const padWidth = 20

type screen struct {
	title  string
	table  [][]string
	footer string
}

func newScreen(title string) *screen {
	return &screen{
		table: make([][]string, 0),
		title: title,
	}
}

func (s *screen) renderN(n int) string {
	parts := []string{}
	if s.title != "" {
		parts = append(parts, s.title+"\n")
	}
	parts = append(parts, renderTableN(s.table, n))
	if s.footer != "" {
		parts = append(parts, "\n"+s.footer)
	}
	return strings.Join(parts, "\n")
}

func (s *screen) render() string {
	return s.renderN(padWidth)
}

func (s *screen) add(name string, value interface{}, icon string) {
	s.table = append(s.table, []string{name, fmt.Sprintf("%s%s", fmt.Sprint(value), icon)})
}

func (s *screen) addLine() {
	s.table = append(s.table, []string{"", ""})
}

func pad(first, last string) string {
	return padN(first, last, padWidth)
}

func padN(first, last string, n int) string {
	firstLen := utf8.RuneCountInString(first)
	lastLen := utf8.RuneCountInString(last)
	if firstLen+lastLen > n-1 {
		r := []rune(first)
		r = r[:n-lastLen-1]
		first = string(r)
		firstLen = utf8.RuneCountInString(first)
	}
	repeatCount := n - firstLen - lastLen
	return first + strings.Repeat(" ", repeatCount) + last
}

func renderTableN(table [][]string, n int) string {
	rows := make([]string, len(table))
	for i, row := range table {
		rows[i] = padN(row[0], row[1], n)
	}
	return strings.Join(rows, "\n")
}

// render Nx2 table
func renderTable(table [][]string) string {
	return renderTableN(table, padWidth)
}
