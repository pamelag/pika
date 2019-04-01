package insight

import (
	"log"
	"time"

	"github.com/pamelag/pika/cab"
)

// Service is the interface that provides methods for an insight into cab trips
type Service interface {
	GetTripCount(medallions []string, tripDate string, ignoreCache bool) ([]TripData, error)
	ClearCache()
}

type service struct {
	trips cab.TripRepository
	cache *cache
}

// NewService creates a trip service with necessary dependencies.
func NewService(trips cab.TripRepository) Service {
	return &service{
		trips: trips,
		cache: newCache(),
	}
}

// GetTripCount service returns the trip count data
func (s *service) GetTripCount(medallions []string, tripDate string, ignoreCache bool) ([]TripData, error) {
	var tripData []TripData
	tripFields, err := cab.NewTripFields(medallions, tripDate)
	if err != nil {
		return tripData, err
	}

	tc, err := s.getTripCountData(tripFields, ignoreCache)
	if err != nil {
		return tripData, err
	}

	return createTripData(tc), nil
}

func (s *service) getTripCountData(tripFields *cab.TripQueryFields, ignoreCache bool) ([]cab.TripCount, error) {
	if ignoreCache {
		log.Println("cache ignored.. fetching fresh data from db")
		return s.getFreshData(tripFields)
	}

	ctc, foundIndices := s.getCachedData(tripFields)
	tripFields = trimQueryFields(foundIndices, tripFields)

	if len(tripFields.Medallions) == 0 {
		return ctc, nil
	}

	tc, err := s.getFreshData(tripFields)
	if err != nil {
		return nil, err
	}

	ctc = append(ctc, tc...)
	return ctc, nil
}

func createTripData(tc []cab.TripCount) []TripData {
	tripData := make([]TripData, 0)
	for _, t := range tc {
		tripData = append(tripData, assemble(t))
	}
	return tripData
}

func (s *service) getFreshData(tripFields *cab.TripQueryFields) ([]cab.TripCount, error) {
	tc, err := s.trips.GetTripCount(tripFields)
	if err != nil {
		return tc, err
	}
	s.addToCache(tc, tripFields.TripDate)
	return tc, nil
}

func (s *service) addToCache(tc []cab.TripCount, tripDate time.Time) {
	for _, v := range tc {
		s.cache.add(v.Medallion, tripDate, v.Rides)
	}
}

func (s *service) getCachedData(tripFields *cab.TripQueryFields) ([]cab.TripCount, map[int]string) {
	cacheIndices := make(map[int]string, 0)
	cabTrips := make([]cab.TripCount, 0)

	for i, v := range tripFields.Medallions {
		count := s.cache.get(v, tripFields.TripDate)
		if count > 0 {
			cacheIndices[i] = v
			tripCount := cab.TripCount{
				Medallion: v,
				Rides:     count,
			}
			cabTrips = append(cabTrips, tripCount)
		}
	}

	return cabTrips, cacheIndices
}

func trimQueryFields(foundIndices map[int]string, tripFields *cab.TripQueryFields) *cab.TripQueryFields {
	if len(foundIndices) == 0 {
		return tripFields
	}

	fields := &cab.TripQueryFields{Medallions: make([]string, 0),
		TripDate: tripFields.TripDate}
	for i, v := range tripFields.Medallions {
		if foundIndices[i] == v {
			continue
		}
		fields.Medallions = append(fields.Medallions, v)
	}
	return fields
}

func (s *service) ClearCache() {
	s.cache.clear()
	return
}

// TripData is the a read model for trip count
type TripData struct {
	Medallion string `json:"medallion"`
	Count     int    `json:"count"`
}

func assemble(t cab.TripCount) TripData {
	return TripData{
		Medallion: t.Medallion,
		Count:     t.Rides,
	}
}
