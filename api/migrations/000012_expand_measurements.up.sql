ALTER TABLE measurements
    ADD COLUMN profile TEXT,
    ADD COLUMN unit TEXT,
    ADD COLUMN values JSONB NOT NULL DEFAULT '{}'::jsonb,
    ADD CONSTRAINT measurements_profile_valid
        CHECK (profile IS NULL OR profile IN ('suit', 'native', 'mixed')),
    ADD CONSTRAINT measurements_unit_valid
        CHECK (unit IS NULL OR unit = 'in'),
    ADD CONSTRAINT measurements_values_object
        CHECK (jsonb_typeof(values) = 'object');
