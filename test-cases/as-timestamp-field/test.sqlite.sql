
DROP TABLE IF EXISTS "as_timestamp";
CREATE TABLE IF NOT EXISTS "as_timestamp" (
    "id" INTEGER,
    "timestamp" TEXT /* stored as RFC3339 string */,
    "timestamp_two" TEXT /* stored as RFC3339 string */,

    PRIMARY KEY (
        "id"
    )
);
