#!/bin/bash

# Continuous Load Test for Go-Insight
# This script generates realistic observability data at high volume

API_URL="http://localhost:8080/api"
API_KEY="your-secure-api-key-here"

# Services to simulate
SERVICES=("user-service" "payment-service" "api-gateway" "notification-service" "order-service" "inventory-service")
LOG_LEVELS=("DEBUG" "INFO" "INFO" "INFO" "WARN" "ERROR" "ERROR")  # Weighted for realistic distribution
HTTP_METHODS=("GET" "POST" "PUT" "DELETE" "PATCH")
ENDPOINTS=("/api/users" "/api/orders" "/api/payments" "/api/notifications" "/api/inventory" "/api/health" "/api/metrics")

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_message() {
    echo -e "${BLUE}[$(date +'%H:%M:%S')]${NC} $1"
}

error_message() {
    echo -e "${RED}[$(date +'%H:%M:%S')] ERROR:${NC} $1"
}

success_message() {
    echo -e "${GREEN}[$(date +'%H:%M:%S')] SUCCESS:${NC} $1"
}

# Generate random log message
generate_log_message() {
    local level=$1
    local service=$2
    
    case $level in
        "ERROR")
            MESSAGES=("Database connection timeout" "Authentication failed" "Payment processing error" "Service unavailable" "Internal server error" "Request validation failed")
            ;;
        "WARN")
            MESSAGES=("High response time detected" "Retry attempt" "Deprecated API usage" "Rate limit approaching" "Cache miss")
            ;;
        "INFO")
            MESSAGES=("Request processed successfully" "User authentication successful" "Cache hit" "Background job completed" "Health check passed")
            ;;
        "DEBUG")
            MESSAGES=("SQL query executed" "Cache lookup" "Method entry" "Variable state" "Configuration loaded")
            ;;
    esac
    
    echo "${MESSAGES[$RANDOM % ${#MESSAGES[@]}]}"
}

# Generate random trace ID
generate_trace_id() {
    echo "$(uuidgen 2>/dev/null || openssl rand -hex 16)"
}

# Generate random span ID  
generate_span_id() {
    echo "$(uuidgen 2>/dev/null || openssl rand -hex 16)"
}

# Send log entry
send_log() {
    local service=${SERVICES[$RANDOM % ${#SERVICES[@]}]}
    local level=${LOG_LEVELS[$RANDOM % ${#LOG_LEVELS[@]}]}
    local message=$(generate_log_message $level $service)
    local trace_id=$(generate_trace_id)
    
    local metadata=""
    case $level in
        "ERROR")
            metadata=',"metadata":{"error_code":"'$(shuf -i 500-599 -n1)'","retry_count":'$(shuf -i 1-5 -n1)',"user_id":"user_'$(shuf -i 1000-9999 -n1)'"}'
            ;;
        "INFO")
            metadata=',"metadata":{"user_id":"user_'$(shuf -i 1000-9999 -n1)'","session_id":"sess_'$(shuf -i 100000-999999 -n1)'","duration_ms":'$(shuf -i 50-500 -n1)'}'
            ;;
        "WARN")
            metadata=',"metadata":{"threshold_value":'$(shuf -i 80-95 -n1)',"current_value":'$(shuf -i 85-100 -n1)'}'
            ;;
    esac
    
    curl -s -X POST "$API_URL/logs" \
        -H "X-API-Key: $API_KEY" \
        -H "Content-Type: application/json" \
        -d "{\"service_name\":\"$service\",\"log_level\":\"$level\",\"message\":\"$message\",\"trace_id\":\"$trace_id\"$metadata}" \
        > /dev/null
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}LOG${NC} [$level] $service: $message"
    else
        error_message "Failed to send log"
    fi
}

# Send metric entry
send_metric() {
    local service=${SERVICES[$RANDOM % ${#SERVICES[@]}]}
    local method=${HTTP_METHODS[$RANDOM % ${#HTTP_METHODS[@]}]}
    local endpoint=${ENDPOINTS[$RANDOM % ${#ENDPOINTS[@]}]}
    
    # Generate realistic status codes (mostly success, some errors)
    local status_codes=(200 200 200 200 200 201 202 400 401 403 404 500 502 503)
    local status_code=${status_codes[$RANDOM % ${#status_codes[@]}]}
    
    # Generate realistic response times based on status code
    local duration_ms
    if [ $status_code -ge 500 ]; then
        duration_ms=$(shuf -i 2000-10000 -n1)  # Slow for server errors
    elif [ $status_code -ge 400 ]; then
        duration_ms=$(shuf -i 100-800 -n1)     # Medium for client errors
    else
        duration_ms=$(shuf -i 20-300 -n1)      # Fast for success
    fi
    
    local frameworks=("gin" "echo" "fiber" "chi")
    local framework=${frameworks[$RANDOM % ${#frameworks[@]}]}
    
    curl -s -X POST "$API_URL/metrics" \
        -H "X-API-Key: $API_KEY" \
        -H "Content-Type: application/json" \
        -d "{\"service_name\":\"$service\",\"path\":\"$endpoint\",\"method\":\"$method\",\"status_code\":$status_code,\"duration_ms\":$duration_ms,\"source\":{\"language\":\"go\",\"framework\":\"$framework\",\"version\":\"1.9.1\"},\"environment\":\"production\"}" \
        > /dev/null
    
    if [ $? -eq 0 ]; then
        local color=$GREEN
        if [ $status_code -ge 500 ]; then
            color=$RED
        elif [ $status_code -ge 400 ]; then
            color=$YELLOW
        fi
        echo -e "${color}METRIC${NC} [$status_code] $service $method $endpoint (${duration_ms}ms)"
    else
        error_message "Failed to send metric"
    fi
}

# Send trace entry
send_trace() {
    local service=${SERVICES[$RANDOM % ${#SERVICES[@]}]}
    local trace_id=$(generate_trace_id)
    
    curl -s -X POST "$API_URL/traces" \
        -H "X-API-Key: $API_KEY" \
        -H "Content-Type: application/json" \
        -d "{\"service_name\":\"$service\",\"id\":\"$trace_id\"}" \
        > /dev/null
    
    if [ $? -eq 0 ]; then
        echo -e "${BLUE}TRACE${NC} Started trace $trace_id for $service"
        
        # 30% chance to end the trace immediately
        if [ $((RANDOM % 10)) -lt 3 ]; then
            sleep $(shuf -i 1-5 -n1)
            curl -s -X POST "$API_URL/traces/$trace_id/end" \
                -H "X-API-Key: $API_KEY" > /dev/null
            echo -e "${BLUE}TRACE${NC} Completed trace $trace_id"
        fi
    else
        error_message "Failed to send trace"
    fi
}

# Statistics tracking
total_requests=0
start_time=$(date +%s)

print_stats() {
    local current_time=$(date +%s)
    local elapsed=$((current_time - start_time))
    local rps=$((total_requests / elapsed))
    
    echo -e "\n${YELLOW}=== LOAD TEST STATS ===${NC}"
    echo -e "Running for: ${elapsed}s"
    echo -e "Total requests: $total_requests"
    echo -e "Requests/second: $rps"
    echo -e "Press Ctrl+C to stop\n"
}

# Trap for cleanup
cleanup() {
    echo -e "\n${YELLOW}Stopping load test...${NC}"
    print_stats
    exit 0
}

trap cleanup SIGINT SIGTERM

# Main load generation loop
main() {
    log_message "Starting continuous load test..."
    log_message "Target: $API_URL"
    log_message "Simulating ${#SERVICES[@]} services with realistic traffic patterns"
    echo ""
    
    # Initial burst to populate dashboard
    log_message "Sending initial data burst..."
    for i in {1..50}; do
        send_log &
        send_metric &
        if [ $((i % 10)) -eq 0 ]; then
            send_trace &
        fi
        total_requests=$((total_requests + 2))
    done
    wait
    success_message "Initial burst complete"
    echo ""
    
    # Continuous load
    while true; do
        # High frequency burst (simulate traffic spike)
        for i in {1..20}; do
            # Weight towards logs and metrics (more common than traces)
            case $((RANDOM % 10)) in
                0|1|2|3|4)  # 50% logs
                    send_log &
                    ;;
                5|6|7|8)    # 40% metrics  
                    send_metric &
                    ;;
                9)          # 10% traces
                    send_trace &
                    ;;
            esac
            total_requests=$((total_requests + 1))
            
            # Small delay to prevent overwhelming
            sleep 0.1
        done
        
        # Wait a bit between bursts
        sleep $(shuf -i 2-8 -n1)
        
        # Print stats every 100 requests
        if [ $((total_requests % 100)) -eq 0 ]; then
            print_stats
        fi
    done
}

# Check if API is accessible
check_api() {
    log_message "Checking API accessibility..."
    if curl -s "$API_URL/health" > /dev/null; then
        success_message "API is accessible"
    else
        error_message "Cannot reach API at $API_URL"
        echo "Make sure Go-Insight is running: docker-compose up"
        exit 1
    fi
}

# Check dependencies
check_dependencies() {
    for cmd in curl uuidgen; do
        if ! command -v $cmd &> /dev/null; then
            error_message "$cmd is required but not installed"
            if [ "$cmd" = "uuidgen" ]; then
                echo "On Ubuntu/Debian: sudo apt-get install uuid-runtime"
                echo "On macOS: uuidgen should be available by default"
            fi
            exit 1
        fi
    done
}

# Run the load test
echo -e "${BLUE}"
echo "╔══════════════════════════════════════╗"
echo "║        Go-Insight Load Tester        ║"
echo "║     Continuous High-Volume Test      ║"
echo "╚══════════════════════════════════════╝"
echo -e "${NC}"

check_dependencies
check_api
main
