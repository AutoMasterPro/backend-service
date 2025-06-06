ALTER TABLE appointments ADD COLUMN attachments text[] DEFAULT ARRAY[]::text[];

ALTER TABLE appointments
    DROP CONSTRAINT IF EXISTS appointments_status_check,
    ADD CONSTRAINT appointments_status_check
        CHECK (status IN ('scheduled', 'in_progress', 'completed', 'cancelled'));