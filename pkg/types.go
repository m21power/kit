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

type Status struct {
	Staged  bool
	Message string
	Hash    string
}

type FileNode struct {
	Path string
	Hash string
	Mode string
	Type string
}

type FileSystemItem struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	Type     string           `json:"type"`
	Content  string           `json:"content,omitempty"`
	Children []FileSystemItem `json:"children,omitempty"`
	Path     string           `json:"path"`
}

type AddRequest struct {
	Username   string         `json:"username"`
	Paths      []string       `json:"files"`
	RootFolder FileSystemItem `json:"rootFolder"`
}

type StatusRequest struct {
	Username   string         `json:"username"`
	RootFolder FileSystemItem `json:"rootFolder"`
}

type RestoreRequest struct {
	Username string   `json:"username"`
	Paths    []string `json:"paths"`
}

type BranchRequest struct {
	Username string `json:"username"`
	Branch   string `json:"branch"`
}

type RestoreResponse struct {
	Restored   map[string]bool `json:"restored"`
	FileSystem FileSystemItem  `json:"fileSystem"`
}

type ResetRequest struct {
	Username   string `json:"username"`
	CommitHash string `json:"hash"`
}
