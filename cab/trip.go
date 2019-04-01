package cab

import (
	"errors"
	"time"
)

// Trip represents a NY cab trip
type Trip struct {
	Medallion           string
	HackLicense         string
	VendorID            string
	RateCode            string
	StoreAndForwardFlag string
	PickupDatetime      time.Time
	DropoffDatetime     time.Time
	PassengerCount      int
	TripTimeInSecs      int
	TripDistance        float64
	PickupLongtitude    float64
	PickupLatitude      float64
	DropoffLongitude    float64
	DropoffLatitude     float64
}

// TripCount is the number of trips a cab has made
type TripCount struct {
	Medallion string
	Rides     int
}

// TripQueryFields are the fields required to fetch information about a trip
type TripQueryFields struct {
	Medallions []string
	TripDate   time.Time
}

// TripRepository provides access to trip store
type TripRepository interface {
	GetTripCount(tripQueryFields *TripQueryFields) ([]TripCount, error)
}

const (
	noMedallios string = "no medallions found in the query"
	invalidDate string = "invalid date"
)

// NewTripFields creates a new TripFields
func NewTripFields(medallions []string, tripDate string) (*TripQueryFields, error) {
	if !hasMedallions(medallions) {
		return nil, errors.New(noMedallios)
	}

	date, err := getTripDate(tripDate)
	if err != nil {
		return nil, errors.New(invalidDate)
	}

	fields := &TripQueryFields{
		Medallions: unique(medallions),
		TripDate:   date,
	}

	return fields, nil
}

func unique(medallions []string) []string {
	medallionIDs := make(map[string]bool)
	uniqueMedallions := []string{}
	for _, medallion := range medallions {
		if _, value := medallionIDs[medallion]; !value {
			medallionIDs[medallion] = true
			uniqueMedallions = append(uniqueMedallions, medallion)
		}
	}
	return uniqueMedallions
}

func getTripDate(tripDate string) (time.Time, error) {
	t, err := time.Parse("02/01/2006", tripDate)
	return t, err
}

func hasMedallions(medallions []string) bool {
	return len(medallions) > 0
}

// Validate method checks the state of TripQueryFields
func (fields *TripQueryFields) Validate() error {
	if !hasMedallions(fields.Medallions) {
		return errors.New(noMedallios)
	}
	return nil
}
