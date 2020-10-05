package wallet

import (
	"github.com/domicio-martins/go-clean-architecture/domain/entity"
	"github.com/domicio-martins/go-clean-architecture/domain/user"
	)

//Monetary interface
type Monetary interface {
	Add(user user.ID, amount float64) (*Wallet, error)
	Retrieve(user user.ID, amount float64) (*Wallet, error)
}

//repository interface
type repository interface {
	Monetary
}

//Manager interface
type Manager interface {
	repository
}