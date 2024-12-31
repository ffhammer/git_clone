package log

import (
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/diff"
	"git_clone/gvc/objectio"
	"git_clone/gvc/pointers"
	"git_clone/gvc/utils"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gobwas/glob"
)

var yellowCommitColor = color.RGB(248, 227, 123)

func datesMatch(commit objectio.CommitMetdata, fromTime, toTime time.Time) (bool, error) {

	date, err := utils.ParseTimeString(commit.Date)
	if err != nil {
		return false, err
	}

	match := fromTime.Before(date) && toTime.After(date)
	return match, nil
}

func authorMatch(commit objectio.CommitMetdata, author string) bool {
	if author == "" {
		return true
	}
	return commit.Author == author
}

func commitMessageGrepMatch(commit objectio.CommitMetdata, pattern string) (bool, error) {
	if pattern == "" {
		return true, nil
	}

	g, err := glob.Compile(pattern)
	if err != nil {
		return false, err
	}
	return g.Match(commit.CommitMessage), nil
}

func StandardLog(patch bool, since, until, author, grep string) (string, error) {

	fromTime := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	if since != "" {
		var err error
		fromTime, err = utils.ParseTimeString(since)
		if err != nil {
			return "", fmt.Errorf("invalid 'since' date format: %w", err)
		}
	}

	// Default to the current time if `until` is not provided
	toTime, _ := utils.ParseTimeString(utils.GetCurrentTimeString())
	if until != "" {
		var err error
		toTime, err = utils.ParseTimeString(until)
		if err != nil {
			return "", fmt.Errorf("invalid 'until' date format: %w", err)
		}
	}

	current, err := pointers.GetLastCommit()
	if err != nil {
		return "", err
	}
	currentHash := config.HEAD

	matchingCommits := make([]objectio.CommitMetdata, 0)
	hashes := make([]string, 0)
	for currentHash != config.DOES_NOT_EXIST_HASH {

		dMatch, err := datesMatch(current, fromTime, toTime)
		if err != nil {
			return "", fmt.Errorf("error checking date for commit %s: %w", currentHash, err)
		}

		mMatch, err := commitMessageGrepMatch(current, grep)
		if err != nil {
			return "", fmt.Errorf("error checking commit message grep %s: %w", currentHash, err)
		}

		if dMatch && authorMatch(current, author) && mMatch {
			matchingCommits = append(matchingCommits, current)
			hashes = append(hashes, currentHash)
		}

		currentHash = current.ParentCommitHash
		current, err = objectio.LoadCommit(current.ParentCommitHash)
		if err != nil {
			return "", err
		}
	}

	var builder strings.Builder
	for i := len(matchingCommits) - 1; i >= 0; i-- {

		current := matchingCommits[i]

		branchBefore := config.STARTING_BRANCH
		if current.ParentCommitHash != config.DOES_NOT_EXIST_HASH {
			commitBefpre, err := objectio.LoadCommit(current.ParentCommitHash)
			if err != nil {
				return "", err
			}
			branchBefore = commitBefpre.BranchName
		}

		branchchange := ""
		if current.BranchName != branchBefore {
			branchchange = fmt.Sprintf("(%s -> %s)", branchBefore, current.BranchName)
		}

		builder.WriteString(fmt.Sprintf("%s %s\nAuthor: %s\nDate:   %s\n\n    %s\n\n\n",
			yellowCommitColor.Sprintf("commit %s", current.CommitMessage),
			branchchange,
			current.Author,
			current.Date,
			current.CommitMessage,
		))
		if patch {
			diff_output, err := diff.CommitToCommit(current.ParentCommitHash, hashes[i], []string{})
			if err != nil {
				return "", fmt.Errorf("")
			}
			builder.WriteString(diff_output)
			builder.WriteString("\n")
		}
	}

	return builder.String(), err
}
