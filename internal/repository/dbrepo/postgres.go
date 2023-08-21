package dbrepo

import (
	"context"
	"github.com/emiljannesson/bed-and-breakfast/internal/models"
	"time"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//"returning" is needed as we're using postgreSQL and want the ID of the reservation
	statement := `insert into reservations (
	                    first_name,
	                    last_name,
	                    email,
	                    phone,
	                    start_date,
	                    end_date,
	                    room_id,
	                    created_at,
	                    updated_at
	                    ) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	                    returning id`

	// We could do this if we used sqlx instead with NamedExec functions
	//statement := `insert into reservations (
	//                      first_name,
	//                      last_name,
	//                      email,
	//                      phone,
	//                      start_date,
	//                      end_date,
	//                      room_id,
	//                      created_at,
	//                      updated_at
	//                      ) values (:first_name, :last_name, :email, :phone, :start_date, :end_date, :room_id, :created_at, :updated_at)
	//                      returning id`

	var newID int
	err := m.DB.QueryRowContext(ctx, statement,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	//err := m.DB.Query(ctx, statement, res).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// "returning" is needed as we're using postgreSQL and want the ID of the reservation
	statement := `insert into room_restrictions (
                          start_date, 
                          end_date, 
                          room_id, 
                          reservation_id,
                          restriction_id,
                          created_at, 
                          updated_at
                          ) values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, statement,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability for roomID exists, and false if no availability
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		select
			count(id)
		from
		    room_restrictions
		where
		    room_id = $1
		    and $2 < end_date and $3 > start_date;
		`

	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)

	var numRows int
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for a given date range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		select
			r.id, r.room_name
		from
		    rooms r
		where r.id not in 
		      (select room_id from room_restrictions rr where $1 < rr.end_date and $2 > start_date);
		`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return nil, err
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

// GetRoomByID gets a room by ID
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		select
			id, room_name, created_at, updated_at
		from
		    rooms r
		where r.id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	var room models.Room
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt)

	if err != nil {
		return room, err
	}

	return room, nil
}
