
DROP TABLE IF EXISTS "updated_at";
CREATE TABLE IF NOT EXISTS "updated_at" (
    "id" INTEGER,
    "data" TEXT,
    "updated_at" TEXT /* stored as RFC3339 string */,

    PRIMARY KEY (
        "id"
    )
);
