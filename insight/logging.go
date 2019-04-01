package insight

import (
	"log"
	"time"
)

type loggingService struct {
	next Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(s Service) Service {
	return &loggingService{s}
}

func (s *loggingService) GetTripCount(medallions []string, tripDate string, ignoreCache bool) (tripData []TripData, err error) {
	defer func(begin time.Time) {
		log.Println("method", "GetTripCount", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.GetTripCount(medallions, tripDate, ignoreCache)
}

func (s *loggingService) ClearCache() {
	defer func(begin time.Time) {
		log.Println("method", "ClearCache", "took", time.Since(begin))
	}(time.Now())
	s.next.ClearCache()
}
