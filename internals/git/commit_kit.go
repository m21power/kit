package git

import (
	util "kit/internals/utils"
)

func CommitGit(message string, username string) (string, error) {
	// go to HEAD and see the branch
	branch, err := util.GetHead(username)
	if err != nil {
		return "", err
	}
	// then go to refs/heads/{branch} and get the hash
	oldTreeHash, err := util.GetCommitTreeHash(branch, username)
	if err != nil {
		return "", err
	}
	indexEntry, err := util.GetIndexEntry(username)
	if err != nil {
		return "", err
	}
	tree, err := util.BuildTree(indexEntry)
	if err != nil {
		return "", err
	}
	newCommitHash, err := util.WriteTree(tree, username)
	if err != nil {
		return "", err
	}

	if newCommitHash == oldTreeHash {
		return "No changes to commit", nil
	}
	oldCommitHash, err := util.GetLastCommitHash(branch, username)
	if err != nil {
		return "", err
	}
	err = util.WriteCommit(branch, newCommitHash, oldCommitHash, message, username)
	if err != nil {
		return "", err
	}
	return "Commit successful", nil
}
