-- add extension for uuid generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


-- create users table 
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(60) UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);