package utils

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultOffset = 0
	defaultLimit  = 10
	perPageParam  = "per_page"
	pageParam     = "page"

	defaultSortParam           = "created_at DESC"
	defaultSellerViewSortParam = "seller_group DESC, rpu DESC, allocated DESC, name ASC"

	multiParamSeparator = ","
	API                 = "api"

	createdAtFromParam = "created_at_from"
	createdAtToParam   = "created_at_to"

	rpuAtFromParam = "rpu_at_from"
	rpuAtToParam   = "rpu_at_to"

	inboundedAtFromParam = "inbounded_at_from"
	inboundedAtToParam   = "inbounded_at_to"
)

type URLParser struct {
	v url.Values
}

func NewURLParser(url *url.URL) *URLParser {
	return &URLParser{url.Query()}
}

// Get returns empty string if key not found
func (up *URLParser) Get(key string) string {
	return up.v.Get(key)
}

// GetList returns nil if key not found
func (up *URLParser) GetList(key string) []string {
	var values []string
	raw := up.Get(key)
	if raw != "" {
		values = strings.Split(raw, multiParamSeparator)
	}
	return values
}

// GetBool returns false if given key doesn't have value 'true'
func (up *URLParser) GetBool(key string) bool {
	return up.Get(key) == "true"
}

// GetPaginationParams returns limit (50 records), offset (starting row from / example 500th record)
func (up *URLParser) GetPaginationParams() (limit int, offset int, err error) {
	limit = defaultLimit
	offset = defaultOffset

	perPageVal := up.v.Get(perPageParam)
	if perPageVal != "" {
		perPageInt, err := strconv.Atoi(perPageVal)
		if err != nil {
			return 0, 0, errors.New("per_page should be a number")
		}
		if perPageInt < 1 {
			return 0, 0, errors.New("per_page should be greater than 0")
		}
		limit = perPageInt
	}

	pageVal := up.v.Get(pageParam)
	if pageVal != "" {
		pageInt, err := strconv.Atoi(pageVal)
		if err != nil {
			return 0, 0, errors.New("page should be a number")
		}
		pageInt = pageInt - 1 // todo check why this is being done? not readable
		if pageInt < 0 {
			return 0, 0, errors.New("page should be greater than 0")
		}
		offset = pageInt * limit
	}

	return limit, offset, nil
}

func (up *URLParser) GetCreatedRangeParams() (*time.Time, *time.Time, error) {
	createdFromVal := up.Get(createdAtFromParam)
	createdToVal := up.Get(createdAtToParam)

	return ToCreatedRange(createdFromVal, createdToVal)
}

func (up *URLParser) GetRpuAtRangeParams() (*time.Time, *time.Time, error) {
	var rpuAtFrom *time.Time
	var rpuAtTo *time.Time

	rpuFromVal := up.Get(rpuAtFromParam)
	rpuToVal := up.Get(rpuAtToParam)

	if rpuFromVal != "" {
		rpuFromEpoch, err := strconv.ParseInt(rpuFromVal, 10, 64)
		if err != nil {
			return nil, nil, errors.New("rpu_at_from should be unix epoch")
		}
		if rpuToVal == "" {
			return nil, nil, errors.New("rpu_at_to is missing")
		}
		startTime := time.Unix(rpuFromEpoch, 0).UTC()
		rpuAtFrom = &startTime
	}

	if rpuToVal != "" {
		endRpuEpoch, err := strconv.ParseInt(rpuToVal, 10, 64)
		if err != nil {
			return nil, nil, errors.New("rpu_at_to should be unix epoch")
		}
		if rpuFromVal == "" {
			return nil, nil, errors.New("rpu_at_from is missing")
		}
		endTime := time.Unix(endRpuEpoch, 0).UTC()
		rpuAtTo = &endTime
	}

	return rpuAtFrom, rpuAtTo, nil
}

func ToCreatedRange(createdFromVal string, createdToVal string) (*time.Time, *time.Time, error) {
	var createdFrom *time.Time
	var createdTo *time.Time

	if createdFromVal != "" {
		createdFromEpoch, err := strconv.ParseInt(createdFromVal, 10, 64)
		if err != nil {
			return nil, nil, errors.New("created_at_from should be unix epoch")
		}
		if createdToVal == "" {
			return nil, nil, errors.New("created_at_to is missing")
		}
		startTime := time.Unix(createdFromEpoch, 0).UTC()
		createdFrom = &startTime
	}

	if createdToVal != "" {
		endCreatedEpoch, err := strconv.ParseInt(createdToVal, 10, 64)
		if err != nil {
			return nil, nil, errors.New("created_at_to should be unix epoch")
		}
		if createdFromVal == "" {
			return nil, nil, errors.New("created_at_from is missing")
		}
		endTime := time.Unix(endCreatedEpoch, 0).UTC()
		createdTo = &endTime
	}

	return createdFrom, createdTo, nil
}

func (up *URLParser) GetInboundedRangeParams() (*time.Time, *time.Time, error) {
	var inboundedFrom *time.Time
	var inboundedTo *time.Time

	inboundedFromVal := up.Get(inboundedAtFromParam)
	inboundedToVal := up.Get(inboundedAtToParam)

	if inboundedFromVal != "" {
		inboundedFromEpoch, err := strconv.ParseInt(inboundedFromVal, 10, 64)
		if err != nil {
			return nil, nil, errors.New("inbounded_at_from should be unix epoch")
		}
		if inboundedToVal == "" {
			return nil, nil, errors.New("inbounded_at_to is missing")
		}
		startTime := time.Unix(inboundedFromEpoch, 0).UTC()
		inboundedFrom = &startTime
	}

	if inboundedToVal != "" {
		endInboundedEpoch, err := strconv.ParseInt(inboundedToVal, 10, 64)
		if err != nil {
			return nil, nil, errors.New("inbounded_at_to should be unix epoch")
		}
		if inboundedFromVal == "" {
			return nil, nil, errors.New("inbounded_at_from is missing")
		}
		endTime := time.Unix(endInboundedEpoch, 0).UTC()
		inboundedTo = &endTime
	}

	return inboundedFrom, inboundedTo, nil
}

func (up *URLParser) GetTimeRangeParams(fromParam, toParam string) (from *time.Time, to *time.Time, err error) {
	fromVal := up.Get(fromParam)
	if fromVal == "" {
		return nil, nil, fmt.Errorf("%s is missing", fromParam)
	}

	toVal := up.Get(toParam)
	if toVal == "" {
		return nil, nil, fmt.Errorf("%s is missing", toParam)
	}

	inboundedFromEpoch, err := strconv.ParseInt(fromVal, 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("%s should be unix epoch", fromParam)
	}
	startTime := time.Unix(inboundedFromEpoch, 0).UTC()
	from = &startTime

	endInboundedEpoch, err := strconv.ParseInt(toVal, 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("%s should be unix epoch", toParam)
	}
	endTime := time.Unix(endInboundedEpoch, 0).UTC()
	to = &endTime

	return from, to, nil
}

func (up *URLParser) GetOptionalTimeRangeParams(fromParam, toParam string) (from *time.Time, to *time.Time, err error) {
	fromVal := up.Get(fromParam)
	if fromVal != "" {
		fromEpoch, err := strconv.ParseInt(fromVal, 10, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("%s should be unix epoch", fromParam)
		}
		fromTime := time.Unix(fromEpoch, 0).UTC()
		from = &fromTime
	}

	toVal := up.Get(toParam)
	if toVal != "" {
		toEpoch, err := strconv.ParseInt(toVal, 10, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("%s should be unix epoch", toParam)
		}
		toTime := time.Unix(toEpoch, 0).UTC()
		to = &toTime
	}

	if (from == nil && to != nil) || (from != nil && to == nil) {
		return nil, nil, fmt.Errorf("%s and %s need to be provided together", fromParam, toParam)
	}

	return from, to, nil
}

func ConvertStringToTime(timeVal string) (*time.Time, error) {
	if timeVal == "" {
		return nil, errors.New("time is missing")
	}
	timeValEpoch, err := strconv.ParseInt(timeVal, 10, 64)
	if err != nil {
		return nil, errors.New("timeValEpoch should be unix epoch")
	}
	timeUTC := time.Unix(timeValEpoch, 0).UTC()
	return &timeUTC, nil
}
