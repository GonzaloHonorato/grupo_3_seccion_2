package persistence

import (
	"context"

	"github.com/gonzalohonorato/servercorego/core/user/domain/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TimescaleUserRepository struct {
	dbPool *pgxpool.Pool
}

func NewTimescaleDBRepository(pool *pgxpool.Pool) *TimescaleUserRepository {
	return &TimescaleUserRepository{
		dbPool: pool,
	}
}

func (r *TimescaleUserRepository) SearchUserByID(id string) (*entities.User, error) {
	ctx := context.Background()
	query := `SELECT * FROM user_details WHERE user_id = $1`
	row := r.dbPool.QueryRow(ctx, query, id)

	var u entities.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Rut, &u.Uid, &u.Type, &u.CreatedAt,
		&u.CustomerType, &u.EmployeeRole)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *TimescaleUserRepository) SearchUsers() (*entities.Users, error) {
	ctx := context.Background()
	query := `SELECT * FROM user_details`
	rows, err := r.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users entities.Users
	for rows.Next() {
		var u entities.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Rut, &u.Uid, &u.Type, &u.CreatedAt,
			&u.CustomerType, &u.EmployeeRole); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &users, nil
}

func (r *TimescaleUserRepository) SearchUsersByType(userType string) (*entities.Users, error) {
	ctx := context.Background()
	query := `SELECT * FROM user_details WHERE user_type = $1`
	rows, err := r.dbPool.Query(ctx, query, userType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users entities.Users
	for rows.Next() {
		var u entities.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Rut, &u.Uid, &u.Type, &u.CreatedAt,
			&u.CustomerType, &u.EmployeeRole); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &users, nil
}

func (r *TimescaleUserRepository) CreateUser(user *entities.User) error {
	ctx := context.Background()
	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	
	userQuery := `
		INSERT INTO "user" (id, name, email, rut, uid, type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
	`
	_, err = tx.Exec(ctx, userQuery, user.ID, user.Name, user.Email, user.Rut, user.Uid, user.Type)
	if err != nil {
		return err
	}

	
	if user.Type == "customer" {
		customerQuery := `
			INSERT INTO customer (id, type)
			VALUES ($1, $2)
		`
		_, err = tx.Exec(ctx, customerQuery, user.ID, user.CustomerType)
		if err != nil {
			return err
		}
	} else if user.Type == "employee" {
		employeeQuery := `
			INSERT INTO employee (id, role)
			VALUES ($1, $2)
		`
		_, err = tx.Exec(ctx, employeeQuery, user.ID, user.EmployeeRole)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *TimescaleUserRepository) UpdateUser(user *entities.User) error {
	ctx := context.Background()
	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	
	userQuery := `
		UPDATE "user" 
		SET name = $2, email = $3, rut = $4, uid = $5, type = $6
		WHERE id = $1
	`
	_, err = tx.Exec(ctx, userQuery, user.ID, user.Name, user.Email, user.Rut, user.Uid, user.Type)
	if err != nil {
		return err
	}

	
	if user.Type == "customer" {
		
		var exists bool
		err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM customer WHERE id = $1)", user.ID).Scan(&exists)
		if err != nil {
			return err
		}

		if exists {
			
			_, err = tx.Exec(ctx, "UPDATE customer SET type = $2 WHERE id = $1", user.ID, user.CustomerType)
		} else {
			
			_, err = tx.Exec(ctx, "INSERT INTO customer (id, type) VALUES ($1, $2)", user.ID, user.CustomerType)
		}
		if err != nil {
			return err
		}

		
		_, err = tx.Exec(ctx, "DELETE FROM employee WHERE id = $1", user.ID)
		if err != nil {
			return err
		}
	} else if user.Type == "employee" {
		
		var exists bool
		err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM employee WHERE id = $1)", user.ID).Scan(&exists)
		if err != nil {
			return err
		}

		if exists {
			
			_, err = tx.Exec(ctx, "UPDATE employee SET role = $2 WHERE id = $1", user.ID, user.EmployeeRole)
		} else {
			
			_, err = tx.Exec(ctx, "INSERT INTO employee (id, role) VALUES ($1, $2)", user.ID, user.EmployeeRole)
		}
		if err != nil {
			return err
		}

		
		_, err = tx.Exec(ctx, "DELETE FROM customer WHERE id = $1", user.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *TimescaleUserRepository) DeleteUserByID(id string) error {
	ctx := context.Background()
	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	
	_, err = tx.Exec(ctx, "DELETE FROM customer WHERE id = $1", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM employee WHERE id = $1", id)
	if err != nil {
		return err
	}

	
	_, err = tx.Exec(ctx, `DELETE FROM "user" WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
