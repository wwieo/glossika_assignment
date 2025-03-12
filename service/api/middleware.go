package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"glossika/service/internal/config"
	boAccount "glossika/service/internal/model/bo/account"
	"net/http"
	"strings"
	"time"
)

const (
	TokenExpiration = 10 * time.Minute
	BearerPrefix    = "Bearer "
)

func ValidateTokenMiddleware(jwtConfig config.JWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
		if tokenString == authHeader {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		claims := &boAccount.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}

			if len(jwtConfig.Secret) == 0 {
				return nil, errors.New("jwt secret is empty")
			}

			return []byte(jwtConfig.Secret), nil
		}, jwt.WithExpirationRequired())

		if err != nil || !token.Valid || claims.Email == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		ctx.Set(boAccount.AccountEmailKey, claims.Email)
		ctx.Next()
	}
}

func GenerateToken(email string, jwtConfig config.JWT) (string, error) {
	now := time.Now()
	claims := &boAccount.JWTClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(TokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtConfig.Secret))
}
