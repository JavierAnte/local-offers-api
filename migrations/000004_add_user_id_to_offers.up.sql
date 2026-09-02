ALTER TABLE offers
    ADD COLUMN user_id UUID REFERENCES users(id);
