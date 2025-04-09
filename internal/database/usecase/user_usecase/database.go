package user_usecase

import (
	"fmt"
)

func (uc *InitUseCase) InitDatabase() error {
	err := uc.repo.RunMigrations()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	return nil
}
