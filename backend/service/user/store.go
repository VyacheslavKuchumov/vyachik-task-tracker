package user

import (
	"VyacheslavKuchumov/test-backend/types"
	"database/sql"
	"fmt"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	row := s.db.QueryRow(
		"SELECT id, first_name, last_name, email, password, created_at FROM users WHERE email = $1",
		email,
	)
	u, err := scanRowIntoUser(row)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	row := s.db.QueryRow(
		"SELECT id, first_name, last_name, email, password, created_at FROM users WHERE id = $1",
		id,
	)
	u, err := scanRowIntoUser(row)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec(
		"INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4)",
		user.FirstName, user.LastName, user.Email, user.Password,
	)
	if err != nil {
		return err
	}

	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanRowIntoUser(row rowScanner) (*types.User, error) {
	user := new(types.User)

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
