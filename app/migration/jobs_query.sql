-- name: GetJobsByUserId :many
SELECT * FROM jobs
WHERE user_id = $1;

-- name: GetJobById :one
SELECT * FROM jobs
WHERE id = $1 LIMIT 1;

-- name: InsertJob :exec
INSERT INTO public.jobs(
	user_id, company, position_desc, remote, contract_type, contacted, general_status, created_at, comments)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE id = $1;

-- name: GetJobsByUserEmail :many
SELECT j.* FROM jobs j
INNER JOIN users u ON j.user_id = u.id
WHERE u.email = $1;

-- name: GetJobByCopmpanyAndDescription :one
SELECT * FROM jobs
WHERE company = $1 AND position_desc = $2;

-- name: UpdateJob :exec
UPDATE jobs SET company=$1,
      position_desc=$2,
      remote=$3,
      contract_type=$4,
      contacted=$5, 
      created_at=$6,
      comments=$7,
      general_status=$8
WHERE id = $9;
