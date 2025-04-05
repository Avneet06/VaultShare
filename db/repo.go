package db

import (
	"errors"
	"fmt"
)

func CreateUser(user *User) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2)`
	_, err := DB.Exec(query, user.Email, user.Password)
	return err
}

func GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, email, password FROM users WHERE email=$1`
	row := DB.QueryRow(query, email)

	var u User
	err := row.Scan(&u.ID, &u.Email, &u.Password)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &u, nil
}

func SaveFileMetadata(file *File) error {
	query := `INSERT INTO files (user_id, filename, size, filetype, storage_path)
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := DB.Exec(query, file.UserID, file.Filename, file.Size, file.FileType, file.StoragePath)
	return err
}

func GetFilesByUserID(userID int) ([]File, error) {
	query := `SELECT id, filename, size, filetype, storage_path, created_at
	          FROM files WHERE user_id=$1 ORDER BY created_at DESC`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var f File
		err := rows.Scan(&f.ID, &f.Filename, &f.Size, &f.FileType, &f.StoragePath, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func GetFileByID(fileID int) (*File, error) {
	query := `SELECT id, filename, size, filetype, storage_path, created_at, user_id
	          FROM files WHERE id = $1`

	row := DB.QueryRow(query, fileID)

	var f File
	err := row.Scan(&f.ID, &f.Filename, &f.Size, &f.FileType, &f.StoragePath, &f.CreatedAt, &f.UserID)
	if err != nil {
		return nil, err
	}
	return &f, nil
}
func SearchFilesByUserID(userID int, name, filetype, date string) ([]File, error) {
	query := `SELECT id, filename, size, filetype, storage_path, created_at
	          FROM files WHERE user_id = $1`
	args := []interface{}{userID}
	i := 2

	if name != "" {
		query += fmt.Sprintf(" AND filename ILIKE $%d", i)
		args = append(args, "%"+name+"%")
		i++
	}
	if filetype != "" {
		query += fmt.Sprintf(" AND filetype ILIKE $%d", i)
		args = append(args, "%"+filetype+"%")
		i++
	}
	if date != "" {
		query += fmt.Sprintf(" AND DATE(created_at) = $%d", i)
		args = append(args, date)
	}

	query += " ORDER BY created_at DESC"

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var f File
		err := rows.Scan(&f.ID, &f.Filename, &f.Size, &f.FileType, &f.StoragePath, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}

func GetExpiredFiles() ([]File, error) {
	query := `SELECT id, user_id, filename, size, filetype, storage_path, created_at, expired_at 
	          FROM files WHERE expired_at IS NOT NULL AND expired_at < NOW()`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var f File
		err := rows.Scan(&f.ID, &f.UserID, &f.Filename, &f.Size, &f.FileType, &f.StoragePath, &f.CreatedAt, &f.ExpiredAt)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func DeleteFileByID(fileID int) error {
	_, err := DB.Exec("DELETE FROM files WHERE id=$1", fileID)
	return err
}
