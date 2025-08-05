package poster

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/pkg/gitmate"
)

type fileName string

const (
	posterBig fileName = "poster.png"
	posterScc fileName = "scc.png"
)

const branchMain = "master"

func toPull(poster model.Poster, event model.Event) gitmate.Pull {
	fileType := posterBig
	if poster.SCC {
		fileType = posterScc
	}

	return gitmate.Pull{
		Title: fmt.Sprintf("%s - %s", event.Name, fileType),
		Body:  "This is an automated action by `events`",
	}
}

func toPath(poster model.Poster, event model.Event) string {
	fileType := posterBig
	if poster.SCC {
		fileType = posterScc
	}

	return fmt.Sprintf("%d-%d/%s/%s", event.Year.Start, event.Year.End, event.Name, fileType)
}

func toBranch(poster model.Poster, event model.Event) string {
	fileType := "poster"
	if poster.SCC {
		fileType = "poster_scc"
	}

	return fmt.Sprintf("feat/%s_%s", sanitizeBranchName(event.Name), fileType)
}

// sanitizeBranchName sanitizes a branch name according to
// https://git-scm.com/docs/git-check-ref-format
func sanitizeBranchName(name string) string {
	// Lowercase
	s := strings.ToLower(name)

	// Replace spaces and illegal characters with underscores
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "`", "")
	s = strings.ReplaceAll(s, "\\", "")
	s = strings.ReplaceAll(s, "~", "")
	s = strings.ReplaceAll(s, "^", "")
	s = strings.ReplaceAll(s, ":", "")
	s = strings.ReplaceAll(s, "?", "")
	s = strings.ReplaceAll(s, "*", "")
	s = strings.ReplaceAll(s, "[", "")
	s = strings.ReplaceAll(s, "]", "")
	s = strings.ReplaceAll(s, "@{", "")
	s = strings.ReplaceAll(s, "..", "-")

	// Replace multiple slashes with a single hyphen
	s = regexp.MustCompile(`/+`).ReplaceAllString(s, "-")

	// Remove trailing `.lock` and leading/trailing `-`, `/`, or `.`
	s = strings.TrimSuffix(s, ".lock")
	s = strings.Trim(s, "-/.")

	// Collapse multiple underscores
	s = regexp.MustCompile(`_+`).ReplaceAllString(s, "_")

	if s == "" {
		s = "default_branch"
	}

	return s
}
