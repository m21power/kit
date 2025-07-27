package git

import (
	util "kit/internals/utils"
)

func CommitGit(message string) (string, error) {
	// go to HEAD and see the branch
	branch, err := util.GetHead()
	if err != nil {
		return "", err
	}
	// then go to refs/heads/{branch} and get the hash
	oldTreeHash, err := util.GetCommitTreeHash(branch)
	if err != nil {
		return "", err
	}
	indexEntry, err := util.GetIndexEntry()
	if err != nil {
		return "", err
	}
	tree, err := util.BuildTree(indexEntry)
	if err != nil {
		return "", err
	}
	newCommitHash, err := util.WriteTree(tree)
	if err != nil {
		return "", err
	}

	if newCommitHash == oldTreeHash {
		return "No changes to commit", nil
	}
	oldCommitHash, err := util.GetLastCommitHash(branch)
	if err != nil {
		return "", err
	}
	err = util.WriteCommit(branch, newCommitHash, oldCommitHash, message)
	if err != nil {
		return "", err
	}
	return "Commit successful", nil
}
