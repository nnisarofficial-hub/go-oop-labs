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
	return fmt.Sprintf("Error: validation error on field '%s': %s", e.Field, e.Message)
}

type DuplicateError struct {
	Resource string
	Value    string
}

func (e *DuplicateError) Error() string {
	return fmt.Sprintf("Error: duplicate '%s': '%s' already exists", e.Resource, e.Value)
}

type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Error: user with id '%s' not found", e.ID)
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
		validation := ValidationError{
			Field:   "username",
			Message: "must be at least 3 characters",
		}
		fmt.Println(validation.Error())
	}
	if !strings.Contains(email, "@") {
		validation := ValidationError{
			Field:   "email",
			Message: "invalid format",
		}
		fmt.Println(validation.Error())
	}
	if age < 18 {
		validation := ValidationError{
			Field:   "age",
			Message: "must be 18 or older",
		}
		fmt.Println(validation.Error())
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

func main() {
	store := NewUserStore()
	fmt.Println(`Register("al", "ali@email.com", 25):`)
	_, err := store.Register("al", "ali@email.com", 25)
	if err != nil {
		fmt.Println("  ", err)
	}
	fmt.Println()
	fmt.Println(`Register("ali", "not-an-email", 25):`)
	_, err = store.Register("ali", "not-an-email", 25)
	if err != nil {
		fmt.Println("  ", err)
	}
	fmt.Println()
	fmt.Println(`Register("ali", "ali@email.com", 15):`)
	_, err = store.Register("ali", "ali@email.com", 15)
	if err != nil {
		fmt.Println("  ", err)
	}
	fmt.Println()
	fmt.Println(`Register("ali", "ali@email.com", 25):`)
	user, err := store.Register("ali", "ali@email.com", 25)
	if err == nil {
		fmt.Printf("   ✓ User created (ID: %d)\n", user.ID)
	}
	fmt.Println()
	fmt.Println(`Register("ali", "ali2@email.com", 30):`)
	_, err = store.Register("ali", "ali2@email.com", 30)
	if err != nil {
		fmt.Println("  ", err)
	}
	fmt.Println()
	fmt.Println(`FindByID(999):`)
	_, err = store.FindByID(999)
	if err != nil {
		fmt.Println("  ", err)
	}
}
