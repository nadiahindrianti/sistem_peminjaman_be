package repositories

import (
	"sistem_peminjaman_be/models"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	UserGetAll(page, limit int, search string) ([]models.User, int, error)
	UserGetDetail(id uint, isDeleted bool) (models.User, error)
	UserGetById(id uint) (models.User, error)
	UserGetById2(id uint) (models.User, error)
	UserGetByEmail(email string) (models.User, error)
	UserGetByEmail2(id uint, email string) (models.User, error)
	UserGetByEmail3(email string) (models.User, error)
	UserCreate(user models.User) (models.User, error)
	UserCreate2(user models.User, isActive bool) (models.User, error)
	UserUpdate(user models.User) (models.User, error)
	UserUpdate2(user models.User, isActive bool) (models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) UserGetAll(page, limit int, search string) ([]models.User, int, error) {
	var users []models.User
	var count int64
	err := r.db.Unscoped().Find(&users).Count(&count).Error

	offset := (page - 1) * limit

	err = r.db.Unscoped().Where("role = 'user' AND full_name LIKE ? OR role = 'user' AND email LIKE ?", "%"+search+"%", "%"+search+"%").Order("id DESC").Limit(limit).Offset(offset).Find(&users).Error

	return users, int(count), err
}

func (r *userRepository) UserGetDetail(id uint, isDeleted bool) (models.User, error) {
	var user models.User
	err := r.db.Unscoped().Where("id = ? AND role = 'user'", id).First(&user).Error
	if isDeleted {
		user.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
		err = r.db.Where("id = ?", id).Save(user).Error
	} else {
		user.DeletedAt = gorm.DeletedAt{}
		err = r.db.Where("id = ?", id).Save(user).Error
	}
	return user, err
}

func (r *userRepository) UserGetById(id uint) (models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return user, err
}

func (r *userRepository) UserGetById2(id uint) (models.User, error) {
	var user models.User
	err := r.db.Unscoped().Where("id = ? AND role = 'user'", id).First(&user).Error
	return user, err
}

func (r *userRepository) UserGetByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *userRepository) UserGetByEmail3(email string) (models.User, error) {
	var user models.User
	err := r.db.Unscoped().Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *userRepository) UserGetByEmail2(id uint, email string) (models.User, error) {
	var user models.User
	err := r.db.Where("id != ? AND email = ?", id, email).First(&user).Error
	return user, err
}

func (r *userRepository) UserCreate(user models.User) (models.User, error) {
	err := r.db.Create(&user).Error
	return user, err
}

func (r *userRepository) UserCreate2(user models.User, isActive bool) (models.User, error) {
	if !isActive {
		user.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	}
	err := r.db.Create(&user).Error
	return user, err
}

func (r *userRepository) UserUpdate(user models.User) (models.User, error) {
	err := r.db.Save(&user).Error
	return user, err
}

func (r *userRepository) UserUpdate2(user models.User, isActive bool) (models.User, error) {
	err := r.db.Unscoped().Save(&user).Error
	if isActive {
		user.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
		err = r.db.Save(user).Error
	} else {
		user.DeletedAt = gorm.DeletedAt{}
		err = r.db.Save(user).Error
	}
	return user, err
}
