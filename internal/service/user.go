package service

import (
	"context"
	"github.com/cespare/xxhash"
	"github.com/google/uuid"
	"github.com/ianaspivac/news-site-go/internal/httperr"
	"github.com/ianaspivac/news-site-go/internal/models"
	"github.com/ianaspivac/news-site-go/internal/store"
	"net/http"
	"strconv"
)

func (p *service) CreateUser(ctx context.Context, u models.User) (*models.User, error) {
	if user, err := p.db.GetUserByMail(ctx, u.Mail); user != nil && err == nil {
		return nil, httperr.New("user with such mail already exists", http.StatusBadRequest)
	}

	u.Password = strconv.FormatUint(xxhash.Sum64String(u.Password), 10)
	u.UUID = uuid.New().String()

	err := p.db.CreateUser(ctx, store.User{
		Mail:     u.Mail,
		Password: u.Password,
		UUID:     u.UUID,
		Type:     store.Basic,
	})

	if err != nil {
		return nil, handleDBError(err, "could not create user")
	}

	return &u, nil
}

func (p *service) GetUserByMail(ctx context.Context, mail string) (*models.User, error) {
	user, err := p.db.GetUserByMail(ctx, mail)
	if err != nil {
		return nil, handleDBError(err, "could not get user by mail")
	}

	return &models.User{
		UUID:     user.UUID,
		Mail:     user.UUID,
		Password: user.Password,
		Type:     string(user.Type),
	}, nil
}

func (p *service) GetUserByUUID(ctx context.Context, uuid string) (*models.User, error) {
	user, err := p.db.GetUserByUUID(ctx, uuid)
	if err != nil {
		return nil, handleDBError(err, "could not get user by uuid")
	}

	return &models.User{
		UUID:     user.UUID,
		Mail:     user.UUID,
		Password: user.Password,
		Type:     string(user.Type),
	}, nil
}

func (p *service) AuthenticateUser(ctx context.Context, u models.User) (*models.User, error) {
	user, err := p.db.GetUserByMail(ctx, u.Mail)
	if err != nil {
		return nil, handleDBError(err, "could not get user by mail")
	}

	if user.Password == strconv.FormatUint(xxhash.Sum64String(u.Password), 10) { // Compare password hashes
		return &models.User{
			UUID:     user.UUID,
			Mail:     user.UUID,
			Password: user.Password,
			Type:     string(user.Type),
		}, nil
	} else {
		return nil, httperr.New("Invalid credentials", http.StatusUnauthorized)
	}
}

func (p *service) PromoteUser(ctx context.Context, uuid string) error {
	user, err := p.db.GetUserByUUID(ctx, uuid)
	if err != nil {
		return handleDBError(err, "could not get user by uuid")
	}

	if user.Type != store.Basic {
		return httperr.New("user is already an editor", http.StatusOK)
	}

	err = p.db.PromoteUserToEditor(ctx, uuid)
	if err != nil {
		return handleDBError(err, "could not promote user")
	}

	return nil
}
