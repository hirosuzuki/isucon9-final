package main

import (
	"fmt"
	"time"
)

func checkAvailableDate(date time.Time) bool {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	t := time.Date(2020, 1, 1, 0, 0, 0, 0, jst)
	t = t.AddDate(0, 0, availableDays)

	return date.Before(t)
}

func getUsableTrainClassList(fromStation Station, toStation Station) []string {
	usable := map[string]string{}

	for key, value := range TrainClassMap {
		usable[key] = value
	}

	if !fromStation.IsStopExpress {
		delete(usable, "express")
	}
	if !fromStation.IsStopSemiExpress {
		delete(usable, "semi_express")
	}
	if !fromStation.IsStopLocal {
		delete(usable, "local")
	}

	if !toStation.IsStopExpress {
		delete(usable, "express")
	}
	if !toStation.IsStopSemiExpress {
		delete(usable, "semi_express")
	}
	if !toStation.IsStopLocal {
		delete(usable, "local")
	}

	ret := []string{}
	for _, v := range usable {
		ret = append(ret, v)
	}

	return ret
}

type ssb struct {
	trainClass    string
	seatClass     string
	isSmokingSeat bool
}

var seatSsbMap map[ssb][]Seat

func initSeatSsbMap() {
	seatSsbMap = map[ssb][]Seat{}

	query := "SELECT * FROM seat_master"

	seatList := []Seat{}
	dbx.Select(&seatList, query)

	for _, seat := range seatList {
		k := ssb{seat.TrainClass, seat.SeatClass, seat.IsSmokingSeat}
		if _, ok := seatSsbMap[k]; !ok {
			seatSsbMap[k] = []Seat{}
		}
		seatSsbMap[k] = append(seatSsbMap[k], seat)
	}
}

func (train Train) getAvailableSeats(seatReservationList []SeatReservation, fromStation Station, toStation Station, seatClass string, isSmokingSeat bool) ([]Seat, error) {
	// 指定種別の空き座席を返す

	k := ssb{train.TrainClass, seatClass, isSmokingSeat}
	seatList := seatSsbMap[k]

	availableSeatMap := map[string]Seat{}
	for _, seat := range seatList {
		availableSeatMap[fmt.Sprintf("%d_%d_%s", seat.CarNumber, seat.SeatRow, seat.SeatColumn)] = seat
	}

	for _, seatReservation := range seatReservationList {
		key := fmt.Sprintf("%d_%d_%s", seatReservation.CarNumber, seatReservation.SeatRow, seatReservation.SeatColumn)
		delete(availableSeatMap, key)
	}

	ret := []Seat{}
	for _, seat := range availableSeatMap {
		ret = append(ret, seat)
	}
	return ret, nil
}

func (train Train) getAvailableSeats4(fromStation Station, toStation Station) (premium_avail_seats []Seat, premium_smoke_avail_seats []Seat, reserved_avail_seats []Seat, reserved_smoke_avail_seats []Seat, err error) {

	// すでに取られている予約を取得する
	query := `
		SELECT sr.reservation_id, sr.car_number, sr.seat_row, sr.seat_column
		FROM seat_reservations sr, reservations r, seat_master s, station_master std, station_master sta
		WHERE
			r.reservation_id=sr.reservation_id AND
			s.train_class=r.train_class AND
			s.car_number=sr.car_number AND
			s.seat_column=sr.seat_column AND
			s.seat_row=sr.seat_row AND
			std.name=r.departure AND
			sta.name=r.arrival
		`

	if train.IsNobori {
		query += "AND ((sta.id < ? AND ? <= std.id) OR (sta.id < ? AND ? <= std.id) OR (? < sta.id AND std.id < ?))"
	} else {
		query += "AND ((std.id <= ? AND ? < sta.id) OR (std.id <= ? AND ? < sta.id) OR (sta.id < ? AND ? < std.id))"
	}

	seatReservationList := []SeatReservation{}
	err = dbx.Select(&seatReservationList, query, fromStation.ID, fromStation.ID, toStation.ID, toStation.ID, fromStation.ID, toStation.ID)
	if err != nil {
		return
	}

	premium_avail_seats, err = train.getAvailableSeats(seatReservationList, fromStation, toStation, "premium", false)
	if err != nil {
		return
	}

	premium_smoke_avail_seats, err = train.getAvailableSeats(seatReservationList, fromStation, toStation, "premium", true)
	if err != nil {
		return
	}

	reserved_avail_seats, err = train.getAvailableSeats(seatReservationList, fromStation, toStation, "reserved", false)
	if err != nil {
		return
	}

	reserved_smoke_avail_seats, err = train.getAvailableSeats(seatReservationList, fromStation, toStation, "reserved", true)
	if err != nil {
		return
	}

	return
}
