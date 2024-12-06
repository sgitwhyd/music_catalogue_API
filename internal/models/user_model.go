package models

import "gorm.io/gorm"

type (
	User struct {
		gorm.Model
		Email string 		`db:"email" gorm:"unique;not null"` 
		Username string `db:"username" gorm:"unique;not null"`
		Password string	`db:"password" gorm:"not null"`
	}

	SignUpRequest struct {
		Username	string		`json:"username" binding:"required"`
		Email		 	string 		`json:"email" binding:"required"` 
		Password 	string			`json:"password" binding:"required"`
	}
	SignInRequest struct {
		Email		 	string 		`json:"email" binding:"required"` 
		Password 	string			`json:"password" binding:"required"`
	}

	LoginResponse struct {
		AccessToken string `json:"access_token"`
	}
)