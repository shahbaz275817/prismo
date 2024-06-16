package transaction

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shahbaz275817/prismo/internal/config"
	"github.com/shahbaz275817/prismo/internal/models"
	"github.com/shahbaz275817/prismo/internal/responder"
	mocks2 "github.com/shahbaz275817/prismo/internal/services/account/mocks"
	mocks3 "github.com/shahbaz275817/prismo/internal/services/operationtype/mocks"
	"github.com/shahbaz275817/prismo/internal/services/transaction/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransactionHandler(t *testing.T) {
	txnSvc := mocks.MockTransactionService{}
	accSvc := mocks2.MockAccountService{}
	optSvc := mocks3.MockOperationtypeService{}

	type response struct {
		body       responder.Response
		statusCode int
	}

	tests := []struct {
		name     string
		request  io.Reader
		mockFunc func()
		want     response
	}{
		{
			name:     "Invalid Request Body",
			request:  getInvalidCreateTxnRequest(),
			mockFunc: func() {},
			want: response{
				body: responder.Response{
					Success: false,
					Data:    nil,
					Errors: []responder.ErrorItem{
						{
							Message:      "Something went wrong while parsing request",
							MessageTitle: "Bad Request",
							Code:         "BAD_REQUEST",
						},
					},
				},
				statusCode: 400,
			},
		},
		{
			name:     "Validation Fails",
			request:  getValidationFailCreateTxnRequest(),
			mockFunc: func() {},
			want: response{
				body: responder.Response{
					Success: false,
					Data:    nil,
					Errors: []responder.ErrorItem{
						{
							Message:      "invalid account_id: must be a positive integer",
							MessageTitle: "Bad Request",
							Code:         "BAD_REQUEST",
						},
					},
				},
				statusCode: 400,
			},
		},
		{
			name:    "Fetch Account Returned Error",
			request: getValidCreateTxnRequest(),
			mockFunc: func() {
				accSvc.On("Get", mock.Anything, &models.Account{AccountID: int64(1)}).Return(nil, errors.New("some error")).Once()
			},
			want: response{
				body: responder.Response{
					Success: false,
					Data:    nil,
					Errors: []responder.ErrorItem{
						{
							Message:      "INTERNAL_SERVER_ERROR",
							MessageTitle: "Internal Server Error",
							Code:         "INTERNAL_SERVER_ERROR",
						},
					},
				},
				statusCode: 500,
			},
		},
		{
			name:    "Account not found",
			request: getValidCreateTxnRequest(),
			mockFunc: func() {
				accSvc.On("Get", mock.Anything, &models.Account{AccountID: int64(1)}).Return(&models.Account{}, nil).Once()
				optSvc.On("Get", mock.Anything, &models.OperationsType{OperationTypeID: int64(1)}).Return(nil, nil).Once()
			},
			want: response{
				body: responder.Response{
					Success: false,
					Data:    nil,
					Errors: []responder.ErrorItem{
						{
							Message:      "operation type not found",
							MessageTitle: "Bad Request",
							Code:         "BAD_REQUEST",
						},
					},
				},
				statusCode: 400,
			},
		},
		{
			name:    "Return Success",
			request: getValidCreateTxnRequest(),
			mockFunc: func() {
				accSvc.On("Get", mock.Anything, &models.Account{AccountID: int64(1)}).Return(&models.Account{}, nil).Once()
				optSvc.On("Get", mock.Anything, &models.OperationsType{OperationTypeID: int64(1)}).Return(&models.OperationsType{}, nil).Once()
				txnSvc.On("Transact", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(context.Context) error)
					err := fn(context.Background())
					assert.NoError(t, err)
				}).Once()
				txnSvc.On("Create", mock.Anything, mock.Anything).Return(nil, nil).Once()
			},
			want: response{
				body: responder.Response{
					Success: true,
					Errors:  []responder.ErrorItem{},
				},
				statusCode: 201,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Load()
			w := httptest.NewRecorder()
			tt.mockFunc()
			r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/transactions"), tt.request)
			assert.NoError(t, err)

			handler := CreateTransactionHandler(&txnSvc, &optSvc, &accSvc)
			handler(w, r)

			assert.Equal(t, tt.want.statusCode, w.Code)
			if w.Code >= 400 {
				assert.Equal(t, tt.want.body.Success, false)
				assert.Equal(t, tt.want.body.Data, nil)
				b, err := json.Marshal(tt.want.body)
				assert.NoError(t, err)
				assert.Equal(t, string(b)+"\n", w.Body.String())
			}
			txnSvc.AssertExpectations(t)
			accSvc.AssertExpectations(t)
			optSvc.AssertExpectations(t)
		})
	}
}

func getValidCreateTxnRequest() io.Reader {
	bodyMap := map[string]interface{}{
		"account_id":        1,
		"operation_type_id": 1,
		"amount":            123.1,
	}
	b, _ := json.Marshal(bodyMap)
	return strings.NewReader(string(b))
}

func getValidationFailCreateTxnRequest() io.Reader {
	bodyMap := map[string]interface{}{
		"operation_type_id": 1,
		"amount":            123.1,
	}
	b, _ := json.Marshal(bodyMap)
	return strings.NewReader(string(b))
}

func getInvalidCreateTxnRequest() io.Reader {
	bodyMap := map[string]interface{}{
		"account_id":        "12",
		"operation_type_id": 1,
		"amount":            123.1,
	}
	b, _ := json.Marshal(bodyMap)
	return strings.NewReader(string(b))
}
