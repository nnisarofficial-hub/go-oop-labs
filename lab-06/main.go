package main

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

type DuplicateError struct {
	Resource string
	Value    string
}

func (e *DuplicateError) Error() string {
	return fmt.Sprintf("duplicate '%s': '%s' already exists", e.Resource, e.Value)
}

type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with id '%s' not found", e.Resource, e.ID)
}

type User struct {
	ID       int
	Username string
	Email    string
	Age      int
}

type UserStore struct {
	users  map[int]*User
	nextID int
}

func NewUserStore() *UserStore {
	return &UserStore{
		users:  make(map[int]*User),
		nextID: 1,
	}
}

func (s *UserStore) Register(username, email string, age int) (*User, error) {
	if len(username) < 3 {
		return nil, &ValidationError{
			Field:   "username",
			Message: "must be at least 3 characters",
		}
	}
	if !strings.Contains(email, "@") {
		return nil, &ValidationError{
			Field:   "email",
			Message: "invalid format",
		}
	}
	if age < 18 {
		return nil, &ValidationError{
			Field:   "age",
			Message: "must be 18 or older",
		}
	}
	for _, existingUser := range s.users {
		if existingUser.Username == username {
			return nil, &DuplicateError{
				Resource: "username",
				Value:    username,
			}
		}
	}
	newUser := &User{
		ID:       s.nextID,
		Username: username,
		Email:    email,
		Age:      age,
	}
	s.users[s.nextID] = newUser
	s.nextID++
	return newUser, nil
}

func (s *UserStore) FindByID(id int) (*User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, &NotFoundError{
			Resource: "user",
			ID:       fmt.Sprintf("%d", id),
		}
	}
	return user, nil
}

func printRegisterResult(store *UserStore, username, email string, age int) {
	fmt.Printf(`Register("%s", "%s", %d):`, username, email, age)
	user, err := store.Register(username, email, age)
	if err != nil {
		fmt.Printf("\n  Error: %v\n\n", err)
		return
	}
	fmt.Printf(" ✓ User created (ID: %d)\n", user.ID)
}

func main() {
	store := NewUserStore()

	printRegisterResult(store, "al", "ali@email.com", 25)
	printRegisterResult(store, "ali", "not-an-email", 25)
	printRegisterResult(store, "ali", "ali@email.com", 15)
	printRegisterResult(store, "ali", "ali@email.com", 25)
	printRegisterResult(store, "ali", "ali2@email.com", 30)

	fmt.Printf(`FindByID(999):`)
	_, err := store.FindByID(999)
	if err != nil {
		fmt.Printf("\n  Error: %v\n", err)
	}
}
