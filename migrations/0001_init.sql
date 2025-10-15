CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $$ BEGIN CREATE TYPE role_enum AS ENUM ('USER','COURIER','ADMIN'); EXCEPTION WHEN duplicate_object THEN null; END $$;
DO $$ BEGIN CREATE TYPE plan_enum AS ENUM ('P7','P15','P30'); EXCEPTION WHEN duplicate_object THEN null; END $$;
DO $$ BEGIN CREATE TYPE sub_status_enum AS ENUM ('ACTIVE','PAUSED','CANCELED','EXPIRED'); EXCEPTION WHEN duplicate_object THEN null; END $$;
DO $$ BEGIN CREATE TYPE order_type_enum AS ENUM ('ONE_TIME','SUBSCRIPTION'); EXCEPTION WHEN duplicate_object THEN null; END $$;
DO $$ BEGIN CREATE TYPE time_option_enum AS ENUM ('ASAP','SCHEDULED'); EXCEPTION WHEN duplicate_object THEN null; END $$;
DO $$ BEGIN CREATE TYPE order_status_enum AS ENUM ('NEW','PAID','ASSIGNED','PICKING_UP','DONE','CANCELED','REFUNDED'); EXCEPTION WHEN duplicate_object THEN null; END $$;
DO $$ BEGIN CREATE TYPE payment_provider_enum AS ENUM ('PAYNETWORKS'); EXCEPTION WHEN duplicate_object THEN null; END $$;
DO $$ BEGIN CREATE TYPE payment_status_enum AS ENUM ('INIT','REQUIRES_ACTION','SUCCEEDED','FAILED','CANCELED'); EXCEPTION WHEN duplicate_object THEN null; END $$;
DO $$ BEGIN CREATE TYPE discount_type_enum AS ENUM ('FIXED','PERCENT'); EXCEPTION WHEN duplicate_object THEN null; END $$;

CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  phone text UNIQUE NOT NULL,
  email text UNIQUE,
  name text NOT NULL,
  password_hash text,
  role role_enum NOT NULL DEFAULT 'USER',
  is_deleted boolean NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS addresses (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id),
  label text,
  lat double precision NOT NULL,
  lng double precision NOT NULL,
  city text NOT NULL,
  street text NOT NULL,
  house text NOT NULL,
  entrance text NOT NULL,
  floor text NOT NULL,
  apartment text NOT NULL,
  intercom text,
  is_default boolean NOT NULL DEFAULT false,
  polygon_id uuid,
  polygon_name text,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_addresses_user ON addresses(user_id);
CREATE INDEX IF NOT EXISTS idx_addresses_polygon ON addresses(polygon_id);

CREATE TABLE IF NOT EXISTS polygons (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text NOT NULL,
  city text NOT NULL,
  geojson jsonb NOT NULL,
  is_active boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS subscriptions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id),
  plan plan_enum NOT NULL,
  total_bags int NOT NULL,
  remaining_bags int NOT NULL,
  price_kzt int NOT NULL,
  status sub_status_enum NOT NULL DEFAULT 'ACTIVE',
  started_at timestamptz NOT NULL DEFAULT now(),
  expires_at timestamptz
);
CREATE INDEX IF NOT EXISTS idx_subscriptions_user ON subscriptions(user_id);

CREATE TABLE IF NOT EXISTS orders (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id),
  address_id uuid NOT NULL REFERENCES addresses(id),
  polygon_id uuid NOT NULL REFERENCES polygons(id),
  type order_type_enum NOT NULL,
  bags_count int NOT NULL,
  price_kzt int NOT NULL,
  comment text NOT NULL DEFAULT '',
  time_option time_option_enum NOT NULL,
  scheduled_at timestamptz,
  courier_id uuid,
  status order_status_enum NOT NULL DEFAULT 'NEW',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_orders_user ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_polygon ON orders(polygon_id);

CREATE TABLE IF NOT EXISTS order_events (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id uuid NOT NULL REFERENCES orders(id),
  from_status order_status_enum,
  to_status order_status_enum NOT NULL,
  at timestamptz NOT NULL DEFAULT now(),
  meta jsonb NOT NULL DEFAULT '{}'::jsonb
);

CREATE TABLE IF NOT EXISTS payments (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id),
  order_id uuid,
  subscription_id uuid,
  amount_kzt int NOT NULL,
  provider payment_provider_enum NOT NULL,
  status payment_status_enum NOT NULL DEFAULT 'INIT',
  provider_intent_id text NOT NULL,
  provider_payload jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_payments_user ON payments(user_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);

CREATE TABLE IF NOT EXISTS promocodes (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  code text UNIQUE NOT NULL,
  discount_type discount_type_enum NOT NULL,
  value int NOT NULL,
  active_from timestamptz NOT NULL,
  active_to timestamptz NOT NULL,
  usage_limit int NOT NULL,
  used_count int NOT NULL DEFAULT 0,
  is_active boolean NOT NULL DEFAULT true
);

CREATE TABLE IF NOT EXISTS user_promocodes (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id),
  promocode_id uuid NOT NULL REFERENCES promocodes(id),
  used_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS device_tokens (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id),
  platform text NOT NULL,
  token text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sessions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id),
  refresh_token_hash text NOT NULL,
  expires_at timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  revoked_at timestamptz
);
