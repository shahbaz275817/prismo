package utils

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/shahbaz275817/prismo/pkg/errors"
	"github.com/shahbaz275817/prismo/pkg/logger"
)

type LatLong struct {
	Latitude  float64
	Longitude float64
}

const (
	latLongTolerance  = 0.00003
	floatPrecision90  = 90
	floatPrecision180 = 180
)

// FromCSV expects comma(,) separated latitude and longitude, example: 84.24655, 90.24255
func (ll *LatLong) FromCSV(csv string) error {
	slice := strings.Split(csv, ",")
	if len(slice) != 2 {
		logger.WithContext(context.Background()).Warnf("Expected csv in 'lat,long' format, but found %s", csv)
		return errors.NewValidationError("invalid_location_format", nil)
	}

	return ll.FromStrings(slice[0], slice[1])
}

// FromStrings expects two string parameters latitude and longitude.
func (ll *LatLong) FromStrings(latStr string, longStr string) error {
	lat, err := ToFloat64(strings.TrimSpace(latStr))
	if err != nil {
		logger.WithContext(context.Background()).Warnf("Latitude %s is not a floating number", latStr)
		return errors.NewValidationError("invalid_location_format", nil)
	}
	lng, err := ToFloat64(strings.TrimSpace(longStr))
	if err != nil {
		logger.WithContext(context.Background()).Warnf("longitude %s is not a floating number", longStr)
		return errors.NewValidationError("invalid_location_format", nil)
	}
	if !valid(lat, lng) {
		return errors.NewValidationError("invalid_location_format", nil)
	}

	ll.Latitude = lat
	ll.Longitude = lng
	return nil
}

func (ll *LatLong) Equals(another *LatLong) bool {
	return math.Abs(ll.Latitude-another.Latitude) < latLongTolerance && math.Abs(ll.Longitude-another.Longitude) < latLongTolerance
}

func (ll *LatLong) String() string {
	return fmt.Sprintf("%f,%f", ll.Latitude, ll.Longitude)
}

func (ll *LatLong) Valid() bool {
	return valid(ll.Latitude, ll.Longitude)
}

func (ll *LatLong) Value() (driver.Value, error) {
	return ll.String(), nil
}

func (ll *LatLong) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to string for lat long failed")
	}
	err := ll.FromCSV(string(b))
	if err != nil {
		return err
	}
	return nil
}

func (ll *LatLong) UnmarshalJSON(b []byte) error {
	unquoted, err := strconv.Unquote(string(b))
	if err != nil {
		return nil
	}
	return ll.Scan([]byte(unquoted))
}

func (ll LatLong) MarshalJSON() (s []byte, err error) {
	sQuoted := strconv.Quote(ll.String())
	return []byte(sQuoted), nil
}

func valid(lat float64, lng float64) bool {
	latDiff := math.Abs(float64(90) - math.Abs(lat))
	lngDiff := math.Abs(float64(180) - math.Abs(lng))
	if latDiff < latLongTolerance || math.Abs(lat) > float64(floatPrecision90) {
		logger.WithContext(context.Background()).Warnf("Invalid latitude %f", lat)
		return false
	}
	if lngDiff < latLongTolerance || math.Abs(lng) > float64(floatPrecision180) {
		logger.WithContext(context.Background()).Warnf("Invalid longitude %f", lng)
		return false
	}
	return true
}
