package store

import (
	"context"
	"github.com/ianaspivac/news-site-go/internal/models"
)

type Post struct {
	ID          uint64
	CreatorUUID string `db:"creator_uuid"`
	models.Post
}

func (d *db) CreatePost(ctx context.Context, post Post) error {
	_, err := d.sql.ExecContext(ctx, `INSERT INTO post (uuid, creator_uuid, headline, summary, preview_img, content, is_protected) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		post.UUID,
		post.CreatorUUID,
		post.Headline,
		post.Summary,
		post.PreviewImage,
		post.Content,
		post.IsProtected,
	)
	return err
}

func (d *db) UpdatePost(ctx context.Context, post Post) error {
	_, err := d.sql.ExecContext(ctx, "UPDATE post SET headline = ?, summary = ?, content = ? WHERE uuid = ?", post.Headline, post.Summary, post.Content, post.UUID)
	if err != nil {
		return err
	}

	return nil
}

func (d *db) GetAllPosts(ctx context.Context) ([]Post, error) {
	var posts []Post

	err := d.sql.SelectContext(ctx, &posts, `SELECT * FROM post`)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (d *db) GetPostByUUID(ctx context.Context, uuid string) (*Post, error) {
	var p Post

	err := d.sql.GetContext(ctx, &p, `SELECT * FROM post WHERE uuid=?`, uuid)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
