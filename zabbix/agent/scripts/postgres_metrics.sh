#!/bin/sh

PG_HOST="pg-dialog-service"
PG_USER="postgres"
PG_PASSWORD="pass"
PG_DB="postgres"

export PGPASSWORD=$PG_PASSWORD

case $1 in
    "connections")
        psql -h $PG_HOST -U $PG_USER -d $PG_DB -t -c "SELECT count(*) FROM pg_stat_activity;" 2>/dev/null | tr -d '[:space:]' | head -1
        ;;
    "active_connections")
        psql -h $PG_HOST -U $PG_USER -d $PG_DB -t -c "SELECT count(*) FROM pg_stat_activity WHERE state = 'active';" 2>/dev/null | tr -d '[:space:]' | head -1
        ;;
    "db_size")
        psql -h $PG_HOST -U $PG_USER -d $PG_DB -t -c "SELECT pg_database_size('$PG_DB');" 2>/dev/null | tr -d '[:space:]' | head -1
        ;;
    "queries_count")
        psql -h $PG_HOST -U $PG_USER -d $PG_DB -t -c "SELECT COALESCE(sum(xact_commit + xact_rollback), 0) FROM pg_stat_database WHERE datname = '$PG_DB';" 2>/dev/null | tr -d '[:space:]' | head -1
        ;;
    "locks_count")
        psql -h $PG_HOST -U $PG_USER -d $PG_DB -t -c "SELECT count(*) FROM pg_locks WHERE granted = false;" 2>/dev/null | tr -d '[:space:]' | head -1
        ;;
    "cache_hitrate")
        psql -h $PG_HOST -U $PG_USER -d $PG_DB -t -c "SELECT COALESCE(round(blks_hit::decimal / greatest(blks_hit + blks_read, 1) * 100, 2), 0) FROM pg_stat_database WHERE datname = '$PG_DB';" 2>/dev/null | tr -d '[:space:]' | head -1
        ;;
    "transactions_rate")
        # Transactions per second (xact_commit + xact_rollback)
        psql -h $PG_HOST -U $PG_USER -d $PG_DB -t -c "SELECT COALESCE(sum(xact_commit + xact_rollback), 0) FROM pg_stat_database WHERE datname = '$PG_DB';" 2>/dev/null | awk 'NR==1{print $1+0}'
        ;;
    *)
        echo "0"
        ;;
esac
