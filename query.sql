-- name: CreateRestaurant :exec
INSERT INTO restaurants ( name ) VALUES ( $name );

-- name: ListRestaurants :many
select *
from restaurants
;

-- name: ListMenuItemsForRestaurant :many
select *
from menu_items
where restaurant_id = $restaurant_id
;

