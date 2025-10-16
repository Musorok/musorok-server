-- Add tables for courier settlements, balances and payout requests
-- Defines the payout_status_enum and tables: order_settlements, courier_balances and payout_requests.

-- payout status enum defines the lifecycle of payout requests
DO $$ BEGIN
    CREATE TYPE payout_status_enum AS ENUM ('REQUESTED','APPROVED','PAID','REJECTED');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

-- order_settlements table records compensation owed to couriers after
-- each completed order. There is a unique constraint on order_id to
-- prevent duplicate settlements per order.
CREATE TABLE IF NOT EXISTS order_settlements (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id uuid NOT NULL REFERENCES orders(id),
    courier_id uuid NOT NULL REFERENCES couriers(id),
    bags_count int NOT NULL CHECK (bags_count > 0),
    amount_kzt int NOT NULL CHECK (amount_kzt >= 0),
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT unique_settlement_per_order UNIQUE(order_id)
);
CREATE INDEX IF NOT EXISTS idx_order_settlements_courier ON order_settlements(courier_id);

-- courier_balances keeps track of the aggregated earnings and withdrawals
-- for a courier. A row is created lazily on first settlement.
CREATE TABLE IF NOT EXISTS courier_balances (
    courier_id uuid PRIMARY KEY REFERENCES couriers(id),
    total_earned_kzt bigint NOT NULL DEFAULT 0,
    total_withdrawn_kzt bigint NOT NULL DEFAULT 0,
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- payout_requests records requests by couriers to withdraw money. The
-- status field uses the payout_status_enum created above. processed_at
-- remains NULL until a request is approved/paid or rejected.
CREATE TABLE IF NOT EXISTS payout_requests (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    courier_id uuid NOT NULL REFERENCES couriers(id),
    amount_kzt int NOT NULL CHECK (amount_kzt > 0),
    status payout_status_enum NOT NULL DEFAULT 'REQUESTED',
    requested_at timestamptz NOT NULL DEFAULT now(),
    processed_at timestamptz
);
CREATE INDEX IF NOT EXISTS idx_payout_requests_courier ON payout_requests(courier_id);