package router

import (
	// "github.com/google/uuid"
	"math/rand"
)

type SignInDetails struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	First_name   string `json:"first_name"`
	Last_name    string `json:"last_name"`
	Role         string `json:"role"`
	Avatar_url   string `json:"avatar_url"`
	Phone_number string `json:"phone_number"`
	Timezone     string `json:"timezone"`
	Language     string `json:"language"`
}

type UserDetails struct {
	// Id                  string `json:"id"`
	Id                  int    `json:"id"`
	Email               string `json:"email"`
	Role                string `json:"role"`
	First_name          string `json:"first_name"`
	Last_name           string `json:"last_name"`
	Avatar_url          string `json:"avatar_url"`
	Phone_number        string `json:"phone_number"`
	Timezone            string `json:"timezone"`
	Dark_mode           bool   `json:"dark_mode"`
	Email_notifications bool   `json:"email_notifications"`
	Language            string `json:"language"`
	Is_verified         bool   `json:"is_verified"`
}

func ConvertSignInDetailsToUserDetails(signInDetails SignInDetails, userDetails *UserDetails) {
	// userDetails.Id = uuid.NewString()
	userDetails.Id = rand.Intn(10000)
	userDetails.Email = signInDetails.Email
	userDetails.First_name = signInDetails.First_name
	userDetails.Last_name = signInDetails.Last_name
	userDetails.Avatar_url = signInDetails.Avatar_url
	userDetails.Timezone = signInDetails.Timezone
	userDetails.Role = signInDetails.Role
	userDetails.Language = signInDetails.Language
	userDetails.Dark_mode = false
	userDetails.Is_verified = false
	userDetails.Email_notifications = true
	userDetails.Phone_number = signInDetails.Phone_number
}
