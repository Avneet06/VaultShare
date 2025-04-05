package main

import (
	"file-sharing-system/db"
	"file-sharing-system/internal/auth"
	"file-sharing-system/internal/file"
	"file-sharing-system/middleware"
	"file-sharing-system/internal/worker"
	"file-sharing-system/cache"

	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting File Server...") // basic startup log (just to check)
	db.ConnectDB() 
	cache.InitRedis()
	worker.StartCleanupWorker() 

	http.HandleFunc("/register", auth.RegisterHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	http.HandleFunc("/upload", middleware.JWTAuthMiddleware(filehandler.UploadHandler))
	http.HandleFunc("/files", middleware.JWTAuthMiddleware(filehandler.ListFilesHandler))
	http.HandleFunc("/share/", filehandler.ShareFileHandler) 



	
http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))





	fmt.Println(" Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
