package repository

import (
	"scanfetcher/internal/db"
	"scanfetcher/internal/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(dbConn *gorm.DB) *UserRepository {
	return &UserRepository{
		db: dbConn,
	}
}

func (r *UserRepository) GetByID(id uint) (*domain.User, error) {
	var dbUser db.User
	if err := r.db.First(&dbUser, id).Error; err != nil {
		return nil, err
	}

	return toDomainUser(dbUser), nil
}

func (r *UserRepository) Create(user domain.User, hashedPassword string) error {
	dbUser := toDBUser(user, hashedPassword)
	return r.db.Create(&dbUser).Error
}

func (r *UserRepository) GetAll() ([]*domain.User, error) {
    var dbUsers []db.User
    if err := r.db.Find(&dbUsers).Error; err != nil {
        return nil, err
    }

    domainUsers := make([]*domain.User, len(dbUsers))
    for i, u := range dbUsers {
        domainUsers[i] = toDomainUser(u)
    }

    return domainUsers, nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
    var dbUser db.User
    if err := r.db.Where("email = ?", email).First(&dbUser).Error; err != nil {
        return nil, err
    }

    return toDomainUser(dbUser), nil
}

// Mapping functions

func toDomainUser(u db.User) *domain.User {
	return &domain.User{
		ID:    int(u.ID),
		Name:  u.Username,
		Email: u.Email,
	}
}

func toDBUser(u domain.User, password string) db.User {
	return db.User{
		Email:    u.Email,
		Username: u.Name,
		Password: password,
	}
}
