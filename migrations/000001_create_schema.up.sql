CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_role AS ENUM ('employee', 'moderator');
CREATE TYPE city_type AS ENUM ('Москва', 'Санкт-Петербург', 'Казань');
CREATE TYPE product_type AS ENUM ('электроника', 'одежда', 'обувь');
CREATE TYPE reception_status AS ENUM ('in_progress', 'close');

CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       email VARCHAR(255) NOT NULL UNIQUE,
                       password_hash VARCHAR(255) NOT NULL,
                       role user_role NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE pvzs (
                      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                      registration_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                      city city_type NOT NULL,
                      created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE receptions (
                            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                            date_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                            pvz_id UUID NOT NULL REFERENCES pvzs(id),
                            status reception_status NOT NULL DEFAULT 'in_progress',
                            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          date_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          type product_type NOT NULL,
                          reception_id UUID NOT NULL REFERENCES receptions(id),
                          created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для повышения производительности
CREATE INDEX idx_receptions_pvz_id ON receptions(pvz_id);
CREATE INDEX idx_receptions_status ON receptions(status);
CREATE INDEX idx_products_reception_id ON products(reception_id);
CREATE INDEX idx_products_date_time ON products(date_time);
