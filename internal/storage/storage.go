package storage

import (
	"context"
	"database/sql"
	"time"

	// импорт для "database/sql"
	_ "github.com/lib/pq"
)

type Storage struct {
	*sql.DB
}

type Post struct {
	Author    string
	Content   string
	CreatedAt string
}

func (s *Storage) GetPosts() (posts []Post, err error) {
	tx, err := s.Begin()
	if err != nil {
		return
	}

	defer tx.Rollback()
	// defer func() {
	// 	if p := recover(); p != nil {
	// 		tx.Rollback()
	// 	} else if err != nil {
	// 		tx.Rollback()
	// 	} else {
	// 		err = tx.Commit()
	// 	}
	// }()

	rows, err := tx.Query("SELECT author, content, created_at FROM users ORDER BY created_at DESC;")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		var rowTime string
		err = rows.Scan(&post.Author, &post.Content, &rowTime)
		if err != nil {
			return
		}

		formatTime, err := time.Parse(time.RFC3339, rowTime)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = formatTime.Format("15:04 02.01.2006")
		posts = append(posts, post)
	}

	err = rows.Err()
	if err != nil {
		return
	}
	tx.Commit()
	return
}

func (s *Storage) SavePost(ctx context.Context, author, content string) error {
	tx, err := s.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	_, err = tx.ExecContext(ctx, "INSERT INTO users (author, content) VALUES ($1, $2);", author, content)

	return err
}
