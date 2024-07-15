package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"testingoauth/database"
	"testingoauth/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

var GoogleConfig *oauth2.Config

var User models.Input_User

var UserInfo struct {
	Email	string	`json:"email"`
}

func Login(c *gin.Context) {
	URL := GoogleConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, URL)
}

func init() {
	err := godotenv.Load()

	if err != nil {
		panic("Failed to load .env file")
	}

    GoogleConfig = &oauth2.Config{
        RedirectURL:  os.Getenv("REDIRECT_URL"),
        ClientID:     os.Getenv("CLIENT_ID"),
        ClientSecret: os.Getenv("CLIENT_SECRET"),
        Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
        Endpoint:     google.Endpoint,
    }
}

func GoogleCallBack(c *gin.Context) {
	Code := c.Query("code")
	Token, err := GoogleConfig.Exchange(context.Background(), Code)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"message": "Cant exchange code to token",
		})
	}

	Client := GoogleConfig.Client(context.Background(), Token)
	Information, err := Client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"message": "Fail to get information user",
		})
	}

	err = json.NewDecoder(Information.Body).Decode(&UserInfo)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "fail",
			"message": "Failed to decode user's information",
		})
	}

	err = database.DB.Where("email = ?", UserInfo.Email).First(&User).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			User = models.Input_User{Email: UserInfo.Email}
			database.DB.Create(&User)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "fail",
				"message": "Failed to find or create user",
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user": User,
	})
}