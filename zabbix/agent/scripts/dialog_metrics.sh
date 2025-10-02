#!/bin/bash

MS_ENDPOINT="ms-dialog:${SERVER_PORT:-8080}"

case $1 in
    "active_users")
        curl -s --connect-timeout 5 "http://$MS_ENDPOINT/metrics" 2>/dev/null | grep "dialog_active_users" | awk '{print $2}' || echo 0
        ;;
    "total_messages")
        curl -s --connect-timeout 5 "http://$MS_ENDPOINT/metrics" 2>/dev/null | grep "dialog_messages_total" | awk '{print $2}' || echo 0
        ;;
    "active_sessions")
        curl -s --connect-timeout 5 "http://$MS_ENDPOINT/metrics" 2>/dev/null | grep "dialog_active_sessions" | awk '{print $2}' || echo 0
        ;;
    "response_time")
        curl -s --connect-timeout 5 "http://$MS_ENDPOINT/metrics" 2>/dev/null | grep "dialog_response_time_seconds" | awk '{print $2}' || echo 0
        ;;
    "health_check")
        curl -s --connect-timeout 5 -f "http://$MS_ENDPOINT/health" >/dev/null 2>&1 && echo 1 || echo 0
        ;;
    *)
        echo 0
        ;;
esac
