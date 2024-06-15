package transaction

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/shahbaz275817/prismo/constants/errcodes"
	"github.com/shahbaz275817/prismo/internal/models"
	"github.com/shahbaz275817/prismo/internal/responder"
	"github.com/shahbaz275817/prismo/internal/services/account"
	"github.com/shahbaz275817/prismo/internal/services/operationtype"
	"github.com/shahbaz275817/prismo/internal/services/transaction"
	"github.com/shahbaz275817/prismo/internal/utils"
	"github.com/shahbaz275817/prismo/internal/wrappers"
	"github.com/shahbaz275817/prismo/pkg/errors"
	"github.com/shahbaz275817/prismo/pkg/logger"
)

func CreateTransactionHandler(txnService transaction.Service, ots operationtype.Service, as account.Service) http.HandlerFunc {
	return wrappers.DefaultWrapper(func(w http.ResponseWriter, r *http.Request) error {

		ctx, lgr := utils.ContextLogger(r)

		var ctReq createTransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&ctReq); err != nil {
			lgr.Errorf("error while decoding create_transaction_request body: %s", err.Error())
			err = errors.NewBadRequestError(errcodes.BadRequest, &errors.ErrDetails{
				Message: "Something went wrong while parsing request",
			})
			responder.WriteError(w, r, err)
			return err
		}

		err := validateCreateTransactionRequest(ctReq)
		if err != nil {
			lgr.Errorf("invalid create transaction request error: %s", err.Error())
			responder.WriteError(w, r, errors.NewBadRequestError(errcodes.BadRequest, &errors.ErrDetails{
				Message: err.Error(),
			}))
			return err
		}

		acc, err := as.Get(ctx, &models.Account{AccountID: ctReq.AccountID})
		if err != nil {
			lgr.Errorf("error in fetching account id error: %s", err.Error())
			responder.WriteError(w, r, errors.NewInternalServerError(errcodes.InternalServerError, &errors.ErrDetails{}))
			return err
		}
		if acc == nil {
			lgr.Errorf("account not found for account id %d", ctReq.AccountID)
			responder.WriteError(w, r, errors.NewBadRequestError(errcodes.BadRequest, &errors.ErrDetails{
				Message: "account not found",
			}))
			return err
		}

		ot, err := ots.Get(ctx, &models.OperationsType{OperationTypeID: ctReq.OperationTypeID})
		if err != nil {
			lgr.Errorf("error in fetching operation type error: %s", err.Error())
			responder.WriteError(w, r, errors.NewInternalServerError(errcodes.InternalServerError, &errors.ErrDetails{}))
			return err
		}
		if ot == nil {
			lgr.Errorf("operation type not found for operation type id %d", ctReq.OperationTypeID)
			responder.WriteError(w, r, errors.NewBadRequestError(errcodes.BadRequest, &errors.ErrDetails{
				Message: "operation type not found",
			}))
			return err
		}

		err = txnService.Transact(ctx, func(ctx context.Context) error {
			_, err = txnService.Create(ctx, models.Transaction{
				AccountID:       ctReq.AccountID,
				OperationTypeID: ctReq.OperationTypeID,
				Amount:          ctReq.Amount,
				EventDate:       time.Now().UTC(),
			})
			return err
		})
		if err != nil {
			logger.Errorf("error in creating transaction error: %s", err.Error())
			responder.WriteError(w, r, errors.NewInternalServerError(errcodes.InternalServerError, &errors.ErrDetails{}))
			return err
		}
		responder.WriteAnyResponse(ctx, w, map[string]string{
			"message": "success",
		}, http.StatusCreated)
		return nil

	})
}

func validateCreateTransactionRequest(req createTransactionRequest) error {
	if req.AccountID <= 0 {
		return errors.New("invalid account_id: must be a positive integer")
	}
	if req.OperationTypeID <= 0 {
		return errors.New("invalid operation_type: must be a positive integer")
	}
	if req.Amount <= 0 {
		return errors.New("invalid amount: must be a positive value")
	}
	return nil
}

type createTransactionRequest struct {
	AccountID       int64   `json:"account_id"`
	OperationTypeID int64   `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
}
