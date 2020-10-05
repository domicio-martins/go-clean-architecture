package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/eminetto/clean-architecture-go-v2/domain"

	"github.com/eminetto/clean-architecture-go-v2/api/presenter"

	"github.com/eminetto/clean-architecture-go-v2/domain/entity"
	"github.com/eminetto/clean-architecture-go-v2/domain/entity/wallet"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func listWallets(manager wallet.Manager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading wallets"
		var data []*wallet.Wallet
		var err error
		title := r.URL.Query().Get("title")
		switch {
		case title == "":
			data, err = manager.List()
		default:
			data, err = manager.Search(title)
		}
		w.Header().Set("Content-Type", "application/json")
		if err != nil && err != domain.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}
		var toJ []*presenter.Wallet
		for _, d := range data {
			toJ = append(toJ, &presenter.Wallet{
				ID:       d.ID,
				Title:    d.Title,
				Author:   d.Author,
				Pages:    d.Pages,
				Quantity: d.Quantity,
			})
		}
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})
}

func createWallet(manager wallet.Manager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error adding wallet"
		var input struct {
			Title    string `json:"title"`
			Author   string `json:"author"`
			Pages    int    `json:"pages"`
			Quantity int    `json:"quantity"`
		}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		b := &wallet.Wallet{
			ID:        entity.NewID(),
			Title:     input.Title,
			Author:    input.Author,
			Pages:     input.Pages,
			Quantity:  input.Quantity,
			CreatedAt: time.Now(),
		}
		b.ID, err = manager.Create(b)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		toJ := &presenter.Wallet{
			ID:       b.ID,
			Title:    b.Title,
			Author:   b.Author,
			Pages:    b.Pages,
			Quantity: b.Quantity,
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

func getWallet(manager wallet.Manager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error reading wallet"
		vars := mux.Vars(r)
		id, err := entity.StringToID(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		data, err := manager.Get(id)
		if err != nil && err != domain.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errorMessage))
			return
		}
		toJ := &presenter.Wallet{
			ID:       data.ID,
			Title:    data.Title,
			Author:   data.Author,
			Pages:    data.Pages,
			Quantity: data.Quantity,
		}
		if err := json.NewEncoder(w).Encode(toJ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
		}
	})
}

func deleteWallet(manager wallet.Manager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorMessage := "Error removing walletmark"
		vars := mux.Vars(r)
		id, err := entity.StringToID(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
		err = manager.Delete(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errorMessage))
			return
		}
	})
}

//MakeWalletHandlers make url handlers
func MakeWalletHandlers(r *mux.Router, n negroni.Negroni, manager wallet.Manager) {
	r.Handle("/v1/wallet", n.With(
		negroni.Wrap(listWallets(manager)),
	)).Methods("GET", "OPTIONS").Name("listWallets")

	r.Handle("/v1/wallet", n.With(
		negroni.Wrap(createWallet(manager)),
	)).Methods("POST", "OPTIONS").Name("createWallet")

	r.Handle("/v1/wallet/{id}", n.With(
		negroni.Wrap(getWallet(manager)),
	)).Methods("GET", "OPTIONS").Name("getWallet")

	r.Handle("/v1/wallet/{id}", n.With(
		negroni.Wrap(deleteWallet(manager)),
	)).Methods("DELETE", "OPTIONS").Name("deleteWallet")
}
