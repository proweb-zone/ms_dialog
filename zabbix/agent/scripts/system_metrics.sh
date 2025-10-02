#!/bin/bash

case $1 in
    "cpu_load")
        cat /proc/loadavg | awk '{print $1}' 2>/dev/null || echo 0
        ;;
    "memory_used")
        free -b | awk 'NR==2{print $3}' 2>/dev/null || echo 0
        ;;
    "memory_total")
        free -b | awk 'NR==2{print $2}' 2>/dev/null || echo 0
        ;;
    "disk_used")
        MOUNT_POINT=$2
        df -B1 $MOUNT_POINT 2>/dev/null | awk 'NR==2{print $3}' || echo 0
        ;;
    "disk_total")
        MOUNT_POINT=$2
        df -B1 $MOUNT_POINT 2>/dev/null | awk 'NR==2{print $2}' || echo 0
        ;;
    *)
        echo 0
        ;;
esac
