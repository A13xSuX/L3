-- +goose Up
CREATE TABLE IF NOT EXISTS sales(
    id UUID primary key DEFAULT gen_random_uuid(),
    title varchar(255) NOT NULL,
    category varchar(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL CHECK ( price > 0.0 ) ,
    quantity INT NOT NULL CHECK ( quantity > 0 ),
    sale_date TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS sales;
