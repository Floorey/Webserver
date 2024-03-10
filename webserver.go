package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type User struct {
	Username string `json:"username"`
	Address  string `json:"address"`
	Email    string `json:"email"`
}

func main() {
	// Gin-Framework
	router := gin.Default()

	//def a route returning the JSON-Token
	router.GET("/token", func(c *gin.Context) {
		user := User{
			Username: "SherlockHolmes",
			Address:  "211bBakerStreet, London, UK",
			Email:    "sherlock.holmes@weee.com",
		}
		//ret. Token as JSON
		tokenString, err := generateToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	})
	//start webserver at port 8080
	router.Run(":8080")
}
func generateToken(user User) (string, error) {
	token := jwt.New(jwt.SigningMethodES256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["address"] = user.Address
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	//signe token with secret key
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
