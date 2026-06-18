
DROP TABLE IF EXISTS "inline_timestamp";
CREATE TABLE IF NOT EXISTS "inline_timestamp" (
    "id" INTEGER,
    "timestamp_seconds" INTEGER,
    "timestamp_nanos" INTEGER,

    PRIMARY KEY (
        "id"
    )
);
