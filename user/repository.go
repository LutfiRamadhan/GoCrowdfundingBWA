package user

import "gorm.io/gorm"

type Repository interface {
	Save(user User) (User, error)
	Get(user User) (User, error)
	Update(user User) (User, error)
	Gets(user User) ([]User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(user User) (User, error) {
	if err := r.db.Save(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *repository) Get(user User) (User, error) {
	var response User
	tx := r.db
	if user.ID != 0 {
		tx = tx.Where("id = ?", user.ID)
	}
	if user.Email != "" {
		tx = tx.Where("email = ?", user.Email)
	}
	if err := tx.Find(&response).Error; err != nil {
		return User{}, err
	}
	return response, nil
}

func (r *repository) Gets(user User) ([]User, error) {
	var response []User
	tx := r.db
	if user.Name != "" {
		tx = tx.Where("name = ?", user.Name)
	}
	if user.Email != "" {
		tx = tx.Where("email =?", user.Email)
	}
	if user.Occupation != "" {
		tx = tx.Where("occupation = ?", user.Occupation)
	}
	if err := tx.Find(&response).Error; err != nil {
		return []User{}, err
	}
	return response, nil
}

func (r *repository) Update(user User) (User, error) {
	if err := r.db.Save(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}
