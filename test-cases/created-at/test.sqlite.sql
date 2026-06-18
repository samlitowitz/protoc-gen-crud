
DROP TABLE IF EXISTS "created_at";
CREATE TABLE IF NOT EXISTS "created_at" (
    "id" INTEGER,
    "data" TEXT,
    "created_at" TEXT /* stored as RFC3339 string */,

    PRIMARY KEY (
        "id"
    )
);
