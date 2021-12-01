SELECT EXISTS (
    SELECT schema_name FROM information_schema.schemata WHERE schema_name = $1
);
