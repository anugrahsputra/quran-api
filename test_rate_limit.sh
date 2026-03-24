#!/bin/bash

# Configuration
API_URL="https://quran-api.downormal.dev/api/v1/surah"
TOTAL_REQUESTS=150
PARALLEL_JOBS=20 # Send 20 at a time

echo "Starting STAGGERED rate limit test against: $API_URL"
echo "Sending $TOTAL_REQUESTS requests (20 at a time)..."
echo "--------------------------------------------------"

# Function to perform a single request
do_request() {
    local i=$1
    # Capture only the status code for cleaner output
    status_code=$(curl -sL -o /dev/null -w "%{http_code}" "$API_URL")
    
    if [ "$status_code" -eq 200 ]; then
        echo -n "."
    elif [ "$status_code" -eq 429 ]; then
        echo -e "\n[Request $i] RATE LIMITED! (429)"
    else
        echo -e "\n[Request $i] Error: $status_code"
    fi
}

export -f do_request
export API_URL

# Using seq and xargs to send requests with a maximum concurrency of 20
seq 1 $TOTAL_REQUESTS | xargs -n 1 -P $PARALLEL_JOBS -I {} bash -c "do_request {}"

echo -e "\n--------------------------------------------------"
echo "Test complete."
