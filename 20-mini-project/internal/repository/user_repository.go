package repository

import (
	"errors"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/model"
	"gorm.io/gorm"
)

// UserRepository handles user data access operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user into database
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByID retrieves user by ID
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil user, no error for not found
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername retrieves user by username
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil user, no error for not found
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail retrieves user by email
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil user, no error for not found
		}
		return nil, err
	}
	return &user, nil
}

// Update updates user data
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes a user
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// ExistsByUsername checks if username already exists
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail checks if email already exists
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
