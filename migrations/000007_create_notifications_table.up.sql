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
