package inmemory

import (
	"sync"

	errs "github.com/hrabalvojta/micro-dvdrental/errors"
	"github.com/hrabalvojta/micro-dvdrental/films"
)

// inMemUserRepository is an implementation of a user repository for storage in local memory
type inMemUserRepository struct {
	mtx   *sync.RWMutex
	users map[int]*films.User
}

// NewInMemUserRepository returns a new user repository for storage in local memory
func NewInMemUserRepository() films.Repository {
	return &inMemUserRepository{
		mtx:   new(sync.RWMutex),
		users: make(map[int]*films.User),
	}
}

// Store inserts a user into the local user map
func (ir *inMemUserRepository) Store(user *films.User) error {
	ir.mtx.Lock()
	ir.users[user.ID] = user
	ir.mtx.Unlock()
	return nil
}

// Find retrieves a single user from the repository
func (ir *inMemUserRepository) Find(id int) (*films.User, error) {
	ir.mtx.RLock()
	u := ir.users[id]
	ir.mtx.RUnlock()

	if u == nil {
		return nil, errs.ErrUserNotFound
	}
	return u, nil
}

// FindAll retrieves all users from memory
func (ir *inMemUserRepository) FindAll() []*films.User {
	ir.mtx.RLock()
	allUsers := []*films.User{}
	for _, v := range ir.users {
		allUsers = append(allUsers, v)
	}
	ir.mtx.RUnlock()
	return allUsers
}
