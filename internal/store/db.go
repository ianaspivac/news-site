package store

import (
	"context"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"log"
)

type DB interface {
	CreateUser(ctx context.Context, user User) error
	GetUserByMail(ctx context.Context, mail string) (*User, error)
	GetUserByUUID(ctx context.Context, UUID string) (*User, error)

	CreatePost(ctx context.Context, post Post) error
	UpdatePost(ctx context.Context, post Post) error
	GetAllPosts(ctx context.Context) ([]Post, error)
	GetPostByUUID(ctx context.Context, uuid string) (*Post, error)
	PromoteUserToEditor(ctx context.Context, UUID string) error
}

type db struct {
	sql *sqlx.DB
}

func initDB(d *db) {
	var (
		dbHost = viper.GetString("DATABASE_HOST")
		dbName = viper.GetString("DATABASE_NAME")
		dbPort = viper.GetString("DATABASE_PORT")
		dbUser = viper.GetString("DATABASE_USER")
		dbPass = viper.GetString("DATABASE_PASSWORD")
		err    error
	)

	if dbHost == "" ||
		dbName == "" ||
		dbPort == "" ||
		dbUser == "" ||
		dbPass == "" {
		log.Panicln("Could not open DB connection. Some of the database env are missing")
	}

	cfg := mysql.Config{
		User:      dbUser,
		Passwd:    dbPass,
		Addr:      fmt.Sprintf("%s:%s", dbHost, dbPort),
		DBName:    dbName,
		ParseTime: true,
	}

	d.sql, err = sqlx.Connect("mysql", cfg.FormatDSN())
	if err != nil {
		log.Panicf("Could not open SQL connection: %v", err)
	}
}

func CreateDB() DB {
	var d db
	initDB(&d)
	return &d
}
