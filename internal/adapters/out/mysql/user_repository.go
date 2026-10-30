package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"

	"github.com/JoelTeoGom/login-service/internal/domain"
	"github.com/JoelTeoGom/login-service/internal/ports"
)

const errDuplicateEntry = 1062

// UserRepository implements ports.UserRepository on top of MySQL.
type UserRepository struct {
	db *sql.DB
}

var _ ports.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO users (email, name, password_hash) VALUES (?, ?, ?)`,
		user.Email, user.Name, user.PasswordHash,
	)
	if isDuplicate(err) {
		return domain.ErrEmailTaken
	}
	if err != nil {
		return fmt.Errorf("mysql: insert user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("mysql: last insert id: %w", err)
	}
	user.ID = id
	user.CreatedAt = time.Now().UTC()
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	return r.findOne(ctx, `SELECT id, email, name, password_hash, created_at FROM users WHERE id = ?`, id)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.findOne(ctx, `SELECT id, email, name, password_hash, created_at FROM users WHERE email = ?`, email)
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	if _, err := r.db.ExecContext(ctx, `UPDATE users SET name = ? WHERE id = ?`, user.Name, user.ID); err != nil {
		return fmt.Errorf("mysql: update user: %w", err)
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("mysql: delete user: %w", err)
	}
	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *UserRepository) findOne(ctx context.Context, query string, arg any) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRowContext(ctx, query, arg).
		Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("mysql: select user: %w", err)
	}
	return &user, nil
}

// isDuplicate translates MySQL's own error into something the domain understands.
func isDuplicate(err error) bool {
	var me *mysqldriver.MySQLError
	return errors.As(err, &me) && me.Number == errDuplicateEntry
}
