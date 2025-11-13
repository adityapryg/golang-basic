package repository

import (
	"github.com/adityapryg/golang-demo/20-mini-project/internal/model"
	"gorm.io/gorm"
)

// TodoRepository handles todo data access
type TodoRepository struct {
	db *gorm.DB
}

// NewTodoRepository creates a new todo repository instance
func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

// Create creates a new todo
func (r *TodoRepository) Create(todo *model.Todo) error {
	return r.db.Create(todo).Error
}

// FindByID finds a todo by ID
func (r *TodoRepository) FindByID(id uint) (*model.Todo, error) {
	var todo model.Todo
	err := r.db.First(&todo, id).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// FindByUserID finds all todos for a specific user
func (r *TodoRepository) FindByUserID(userID uint) ([]model.Todo, error) {
	var todos []model.Todo
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&todos).Error
	return todos, err
}

// FindByUserIDWithFilters finds todos with filters (status, priority)
func (r *TodoRepository) FindByUserIDWithFilters(userID uint, status, priority string) ([]model.Todo, error) {
	query := r.db.Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	var todos []model.Todo
	err := query.Order("created_at DESC").Find(&todos).Error
	return todos, err
}

// Update updates a todo
func (r *TodoRepository) Update(todo *model.Todo) error {
	return r.db.Save(todo).Error
}

// Delete soft deletes a todo
func (r *TodoRepository) Delete(id uint) error {
	return r.db.Delete(&model.Todo{}, id).Error
}

// ExistsByID checks if a todo exists by ID
func (r *TodoRepository) ExistsByID(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Todo{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// IsOwnedByUser checks if a todo belongs to a specific user
func (r *TodoRepository) IsOwnedByUser(todoID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Todo{}).Where("id = ? AND user_id = ?", todoID, userID).Count(&count).Error
	return count > 0, err
}
