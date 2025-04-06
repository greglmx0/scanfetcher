package usecase

import (
	"scanfetcher/internal/domain"
	"scanfetcher/internal/repository"
)

type UserUseCase struct {
    userRepo *repository.UserRepository
}

func NewUserUseCase(repo *repository.UserRepository) *UserUseCase {
    return &UserUseCase{userRepo: repo}
}

func (u *UserUseCase) GetUserByID(id int) (*domain.User, error) {
    return u.userRepo.GetByID(uint(id)) // Conversion int -> uint
}

func (u *UserUseCase) GetAllUsers() ([]*domain.User, error) {
    return u.userRepo.GetAll()
}

func (u *UserUseCase) Create(user domain.User, hashedPassword string) error {
    return u.userRepo.Create(user, hashedPassword)
}

func (u *UserUseCase) GetUserByEmail(email string) (*domain.User, error) {
    return u.userRepo.GetByEmail(email)
}