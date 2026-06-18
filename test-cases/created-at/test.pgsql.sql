
DROP TABLE IF EXISTS "created_at";
CREATE TABLE IF NOT EXISTS "created_at" (
    "id" INTEGER,
    "data" TEXT,
    "created_at" TIMESTAMP WITH TIME ZONE,

    PRIMARY KEY (
        "id"
    )
);
