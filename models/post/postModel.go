package postModel

import (
	"database/sql"
	"time"
)

type Post struct {
	ID        int
	Title     string
	Content   string
	CreatedBy int
	CreatedAt time.Time
}

func GetAll(db *sql.DB) ([]Post, error) {
	res, err := db.Query("SELECT * from post")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var posts []Post
	for res.Next() {
		var post Post
		if err := res.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedBy, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (p *Post) Get(db *sql.DB) (Post, error) {
	err := db.QueryRow("SELECT * FROM post WHERE id = ?", p.ID).Scan(&p.ID, &p.Title, &p.Content, &p.CreatedBy, &p.CreatedAt)
	if err != nil {
		return Post{}, err
	}
	return *p, nil
}

func (p *Post) Create(db *sql.DB) (int64, error) {
	res, err := db.Exec("INSERT INTO post (title, content, created_by) VALUES (?, ?, ?)", p.Title, p.Content, p.CreatedBy)
	if err != nil {
		return 0, err
	}
	resID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return resID, nil
}

func (p *Post) Update(db *sql.DB, data *Post) (int64, error) {
	query := "UPDATE post SET "
	args := []interface{}{}
	if data.Title != "" {
		query += "title = ?, "
		args = append(args, data.Title)
	}
	if data.Content != "" {
		query += "content = ?, "
		args = append(args, data.Content)
	}
	if data.CreatedBy != 0 {
		query += "created_by = ?, "
		args = append(args, data.CreatedBy)
	}
	if !data.CreatedAt.IsZero() {
		query += "created_at = ?, "
		args = append(args, data.CreatedAt)
	}
	query = query[:len(query)-2] // Remove the trailing comma and space
	query += " WHERE id = ?"
	args = append(args, p.ID)

	res, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	resID, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return resID, nil
}

func (p *Post) Delete(db *sql.DB) (int64, error) {
	res, err := db.Exec("DELETE FROM post WHERE id = ?", p.ID)
	if err != nil {
		return 0, err
	}
	resID, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return resID, nil
}
