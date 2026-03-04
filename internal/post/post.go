package post

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

type Post struct {
	Author    string
	Content   string
	CreatedAt string
}

func GetPosts(ctx context.Context) (posts []Post, err error) {
	err = godotenv.Load()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	database, err := sql.Open("pgx", os.Getenv("credentials"))
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	defer database.Close()

	if err = database.Ping(); err != nil {
		fmt.Printf("%v", err)
		return
	}

	tx, err := database.Begin()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, "SELECT author, content, created_at FROM users ORDER BY created_at DESC;")
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		var rowTime string
		err = rows.Scan(&post.Author, &post.Content, &rowTime)
		if err != nil {
            return nil, err
        }

		formatTime, err := time.Parse(time.RFC3339, rowTime)
		if err != nil {
			fmt.Print(err)
            return nil, err
        }
		post.CreatedAt = formatTime.Format("15:04 02.01.2006")
		posts = append(posts, post)
	}

	err = rows.Err()
	if err != nil {
        return nil, err
    }

	tx.Commit()
	return
}

func SavePost(ctx context.Context, author, content string) error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	database, err := sql.Open("pgx", os.Getenv("credentials"))
	if err != nil {
		return err
	}
	defer database.Close()

	if err = database.Ping(); err != nil {
		return err
	}

	tx, err := database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "INSERT INTO users (author, content) VALUES ($1, $2);", author, content)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return tx.Commit()
}
