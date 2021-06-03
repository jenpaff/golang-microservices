CREATE TABLE IF NOT EXISTS users(
   id serial PRIMARY KEY,
   username VARCHAR (50) UNIQUE NOT NULL,
   email VARCHAR (300) UNIQUE,
   phone_number VARCHAR (20) UNIQUE
);
