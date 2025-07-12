-- name: CreateCategory :exec
insert into categories (
  name,
  type,
  user_id
)
values ($1, $2, $3);

-- name: GetCategory :one
select *
  from categories
 where id = $1;

-- name: GetAllCategoriesByUserId :many
select *
  from categories
 where user_id = $1;

-- name: UpdateCategory :exec
update categories
   set name = $2,
       type = $3,
       updated_at = now()
 where id = $1;

-- name: DeleteCategory :exec
delete from categories
 where id = $1;