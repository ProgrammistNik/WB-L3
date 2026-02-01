package middlerware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("secret")

func Auth(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		role := claims["role"].(string)
		username := claims["username"].(string)

		c.Set("user", username)
		c.Set("role", role)

		if len(requiredRoles) > 0 {
			allowed := false
			for _, r := range requiredRoles {
				if r == role {
					allowed = true
				}
			}
			if !allowed {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		c.Set("sqlSetUser", "SET app.user = '"+username+"'")

		c.Next()
	}
}

func GenerateToken(username, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}
