package account

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/shahbaz275817/prismo/constants/errcodes"
	"github.com/shahbaz275817/prismo/internal/models"
	"github.com/shahbaz275817/prismo/internal/responder"
	"github.com/shahbaz275817/prismo/internal/services/account"
	"github.com/shahbaz275817/prismo/internal/utils"
	"github.com/shahbaz275817/prismo/internal/wrappers"
	"github.com/shahbaz275817/prismo/pkg/errors"
)

func GetAccountHandler(accountService account.Service) http.HandlerFunc {
	return wrappers.DefaultWrapper(func(w http.ResponseWriter, r *http.Request) error {
		ctx, lgr := utils.ContextLogger(r)

		accountIDStr := mux.Vars(r)["account_id"]

		accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
		if err != nil {
			lgr.Errorf("invalid account_id: %s, error: %s", accountIDStr, err.Error())
			responder.WriteError(w, r, errors.NewBadRequestError(errcodes.BadRequest, &errors.ErrDetails{
				Message: "Invalid account ID",
			}))
			return err
		}

		accountData, err := accountService.Get(ctx, &models.Account{AccountID: accountID})
		if err != nil {
			lgr.Errorf("error in fetching account id error: %s", err.Error())
			responder.WriteError(w, r, errors.NewInternalServerError(errcodes.InternalServerError, &errors.ErrDetails{}))
			return err
		}

		if accountData == nil {
			lgr.Errorf("account not found with account id")
			responder.WriteError(w, r, errors.NewNotFoundError("account_not_found", &errors.ErrDetails{}))
			return errors.Errorf("account not found")
		}

		res := getAccountResponse{
			AccountID:      accountData.AccountID,
			DocumentNumber: accountData.DocumentNumber,
		}
		responder.WriteAnyResponse(ctx, w, res)
		return nil
	})
}

type getAccountResponse struct {
	AccountID      int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}
