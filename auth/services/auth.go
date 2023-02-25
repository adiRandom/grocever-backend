package services

import (
	"auth/data/models"
	"auth/data/repository"
	"auth/services/crypto/jwt"
	dto "lib/data/dto/auth"
	"lib/data/models/auth"
	"lib/helpers"
	"lib/network/http"
)

type LoginDetails struct {
	body *dto.LoginRequest
	user *auth.User
}

func NewLoginDetails(body *dto.LoginRequest, user *auth.User) LoginDetails {
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
		Body: dto.AuthResponse{
			AccessToken:  token,
			RefreshToken: refreshToken,
			User:         dto.NewUserFromModel(*user),
		},
	}
}

func HandleRegister(body dto.RegisterRequest, userRepository *repository.User) http.Response[any] {
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

func HandleRefresh(body dto.RefreshRequest, userRepository *repository.User) http.Response[any] {
	newAccessToken, err := jwt.RefreshJwtToken(body.LastValidAccessToken, body.RefreshToken, userRepository)
	if err != nil {
		return http.Response[any]{
			StatusCode: 403,
			Err:        err.Error(),
			Body:       helpers.None{},
		}
	}

	return http.Response[any]{
		StatusCode: 200,
		Body: dto.RefreshResponse{
			AccessToken: newAccessToken,
		},
	}
}

func HandleValidate(body dto.ValidateRequest) http.Response[any] {
	claims, err := jwt.VerifyJwtToken(body.AccessToken)
	if err != nil {
		return http.Response[any]{
			StatusCode: 400,
			Err:        err.Error(),
			Body:       helpers.None{},
		}
	}

	return http.Response[any]{
		StatusCode: 200,
		Body: dto.ValidateResponse{
			UserId: int(claims.UserId),
		},
	}
}
