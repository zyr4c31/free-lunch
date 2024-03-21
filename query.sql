-- name: CreateRestaurant :exec
INSERT INTO restaurants ( name ) VALUES ( $name );

-- name: ListRestaurants :many
select *
from restaurants
;

