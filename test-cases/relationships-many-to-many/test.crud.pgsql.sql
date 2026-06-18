
DROP TABLE IF EXISTS "primary_key_enum";
CREATE TABLE IF NOT EXISTS "primary_key_enum" (
    "id" INTEGER PRIMARY KEY,
    "value" TEXT
);

INSERT INTO "primary_key_enum" ("id", "value") VALUES
    (0, 'ZERO'),
    (1, 'ONE'),
    (2, 'TWO'),
    (3, 'THREE')
;

DROP TABLE IF EXISTS "sa_int_32_ma_all";
CREATE TABLE IF NOT EXISTS "sa_int_32_ma_all" (
    "maall_id_uint_32" INTEGER,
    "maall_id_uint_64" INTEGER,
    "maall_id_string" TEXT,
    "saint_32_id" INTEGER,
    "maall_id_enum" INTEGER /* references "primary_key_enum"."id" */,
    "maall_id_int_32" INTEGER,
    "maall_id_int_64" INTEGER,

    PRIMARY KEY (
        "maall_id_uint_32",
        "maall_id_uint_64",
        "maall_id_string",
        "saint_32_id",
        "maall_id_enum",
        "maall_id_int_32",
        "maall_id_int_64"
    )
);
