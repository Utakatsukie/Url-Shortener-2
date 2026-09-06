package di

import "Url-Shortener-2/internal/user"

type IStatRepository interface {
	AddClick(lickId uint)
}

type IUserRepository interface {
	Create(user *user.User) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
}
