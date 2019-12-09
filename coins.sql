CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

CREATE TABLE public.coins
(
    id                      serial                                PRIMARY KEY,
    volume                  numeric(70, 0),
    reserve_balance         numeric(70, 0),
    price                   numeric(100, 0)          default 0.0,
    capitalization          numeric(100, 0)          default 0.0,
    symbol                  character varying(20)                  NOT NULL,
    created_at              timestamp with time zone DEFAULT now() NOT NULL
);