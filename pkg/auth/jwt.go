package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jixlox0/studoto-backend/pkg/cache"
)

type JWTAuth struct {
	secretKey       string
	expirationHours int
	tokenCache      cache.TokenCache
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewJWTAuth(secretKey string, expirationHours int, tokenCache cache.TokenCache) *JWTAuth {
	return &JWTAuth{
		secretKey:       secretKey,
		expirationHours: expirationHours,
		tokenCache:      tokenCache,
	}
}

func (j *JWTAuth) GenerateToken(ctx context.Context, userID uint, email string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(j.expirationHours) * time.Hour)
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	// Cache the token in Redis
	if j.tokenCache != nil {
		expiration := time.Until(expirationTime)
		if err := j.tokenCache.SetToken(ctx, tokenString, userID, expiration); err != nil {
			// Log error but don't fail token generation if cache fails
			// Token is still valid, just not cached
		}
	}

	return tokenString, nil
}

func (j *JWTAuth) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
	// First check if token exists in cache (for faster validation and logout support)
	if j.tokenCache != nil {
		cachedUserID, err := j.tokenCache.GetToken(ctx, tokenString)
		if err == nil {
			// Token is in cache, validate JWT structure
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("invalid signing method")
				}
				return []byte(j.secretKey), nil
			})

			if err != nil {
				return nil, err
			}

			if !token.Valid {
				return nil, errors.New("invalid token")
			}

			// Verify cached userID matches token userID
			if claims.UserID != cachedUserID {
				return nil, errors.New("token user mismatch")
			}

			return claims, nil
		}
		// If token not in cache but JWT is valid, it might have expired from cache
		// Continue with JWT validation
	}

	// Validate JWT token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// If cache is available and token is valid, cache it
	if j.tokenCache != nil && claims.ExpiresAt != nil {
		expirationTime := claims.ExpiresAt.Time
		if expirationTime.After(time.Now()) {
			expiration := time.Until(expirationTime)
			j.tokenCache.SetToken(ctx, tokenString, claims.UserID, expiration)
		}
	}

	return claims, nil
}

// InvalidateToken removes a token from the cache (for logout)
func (j *JWTAuth) InvalidateToken(ctx context.Context, tokenString string) error {
	if j.tokenCache != nil {
		return j.tokenCache.DeleteToken(ctx, tokenString)
	}
	return nil
}

// InvalidateUserTokens removes all tokens for a user (for logout all devices)
func (j *JWTAuth) InvalidateUserTokens(ctx context.Context, userID uint) error {
	if j.tokenCache != nil {
		return j.tokenCache.DeleteUserTokens(ctx, userID)
	}
	return nil
}
