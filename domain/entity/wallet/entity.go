package wallet

import (
	"time"
	"github.com/domicio-martins/go-clean-architecture/domain/user"
	"github.com/domicio-martins/go-clean-architecture/domain/entity"
)

type Wallet struct {
	ID entity.ID
	Amount float64
	User user.ID
	CreatedAt time.Time
	UpdatedAt time.Time
}