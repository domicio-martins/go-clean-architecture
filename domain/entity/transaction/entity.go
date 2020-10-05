package transaction

import (
	"time"
	"github.com/domicio/go-clean-architecture/domain/entity"
	"github.com/domicio/go-clean-architecture/domain/entity/user"
)

type Transaction struct {
	ID entity.ID
	Payer User.ID
	Payee User.ID
	Status string
	CreatedAt time.Time
	UpdatedAt time.Time
}