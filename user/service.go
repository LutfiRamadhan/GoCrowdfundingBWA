package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	LoginUser(input LoginUserInput) (User, error)
	ValidateEmail(input ValidateEmailInput) (bool, error)
	SaveAvatar(id int, fileLoc string) (User, error)
	FindByID(id int) (User, error)
	Get(criteria User) ([]User, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return User{}, err
	}
	user := User{
		Name:       input.Name,
		Email:      input.Email,
		Password:   string(password),
		Occupation: input.Occupation,
		ProfilePic: "",
		Role:       "user",
	}

	newUser, err := s.repository.Save(user)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (s *service) LoginUser(input LoginUserInput) (User, error) {
	user := User{
		Email: input.Email,
	}
	userData, err := s.repository.Get(user)
	if err != nil {
		return User{}, err
	}
	if userData.ID == 0 {
		return User{}, errors.New("Email/Password not match")
	}
	matchPassword := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(input.Password))
	if matchPassword != nil {
		return User{}, errors.New("Email/Password not match")
	}
	return userData, nil
}

func (s *service) ValidateEmail(input ValidateEmailInput) (bool, error) {
	user := User{
		Email: input.Email,
	}
	userData, err := s.repository.Get(user)
	if err != nil {
		return false, err
	}
	if userData.ID != 0 {
		return false, nil
	}
	return true, nil
}

func (s *service) SaveAvatar(id int, fileLoc string) (User, error) {
	user := User{
		ID: id,
	}
	userData, err := s.repository.Get(user)
	if err != nil {
		return User{}, err
	}

	userData.ProfilePic = fileLoc
	updatedUser, err := s.repository.Update(userData)
	if err != nil {
		return User{}, err
	}
	return updatedUser, nil
}

func (s *service) FindByID(id int) (User, error) {
	user := User{
		ID: id,
	}
	userData, err := s.repository.Get(user)
	if err != nil {
		return User{}, err
	}

	if userData.ID == 0 {
		return User{}, errors.New("User not found with that ID")
	}
	return userData, nil
}

func (s *service) Get(criteria User) ([]User, error) {
	user := User{
		Name:       criteria.Name,
		Occupation: criteria.Occupation,
		Email:      criteria.Email,
	}
	userData, err := s.repository.Gets(user)
	if err != nil {
		return []User{}, err
	}

	return userData, nil
}
