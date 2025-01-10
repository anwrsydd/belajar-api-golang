package postHandler

import (
	postModel "belajar-api-golang/models/post"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAll(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		posts, err := postModel.GetAll(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "There was an error while fetching posts"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Success",
			"data":    posts,
		})
	}
}

func Get(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("id")
		id, err := strconv.Atoi(postID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be a number"})
			return
		}

		p := postModel.Post{
			ID: id,
		}
		post, err := p.Get(db)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "There was an error while fetching post"})
				fmt.Println("Error while fetching post:", err.Error())
				return
			}

		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Post found",
			"data":    post,
		})
	}
}

func Create(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		type PostForm struct {
			Title     string `form:"title" json:"title" binding:"required"`
			Content   string `form:"content" json:"content" binding:"required"`
			CreatedBy int    `form:"created_by" json:"created_by" binding:"required"`
		}
		var p PostForm
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		post := postModel.Post{
			Title:     p.Title,
			Content:   p.Content,
			CreatedBy: p.CreatedBy,
		}
		res, err := post.Create(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "There was an error while creating post"})
			return
		}
		post = postModel.Post{
			ID: int(res),
		}
		post, err = post.Get(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "There was an error while fetching post"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Post created",
			"data":    post,
		})
	}
}

func Update(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		type PostFormPut struct {
			Title   string `form:"title" json:"title"`
			Content string `form:"content" json:"content"`
		}
		var p PostFormPut
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		postIDKey := c.Param("id")
		postID, err := strconv.Atoi(postIDKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be a number"})
			return
		}
		post := postModel.Post{
			ID:      postID,
			Title:   p.Title,
			Content: p.Content,
		}
		rowAffected, err := post.Update(db, &post)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "There was an error while updating post"})
			return
		}
		if rowAffected == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No post updated"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Post updated"})
	}
}

func Delete(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postIDKey := c.Param("id")
		postID, err := strconv.Atoi(postIDKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be a number"})
			return
		}
		post := postModel.Post{
			ID: postID,
		}
		rowAffected, err := post.Delete(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "There was an error while deleting post"})
			return
		}
		if rowAffected == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No post deleted"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
	}
}
