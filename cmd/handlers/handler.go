package handlers

import (
	git "kit/internals/git"
	"kit/pkg"
	"net/http"
)

func InitKit(w http.ResponseWriter, r *http.Request) {
	if err := git.InitKit(); err != nil {
		pkg.WriteError(w, err, http.StatusBadRequest)
		return
	}
	pkg.WriteJSON(w, http.StatusOK, map[string]string{"message": "kit initialized successfully"})
}
