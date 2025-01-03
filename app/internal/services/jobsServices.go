package services

import (
	"context"
	"errors"
	"time"

	"github.com/galrub/go/jobSearch/internal/database"
	"github.com/galrub/go/jobSearch/internal/pool"
	"github.com/gofrs/uuid"
)

type Jobs []database.Job

func GetJob(id uuid.UUID) (database.Job, error) {
	if id.IsNil() {
		return database.Job{}, errors.New("id provided is not valid")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pool.GetConnection(ctx)
	if err != nil {
		return database.Job{}, err
	}
	defer conn.Release()

	return database.New(conn).GetJobById(ctx, id)
}

func GetJobs(email string) (Jobs, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pool.GetConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	return database.New(conn).GetJobsByUserEmail(ctx, email)
}

func DeleteJob(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pool.GetConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	return database.New(conn).DeleteJob(ctx, id)
}

func SaveJob(job *database.InsertJobParams, email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pool.GetConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	db := database.New(conn)

	user, err := db.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	job.UserID = user.ID

	err = db.InsertJob(ctx, *job)
	return err
}

func GetJobByCompanyAndDescription(company string, positionDescription string) (database.Job, error) {
	var job database.Job
	if company == "" || positionDescription == "" {
		return job, errors.New("inputs cannot be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pool.GetConnection(ctx)
	if err != nil {
		return job, err
	}
	defer conn.Release()

	params := database.GetJobByCopmpanyAndDescriptionParams{
		Company:      company,
		PositionDesc: positionDescription,
	}

	return database.New(conn).GetJobByCopmpanyAndDescription(ctx, params)
}

func UpdateJob(params database.UpdateJobParams) error {
	if params.Company == "" || params.PositionDesc == "" {
		return errors.New("company and postion cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pool.GetConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	return database.New(conn).UpdateJob(ctx, params)
}
