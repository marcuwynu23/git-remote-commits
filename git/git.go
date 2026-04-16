package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Commit struct {
	Hash    string
	Author  string
	Message string
	Time    string
	When    time.Time
}

type Snapshot struct {
	Branch          string
	Status          string
	Commits         []Commit
	RemoteStatus    string
	CommitsBehind   int
	LastRefresh     time.Time
	RepoError       string
	SelectedDiff    string
	SelectedHash    string
	RemoteTrackName string
}

func IsGitRepo(repoPath string) error {
	_, err := run(repoPath, "rev-parse", "--is-inside-work-tree")
	if err != nil {
		return errors.New("not a git repository")
	}
	return nil
}

func CollectSnapshot(repoPath string, limit int) Snapshot {
	s := Snapshot{
		LastRefresh:     time.Now(),
		RemoteTrackName: "origin/main",
	}

	if err := IsGitRepo(repoPath); err != nil {
		s.RepoError = err.Error()
		return s
	}

	branch, err := run(repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		s.RepoError = "failed to read branch"
		return s
	}
	s.Branch = strings.TrimSpace(branch)

	porcelain, err := run(repoPath, "status", "--porcelain")
	if err != nil {
		s.RepoError = "failed to read git status"
		return s
	}
	if strings.TrimSpace(porcelain) == "" {
		s.Status = "clean"
	} else {
		s.Status = "dirty"
	}

	commits, err := listCommits(repoPath, limit)
	if err != nil {
		s.RepoError = err.Error()
		return s
	}
	s.Commits = commits

	_ = fetchRemote(repoPath)
	behind, err := commitsBehind(repoPath)
	if err != nil {
		s.RemoteStatus = "remote unavailable"
	} else {
		s.CommitsBehind = behind
		switch {
		case behind == 0:
			s.RemoteStatus = "up to date"
		case behind == 1:
			s.RemoteStatus = "1 commit behind remote"
		default:
			s.RemoteStatus = fmt.Sprintf("%d commits behind remote", behind)
		}
	}

	return s
}

func ShowCommit(repoPath, hash string) string {
	if strings.TrimSpace(hash) == "" {
		return "No commit selected."
	}
	out, err := run(repoPath, "show", "--stat", "--patch", "--color=never", hash)
	if err != nil {
		return "Unable to load diff preview."
	}
	if strings.TrimSpace(out) == "" {
		return "No diff output for selected commit."
	}
	return out
}

func listCommits(repoPath string, limit int) ([]Commit, error) {
	out, err := run(
		repoPath,
		"log",
		fmt.Sprintf("-%d", limit),
		"--pretty=format:%h|%an|%ar|%at|%s",
	)
	if err != nil {
		if strings.Contains(err.Error(), "does not have any commits") {
			return []Commit{}, nil
		}
		return nil, errors.New("failed to read git log")
	}

	if strings.TrimSpace(out) == "" {
		return []Commit{}, nil
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	commits := make([]Commit, 0, len(lines))
	for _, line := range lines {
		parts := strings.SplitN(line, "|", 5)
		if len(parts) != 5 {
			continue
		}
		epoch, _ := strconv.ParseInt(parts[3], 10, 64)
		commits = append(commits, Commit{
			Hash:    parts[0],
			Author:  parts[1],
			Time:    parts[2],
			When:    time.Unix(epoch, 0),
			Message: parts[4],
		})
	}
	return commits, nil
}

func fetchRemote(repoPath string) error {
	_, err := run(repoPath, "fetch", "origin", "--quiet")
	return err
}

func commitsBehind(repoPath string) (int, error) {
	out, err := run(repoPath, "rev-list", "--count", "HEAD..origin/main")
	if err != nil {
		return 0, err
	}
	out = strings.TrimSpace(out)
	if out == "" {
		return 0, nil
	}
	n, err := strconv.Atoi(out)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func run(repoPath string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("%s", msg)
	}
	return stdout.String(), nil
}
