
DROP TABLE IF EXISTS "updated_at";
CREATE TABLE IF NOT EXISTS "updated_at" (
    "id" INTEGER,
    "data" TEXT,
    "updated_at" TIMESTAMP WITH TIME ZONE,

    PRIMARY KEY (
        "id"
    )
);
