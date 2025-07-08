-- name: CreateAccount :exec
INSERT INTO accounts (
  user_id,
  name,
  type,
  balance
) VALUES ($1, $2, $3, $4);

-- name: GetAccount :one
SELECT * FROM accounts WHERE id = $1;

-- name: GetAccountsByUserId :many
SELECT * FROM accounts WHERE user_id = $1;

-- name: UpdateAccount :exec
UPDATE accounts
SET name = $2,
    type = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;