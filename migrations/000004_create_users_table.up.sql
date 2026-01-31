CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    -- citext means case-insensitive text
    -- This type stores text data
    -- exactly as it is inputted — without changing the case in any way — but comparisons
    -- against the data are always case-insensitive… including lookups on associated indexes.
    email citext UNIQUE NOT NULL,
    -- bytea means binary string
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    version integer NOT NULL DEFAULT 1
);