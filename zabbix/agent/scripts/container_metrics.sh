#!/bin/bash

CONTAINER_NAME=$2

case $1 in
    "cpu_usage")
        docker stats --no-stream --format "{{.CPUPerc}}" $CONTAINER_NAME 2>/dev/null | head -1 | sed 's/%//' || echo 0
        ;;
    "memory_usage")
        docker stats --no-stream --format "{{.MemUsage}}" $CONTAINER_NAME 2>/dev/null | head -1 | cut -d'/' -f1 | sed 's/[^0-9.]//g' || echo 0
        ;;
    "memory_limit")
        docker stats --no-stream --format "{{.MemUsage}}" $CONTAINER_NAME 2>/dev/null | head -1 | cut -d'/' -f2 | sed 's/[^0-9.]//g' || echo 0
        ;;
    "disk_io")
        docker stats --no-stream --format "{{.BlockIO}}" $CONTAINER_NAME 2>/dev/null | head -1 || echo "0B / 0B"
        ;;
    *)
        echo 0
        ;;
esac
