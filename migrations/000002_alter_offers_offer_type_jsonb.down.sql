ALTER TABLE offers
    ALTER COLUMN offer_type TYPE TEXT USING offer_type::TEXT;
