package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/ianaspivac/news-site-go/internal/httperr"
	"github.com/ianaspivac/news-site-go/internal/models"
	"github.com/ianaspivac/news-site-go/internal/store"
	"net/http"
)

type Service interface {
	CreateUser(ctx context.Context, u models.User) (*models.User, error)
	GetUserByUUID(ctx context.Context, uuid string) (*models.User, error)
	GetUserByMail(ctx context.Context, mail string) (*models.User, error)
	AuthenticateUser(ctx context.Context, u models.User) (*models.User, error)

	CreatePost(ctx context.Context, creator string, u models.Post) (*models.Post, error)
	UpdatePost(ctx context.Context, newPost models.Post) (*models.Post, error)
	GetAllPosts(ctx context.Context) ([]models.Post, error)
	GetPostByUUID(ctx context.Context, uuid string) (*models.Post, error)
	PromoteUser(ctx context.Context, uuid string) error
}

type service struct {
	db store.DB
}

func New(datastore store.DB) Service {
	return &service{
		db: datastore,
	}
}

func handleDBError(err error, msg string) error {
	if err == sql.ErrNoRows {
		return httperr.New(fmt.Sprintf("%s: not found", msg), http.StatusNotFound)
	}

	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		switch mysqlErr.Number {
		case 1062:
			return httperr.New(fmt.Sprintf("%s: entry already exists", msg), http.StatusBadRequest)
		case 1741:
			return httperr.New(fmt.Sprintf("%s: key not found", msg), http.StatusNotFound)
		}
	}
	// TODO: Change in live environment
	return httperr.New(fmt.Sprintf("%s: unknown internal error: %v", msg, err), http.StatusInternalServerError)
}
