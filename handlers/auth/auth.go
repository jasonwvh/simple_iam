package auth

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jasonwvh/simple_iam/handlers"
	"github.com/jasonwvh/simple_iam/storage"
)

func LoginHandler(env *handlers.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		log.Printf("found auth header %s", tokenString)
		ok, _ := VerifyJWT(tokenString)
		if ok {
			handlers.HandleSuccess(w, http.StatusOK, "login success")
			return
		}

		var user *storage.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			handlers.HandleError(w, http.StatusBadRequest, err.Error())
			return
		}

		if user.Username == "" || user.Password == "" {
			handlers.HandleError(w, http.StatusBadRequest, "Username or password cannot be empty")
			return
		}

		authSuccess, err := Authenticate(env, user.Username, user.Password)
		if err != nil {
			handlers.HandleError(w, http.StatusInternalServerError, err.Error())
			return
		}

		var token string
		if authSuccess {
			token, _ = Authorize(env, user.Username)
		} else {
			token = ""
		}

		handlers.HandleSuccess(w, http.StatusOK, token)
	}
}

func VerifyJWT(tokenString string) (bool, error) {
	if tokenString == "" {
		return false, fmt.Errorf("empty token")
	}

	secret := []byte("super-secret")

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return secret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["user"], claims["nbf"])
		return true, nil
	} else {
		fmt.Println(err)
		return false, err
	}
}

func Authenticate(env *handlers.Env, username string, password string) (bool, error) {
	dbUser, err := storage.GetUser(env.DB, username)
	if err != nil {
		return false, err
	}

	salt := dbUser.Salt
	saltPass := password + salt

	hasher := sha512.New512_256()
	hasher.Write([]byte(saltPass))
	hashedPass := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	log.Printf("found user: %v", dbUser)
	if hashedPass != dbUser.Password {
		log.Printf("password mismatch %s and %s", hashedPass, dbUser.Password)
		return false, err
	}

	return true, nil
}

func Authorize(env *handlers.Env, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": username,
		"nbf":  time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	tokenString, err := token.SignedString([]byte("super-secret"))
	log.Printf("jwt token %v (%s)", token, tokenString)

	return tokenString, err
}
