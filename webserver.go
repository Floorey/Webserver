package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"

	"os"
	"os/signal"
	"syscall"

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

	// add JSON-USER-DATA
	router.POST("/user", addUserHandler)

	//start webserver at port 8080
	go func() {
		if err := router.Run(":8080"); err != nil {
			fmt.Println("Server error", err)
		}
	}()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		go worker(i)
	}

	for {
		select {
		case <-sig:
			fmt.Println("Exit...")
			return
		default:
			time.Sleep(time.Second)
			fmt.Println("looping..")
		}
	}
}

func generateToken(user User) (string, error) {
	token := jwt.New(jwt.SigningMethodES256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["address"] = user.Address
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	fmt.Println("Generated token:", token)

	//signe token with secret key
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func addUserHandler(c *gin.Context) {
	// Parse the JSON request body
	var newUser User
	err := json.NewDecoder(c.Request.Body).Decode(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the user data (assuming you have a function to handle this)
	// saveUser(newUser)

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{"message": "User successfully added"})
}

func worker(id int) {
	for {
		fmt.Printf("Worker %d started\n", id)
		time.Sleep(time.Second)
		fmt.Printf("Worker %d beendet\n", id)
	}
}
