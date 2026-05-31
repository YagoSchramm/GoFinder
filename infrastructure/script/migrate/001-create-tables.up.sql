CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS searches (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users(id),
	query TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS products (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	search_id UUID NOT NULL REFERENCES searches(id),
	title TEXT NOT NULL,
	price NUMERIC(12, 2) NOT NULL,
	store TEXT NOT NULL,
	url TEXT NOT NULL,
	thumbnail TEXT,
	found_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS price_history (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	product_url TEXT NOT NULL,
	store TEXT NOT NULL,
	price NUMERIC(12, 2) NOT NULL,
	captured_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS alerts (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users(id),
	product_url TEXT NOT NULL,
	store TEXT NOT NULL,
	target_price NUMERIC(12, 2) NOT NULL,
	active BOOLEAN NOT NULL DEFAULT true,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
