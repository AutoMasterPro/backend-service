package storages

import (
	"backend-service/internal/entity"
	"backend-service/pkg/database"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (userID uuid.UUID, err error)
	GetById(ctx context.Context, id uuid.UUID) (user *entity.User, err error)
	GetByEmail(ctx context.Context, email string) (user *entity.User, err error)
	GetAllClients(ctx context.Context) ([]*entity.User, error)
}

type userStorage struct {
	pg *database.PostgresDB
}

func NewUserStorage(deps StorageDeps) UserRepository {
	return &userStorage{
		pg: deps.PostgresDB,
	}
}

func (s *userStorage) Create(ctx context.Context, user *entity.User) (uuid.UUID, error) {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	const query = `
		INSERT INTO users (id, full_name, phone, email, password_hash, is_admin)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`

	row := s.pg.DB.QueryRowContext(ctx, query,
		user.ID, user.FullName, user.Phone, user.Email, user.PasswordHash, user.IsAdmin,
	)

	if err := row.Scan(&user.ID); err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return user.ID, nil
}

func (s *userStorage) GetById(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	const query = `
		SELECT id, full_name, phone, email, password_hash, is_admin, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL;
	`

	row := s.pg.DB.QueryRowContext(ctx, query, id)

	var user entity.User
	if err := row.Scan(
		&user.ID, &user.FullName, &user.Phone, &user.Email,
		&user.PasswordHash, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (s *userStorage) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	const query = `
		SELECT id, full_name, phone, email, password_hash, is_admin
		FROM users
		WHERE email = $1 AND deleted_at IS NULL;
	`

	row := s.pg.DB.QueryRowContext(ctx, query, email)

	var user entity.User
	if err := row.Scan(&user.ID, &user.FullName, &user.Phone, &user.Email, &user.PasswordHash, &user.IsAdmin); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userStorage) GetAllClients(ctx context.Context) ([]*entity.User, error) {
	const query = `
		SELECT id, full_name, phone, email, is_admin, created_at, updated_at
		FROM users
		WHERE is_admin = false AND deleted_at IS NULL;
	`
	rows, err := s.pg.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(
			&user.ID, &user.FullName, &user.Phone, &user.Email,
			&user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
