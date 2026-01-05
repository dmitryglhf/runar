package git

import (
	"os/exec"
	"strings"
)

type Info struct {
	Commit string
	Branch string
	Dirty  bool
}

func GetInfo() *Info {
	commit := runGit("rev-parse", "--short", "HEAD")
	// Return if no git repo found
	if commit == "" {
		return nil
	}

	branch := runGit("rev-parse", "--abbrev-ref", "HEAD")
	dirty := runGit("status", "--porcelain") != ""

	return &Info{
		Commit: commit,
		Branch: branch,
		Dirty:  dirty,
	}
}

func runGit(args ...string) string {
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
