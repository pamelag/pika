package insight

import (
	"testing"
	"time"

	"github.com/pamelag/pika/cab"
	"github.com/stretchr/testify/assert"
)

var (
	queryCounter int
	trips        cab.TripRepository
	is           Service
)

func TestGetTripCount(t *testing.T) {
	trips = NewMockRepository()
	is = NewService(trips)

	is.GetTripCount([]string{"ABC123"}, "12/02/2019", false)
	is.GetTripCount([]string{"ABC123"}, "12/02/2019", false)
	assert.Equal(t, 1, queryCounter, "they should be equal")

	queryCounter = 0
	is.GetTripCount([]string{"ABC123"}, "12/02/2019", true)
	assert.Equal(t, 1, queryCounter, "they should be equal")

	is.GetTripCount([]string{"ABC123"}, "12/02/2019", false)
	is.GetTripCount([]string{"DEF456"}, "12/02/2019", false)
	is.GetTripCount([]string{"789GHI"}, "12/02/2019", false)

	queryCounter = 0
	tripCount, err := is.GetTripCount([]string{"ABC123", "DEF456", "789GHI"}, "12/02/2019", false)
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, 0, queryCounter, "they should be equal")

	assert.Equal(t, 5, tripCount[0].Count, "they should be equal")
	assert.Equal(t, 1, tripCount[1].Count, "they should be equal")
	assert.Equal(t, 2, tripCount[2].Count, "they should be equal")

}

type mockRepository struct {
	data map[mockTrip]int
}

type mockTrip struct {
	medallion string
	date      time.Time
}

func NewMockRepository() cab.TripRepository {
	mockData := make(map[mockTrip]int)
	t, _ := time.Parse("02/01/2006", "12/02/2019")
	row1 := mockTrip{medallion: "ABC123", date: t}
	row2 := mockTrip{medallion: "DEF456", date: t}
	row3 := mockTrip{medallion: "789GHI", date: t}
	mockData[row1] = 5
	mockData[row2] = 1
	mockData[row3] = 2
	return &mockRepository{
		data: mockData,
	}
}

func (m *mockRepository) GetTripCount(tripQueryFields *cab.TripQueryFields) ([]cab.TripCount, error) {
	tripData := make([]cab.TripCount, 0)
	for _, v := range tripQueryFields.Medallions {
		queryCounter++
		row := mockTrip{medallion: v, date: tripQueryFields.TripDate}
		count, _ := m.data[row]
		tripData = append(tripData, cab.TripCount{Medallion: v, Rides: count})
	}
	return tripData, nil
}
