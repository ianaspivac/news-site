package server

import (
	"github.com/gin-gonic/gin"
	"github.com/ianaspivac/news-site-go/internal/httperr"
	"github.com/ianaspivac/news-site-go/internal/models"
	"github.com/ianaspivac/news-site-go/internal/util"
	"net/http"
	"time"
)

func (n *news) CreateUser(c *gin.Context) {
	var u models.User

	if err := c.ShouldBindJSON(&u); err != nil {
		httperr.Handle(c, httperr.New(err.Error(), http.StatusBadRequest))
		return
	}

	if err := util.ValidatePassword(u.Password); err != nil {
		httperr.Handle(c, err)
		return
	} else if err = util.ValidateMail(u.Mail); err != nil {
		httperr.Handle(c, err)
		return
	}

	createdUser, err := n.service.CreateUser(c, u)
	if err != nil {
		httperr.Handle(c, err)
		return
	}

	c.JSON(http.StatusOK, models.User{
		UUID: createdUser.UUID,
		Mail: createdUser.Mail,
	})
}

//func (n *news) GetUserByMail(c *gin.Context) {
//	var u models.User
//
//	if err := c.ShouldBindJSON(&u); err != nil {
//		httperr.Handle(c, httperr.New(err.Error(), http.StatusBadRequest))
//		return
//	}
//
//	if err := u.ValidateMail(); err != nil {
//		httperr.Handle(c, err)
//		return
//	}
//
//	responseUser, err := n.service.GetUserByMail(c, u.Mail)
//	if err != nil {
//		httperr.Handle(c, err)
//		return
//	}
//
//	c.JSON(http.StatusOK, responseUser)
//}

func (n *news) LoginUser(c *gin.Context) {
	var u models.User

	if err := c.ShouldBindJSON(&u); err != nil {
		httperr.Handle(c, httperr.New(err.Error(), http.StatusBadRequest))
		return
	}

	if err := util.ValidateMail(u.Mail); err != nil {
		httperr.Handle(c, err)
		return
	}

	responseUser, err := n.service.AuthenticateUser(c, u)
	if err != nil {
		httperr.Handle(c, err)
		return
	}

	tokenString, err := n.jwt.CreateSignedToken(responseUser.UUID, time.Now().Add(tokenLife).Unix(), tokenIssuer, []byte(tokenKey))
	if err != nil {
		httperr.Handle(c, httperr.WrapHttp(err, "failed to create JWT token", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"role":  responseUser.Type,
	})
}

func (n *news) PromoteUser(c *gin.Context) {
	uuid := c.Param("user_uuid")
	if uuid == "" {
		httperr.Handle(c, httperr.New("user_uuid parameter is missing", http.StatusBadRequest))
		return
	}

	err := n.service.PromoteUser(c, uuid)
	if err != nil {
		httperr.Handle(c, err)
		return
	}

	c.Status(http.StatusOK)
}
