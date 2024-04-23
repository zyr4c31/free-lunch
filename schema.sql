CREATE TABLE IF NOT EXISTS restaurants(
    id INTEGER PRIMARY KEY NOT NULL,
    name TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS menu_items(
    id INTEGER PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    price DOUBLE NOT NULL,
    restaurant_id INTEGER NOT NULL,
    FOREIGN KEY(restaurant_id) REFERENCES restaurants(id)
);

