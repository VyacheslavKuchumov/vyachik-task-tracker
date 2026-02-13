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

func (s *Store) UpdateUserProfile(userID int, payload types.UpdateProfilePayload) (*types.User, error) {
	row := s.db.QueryRow(
		`UPDATE users
		 SET first_name = $1, last_name = $2
		 WHERE id = $3
		 RETURNING id, first_name, last_name, email, password, created_at`,
		payload.FirstName,
		payload.LastName,
		userID,
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

func (s *Store) UpdateUserPassword(userID int, hashedPassword string) error {
	result, err := s.db.Exec(
		`UPDATE users
		 SET password = $1
		 WHERE id = $2`,
		hashedPassword,
		userID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (s *Store) ListUsers() ([]*types.UserLookup, error) {
	rows, err := s.db.Query(
		`SELECT id, TRIM(CONCAT(first_name, ' ', last_name)) AS full_name
		 FROM users
		 ORDER BY first_name, last_name, id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*types.UserLookup, 0)
	for rows.Next() {
		user := new(types.UserLookup)
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
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
