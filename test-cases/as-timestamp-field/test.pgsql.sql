
DROP TABLE IF EXISTS "as_timestamp";
CREATE TABLE IF NOT EXISTS "as_timestamp" (
    "id" INTEGER,
    "timestamp" TIMESTAMP WITH TIME ZONE,
    "timestamp_two" TIMESTAMP WITH TIME ZONE,

    PRIMARY KEY (
        "id"
    )
);
