package main

import (
	"bloom/read"
	"bloom/repo"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func getUserTags(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	u64, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	userIDUint := uint(u64)

	tags := repo.GetUserTagByID(userIDUint)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

func main() {
	credits := read.ReadData()

	repo.SetupDB(credits)
	http.HandleFunc("/user-tags", getUserTags)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

}
