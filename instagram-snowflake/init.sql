-- ID Structure:
-- 
-- +------------------+-----------+------------+
-- |    Timestamp     | Shard ID  |  Sequence  |
-- |     41 bits      |  13 bits  |   10 bits  |
-- +------------------+-----------+------------+
-- 
--  63                22          9            0
-- +------------------+-----------+------------+
-- |    Timestamp     | Shard ID  |  Sequence  |
-- +------------------+-----------+------------+
--
-- - Timestamp: Milliseconds since custom epoch (1314220021721)
-- - Shard ID: 13 bits, allowing for 8,192 shards
-- - Sequence: 10 bits, allowing 1,024 IDs per millisecond per shard
--
-- This structure allows for:
-- - ~69 years of IDs from the custom epoch (2011-08-24 19:07:01.721 UTC)
-- - 8,192 shards
-- - 1,024 IDs per millisecond per shard
CREATE SEQUENCE IF NOT EXISTS public.global_id_seq;

CREATE OR REPLACE FUNCTION public.id_generator() 
RETURNS bigint 
LANGUAGE 'plpgsql' 
AS $BODY$
DECLARE
    our_epoch bigint := 1314220021721; -- Custom epoch timestamp
    seq_id bigint;
    now_millis bigint;
    shard_id int := 1; -- Set this for each schema shard
    result bigint := 0;
BEGIN
    SELECT nextval('public.global_id_seq') % 1024 INTO seq_id; -- Get sequence number
    SELECT FLOOR(EXTRACT(EPOCH FROM clock_timestamp()) * 1000) INTO now_millis; -- Current time in milliseconds
    result := (now_millis - our_epoch) << 23; -- Shift timestamp left by 23 bits
    result := result | (shard_id << 10); -- Shift shard ID left by 10 bits
    result := result | seq_id; -- Combine with sequence ID
    RETURN result; -- Return the generated ID
END;
$BODY$;

ALTER FUNCTION public.id_generator() OWNER TO postgres;

CREATE TABLE IF NOT EXISTS public.users (
    id bigint NOT NULL DEFAULT id_generator(),
    email varchar(255) NOT NULL UNIQUE,
    first varchar(50),
    last varchar(50)
);