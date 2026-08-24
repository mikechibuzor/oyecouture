ALTER TABLE measurements
    DROP CONSTRAINT IF EXISTS measurements_values_object,
    DROP CONSTRAINT IF EXISTS measurements_unit_valid,
    DROP CONSTRAINT IF EXISTS measurements_profile_valid,
    DROP COLUMN IF EXISTS values,
    DROP COLUMN IF EXISTS unit,
    DROP COLUMN IF EXISTS profile;
