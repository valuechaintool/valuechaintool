package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/valuechaintool/valuechaintool/models"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := getUserIDFromToken(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Set("userID", *userID)
		capabilities, err := models.ListCapabilitiesByUser(*userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Set("capabilities", capabilities)
		c.Next()
	}
}

func getUserIDFromToken(r *http.Request) (*uuid.UUID, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}
	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return nil, err
	}
	return &userID, nil
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	authorization := r.Header.Get("Authorization")
	if len(authorization) < 8 {
		return nil, errors.New("unrecognized authorization method")
	}
	if authorization[:7] != "Bearer " {
		return nil, errors.New("unrecognized authorization method")
	}
	token, err := jwt.Parse(authorization[7:], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString("accessSecret")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
