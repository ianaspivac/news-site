package server

import (
	"github.com/gin-gonic/gin"
	"github.com/ianaspivac/news-site-go/internal/httperr"
	"github.com/ianaspivac/news-site-go/internal/models"
	"net/http"
)

func (n *news) CreatePost(c *gin.Context) {
	var p models.Post

	if err := c.ShouldBindJSON(&p); err != nil {
		httperr.Handle(c, httperr.New(err.Error(), http.StatusBadRequest))
		return
	}

	if p.Headline == "" ||
		p.Summary == "" ||
		p.Content == "" {
		httperr.Handle(c, httperr.New("some fields are missing", http.StatusBadRequest))
		return
	}

	editorUUID := c.GetString(UUIDKey)
	if editorUUID == "" {
		httperr.Handle(c, httperr.New("uuid missing in context", http.StatusUnauthorized))
		return
	}

	createdPost, err := n.service.CreatePost(c, editorUUID, p)
	if err != nil {
		httperr.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, models.Post{
		UUID:         createdPost.UUID,
		CreatorUUID:  createdPost.CreatorUUID,
		Headline:     createdPost.Headline,
		Summary:      createdPost.Summary,
		PreviewImage: createdPost.PreviewImage,
		Content:      createdPost.Content,
		IsProtected:  createdPost.IsProtected,
	})
}

func (n *news) GetPostByUUID(c *gin.Context) {
	uuid := c.Param("post_uuid")
	if uuid == "" {
		httperr.Handle(c, httperr.New("post_uuid parameter is missing", http.StatusBadRequest))
		return
	}

	post, err := n.service.GetPostByUUID(c, uuid)
	if err != nil {
		httperr.Handle(c, err)
		return
	}

	if !post.IsProtected {
		c.JSON(http.StatusOK, post)
		return
	} else {
		_, err := n.required(c)
		if err != nil {
			httperr.Handle(c, err)
			return
		}

		c.JSON(http.StatusOK, post)
		return
	}
}

func (n *news) GetAllPosts(c *gin.Context) {
	posts, err := n.service.GetAllPosts(c)
	if err != nil {
		httperr.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (n *news) EditPost(c *gin.Context) {
	uuid := c.Param("post_uuid")
	if uuid == "" {
		httperr.Handle(c, httperr.New("post_uuid parameter is missing", http.StatusBadRequest))
		return
	}

	var p models.Post

	if err := c.ShouldBindJSON(&p); err != nil {
		httperr.Handle(c, httperr.New(err.Error(), http.StatusBadRequest))
		return
	}

	p.UUID = uuid

	newPost, err := n.service.UpdatePost(c, p)
	if err != nil {
		httperr.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, newPost)
}
