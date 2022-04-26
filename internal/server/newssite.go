package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ianaspivac/news-site-go/internal/httperr"
	"github.com/ianaspivac/news-site-go/internal/models"
	"github.com/ianaspivac/news-site-go/internal/service"
	"github.com/ianaspivac/news-site-go/internal/store"
	"github.com/ianaspivac/news-site-go/internal/util"
	"net/http"
	"time"
)

const (
	tokenIssuer = "news_dev"
	tokenLife   = time.Hour * 4
	tokenKey    = "bebra"
	AuthSchema  = "Bearer " // Space if required by auth header standard
	UUIDKey     = "uuid"
)

type News interface {
	Run(hostname string) error
}

type news struct {
	mux *gin.Engine

	jwt     util.JWT
	service service.Service
}

func New(
	jwtProvider util.JWT,
	serviceProvider service.Service,
) News {
	muxProvider := gin.New()
	muxProvider.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders: []string{"*"},
		MaxAge:       12 * time.Hour,
	}))
	serv := &news{
		mux:     muxProvider,
		jwt:     jwtProvider,
		service: serviceProvider,
	}

	v1 := muxProvider.Group("/v1")
	{
		posts := v1.Group("/posts")
		{
			posts.GET("/", serv.GetAllPosts)                      // Get All posts (return just headline, summary and preview pic)
			posts.GET("/:post_uuid", serv.GetPostByUUID)          // Get specific post (if post if protected - token required)
			posts.POST("/", serv.editorRequired, serv.CreatePost) // Create a new post
			posts.PATCH("/:post_uuid", serv.editorRequired, serv.EditPost)
		}
		user := v1.Group("/user")
		{
			user.POST("/register", serv.CreateUser)                                // Register a new user (non-anonymous a.k.a subscriber)
			user.POST("/login", serv.LoginUser)                                    // Authenticate user and return a JWT token
			user.POST("/promote/:user_uuid", serv.adminRequired, serv.PromoteUser) // Promote basic user to an editor
		}
	}

	return serv
}

func (n *news) Run(hostname string) error {
	return n.mux.Run(hostname)
}

func (n *news) basicRequired(c *gin.Context) {
	if _, err := n.required(c); err != nil {
		httperr.Handle(c, err)
		c.Abort()
		return
	}
}

func (n *news) editorRequired(c *gin.Context) {
	user, err := n.required(c)
	if err != nil {
		httperr.Handle(c, err)
		c.Abort()
		return
	}

	if user.Type != string(store.Editor) && user.Type != string(store.Admin) {
		httperr.Handle(c, httperr.New("user is not an editor", http.StatusUnauthorized))
		c.Abort()
		return
	}
}

func (n *news) adminRequired(c *gin.Context) {
	user, err := n.required(c)
	if err != nil {
		httperr.Handle(c, err)
		c.Abort()
		return
	}

	if user.Type != string(store.Admin) {
		httperr.Handle(c, httperr.New("user is not an admin", http.StatusUnauthorized))
		c.Abort()
		return
	}
}

func (n *news) required(c *gin.Context) (*models.User, error) {
	authHeader := c.GetHeader("Authorization")

	if len(authHeader) <= len(AuthSchema) {
		return nil, httperr.New("missing or invalid JWT token", http.StatusUnauthorized)
	}

	token := authHeader[len(AuthSchema):]

	uuid, err := n.jwt.ValidateToken(token, tokenIssuer, []byte(tokenKey))
	if err != nil {
		return nil, httperr.WrapHttp(err, "failed to validate JWT token", http.StatusUnauthorized)

	}

	user, err := n.service.GetUserByUUID(c, uuid)
	if err != nil {
		return nil, err
	}

	c.Set(UUIDKey, uuid)

	return user, nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
