package repositories

import "github.com/gonzalohonorato/servercorego/core/user/domain/entities"

type UserRepository interface {
	SearchUserByID(id string) (*entities.User, error)
	CreateUser(user *entities.User) error
	UpdateUser(user *entities.User) error
	SearchUsers() (*entities.Users, error)
	SearchUsersByType(userType string) (*entities.Users, error)
}
