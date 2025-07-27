package git

import (
	util "kit/internals/utils"
	"kit/pkg"
)

func LogKit(count int) ([]pkg.CommitObject, error) {
	result := make([]pkg.CommitObject, 0)
	branch, err := util.GetHead()
	if err != nil {
		return nil, err
	}
	// then go to refs/heads/{branch} and get the hash
	oldCommitHash, err := util.GetLastCommitHash(branch)
	if err != nil {
		return nil, err
	}
	for range count {
		commit, err := util.GetLogs(oldCommitHash)
		if err != nil {
			return nil, err
		}
		result = append(result, commit)
		oldCommitHash = commit.Parent
		if oldCommitHash == "" {
			break
		}
	}
	return result, nil
}
