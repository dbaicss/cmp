

package auth

import (
"fmt"
"net/http"

"cmp-server/model"
c "cmp-server/common"
jwt "github.com/dgrijalva/jwt-go"
"time"

)




func GenerateToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.UserName,
		"exp":      time.Now().Add(time.Minute * 60).Unix(),
	})

	return token.SignedString([]byte("secret"))
}
func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("authorization")
		if tokenStr == "" {
			c.ResponseWithJson(w, http.StatusUnauthorized,
				c.Response{ErrNo: http.StatusUnauthorized, Msg: "not authorized"})
		} else {
			token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					c.ResponseWithJson(w, http.StatusUnauthorized,
						c.Response{ErrNo: http.StatusUnauthorized, Msg: "not authorized"})
					return nil, fmt.Errorf("not authorization")
				}
				return []byte("secret"), nil
			})
			if !token.Valid {
				c.ResponseWithJson(w, http.StatusUnauthorized,
					c.Response{ErrNo: http.StatusUnauthorized, Msg: "not authorized"})
			} else {
				next.ServeHTTP(w, r)
			}
		}
	})
}

