package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	contextWrapper "github.com/shahbaz275817/prismo/pkg/context"
	"github.com/shahbaz275817/prismo/pkg/logger"
)

func ContextLogger(r *http.Request) (context.Context, logger.LgrWrapper) {
	ctx := r.Context()
	endpoint := contextWrapper.GetAPIEndpoint(ctx)
	lgr := logger.WithFields(logger.Fields{API: endpoint}).WithContext(r.Context())
	return ctx, lgr
}

func ContextLoggerURLParser(r *http.Request) (context.Context, logger.LgrWrapper, *URLParser) {
	ctx, lgr := ContextLogger(r)
	up := NewURLParser(r.URL)
	return ctx, lgr, up
}

func GetRouteIndexFromParam(routeSeqArray []string) ([]int, error) {
	var routeIndex []int
	for i := 0; i < len(routeSeqArray); i++ {
		seq, err := strconv.Atoi(routeSeqArray[i])
		if err != nil {
			return nil, errors.New("route_index should be an array of numbers - " + err.Error())
		}
		routeIndex = append(routeIndex, seq)
	}
	return routeIndex, nil
}

func ValidateSortVerb(sort string, allowedValues []string) (string, error) {
	if sort == "" {
		return defaultSortParam, nil
	}
	splits := strings.Split(sort, ":")
	if len(splits) != 2 {
		return "", errors.New("sort should be in this format field.orderdirection")
	}
	field, order := splits[0], splits[1]
	if order != "desc" && order != "asc" {
		return "", errors.New("Order direction should be asc or desc")
	}

	if !StringInSlice(allowedValues, field) {
		return "", errors.New("Unknown field in sort query parameter")
	}

	return fmt.Sprintf("%s %s", field, strings.ToUpper(order)), nil
}

func ValidateSellerViewSortVerb(sort string, allowedValues []string) (string, error) {
	if sort == "" {
		return defaultSellerViewSortParam, nil
	}
	return ValidateSortVerb(sort, allowedValues)
}

func FromDaysToHours(days int) int {
	return 24 * days
}
