package main

import (
	"database/sql"
	"errors"
	"time"
)

type Parcel struct {
	ID        int
	Client    int
	Status    string
	Address   string
	CreatedAt time.Time
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) AddParcel(p Parcel) (int, error) {
	result, err := s.db.Exec("INSERT INTO parcel (client, status, address) VALUES (?, ?, ?)", p.Client, p.Status, p.Address)
	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastID), nil
}

func (s ParcelStore) GetParcel(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = ?", number)
	var id, client int
	var status, address string
	err := row.Scan(&id, &client, &status, &address)
	if err != nil {
		return Parcel{}, err
	}

	// заполните объект Parcel данными из таблицы
	p := Parcel{
		ID:        id,
		Client:    client,
		Status:    status,
		Address:   address,
		CreatedAt: time.Now().UTC(),
	}

	return p, nil
}

func (s ParcelStore) GetByClientParcel(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = ?", client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	for rows.Next() {
		var id, client int
		var status, address string

		err := rows.Scan(&id, &client, &status, &address)
		if err != nil {
			return nil, err
		}

		p := Parcel{
			ID:        id,
			Client:    client,
			Status:    status,
			Address:   address,
			CreatedAt: time.Now().UTC(),
		}
		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatusParcel(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = ? WHERE number = ?", status, number)
	return err
}

func (s ParcelStore) SetAddressParcel(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	var currentStatus string
	err := s.db.QueryRow("SELECT status FROM parcel WHERE number = ?", number).Scan(&currentStatus)
	if err != nil {
		return err
	}

	if currentStatus != "registered" {
		return errors.New("нельзя менять адрес для посылки со статусом " + currentStatus)
	}

	_, err = s.db.Exec("UPDATE parcel SET address = ? WHERE number = ?", address, number)
	return err
}

func (s ParcelStore) DeleteParcel(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	var currentStatus string
	err := s.db.QueryRow("SELECT status FROM parcel WHERE number = ?", number).Scan(&currentStatus)
	if err != nil {
		return err
	}

	if currentStatus != "registered" {
		return errors.New("нельзя удалять посылку со статусом " + currentStatus)
	}

	_, err = s.db.Exec("DELETE FROM parcel WHERE number = ?", number)
	return err
}
