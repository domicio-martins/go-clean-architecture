package transfer

import "github.com/domicio-martins/go-clean-architeture/domain/user"

type UseCase interface {
	Credit(payer *user.User, payee *user.User, value float64)
}
