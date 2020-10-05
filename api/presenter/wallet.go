package presenter

import (
	"github.com/domicio-martins/go-clean-architecture/domain/entity"
	"time"
)

//Book data
type Wallet struct {
	ID       entity.ID `json:"id"`
	Amount    float64    `json:"amount"`
	User   string    `json:"user"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
