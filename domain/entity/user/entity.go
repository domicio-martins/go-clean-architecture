package user

import (
	"time"
	"github.com/domicio/go-clean-architecture/domain/entity"
)

type User struct {
	ID entity.ID
	Name string
	Surname string
	Type string
	CreatedAt time.Time
	UpdatedAt time.Time
}