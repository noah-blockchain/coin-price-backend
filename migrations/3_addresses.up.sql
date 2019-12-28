CREATE TABLE IF NOT EXISTS public.addresses
(
    id             serial PRIMARY KEY,
    address        character varying(40)                  NOT NULL,
    symbol         character varying(20)                  NOT NULL,
    amount         numeric(70, 0)           default 0.0,
    created_at     timestamp with time zone DEFAULT now() NOT NULL
);

create index addresses_address_index
    on addresses (address);
