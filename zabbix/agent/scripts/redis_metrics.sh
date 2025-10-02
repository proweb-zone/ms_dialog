#!/bin/bash

REDIS_HOST="redis-service.shared"
REDIS_PASSWORD="123123vv"

case $1 in
    "memory_used")
        redis-cli -h $REDIS_HOST -a $REDIS_PASSWORD info memory 2>/dev/null | grep "used_memory:" | cut -d: -f2 | tr -d '\r' || echo 0
        ;;
    "connected_clients")
        redis-cli -h $REDIS_HOST -a $REDIS_PASSWORD info clients 2>/dev/null | grep "connected_clients:" | cut -d: -f2 | tr -d '\r' || echo 0
        ;;
    "keys_count")
        redis-cli -h $REDIS_HOST -a $REDIS_PASSWORD dbsize 2>/dev/null | tr -d '\r' || echo 0
        ;;
    "hit_rate")
        redis-cli -h $REDIS_HOST -a $REDIS_PASSWORD info stats 2>/dev/null | grep "keyspace_hits" | cut -d: -f2 | tr -d '\r' || echo 0
        ;;
    "ops_persecond")
        redis-cli -h $REDIS_HOST -a $REDIS_PASSWORD info stats 2>/dev/null | grep "instantaneous_ops_per_sec" | cut -d: -f2 | tr -d '\r' || echo 0
        ;;
    "used_memory_peak")
        redis-cli -h $REDIS_HOST -a $REDIS_PASSWORD info memory 2>/dev/null | grep "used_memory_peak:" | cut -d: -f2 | tr -d '\r' || echo 0
        ;;
    *)
        echo 0
        ;;
esac
