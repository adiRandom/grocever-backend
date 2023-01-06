package services

import (
	"auth/data/models"
	"auth/data/repository"
	"auth/services/crypto/jwt"
	auth2 "lib/data/dto/auth"
	"lib/data/models/auth"
	"lib/helpers"
	"lib/network/http"
)

type LoginDetails struct {
	body *auth2.LoginRequest
	user *auth.User
}

func NewLoginDetails(body *auth2.LoginRequest, user *auth.User) LoginDetails {
	return LoginDetails{
		body: body,
		user: user,
	}
}

func HandleLogin(details LoginDetails, userRepository *repository.User) http.Response[any] {
	var user *auth.User
	if details.body != nil {
		user_, err := userRepository.GetByUsernameAndPwd(details.body.Username, details.body.Password)
		if err != nil {
			return http.Response[any]{
				StatusCode: 400,
				Err:        err.Error(),
				Body:       helpers.None{},
			}
		}
		user = user_
	} else if details.user != nil {
		user = details.user
	} else {
		return http.Response[any]{
			StatusCode: 400,
			Err:        "Invalid login details",
			Body:       helpers.None{},
		}
	}

	token, err := jwt.GenerateJwtToken(*user)
	if err != nil {
		println(err.Error())

		return http.Response[any]{
			StatusCode: 500,
			Err:        "Internal server error",
			Body:       helpers.None{},
		}
	}

	refreshToken, err := jwt.GenerateRefreshToken(*user)
	if err != nil {
		println(err.Error())

		return http.Response[any]{
			StatusCode: 500,
			Err:        "Internal server error",
			Body:       helpers.None{},
		}
	}

	return http.Response[any]{
		StatusCode: 200,
		Body: auth2.AuthResponse{
			AccessToken:  token,
			RefreshToken: refreshToken,
			User:         auth2.NewUserFromModel(*user),
		},
	}
}

func HandleRegister(body auth2.RegisterRequest, userRepository *repository.User) http.Response[any] {
	user, err := userRepository.CreateFromAuth(models.NewRegisterFromDto(body))
	if err != nil {
		return http.Response[any]{
			StatusCode: 400,
			Err:        err.Error(),
			Body:       helpers.None{},
		}
	}

	return HandleLogin(NewLoginDetails(nil, user), userRepository)
}

func HandleRefresh(body auth2.RefreshRequest, userRepository *repository.User) http.Response[any] {
	newAccessToken, err := jwt.RefreshJwtToken(body.LastValidAccessToken, body.RefreshToken, userRepository)
	if err != nil {
		return http.Response[any]{
			StatusCode: 400,
			Err:        err.Error(),
			Body:       helpers.None{},
		}
	}

	return http.Response[any]{
		StatusCode: 200,
		Body: auth2.RefreshResponse{
			AccessToken: newAccessToken,
		},
	}
}
