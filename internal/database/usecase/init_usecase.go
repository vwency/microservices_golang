package usecase

import (
	"github.com/vwency/microservices_golang/internal/database/models"
	"github.com/vwency/microservices_golang/internal/database/repository"
)

type InitUseCase struct {
	repo repository.UserRepository
}

func NewInitUseCase(repo repository.UserRepository) *InitUseCase {
	return &InitUseCase{repo: repo}
}

func (uc *InitUseCase) InitDatabase() error {
	return uc.repo.RunMigrations()
}

func (uc *InitUseCase) GetUser(username, email string) (*models.User, error) {
	return uc.repo.GetUserByUsernameOrEmail(username, email)
}

func (uc *InitUseCase) AddUser(username, password, hashedRt, accessRt string) error {
	user := &models.User{
		Username: username,
		Password: password,
		HashedRt: hashedRt,
		HashedAt: accessRt,
	}
	return uc.repo.AddUser(user)
}

func (uc *InitUseCase) UpdateUserTokens(username, hashedRt, accessRt string) error {
	return uc.repo.UpdateUserTokens(username, hashedRt, accessRt)
}

type DatabaseInit interface {
	InitDatabase() error
	GetUser(username, email string) (*models.User, error)
	AddUser(username, password, hashedRt, accessRt string) error
	UpdateUserTokens(username, hashedRt, accessRt string) error
}
