package main

import (
	"Url-Shortener-2/internal/auth"
	"Url-Shortener-2/internal/user"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Email:    "a4@0b.by",
		Password: "$2a$10$2hOlfZQj0TzSnbZzD3z78e881L5yMt9KCuAIxp88Dm2OmHYWyTwa.",
		Name:     "Ivan",
	})
}

func TestLoginSuccess(t *testing.T) {
	// Prepare
	db := initDb()
	initData(db)
	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "a4@0b.by",
		Password: "4",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("Expected %d got %d", 200, res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	var resData auth.LoginResponse
	err = json.Unmarshal(body, &resData)
	if err != nil {
		t.Fatal(err)
	}
	if resData.Token == "" {
		t.Fatal("Token empty")
	}
}

func TestLoginFail(t *testing.T) {
	// Prepare
	db := initDb()
	initData(db)
	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "a4@0b.by",
		Password: "1",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 401 {
		t.Fatalf("Expected %d got %d", 401, res.StatusCode)
	}
}
