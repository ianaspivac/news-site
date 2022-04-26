package store

import (
	"context"
)

type UserType string

const (
	Basic  UserType = "BASIC"
	Editor UserType = "EDITOR"
	Admin  UserType = "ADMIN"
)

type User struct {
	ID       uint64
	Mail     string
	Password string
	UUID     string
	Type     UserType
}

func (d *db) CreateUser(ctx context.Context, user User) error {
	_, err := d.sql.ExecContext(ctx, `INSERT INTO user (mail, password, uuid, type) VALUES (?, ?, ?, ?)`, user.Mail, user.Password, user.UUID, user.Type)
	return err
}

func (d *db) GetUserByMail(ctx context.Context, mail string) (*User, error) {
	var u User

	err := d.sql.GetContext(ctx, &u, `SELECT * FROM user WHERE mail=?`, mail)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (d *db) GetUserByUUID(ctx context.Context, UUID string) (*User, error) {
	var u User

	err := d.sql.GetContext(ctx, &u, `SELECT * FROM user WHERE uuid=?`, UUID)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (d *db) PromoteUserToEditor(ctx context.Context, UUID string) error {
	_, err := d.sql.ExecContext(ctx, `UPDATE user SET type = ? WHERE uuid = ?`, Editor, UUID)
	if err != nil {
		return err
	}

	return err
}
