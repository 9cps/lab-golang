package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// CorsMiddleware returns a CORS handler that reads allowed origins from the
// CORS_ALLOWED_ORIGINS env var (comma-separated). If the list is empty, no
// origins are allowed. Using "*" with credentials violates the CORS spec, so
// this middleware refuses to set credentials when a wildcard is configured.
func CorsMiddleware() gin.HandlerFunc {
	raw := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	var allowed []string
	for _, o := range strings.Split(raw, ",") {
		if o = strings.TrimSpace(o); o != "" {
			allowed = append(allowed, o)
		}
	}
	allowCredentials := len(allowed) > 0 && !contains(allowed, "*")

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && contains(allowed, origin) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			if allowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		// Allow any request header. We reflect the headers the browser asks for
		// in its preflight (Access-Control-Request-Headers) instead of using a
		// literal "*", because "*" is interpreted literally when credentials are
		// allowed and would reject real headers such as Authorization. When the
		// preflight does not advertise any headers, fall back to "*".
		if reqHeaders := c.GetHeader("Access-Control-Request-Headers"); reqHeaders != "" {
			c.Header("Access-Control-Allow-Headers", reqHeaders)
		} else {
			c.Header("Access-Control-Allow-Headers", "*")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// AuthEnabled reports whether JWT authentication should be enforced. It reads
// the AUTH_ENABLED env var and defaults to true (auth on) when unset. The
// values "false", "0", "no", and "off" (case-insensitive) disable auth.
func AuthEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("AUTH_ENABLED"))) {
	case "false", "0", "no", "off":
		return false
	default:
		return true
	}
}

// AuthMiddleware verifies a JWT bearer token using the HS256 algorithm and a
// secret loaded from the JWT_SECRET env var. It protects against alg-confusion
// attacks by pinning the allowed signing methods.
//
// When AUTH_ENABLED is set to a falsy value the middleware becomes a no-op
// pass-through, allowing every request without a token. This is intended for
// local development only.
func AuthMiddleware() gin.HandlerFunc {
	if !AuthEnabled() {
		return func(c *gin.Context) { c.Next() }
	}

	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		panic("JWT_SECRET is not set")
	}

	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header must be a Bearer token"})
			return
		}
		tokenString := strings.TrimPrefix(header, "Bearer ")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header must be a Bearer token"})
			return
		}

		token, err := jwt.Parse(
			tokenString,
			func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return secret, nil
			},
			jwt.WithValidMethods([]string{"HS256"}),
		)
		if err != nil || !token.Valid {
			msg := "invalid token"
			if errors.Is(err, jwt.ErrTokenExpired) {
				msg = "token expired"
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if _, hasExp := claims["exp"]; !hasExp {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token missing exp claim"})
				return
			}
		}

		c.Set("user", token.Claims)
		c.Next()
	}
}

func contains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}
