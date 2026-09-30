CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE users (
    id UUID PRIMARY KEY,

    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE offers (
    id UUID PRIMARY KEY,

    headline TEXT NOT NULL,
    description TEXT,

    business_name TEXT NOT NULL,
    category TEXT NOT NULL,

    image_url TEXT,
    offer_type JSONB NOT NULL,

    location GEOGRAPHY(POINT, 4326) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),

    expires_at TIMESTAMPTZ,

    confirmations_count INTEGER NOT NULL DEFAULT 0,
    invalidations_count INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX offers_location_idx
    ON offers
    USING GIST(location);

CREATE INDEX idx_offers_user_created_at
    ON offers(user_id, created_at DESC);

CREATE TABLE comments (
    id UUID PRIMARY KEY,

    offer_id UUID NOT NULL REFERENCES offers(id),
    user_id UUID NOT NULL REFERENCES users(id),

    body TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_comments_offer_id ON comments(offer_id);

CREATE TABLE offer_votes (
    id UUID PRIMARY KEY,

    offer_id UUID NOT NULL REFERENCES offers(id),
    user_id UUID NOT NULL REFERENCES users(id),

    type TEXT NOT NULL CHECK (type IN ('validate', 'invalidate')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (offer_id, user_id)
);

CREATE INDEX idx_offer_votes_offer_id ON offer_votes(offer_id);

CREATE TABLE notifications (
    id UUID PRIMARY KEY,

    recipient_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    actor_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    offer_id UUID NOT NULL REFERENCES offers(id) ON DELETE CASCADE,

    type TEXT NOT NULL CHECK (type IN ('comment_received', 'offer_validated', 'offer_invalidated')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notifications_recipient_created_at
    ON notifications(recipient_user_id, created_at DESC);
