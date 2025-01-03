package services

import (
	"context"
	"time"

	"github.com/galrub/go/jobSearch/internal/pool"
	"github.com/galrub/go/jobSearch/internal/database"
)

func GetUserByEmail(email string) (database.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user database.User

	conn, err := pool.GetConnection(ctx)
	if err != nil {
		return user, err
	}
	defer conn.Release()

	user, err = database.New(conn).GetUserByEmail(ctx, email)
	return user, nil
}

func LoginUser(u *database.User, pw string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pool.GetConnection(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	var res bool

	rows := conn.QueryRow(ctx, "SELECT verify_user AS res FROM verify_user($1, $2)", u.ID, pw)

	err = rows.Scan(&res)
	if err != nil {
		return false, err
	}

	return res, nil
}

func LoginUserByEmail(email string, pw string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pool.GetConnection(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	var u database.User
	u, err = GetUserByEmail(email)

	var res bool
	rows := conn.QueryRow(ctx, "SELECT verify_user AS res FROM verify_user($1, $2)", u.ID, pw)
	err = rows.Scan(&res)
	if err != nil {
		return false, err
	}

	return res, nil
}
