CREATE TABLE offer_votes (
    id UUID PRIMARY KEY,

    offer_id UUID NOT NULL REFERENCES offers(id),
    user_id UUID NOT NULL REFERENCES users(id),

    type TEXT NOT NULL CHECK (type IN ('validate', 'invalidate')),

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    UNIQUE (offer_id, user_id)
);

CREATE INDEX idx_offer_votes_offer_id ON offer_votes(offer_id);
