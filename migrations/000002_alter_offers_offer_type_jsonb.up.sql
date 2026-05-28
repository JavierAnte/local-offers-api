ALTER TABLE offers
    ALTER COLUMN offer_type TYPE JSONB USING offer_type::JSONB;
