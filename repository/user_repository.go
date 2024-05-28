package repository

import (
	"fmt"
	"sync"

	"github.com/grpc-user-service/models"
)

type UserRepository interface {
	CreateUser(user models.User) (int32, error)
	GetUser(id int32) (models.User, error)
	GetUsers(ids []int32) ([]models.User, error)
	SearchUsersByField(field string, value interface{}) ([]models.User, error)
}

type InMemoryUserRepository struct {
	users  map[int32]models.User
	lastID int32
	mu     sync.Mutex
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:  make(map[int32]models.User),
		lastID: 0,
	}
}

func (r *InMemoryUserRepository) CreateUser(user models.User) (int32, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.lastID++
	user.ID = r.lastID
	r.users[r.lastID] = user
	return r.lastID, nil
}

func (r *InMemoryUserRepository) GetUser(id int32) (models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return models.User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *InMemoryUserRepository) GetUsers(ids []int32) ([]models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var users []models.User
	for _, id := range ids {
		if user, ok := r.users[id]; ok {
			users = append(users, user)
		}
	}
	return users, nil
}

func (r *InMemoryUserRepository) SearchUsersByField(field string, value interface{}) ([]models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var users []models.User
	for _, user := range r.users {
		switch field {
		case "ID":
			if id, ok := value.(int32); ok && user.ID == id {
				users = append(users, user)
			}
		case "FName":
			if fname, ok := value.(string); ok && user.FName == fname {
				users = append(users, user)
			}
		case "City":
			if city, ok := value.(string); ok && user.City == city {
				users = append(users, user)
			}
		case "Phone":
			if phone, ok := value.(int64); ok && user.Phone == phone {
				users = append(users, user)
			}
		case "Height":
			if height, ok := value.(float32); ok && user.Height == height {
				users = append(users, user)
			}
		case "Married":
			if married, ok := value.(bool); ok && user.Married == married {
				users = append(users, user)
			}
		}
	}

	return users, nil
}
