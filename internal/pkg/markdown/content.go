package markdown

import (
	"fmt"
	"strings"
)

type contentSection struct {
	Title    string
	Lines    []string
	Numbered bool
}

func (s *contentSection) Markdown() string {
	lines := []string{
		fmt.Sprintf("### %s", s.Title),
	}

	for _, l := range s.Lines {
		if s.Numbered && len(s.Lines) > 1 {
			lines = append(lines, fmt.Sprintf("1. %s", l))
		} else {
			lines = append(lines, fmt.Sprintf("- %s", l))
		}
	}

	return strings.Join(lines, "\n")
}
