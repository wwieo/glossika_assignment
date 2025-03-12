package boAccount

import "github.com/golang-jwt/jwt/v5"

const (
	AccountEmailKey = "account_email"
)

type Account struct {
	Email    string `json:"email" valid:"required"`
	Password string `json:"password" valid:"required"`
}

type RegisterArgs struct {
	Account
}

type RegisterReply struct {
	VerifyURL string `json:"verify_url"`
}

type JWTClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type LoginArgs struct {
	Account
}

type LoginReply struct {
	Token string `json:"token"`
}
