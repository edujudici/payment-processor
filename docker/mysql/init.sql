CREATE TABLE IF NOT EXISTS payments (
    id VARCHAR(36) PRIMARY KEY,
    protocol VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    quantity INT NOT NULL,
    total DECIMAL(10, 2) NOT NULL,
    subtotal DECIMAL(10, 2) NOT NULL,
    description TEXT,
    preference_id VARCHAR(255),
    external_reference VARCHAR(255),
    preference_init_point TEXT,
    preference_sandbox_init_point TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);