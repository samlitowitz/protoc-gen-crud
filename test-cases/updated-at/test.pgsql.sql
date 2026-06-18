
DROP TABLE IF EXISTS "updated_at";
CREATE TABLE IF NOT EXISTS "updated_at" (
    "id" INTEGER,
    "updated_at" TIMESTAMP WITH TIME ZONE,
    "data" TEXT,

    PRIMARY KEY (
        "id"
    )
);
