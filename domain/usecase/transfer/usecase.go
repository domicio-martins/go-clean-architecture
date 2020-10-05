package transfer

import (
	"github.com/domicio-martins/go-clean-architeture/domain/user"
	"github.com/domicio-martins/go-clean-architeture/domain/transaction"
	"github.com/domicio-martins/go-clean-architeture/domain/wallet"
)

type usecase struct {
	userManager user.Manager,
	walletManager wallet.Manager
}

func New(user user.manager, wallet wallet.Manager) *usecase {
	return &usecase{
		userManager: user,
		walletManager: wallet,
	}
}

func (transfer *usecase) Credit(payer , payee *user.User, wallet *wallet.Manager) error {
	//TODO fazer a logica da transferencia
}