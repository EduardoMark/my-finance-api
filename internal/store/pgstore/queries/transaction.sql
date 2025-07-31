-- name: CreateTransaction :exec
insert into transactions (
  description,
  amount,
  date,
  type,
  user_id,
  account_id,
  category_id
)
values ($1, $2, $3, $4, $5, $6, $7);

-- name: GetTrasaction :one
select *
  from transactions
 where id = $1;

-- name: GetAllTransactions :many
select *
  from transactions
 where user_id = $1;

-- name: GetAllTransactionsByAccount :many
select *
  from transactions
 where account_id = $1;

-- name: GetAllTransactionsByCategory :many
select *
  from transactions
 where category_id = $1;

-- name: UpdateTransaction :exec
update transactions
   set description = $2,
   amount = $3,
   date = $4,
   type = $5,
   updated_at = now()
 where id = $1;

-- name: DeleteTransaction :exec
delete from transactions
 where id = $1;