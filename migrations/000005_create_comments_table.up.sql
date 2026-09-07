CREATE TABLE comments (
    id UUID PRIMARY KEY,

    offer_id UUID NOT NULL REFERENCES offers(id),
    user_id UUID NOT NULL REFERENCES users(id),

    body TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_comments_offer_id ON comments(offer_id);
