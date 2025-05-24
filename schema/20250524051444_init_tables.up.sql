-- Создание таблицы пользователей
CREATE TABLE users
(
    id            UUID PRIMARY KEY,
    full_name     TEXT NOT NULL,
    phone         TEXT NOT NULL UNIQUE,
    email         TEXT UNIQUE,
    password_hash TEXT NOT NULL,
    is_admin      BOOLEAN   DEFAULT False,
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW(),
    deleted_at    TIMESTAMP
);

-- Создание таблицы услуг
CREATE TABLE services
(
    id           UUID PRIMARY KEY,
    name         TEXT           NOT NULL,
    description  TEXT,
    price        NUMERIC(10, 2) NOT NULL,
    duration_min INT            NOT NULL,
    created_at   TIMESTAMP DEFAULT NOW(),
    updated_at   TIMESTAMP DEFAULT NOW(),
    deleted_at   TIMESTAMP
);

-- Создание таблицы автомобилей
CREATE TABLE vehicles
(
    id            UUID PRIMARY KEY,
    user_id       UUID NOT NULL,
    brand         TEXT NOT NULL,
    model         TEXT NOT NULL,
    license_plate TEXT NOT NULL UNIQUE,
    year          INT,
    vin           TEXT UNIQUE,
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW(),
    deleted_at    TIMESTAMP
);

-- Создание таблицы записей на осмотр
CREATE TABLE appointments
(
    id               UUID PRIMARY KEY,
    user_id          UUID      NOT NULL,
    vehicle_id       UUID      NOT NULL,
    appointment_time TIMESTAMP NOT NULL,
    status           TEXT      NOT NULL CHECK (status IN ('scheduled', 'completed', 'cancelled')),
    created_at       TIMESTAMP DEFAULT NOW(),
    updated_at       TIMESTAMP DEFAULT NOW(),
    deleted_at       TIMESTAMP
);

-- Создание связующей таблицы заказанных услуг
CREATE TABLE appointment_services
(
    id             UUID PRIMARY KEY,
    appointment_id UUID           NOT NULL REFERENCES appointments (id) ON DELETE CASCADE,
    service_id     UUID           NOT NULL REFERENCES services (id) ON DELETE CASCADE,
    price          NUMERIC(10, 2) NOT NULL,
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP DEFAULT NOW(),
    deleted_at     TIMESTAMP
);