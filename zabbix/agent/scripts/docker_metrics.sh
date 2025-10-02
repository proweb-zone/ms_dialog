#!/bin/bash

case $1 in
    "container_status")
        CONTAINER=$2
        docker inspect --format='{{.State.Status}}' $CONTAINER 2>/dev/null || echo "not_found"
        ;;
    "container_cpu")
        CONTAINER=$2
        docker stats --no-stream --format "{{.CPUPerc}}" $CONTAINER 2>/dev/null | head -1 | sed 's/%//' || echo 0
        ;;
    "container_memory")
        CONTAINER=$2
        docker stats --no-stream --format "{{.MemUsage}}" $CONTAINER 2>/dev/null | head -1 | cut -d'/' -f1 | sed 's/[^0-9.]//g' || echo 0
        ;;
    "container_memory_limit")
        CONTAINER=$2
        docker stats --no-stream --format "{{.MemUsage}}" $CONTAINER 2>/dev/null | head -1 | cut -d'/' -f2 | sed 's/[^0-9.]//g' || echo 0
        ;;
    "container_network_io")
        CONTAINER=$2
        docker stats --no-stream --format "{{.NetIO}}" $CONTAINER 2>/dev/null | head -1 || echo "0B / 0B"
        ;;
    "containers_running")
        docker ps -q | wc -l | tr -d ' '
        ;;
    "containers_total")
        docker ps -a -q | wc -l | tr -d ' '
        ;;
    *)
        echo 0
        ;;
esac
