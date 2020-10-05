package handler

import (
	"fmt"
	"net/http"

	"github.com/domicio-martins/go-clean-architecture/domain/entity/user"

	"github.com/domicio-martins/go-clean-architecture/domain/usecase/transfer"

	"github.com/domicio-martins/go-clean-architecture/domain"

	"github.com/domicio-martins/go-clean-architecture/domain/entity"
	"github.com/domicio-martins/go-clean-architecture/domain/entity/wallet"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func Credit(userManager user.Manager, transferUseCase transfer.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		payer, payee, amount, err := validateTransaction(userManager)

		err = transferUseCase.Credit(payer, payee, amount)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		w.WriteHeader(http.StatusCreated)
	})
}

func validateTransaction(userManager user.Manager) (payer, payee *user.Manager, amount *float64, err error) {
	errorMessage := "Error borrowing book"
	vars := mux.Vars(r)
	amount = vars["amount"]
	payerID, err := entity.StringToID(vars["payer"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}
	payer, err = userManager.Get(payerID)
	if err != nil && err != domain.ErrNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}
	if payer == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(errorMessage))
		return
	}
	payeeID, err := entity.StringToID(vars["payee"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}
	payee, err = userManager.Get(payeeID)
	if err != nil && err != domain.ErrNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errorMessage))
		return
	}
	if payee == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(errorMessage))
		return
	}

	return
}

//MakeLoanHandlers make url handlers
func ₢₢ß(r *mux.Router, n negroni.Negroni, userManager user.Manager, transferUseCase transfer.UseCase) {
	r.Handle("/v1/transfer", n.With(
		negroni.Wrap(Credit(userManager, transferUseCase)),
	)).
		Methods("POST", "OPTIONS").
		Name("Credit")
}