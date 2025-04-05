package worker

import (
	"file-sharing-system/db"
	"fmt"
	"os"
	"time"
)

func StartCleanupWorker() {
	go func() {
		for {
			fmt.Println("Checking for expired files...") // background log

			files, err := db.GetExpiredFiles()
			if err != nil {
				fmt.Println(" Failed to fetch expired files:", err)
				time.Sleep(1 * time.Minute)
				continue
			}

			for _, f := range files {
				err := os.Remove(f.StoragePath)
				if err != nil {
					fmt.Println(" Failed to delete file from disk:", f.StoragePath)
					continue
				}
				db.DeleteFileByID(f.ID)
				fmt.Println("Deleted:", f.StoragePath)
			}

			time.Sleep(1 * time.Minute)
		}
	}()
}
