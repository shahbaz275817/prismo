package account

import (
	"encoding/json"
	"net/http"

	"github.com/shahbaz275817/prismo/constants/errcodes"
	"github.com/shahbaz275817/prismo/internal/models"
	"github.com/shahbaz275817/prismo/internal/responder"
	"github.com/shahbaz275817/prismo/internal/services/account"
	"github.com/shahbaz275817/prismo/internal/utils"
	"github.com/shahbaz275817/prismo/internal/wrappers"
	"github.com/shahbaz275817/prismo/pkg/errors"
	"github.com/shahbaz275817/prismo/pkg/locks"
	"github.com/shahbaz275817/prismo/pkg/logger"
)

func CreateAccountHandler(accountService account.Service, lock *locks.AtomicLock) http.HandlerFunc {
	return wrappers.DefaultWrapper(func(w http.ResponseWriter, r *http.Request) error {

		ctx, lgr := utils.ContextLogger(r)

		var caReq createAccountRequest
		if err := json.NewDecoder(r.Body).Decode(&caReq); err != nil {
			lgr.Errorf("error while decoding create_account_request body: %s", err.Error())
			err = errors.NewBadRequestError(errcodes.BadRequest, &errors.ErrDetails{
				Message: "Something went wrong while parsing request",
			})
			responder.WriteError(w, r, err)
			return err
		}

		err := validateCreateAccountRequest(caReq)
		if err != nil {
			logger.Errorf("invalid create account request error: %s", err.Error())
			responder.WriteError(w, r, errors.NewBadRequestError(errcodes.BadRequest, &errors.ErrDetails{
				Message: err.Error(),
			}))
			return err
		}

		lockState := locks.LockState{}
		_, err = lock.Execute(ctx, utils.BuildLockKey("doc_number", caReq.DocumentNumber), locks.Def, func(lockState *locks.LockState) ([]interface{}, error) {
			return nil, accountService.Create(ctx, models.Account{DocumentNumber: caReq.DocumentNumber})
		}, &lockState)
		if err != nil {
			lgr.Errorf("Error while creating account : %s", err.Error())
			responder.WriteError(w, r, err)
			return err
		}

		responder.WriteAnyResponse(ctx, w, map[string]string{
			"message": "success",
		}, http.StatusCreated)
		return nil
	})
}

func validateCreateAccountRequest(req createAccountRequest) error {
	if req.DocumentNumber == "" {
		return errors.New("invalid document number")
	}

	return nil
}

type createAccountRequest struct {
	DocumentNumber string `json:"document_number"`
}
