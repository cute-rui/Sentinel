package auth

import (
	"Sentinel/dao/models"
	utils "Sentinel/utils/string"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

type AuthKey struct{}

const (
	DefaultExpireTime = time.Hour * 24 * 7
)

type JWTClaim struct {
	UserID   int
	Username string
	MFA      bool
	Remember bool
	jwt.RegisteredClaims
}

type JWTClaimOption func(*JWTClaim)

func NewJWTClaim(u *models.User, opt ...JWTClaimOption) *JWTClaim {
	if u == nil {
		return &JWTClaim{}
	}

	c := &JWTClaim{
		UserID:   u.ID,
		Username: u.Username,
		MFA:      false,
		Remember: false,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Sentinel",
			Subject:   u.Username,
			Audience:  jwt.ClaimStrings{"Sentinel_AUTH"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(DefaultExpireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        utils.RandString(16),
		},
	}

	for i := range opt {
		opt[i](c)
	}

	return c

}

func WithUsername(username string) JWTClaimOption {
	return func(c *JWTClaim) {
		c.Username = username
	}
}

func WithUID(id int) JWTClaimOption {
	return func(c *JWTClaim) {
		c.UserID = id
	}
}

func WithMFAStatus(status bool) JWTClaimOption {
	return func(c *JWTClaim) {
		c.MFA = status
	}
}

func WithPermission(permission string) JWTClaimOption {
	return func(c *JWTClaim) {

	}
}

func WithRemember(remember bool) JWTClaimOption {
	return func(c *JWTClaim) {
		c.Remember = remember
		c.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 128))
	}
}

func (c *JWTClaim) Sign(key string) (string, error) {
	signedString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(key))
	if err != nil {
		return "", errors.New("generate token failed:" + err.Error())
	}
	return signedString, nil
}

func GetClaimsFromContext(ctx context.Context) (*JWTClaim, error) {
	c, ok := fromContext(ctx)
	if !ok {
		return nil, errors.New(`invalid claim`)
	}

	output, ok := c.(*JWTClaim)
	if !ok {
		return nil, errors.New(`invalid claim`)
	}

	return output, nil
}

// NewContext put authentication info into context
func NewContext(ctx context.Context, info jwt.Claims) context.Context {
	return context.WithValue(ctx, AuthKey{}, info)
}

func bytesToJWTClaim(bytes []byte, key string) (*JWTClaim, error) {
	return StringToJWTClaim(utils.ByteToString(bytes), key)
}

func TrimBearerScheme(t string) string {
	return strings.TrimPrefix(t, `Bearer `)
}

func StringToJWTClaim(tokenString string, key string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	//Todo: refactor error output
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("claim invalid")
	}

	claim, ok := token.Claims.(*JWTClaim)
	if !ok {
		return nil, errors.New("invalid claim type")
	}

	return claim, nil
}

// FromContext extract authentication info from context
func fromContext(ctx context.Context) (token jwt.Claims, ok bool) {
	token, ok = ctx.Value(AuthKey{}).(jwt.Claims)
	return
}
