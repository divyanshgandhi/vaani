#!/bin/bash

# Vaani Cleanup Script
echo "🧹 Cleaning up Vaani processes and files"
echo "========================================"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Load test session to get PID
if [ -f "/tmp/vaani_test_session" ]; then
    source /tmp/vaani_test_session
    print_status $YELLOW "📂 Found test session with API PID: $API_PID"
    
    if [ -n "$API_PID" ]; then
        print_status $YELLOW "🛑 Stopping API Gateway (PID: $API_PID)..."
        kill $API_PID 2>/dev/null && print_status $GREEN "✅ API Gateway stopped" || print_status $YELLOW "⚠️ API Gateway might already be stopped"
    fi
fi

# Kill any remaining processes
print_status $YELLOW "🔍 Killing any remaining processes..."
killall api-gateway 2>/dev/null && print_status $GREEN "✅ api-gateway processes killed" || true
killall flutter 2>/dev/null && print_status $GREEN "✅ flutter processes killed" || true
killall node 2>/dev/null && print_status $GREEN "✅ node processes killed" || true

# Clean up temporary files
print_status $YELLOW "🗑️ Cleaning up temporary files..."
rm -f /tmp/api-gateway.log && print_status $GREEN "✅ API Gateway logs removed" || true
rm -f /tmp/vaani_test_session && print_status $GREEN "✅ Test session file removed" || true

# Check if any processes are still running
print_status $YELLOW "🔍 Checking for remaining processes..."
REMAINING=$(ps aux | grep -E "(api-gateway|flutter|vaani)" | grep -v grep | grep -v cleanup.sh | wc -l)

if [ "$REMAINING" -gt 0 ]; then
    print_status $YELLOW "⚠️ Some processes might still be running:"
    ps aux | grep -E "(api-gateway|flutter|vaani)" | grep -v grep | grep -v cleanup.sh
    print_status $YELLOW "You may need to kill them manually"
else
    print_status $GREEN "✅ All processes cleaned up successfully"
fi

print_status $GREEN "🎉 Cleanup complete!"
print_status $YELLOW "💡 To start testing again, run: ./test_backend.sh" 