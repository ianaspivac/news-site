package service

import (
	"bytes"
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/ianaspivac/news-site-go/internal/httperr"
	"github.com/ianaspivac/news-site-go/internal/models"
	"github.com/ianaspivac/news-site-go/internal/store"
	"net/http"
)

func (p *service) CreatePost(ctx context.Context, creator string, u models.Post) (*models.Post, error) {
	u.UUID = uuid.New().String()

	err := p.db.CreatePost(ctx, store.Post{
		CreatorUUID: creator,
		Post:        u,
	})

	if err != nil {
		return nil, handleDBError(err, "could not create user")
	}

	u.CreatorUUID = creator

	return &u, nil
}

func (p *service) UpdatePost(ctx context.Context, newPost models.Post) (*models.Post, error) {
	post, err := p.db.GetPostByUUID(ctx, newPost.UUID)
	if err != nil {
		return nil, handleDBError(err, "could not get post by uuid")
	}

	if newPost.Summary != "" && newPost.Summary != post.Summary {
		post.Summary = newPost.Summary
	}
	if newPost.Headline != "" && newPost.Headline != post.Headline {
		post.Headline = newPost.Headline
	}
	if newPost.Content != "" && newPost.Content != post.Content {
		post.Content = newPost.Content
	}
	if newPost.PreviewImage != nil && !bytes.Equal(newPost.PreviewImage, post.PreviewImage) {
		post.PreviewImage = newPost.PreviewImage
	}

	if cmp.Equal(post, newPost) {
		return nil, httperr.New("no new data in request", http.StatusBadRequest)
	}

	err = p.db.UpdatePost(ctx, *post)
	if err != nil {
		return nil, handleDBError(err, "could not update post")
	}

	return &models.Post{
		UUID:         post.UUID,
		CreatorUUID:  post.CreatorUUID,
		Headline:     post.Headline,
		Summary:      post.Summary,
		PreviewImage: post.PreviewImage,
		Content:      post.Content,
		IsProtected:  post.IsProtected,
	}, nil
}

func (p *service) GetAllPosts(ctx context.Context) ([]models.Post, error) {
	posts, err := p.db.GetAllPosts(ctx)
	if err != nil {
		return nil, handleDBError(err, "could not get all posts")
	}

	var response []models.Post
	for _, post := range posts {
		response = append(response, models.Post{
			UUID:         post.UUID,
			CreatorUUID:  post.CreatorUUID,
			Headline:     post.Headline,
			Summary:      post.Summary,
			PreviewImage: post.PreviewImage,
			IsProtected:  post.IsProtected,
			CreatedAt:    post.CreatedAt,
			UpdatedAt:    post.UpdatedAt,
		})
	}

	return response, nil
}

func (p *service) GetPostByUUID(ctx context.Context, uuid string) (*models.Post, error) {
	post, err := p.db.GetPostByUUID(ctx, uuid)
	if err != nil {
		return nil, handleDBError(err, "could not get post by uuid")
	}

	return &models.Post{
		UUID:         post.UUID,
		CreatorUUID:  post.CreatorUUID,
		Headline:     post.Headline,
		Summary:      post.Summary,
		PreviewImage: post.PreviewImage,
		Content:      post.Content,
		IsProtected:  post.IsProtected,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
	}, nil
}
