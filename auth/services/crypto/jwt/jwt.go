package jwt

import (
	"auth/data/models"
	"auth/data/repository"
	"github.com/golang-jwt/jwt/v4"
	"lib/data/models/auth"
	"os"
	"time"
)

const jwtExp = 24 * time.Hour
const refreshExp = 14 * 24 * time.Hour
const signKeyArg = "SIGN_KEY"
const issuerArg = "ISSUER"

func getSignKey() []byte {
	return []byte(os.Getenv(signKeyArg))
}

func GenerateJwtToken(user auth.User) (string, error) {
	issuer := os.Getenv(issuerArg)
	claims := models.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(jwt.TimeFunc().Add(jwtExp)),
		},
		Username: user.Username,
		Email:    user.Email,
		UserId:   user.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getSignKey())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(user auth.User) (string, error) {
	issuer := os.Getenv(issuerArg)
	claims := models.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(jwt.TimeFunc().Add(refreshExp)),
		},
		Username: user.Username,
		UserId:   user.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getSignKey())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJwtToken(tokenString string) (*models.JwtClaims, error) {
	issuer := os.Getenv(issuerArg)
	token, err := jwt.ParseWithClaims(tokenString, &models.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getSignKey(), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.JwtClaims)
	if !ok {
		return nil, err
	}

	if claims.Issuer != issuer {
		return nil, err
	}

	if err := claims.Valid(); err != nil {
		return nil, err
	}

	return claims, nil
}

func VerifyRefreshToken(oldJwt string, refreshString string) (*models.JwtClaims, error) {
	issuer := os.Getenv(issuerArg)
	token, err := jwt.ParseWithClaims(refreshString, &models.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getSignKey(), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.JwtClaims)
	if !ok {
		return nil, err
	}

	if claims.Issuer != issuer {
		return nil, err
	}

	if err := claims.Valid(); err != nil {
		return nil, err
	}

	oldClaims, err := VerifyJwtToken(oldJwt)
	if err != nil {
		return nil, err
	}

	if oldClaims.UserId != claims.UserId {
		return nil, err
	}

	return claims, nil
}

func RefreshJwtToken(oldJwt string, refreshString string, userRepository *repository.User) (string, error) {
	claims, err := VerifyRefreshToken(oldJwt, refreshString)
	if err != nil {
		return "", err
	}

	user, err := userRepository.GetById(claims.UserId)
	if err != nil {
		return "", err
	}

	return GenerateJwtToken(*user)
}
