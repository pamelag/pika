package mysql

import (
	"bytes"
	"database/sql"
	"reflect"
	"time"

	"github.com/pamelag/pika/cab"
)

type tripRepository struct {
	db *sql.DB
}

// NewTripRepository creates a new instance of tripRepository
func NewTripRepository(db *sql.DB) cab.TripRepository {
	r := &tripRepository{
		db: db,
	}
	return r
}

func (t *tripRepository) GetTripCount(fields *cab.TripQueryFields) ([]cab.TripCount, error) {

	if err := fields.Validate(); err != nil {
		return nil, err
	}

	query := buildQuery(len(fields.Medallions))

	stmt, err := t.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	rows, err := execDynamicStmt(stmt, fields)
	if err != nil {
		return nil, err
	}

	var counter int
	cabTrips := make([]cab.TripCount, 0)

	defer rows.Close()
	for rows.Next() {
		var count int

		err = rows.Scan(&count)
		if err != nil {
			return nil, err
		}

		tripCount := cab.TripCount{
			Medallion: fields.Medallions[counter],
			Rides:     count,
		}

		cabTrips = append(cabTrips, tripCount)
		counter++

	}

	return cabTrips, nil

}

func getSQLDateFormat(tripDate time.Time) string {
	return tripDate.Format("2006-01-02")
}

func buildQuery(medallionCount int) string {
	var query bytes.Buffer
	var counter int

	for counter < medallionCount {
		if counter == medallionCount-1 {
			query.WriteString("SELECT count(*) FROM cab_trip_data where medallion = ? and DATE(pickup_datetime) = ?")
			break
		}
		query.WriteString("SELECT count(*) FROM cab_trip_data where medallion = ? and DATE(pickup_datetime) = ? UNION ALL ")
		counter++
	}

	return query.String()
}

func execDynamicStmt(stmt *sql.Stmt, fields *cab.TripQueryFields) (*sql.Rows, error) {
	relectstmt := reflect.ValueOf(stmt)
	values := make([]reflect.Value, 0)

	for i := 0; i < len(fields.Medallions); i++ {
		values = append(values, reflect.ValueOf(fields.Medallions[i]))
		values = append(values, reflect.ValueOf(getSQLDateFormat(fields.TripDate)))
	}
	reflectQuery := relectstmt.MethodByName("Query")
	reflectReturns := reflectQuery.Call(values)

	var rows *sql.Rows
	var err error
	if reflectReturns[0].Interface() != nil {
		rows = reflectReturns[0].Interface().(*sql.Rows)
	}
	if reflectReturns[1].Interface() != nil {
		err = reflectReturns[1].Interface().(error)
	}

	return rows, err
}
