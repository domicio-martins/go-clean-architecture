package presenter

import (
	"github.com/domicio-martins/go-clean-architecture/domain/entity"
)

//User data
type User struct {
	ID        entity.ID `json:"id"`
	Name     string    `json:"email"`
	SurName string    `json:"first_name"`
}
