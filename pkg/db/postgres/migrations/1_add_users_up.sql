CREATE TABLE IF NOT EXISTS "users"
(
    id SERIAL UNIQUE NOT NULL,
    name text,
    age text,
    CONSTRAINT "pk_Users" PRIMARY KEY ("id")
);




