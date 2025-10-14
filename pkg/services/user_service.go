package services

import (
	"errors"
	"fmt"
	"regexp"
	"share-docs/pkg/app/domain/userapp"
	"share-docs/pkg/db/models"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrFailedToCreateUser = errors.New("failed to create user")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrAccountLocked      = errors.New("account is temporarily locked")
	ErrAccountNotVerified = errors.New("account is not verified")
	ErrAccountInactive    = errors.New("account is inactive")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrWeakPassword       = errors.New("password does not meet requirements")
)

type UserServiceInterface interface {
	CreateUser(email, password, firstName, lastName string, birthDate *time.Time) (*userapp.User, error)
	GetUserByID(userID string) (*userapp.User, error)
	GetUserByEmail(email string) (*userapp.User, error)
	LoginWithEmailPassword(email, password string) (*userapp.User, error)
}

type UserService struct {
	db                *gorm.DB
	emailRegex        *regexp.Regexp
	passwordMinLength int
	bcryptCost        int
}

func NewUserService(db *gorm.DB) *UserService {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	return &UserService{
		db:                db,
		emailRegex:        emailRegex,
		passwordMinLength: 8,
		bcryptCost:        5,
	}
}

func (s *UserService) CreateUser(email, password, firstName, lastName string, birthDate *time.Time) (*userapp.User, error) {
	if err := s.validateEmail(email); err != nil {
		return nil, err
	}

	if err := s.validatePassword(password); err != nil {
		return nil, err
	}

	var existingUser models.User
	res := s.db.Where("email = ?", email).First(&existingUser)

	if res.Error == nil {
		return nil, ErrEmailAlreadyExists
	}

	if res.Error != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("error querying for a user with %s email: %w", email, res.Error)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)

	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	modelUser := &models.User{
		Email:      email,
		Password:   string(hashedPassword),
		IsActive:   true,
		IsVerified: false,
		FirstName:  firstName,
		LastName:   lastName,
		BirthDate:  birthDate,
	}

	if result := s.db.Create(modelUser); result.Error != nil {
		return nil, ErrFailedToCreateUser
	}

	user := userapp.ToAppUser(*modelUser)
	return &user, nil
}

func (s *UserService) GetUserByID(userIDString string) (*userapp.User, error) {
	userID, err := uuid.Parse(userIDString)

	if err != nil {
		return nil, err
	}

	var modelUser *models.User

	result := s.db.First(&modelUser, userID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}

		return nil, result.Error
	}

	user := userapp.ToAppUser(*modelUser)
	return &user, nil
}

func (s *UserService) GetUserByEmail(email string) (*userapp.User, error) {
	if err := s.validateEmail(email); err != nil {
		return nil, ErrInvalidEmail
	}

	var modelUser *models.User

	result := s.db.Where("email = ?", email).First(&modelUser)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}

		return nil, result.Error
	}

	user := userapp.ToAppUser(*modelUser)
	return &user, nil
}

func (s *UserService) LoginWithEmailPassword(email, password string) (*userapp.User, error) {
	var modelUser *models.User

	result := s.db.Where("email = ?", email).First(&modelUser)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}

		return nil, result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(modelUser.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	user := userapp.ToAppUser(*modelUser)
	return &user, nil
}

func (s *UserService) ValidatePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (s *UserService) validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return ErrInvalidEmail
	}

	if !s.emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}

	return nil
}

func (s *UserService) validatePassword(password string) error {
	if len(password) < s.passwordMinLength {
		return ErrWeakPassword
	}

	return nil
}
