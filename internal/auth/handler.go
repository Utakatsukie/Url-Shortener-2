package auth

import (
	"Url-Shortener-2/configs"
	"Url-Shortener-2/pkg/res"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
)

type AuthHandler struct {
	*configs.Config
}

type AuthHandlerDeps struct {
	*configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println(handler.Config.Auth.Secret)
		// read body
		var payload LoginRequest
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			res.Json(w, err.Error(), 402)
			return
		}
		if payload.Email == "" {
			res.Json(w, "Email required", 402)
			return
		}
		//reg, _ := regexp.Compile(`[A-Za-z0-9\._%+\-]+@[A-Za-z0-9\.\-]+\.[A-Za-z]{2,}`)
		//ещё есть варианты проверки почты, видос пересмотри

		if _, err := mail.ParseAddress(payload.Email); err != nil {
			res.Json(w, "Wrong email", 402)
			return
		}
		if payload.Password == "" {
			res.Json(w, "Password required", 402)
			return
		}
		fmt.Println(payload)
		data := LoginResponse{
			Token: "123",
		}
		res.Json(w, data, 200)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Register")
	}
}
