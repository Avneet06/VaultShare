package db

type File struct {
	ID          int
	UserID      int
	Filename    string
	Size        int64
	FileType    string
	StoragePath string
	CreatedAt   string
	ExpiredAt   *string 
}
