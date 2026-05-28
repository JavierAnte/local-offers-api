CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE offers (
    id UUID PRIMARY KEY,

    headline TEXT NOT NULL,
    description TEXT,

    business_name TEXT NOT NULL,
    category TEXT NOT NULL,

    image_url TEXT,

    offer_type JSONB NOT NULL,

    location GEOGRAPHY(POINT, 4326) NOT NULL,

    expires_at TIMESTAMP,

    confirmations_count INTEGER NOT NULL DEFAULT 0,
    invalidations_count INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX offers_location_idx
ON offers
USING GIST(location);