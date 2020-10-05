package transaction

import (
	"strings"
	"time"

	"github.com/domicio-martins/go-clean-architecture/domain/entity"
	"github.com/domicio-martins/go-clean-architecture/domain/wallet"
	"github.com/domicio-martins/go-clean-architecture/domain/user"
)

type manager struct {
	repo repository
}

//NewManager create new manager
func NewManager(r repository) *manager {
	return &manager{
		repo: r,
	}
}

//Create a book
func (s *manager) Add(user *user.User, amount float64) (*Wallet, error) {
	return s.updateAmount(user, amount)
}

//Get a book
func (s *manager) Retrieve(user *user.User, amount float64) (*Wallet, error) {
	minusAmount := - amount
	return s.updateAmount(user, minusAmount)
}

func (s *manager) updateAmount(user *user.User, amount float64) (*Wallet, error) {
	return s.repo.Update(user, amount)
}
