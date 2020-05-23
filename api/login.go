package api

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/valuechaintool/valuechaintool/models"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	type login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var l login
	if err := c.ShouldBindJSON(&l); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid JSON provided")
		return
	}
	user, err := models.GetUserByName(l.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Invalid login credentials")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(l.Password)) != nil {
		c.JSON(http.StatusUnauthorized, "Invalid login credentials")
		return
	}
	token, err := CreateToken(*user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
}

func CreateToken(user models.User) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = user.ID.String()
	atClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(viper.GetString("accessSecret")))
	if err != nil {
		return "", err
	}
	return token, nil
}
