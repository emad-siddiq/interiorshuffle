-- Create the database
CREATE DATABASE realestate;

-- Create a database user with limited privileges
CREATE USER realestate_user WITH ENCRYPTED PASSWORD 'securepassword';
GRANT CONNECT ON DATABASE realestate TO realestate_user;
GRANT USAGE ON SCHEMA public TO realestate_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO realestate_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO realestate_user;


-- Create the properties table
CREATE TABLE IF NOT EXISTS properties (
    id SERIAL PRIMARY KEY,
    data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Create the property_details table
CREATE TABLE IF NOT EXISTS property_details (
    property_id VARCHAR(255) PRIMARY KEY,  -- Unique identifier for the property
    address TEXT NOT NULL,                 -- Full address of the property
    price DECIMAL(15, 2) NOT NULL,         -- Price of the property
    fetched_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP  -- Timestamp when the property data was fetched
);
