export $(grep -v '^#' .env | xargs)

echo "Applying database migrations..."

for file in internal/db/migrations/*.sql; do
    echo "📂 Running migration: $file"
    PGPASSWORD=$DB_PASS psql -U $DB_USER -h $DB_HOST -p $DB_PORT -d $DB_NAME -f "$file"
done

echo "Database setup done."