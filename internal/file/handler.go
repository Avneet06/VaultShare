package filehandler

import (
	"file-sharing-system/db"
	"file-sharing-system/middleware"
	"file-sharing-system/cache"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"strconv"
	"strings"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File missing", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// get email from JWT context
	email := r.Context().Value(middleware.UserEmailKey).(string)
	

	
	user, err := db.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	filename := timestamp + "_" + handler.Filename
	filePath := filepath.Join("uploads", filename)

	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	
	done := make(chan error)
	go func() {
		_, err := io.Copy(out, file)
		done <- err
	}()

	if err := <-done; err != nil {
		http.Error(w, "File write failed", http.StatusInternalServerError)
		return
	}

	
	meta := &db.File{
		UserID:     user.ID,
		Filename:   handler.Filename,
		Size:       handler.Size,
		FileType:   handler.Header.Get("Content-Type"),
		StoragePath: filePath,
	}

	err = db.SaveFileMetadata(meta)
	if err != nil {
		http.Error(w, "Failed to save metadata", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("http://localhost:8080/uploads/%s", filename)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("File uploaded successfully. URL: %s", url)))
}

func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(middleware.UserEmailKey).(string)
	user, err := db.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	
	name := r.URL.Query().Get("name")
	filetype := r.URL.Query().Get("type")
	date := r.URL.Query().Get("date")

	cacheKey := fmt.Sprintf("files:user:%d:name:%s:type:%s:date:%s", user.ID, name, filetype, date)


	val, err := cache.RDB.Get(cache.Ctx, cacheKey).Result()
	if err == nil {
	
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(val))
		return
	}

	
	files, err := db.SearchFilesByUserID(user.ID, name, filetype, date)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	jsonResp, _ := json.Marshal(files)

	
	cache.RDB.Set(cache.Ctx, cacheKey, jsonResp, 300_000_000_000)

	// 4. Return
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}



func ShareFileHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	fileID := parts[2]
	cacheKey := "file:" + fileID


	val, err := cache.RDB.Get(cache.Ctx, cacheKey).Result()
	if err == nil {
		
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(val))
		return
	}

	
	id, err := strconv.Atoi(fileID)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	file, err := db.GetFileByID(id)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	url := fmt.Sprintf("http://localhost:8080/%s", file.StoragePath)
	resp := map[string]string{"url": url}
	jsonResp, _ := json.Marshal(resp)

	
	cache.RDB.Set(cache.Ctx, cacheKey, jsonResp, 300_000_000_000) 

	
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}


