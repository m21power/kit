package pkg

import "time"

type IndexEntry struct {
	Mode string
	Path string
	Hash string
}
type TreeNode struct {
	Blobs map[string]BlobEntry // files: name -> BlobEntry
	Trees map[string]*TreeNode // subdirs: name -> TreeNode
}

type BlobEntry struct {
	Mode string
	Hash string
}

type TreeEntry struct {
	Mode string
	Hash string
	Name string
	Type string
}

type CommitObject struct {
	Hash    string
	Author  string
	Message string
	Email   string
	Date    time.Time
	Parent  string
}
