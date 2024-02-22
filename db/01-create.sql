CREATE TABLE products(
    id INT AUTO_INCREMENT PRIMARY KEY
    , name VARCHAR(255) NOT NULL
    , description TEXT
    , price FLOAT
    , category VARCHAR(255)
);

CREATE TABLE orders (
    order_id VARCHAR(36)
    , product_id INT
    , quantity INT
    , customer_info VARCHAR(255)
    , FOREIGN KEY (product_id) REFERENCES products(id)
);
