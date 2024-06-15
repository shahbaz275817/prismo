package server

import (
	"net/http"

	"github.com/shahbaz275817/prismo/internal/services/account"
	"github.com/shahbaz275817/prismo/internal/services/operationtype"
	"github.com/shahbaz275817/prismo/internal/services/transaction"
	"github.com/shahbaz275817/prismo/pkg/locks"
)

type Dependencies struct {
	TransactionService    transaction.Service
	AccountService        account.Service
	OperationTypesService operationtype.Service
	AtomicLock            *locks.AtomicLock
}

type ExternalDependencies struct {
	HTTPClient *http.Client
}
