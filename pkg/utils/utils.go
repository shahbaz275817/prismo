package utils

import "time"

const (
	HeaderContentType    = "Content-Type"
	HeaderAccept         = "Accept"
	HeaderAcceptLanguage = "Accept-Language"
	ClientID             = "client-id"
	PassKey              = "pass-key"
	ApplicationJSON      = "application/json"
	AuthKey              = "auth-key"
	jakartaTimeZone      = "Asia/Jakarta"
	TextHTML             = "text/html"
	Authorization        = "Authorization"
)

func ContainsInt(haystack []int, needle int) bool {
	for _, b := range haystack {
		if b == needle {
			return true
		}
	}
	return false
}

func GetTimeBeginningOfTheDay(timezone *string) (*time.Time, error) {
	tz := "UTC"
	if timezone != nil {
		tz = *timezone
	}
	timeLocation, err := time.LoadLocation(tz) // load time zone check if correct timezone
	if err != nil {
		return nil, err
	}
	year, month, day := time.Now().Date()                                       // real time-server date
	bod := time.Date(year, month, day, 0, 0, 0, 0, time.Local).In(timeLocation) // real beginning day server time in UTC
	return &bod, nil
}

func getTimeBaseOnTimeZone(datetime *time.Time, timezone string) (*time.Time, error) {
	timeLocation, err := time.LoadLocation(timezone) // load time zone check if correct timezone
	if err != nil {
		return nil, err
	}

	dateTimeWithTimeZone := datetime.In(timeLocation)
	return &dateTimeWithTimeZone, nil
}

func GetJakartaTimeFor(t *time.Time) (*time.Time, error) {
	return getTimeBaseOnTimeZone(t, jakartaTimeZone)
}
