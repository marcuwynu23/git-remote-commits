package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Commit struct {
	Hash        string
	Author      string
	AuthorEmail string
	Refs        string
	Message     string
	Body        string
	Time        string
	When        time.Time
}

type Snapshot struct {
	RepoName        string
	Branch          string
	Status          string
	Commits         []Commit
	CurrentAuthor   string
	RemoteStatus    string
	CommitsBehind   int
	CommitsAhead    int
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

func CollectSnapshot(repoPath string, remoteName string, limit int) Snapshot {
	s := Snapshot{
		LastRefresh: time.Now(),
	}
	if strings.TrimSpace(remoteName) == "" {
		remoteName = "origin"
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

	remoteOnline := isRemoteReachable(repoPath, remoteName)
	var pullErr error
	if remoteOnline {
		pullErr = pullRemoteBranch(repoPath, remoteName, s.Branch)
	}

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
	s.CurrentAuthor = currentAuthor(repoPath)
	s.RepoName = remoteRepoName(repoPath, remoteName)

	s.RemoteTrackName = remoteName + "/" + s.Branch
	ahead, behind, err := aheadBehind(repoPath, s.RemoteTrackName)
	if !remoteOnline {
		s.RemoteStatus = "offline"
	} else if pullErr != nil {
		s.RemoteStatus = "pull failed: " + pullErr.Error()
	} else if err != nil {
		s.RemoteStatus = "remote unavailable"
	} else {
		s.CommitsBehind = behind
		s.CommitsAhead = ahead
		switch {
		case behind == 0 && ahead == 0:
			s.RemoteStatus = "up to date"
		case behind > 0 && ahead == 0:
			if behind == 1 {
				s.RemoteStatus = "1 commit behind remote"
			} else {
				s.RemoteStatus = fmt.Sprintf("%d commits behind remote", behind)
			}
		case behind == 0 && ahead > 0:
			if ahead == 1 {
				s.RemoteStatus = "1 commit ahead of remote"
			} else {
				s.RemoteStatus = fmt.Sprintf("%d commits ahead of remote", ahead)
			}
		default:
			s.RemoteStatus = fmt.Sprintf("%d behind / %d ahead", behind, ahead)
		}
	}

	return s
}

func remoteRepoName(repoPath, remoteName string) string {
	if strings.TrimSpace(remoteName) == "" {
		return ""
	}
	out, err := run(repoPath, "remote", "get-url", remoteName)
	if err != nil {
		return ""
	}
	url := strings.TrimSpace(out)
	if url == "" {
		return ""
	}
	url = strings.TrimSuffix(url, ".git")

	// SSH format: git@github.com:owner/repo
	if strings.Contains(url, "@") && strings.Contains(url, ":") && !strings.Contains(url, "://") {
		parts := strings.SplitN(url, ":", 2)
		host := parts[0]
		path := parts[1]
		host = strings.TrimPrefix(host, "git@")
		path = strings.TrimPrefix(path, "/")
		if host != "" && path != "" {
			return host + "/" + path
		}
	}

	// HTTPS/SSH URL format: scheme://host/owner/repo
	if i := strings.Index(url, "://"); i >= 0 {
		rest := url[i+3:]
		slash := strings.Index(rest, "/")
		if slash > 0 && slash < len(rest)-1 {
			host := rest[:slash]
			path := strings.TrimPrefix(rest[slash+1:], "/")
			if host != "" && path != "" {
				return host + "/" + path
			}
		}
	}

	return url
}

func currentAuthor(repoPath string) string {
	out, err := run(repoPath, "config", "user.name")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}

func fetchRemote(repoPath, remoteName string) error {
	if strings.TrimSpace(remoteName) == "" {
		return errors.New("remote name required")
	}
	_, err := run(repoPath, "fetch", remoteName, "--quiet")
	return err
}

func pullRemoteBranch(repoPath, remoteName, branch string) error {
	if strings.TrimSpace(remoteName) == "" {
		return errors.New("remote name required")
	}
	if strings.TrimSpace(branch) == "" {
		return errors.New("branch name required")
	}
	_, err := runWithTimeout(8*time.Second, repoPath, "pull", remoteName, branch)
	return err
}

func isRemoteReachable(repoPath, remoteName string) bool {
	if strings.TrimSpace(remoteName) == "" {
		return false
	}
	_, err := runWithTimeout(3*time.Second, repoPath, "ls-remote", "--exit-code", remoteName, "HEAD")
	return err == nil
}

func aheadBehind(repoPath, upstream string) (ahead int, behind int, err error) {
	out, err := run(repoPath, "rev-list", "--left-right", "--count", upstream+"...HEAD")
	if err != nil {
		return 0, 0, err
	}
	parts := strings.Fields(strings.TrimSpace(out))
	if len(parts) != 2 {
		return 0, 0, errors.New("unexpected rev-list output")
	}
	behind, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	ahead, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}
	return ahead, behind, nil
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
	args := []string{"log", "--decorate=short", "--pretty=format:%h%x1f%an%x1f%ae%x1f%ar%x1f%at%x1f%d%x1f%s%x1f%b%x1e"}
	if limit > 0 {
		args = append(args, fmt.Sprintf("-%d", limit))
	}
	out, err := run(repoPath, args...)
	if err != nil {
		if strings.Contains(err.Error(), "does not have any commits") {
			return []Commit{}, nil
		}
		return nil, errors.New("failed to read git log")
	}

	if strings.TrimSpace(out) == "" {
		return []Commit{}, nil
	}

	records := strings.Split(out, "\x1e")
	commits := make([]Commit, 0, len(records))
	for _, record := range records {
		record = strings.TrimSpace(record)
		if record == "" {
			continue
		}

		parts := strings.SplitN(record, "\x1f", 8)
		if len(parts) != 8 {
			continue
		}
		epoch, _ := strconv.ParseInt(parts[4], 10, 64)
		commits = append(commits, Commit{
			Hash:        parts[0],
			Author:      parts[1],
			AuthorEmail: parts[2],
			Time:        parts[3],
			When:        time.Unix(epoch, 0),
			Refs:        strings.TrimSpace(parts[5]),
			Message:     parts[6],
			Body:        strings.TrimSpace(parts[7]),
		})
	}
	return commits, nil
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

func runWithTimeout(timeout time.Duration, repoPath string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = repoPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return "", errors.New("timed out while contacting remote")
		}
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("%s", msg)
	}
	return stdout.String(), nil
}
