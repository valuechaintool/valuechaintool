package web

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/valuechaintool/valuechaintool/models"
)

func Register(c *gin.Context) {
	d := map[string]string{
		"PageTitle": "ValueChain Registration",
	}
	c.HTML(http.StatusOK, "register.html", d)
}

func RegisterPost(c *gin.Context) {
	if c.PostForm("password") != c.PostForm("password2") {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, errors.New("passwords do not match"))
		return
	}

	user := models.User{
		Username: c.PostForm("username"),
		RealName: c.PostForm("realname"),
		Email:    c.PostForm("email"),
		Password: c.PostForm("password"),
	}

	if err := models.NewUser(&user); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusFound, "/")
}
