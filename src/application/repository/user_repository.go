package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type IUserRepository interface {
	CreateUser(firstName string, lastName string) (*uuid.UUID, error)
	UpdateUser(id uuid.UUID, firstName string, lastName string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &userRepository{db}
}

func (u *userRepository) CreateUser(firstName string, lastName string) (*uuid.UUID, error) {

	newUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("échec de la génération UUID: %w", err)
	}
	userIDString := newUUID.String()

	query := `INSERT INTO main.user (id, first_name, last_name)
              VALUES (?, ?, ?)
              RETURNING id;`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnedIDString string

	row := u.db.QueryRowContext(ctx, query, userIDString, firstName, lastName)

	err = row.Scan(&returnedIDString)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("l'insertion a échoué ou n'a pas retourné d'ID")
		}
		return nil, fmt.Errorf("échec de l'exécution ou du scan de l'ID retourné: %w", err)
	}
	err = row.Err()
	if err != nil {
		return nil, fmt.Errorf("erreur différée sur la ligne: %w", err)
	}

	parsedUUID, err := uuid.Parse(returnedIDString)
	if err != nil {
		return nil, fmt.Errorf("échec du parsing de l'UUID retourné '%s': %w", returnedIDString, err)
	}
	return &parsedUUID, nil
}

func (r *userRepository) UpdateUser(id uuid.UUID, firstName string, lastName string) error {
	query := `UPDATE user
              SET first_name = ?, last_name = ?
              WHERE id = ?;`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	idString := id.String()

	result, err := r.db.ExecContext(ctx, query, firstName, lastName, idString)
	if err != nil {
		return fmt.Errorf("error executing user update: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Unable to get rows affected, but update might have worked
		// Consider logging this error
	} else if rowsAffected == 0 {
		// If 0 rows affected, the user ID did not match any record
		return fmt.Errorf("user with ID %s not found or no changes made", idString)
	}

	return nil
}
