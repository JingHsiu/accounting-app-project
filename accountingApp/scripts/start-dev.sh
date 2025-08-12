#!/bin/bash

# Development startup script
set -e

echo "🐳 Starting PostgreSQL database..."
docker-compose up -d postgres

echo "⏳ Waiting for database to be ready..."
until docker-compose exec postgres pg_isready -U postgres -d accountingdb > /dev/null 2>&1; do
  echo "   Database is not ready yet, waiting 2 seconds..."
  sleep 2
done

echo "✅ Database is ready!"
echo "📊 Database connection: postgres://postgres:password@localhost:5432/accountingdb"

echo ""
echo "🚀 You can now start the application with:"
echo "   go run cmd/accoountingApp/main.go"
echo ""
echo "🔧 Optional: Start pgAdmin with:"
echo "   docker-compose --profile admin up -d pgadmin"
echo "   Then visit: http://localhost:8081"

echo ""
echo "🛑 To stop database:"
echo "   docker-compose down"