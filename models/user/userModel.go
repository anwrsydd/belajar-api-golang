package userModel

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int
	Username  string
	Name      string
	Password  string
	Email     string
	CreatedAt time.Time
}

func GetAll(db *sql.DB) ([]User, error) {
	res, err := db.Query("SELECT * from user")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var users []User
	for res.Next() {
		var user User
		if err := res.Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Password, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
func (u *User) Get(db *sql.DB) (User, error) {
	err := db.QueryRow("SELECT * FROM user WHERE id = ?", u.ID).Scan(&u.ID, &u.Username, &u.Name, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		return User{}, err
	}
	return *u, nil
}

func (u *User) Create(db *sql.DB) (int64, error) {
	res, err := db.Exec("INSERT INTO user (username, name, email, password) VALUES (?, ?, ?, ?)", u.Username, u.Name, u.Email, u.Password)
	if err != nil {
		return 0, err
	}
	resID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return resID, nil
}

func (u *User) Update(db *sql.DB, data *User) (int64, error) {
	query := "UPDATE user SET "
	args := []interface{}{}
	if data.Username != "" {
		query += "username = ?, "
		args = append(args, data.Username)
	}
	if data.Name != "" {
		query += "name = ?, "
		args = append(args, data.Name)
	}
	if data.Email != "" {
		query += "email = ?, "
		args = append(args, data.Email)
	}
	if data.Password != "" {
		query += "password = ?, "
		args = append(args, data.Password)
	}
	if !data.CreatedAt.IsZero() {
		query += "created_at = ?, "
		args = append(args, data.CreatedAt)
	}
	query = query[:len(query)-2] // Remove the trailing comma and space
	query += " WHERE id = ?"
	args = append(args, u.ID)

	res, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowAffected, nil
}

func (u *User) Delete(db *sql.DB) (int64, error) {
	res, err := db.Exec("DELETE FROM user WHERE id = ?", u.ID)
	if err != nil {
		return 0, err
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowAffected, nil
}
