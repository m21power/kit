package handlers

import (
	"encoding/json"
	"fmt"
	git "kit/internals/git"
	"kit/internals/utils"
	"kit/pkg"
	"net/http"
	"os"
	"path/filepath"
)

func InitHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Username == "" {
		pkg.WriteError(w, fmt.Errorf("username is required"), http.StatusBadRequest)
		return
	}
	if err := git.InitKit(req.Username); err != nil {
		pkg.WriteError(w, err, http.StatusBadRequest)
		return
	}
	pkg.WriteJSON(w, http.StatusOK, map[string]string{"message": "kit initialized successfully"})
}
func AddHandler(w http.ResponseWriter, r *http.Request) {
	var req pkg.AddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteError(w, err, http.StatusBadRequest)
		return
	}
	workspace := filepath.Join("workspaces", req.Username)
	if err := os.MkdirAll(workspace, 0755); err != nil {
		pkg.WriteError(w, fmt.Errorf("failed to create workspace: %v", err), http.StatusInternalServerError)
		return
	}

	if err := utils.WriteStructure(workspace, req.RootFolder); err != nil {
		pkg.WriteError(w, fmt.Errorf("failed to write structure: %v", err), http.StatusInternalServerError)
		return
	}

	// "kit add ."
	if len(req.Paths) == 1 && req.Paths[0] == "." {
		staged, err := git.AddKit(workspace, req.Username)
		if err != nil {
			pkg.WriteError(w, fmt.Errorf("add failed: %v", err), http.StatusInternalServerError)
			return
		}

		if len(staged) == 0 {
			pkg.WriteJSON(w, http.StatusAccepted, map[string]any{
				"message": "Nothing to stage",
			})
		} else {
			pkg.WriteJSON(w, http.StatusAccepted, map[string]any{
				"message": "Files staged successfully",
				"staged":  staged,
			})
		}
	} else {
		var allStaged []string

		for _, file := range req.Paths {
			filePath := filepath.Join(workspace, file)
			staged, err := git.AddKit(filePath, req.Username)
			if err != nil {
				pkg.WriteError(w, fmt.Errorf("add failed: %v", err), http.StatusInternalServerError)
				return
			}
			allStaged = append(allStaged, staged...)
		}

		if len(allStaged) == 0 {
			pkg.WriteJSON(w, http.StatusAccepted, map[string]any{
				"message": "Nothing to stage",
			})
		} else {
			pkg.WriteJSON(w, http.StatusAccepted, map[string]any{
				"message": "Files staged successfully",
				"staged":  allStaged,
			})
		}
	}
}

func CommitHandler(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Message  string `json:"message"`
		Username string `json:"username"`
	}
	var payload req
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		pkg.WriteError(w, err, http.StatusBadRequest)
		return
	}
	response, err := git.CommitGit(payload.Message, payload.Username)
	if err != nil {
		pkg.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	pkg.WriteJSON(w, http.StatusOK, map[string]string{"message": response})
}

func LogHandler(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Count    int    `json:"count"`
		Username string `json:"username"`
	}
	var payload req

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		pkg.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if payload.Count <= 0 {
		payload.Count = 5
	}

	response, err := git.LogKit(payload.Count, payload.Username)
	if err != nil {
		pkg.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	pkg.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Commit logs retrieved successfully",
		"logs":    response,
	})
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	var req pkg.StatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteError(w, err, http.StatusBadRequest)
		return
	}
	workspace := filepath.Join("workspaces", req.Username)
	err := utils.WriteStructure(workspace, req.RootFolder)
	if err != nil {
		pkg.WriteError(w, fmt.Errorf("failed to write structure: %v", err), http.StatusInternalServerError)
		return
	}
	if req.Username == "" {
		pkg.WriteError(w, fmt.Errorf("username is required"), http.StatusBadRequest)
		return
	}

	branch, err := utils.GetHead(req.Username)
	if err != nil {
		pkg.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	treeComm, err := utils.GetCommitTreeHash(branch, req.Username)
	if err != nil {
		pkg.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	result, err := git.StatusKit(treeComm, "", req.Username, make(map[string]bool))
	if err != nil {
		pkg.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	// fmt.Println("from tree")
	// for _, entry := range result {
	// 	fmt.Println(entry.Path, entry.Hash)
	// }

	res, err := git.IsChanged(result, ".", req.Username)
	// fmt.Println("from index")
	// for path, status := range res {
	// 	fmt.Println(path, status.Hash, status.Message)
	// }
	if err != nil {
		pkg.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	var stagedFiles []map[string]string
	var unstagedFiles []map[string]string

	for path, status := range res {
		if status.Staged {
			stagedFiles = append(stagedFiles, map[string]string{"path": path, "status": status.Message})
		} else {
			unstagedFiles = append(unstagedFiles, map[string]string{"path": path, "status": status.Message})
		}
	}

	if len(stagedFiles) == 0 && len(unstagedFiles) == 0 {
		pkg.WriteJSON(w, http.StatusOK, map[string]any{
			"message": "Nothing to commit, working tree clean",
		})
		return
	}

	pkg.WriteJSON(w, http.StatusOK, map[string]any{
		"branch":         branch,
		"staged":         stagedFiles,
		"unstaged":       unstagedFiles,
		"staged_count":   len(stagedFiles),
		"unstaged_count": len(unstagedFiles),
	})
}

func RestoreHandler(w http.ResponseWriter, r *http.Request) {
	var req pkg.RestoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteError(w, err, http.StatusBadRequest)
		return
	}
	if req.Username == "" {
		pkg.WriteError(w, fmt.Errorf("username is required"), http.StatusBadRequest)
		return
	}
	branch, err := utils.GetHead(req.Username)
	if err != nil {
		pkg.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	treeComm, err := utils.GetCommitTreeHash(branch, req.Username)
	if err != nil {
		pkg.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	oldTree, err := git.StatusKit(treeComm, "", req.Username, make(map[string]bool))
	if err != nil {
		pkg.WriteError(w, fmt.Errorf("failed to get index entry: %w", err), http.StatusInternalServerError)
		return
	}

	result, err := git.KitRestore(req.Username, req.Paths, oldTree)
	if err != nil {
		pkg.WriteError(w, fmt.Errorf("restore failed: %w", err), http.StatusInternalServerError)
		return
	}

	pkg.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Files restored successfully",
		"result":  result,
	})
}

func CreateBranch(w http.ResponseWriter, r *http.Request) {
	var payload pkg.BranchRequest
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		pkg.WriteError(w, err, http.StatusBadRequest)
		return
	}
	err = git.CreateBranch(payload.Branch, payload.Username)
	if err != nil {
		pkg.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	pkg.WriteJSON(w, http.StatusOK, map[string]string{"message": "branch created successfully!"})
}

func CheckoutBranch(w http.ResponseWriter, r *http.Request) {
	var payload pkg.BranchRequest
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		pkg.WriteError(w, err, http.StatusBadRequest)
		return
	}
	result, err := git.CheckoutBranch(payload.Branch, payload.Username)
	if err != nil {
		pkg.WriteError(w, err, http.StatusInternalServerError)
		return
	}
	pkg.WriteJSON(w, http.StatusOK, map[string]pkg.FileSystemItem{"data": *result})
}
