package userHandler

import (
	userModel "belajar-api-golang/models/user"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserForm struct {
	Username string `form:"username" json:"username" binding:"required"`
	Name     string `form:"name" json:"name" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

func GetAll(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := userModel.GetAll(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": "There was an error while fetching users"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Success",
			"data":    users,
		})
	}
}
func Get(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		id, err := strconv.Atoi(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be a number"})
			return
		}
		user := userModel.User{
			ID: id,
		}
		user, err = user.Get(db)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "There was an error while fetching user"})
				fmt.Println("Error while fetching user:", err.Error())
				return
			}

		}
		c.JSON(http.StatusOK, gin.H{
			"message": "User found",
			"data":    user,
		})
	}
}

func Create(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var u UserForm
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user := userModel.User{
			Username: u.Username,
			Name:     u.Name,
			Email:    u.Email,
			Password: u.Password,
		}
		res, err := user.Create(db)
		if err != nil {
			fmt.Println("Error while creating user:", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "There was an error while creating user",
			})
			return
		}

		user = userModel.User{
			ID: int(res),
		}
		user, err = user.Get(db)
		if err != nil {
			fmt.Println("Error while fetching user:", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "There was an error while fetching user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "User created successfully",
			"data":    user,
		})
	}
}

func Update(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		type UserFormPut struct {
			Username string `form:"username" json:"username"`
			Name     string `form:"name" json:"name"`
			Email    string `form:"email" json:"email" binding:"omitempty,email"`
			Password string `form:"password" json:"password"`
		}
		var u UserFormPut
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userIDKey := c.Param("id")
		userID, err := strconv.Atoi(userIDKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be a number"})
			return
		}
		user := userModel.User{
			ID:       userID,
			Username: u.Username,
			Name:     u.Name,
			Email:    u.Email,
			Password: u.Password,
		}
		rowAffected, err := user.Update(db, &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "There was an error while updating user"})
			fmt.Println("Error while updating user:", err.Error())
			return
		}
		if rowAffected == 0 { // if there is no row affected, it's caused by the user not found
			c.JSON(http.StatusOK, gin.H{"message": "No user updated"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
	}
}

func Delete(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDKey := c.Param("id")
		userID, err := strconv.Atoi(userIDKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be a number"})
			return
		}
		user := userModel.User{
			ID: userID,
		}

		rowAffected, err := user.Delete(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "There was an error while deleting user"})
			fmt.Println("Error while deleting user:", err.Error())
			return
		}
		if rowAffected == 0 { // if there is no row affected, it's caused by the user not found
			c.JSON(http.StatusOK, gin.H{"message": "No user deleted"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}
