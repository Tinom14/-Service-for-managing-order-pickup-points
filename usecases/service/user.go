package service

import (
	"avito_test/config"
	"avito_test/domain"
	"avito_test/repository"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type User struct {
	repo repository.User
}

func NewUserService(repo repository.User) *User {
	return &User{repo: repo}
}

func (u *User) GetToken(id string, role string) (string, error) {
	payload := jwt.MapClaims{
		"id":   id,
		"role": role,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return "", errors.New("token create error")
	}
	return tokenString, nil
}

func (u *User) Register(email string, password string, role string) (domain.User, error) {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashPasswordStr := string(hashPassword)
	return u.repo.Register(email, hashPasswordStr, role)
}

func (u *User) Login(email string, password string) (string, error) {
	user, err := u.repo.Login(email)
	if err != nil {
		return "", errors.New("wrong email")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("wrong password")
	}
	token, err := u.GetToken(strconv.Itoa(user.Id), user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}
