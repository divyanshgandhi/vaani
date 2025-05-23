#!/bin/bash

# Vaani Frontend Testing Script
echo "📱 Starting Vaani Frontend Testing"
echo "=================================="

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Load test session if available
if [ -f "/tmp/vaani_test_session" ]; then
    source /tmp/vaani_test_session
    print_status $GREEN "📂 Loaded test session:"
    print_status $BLUE "   📞 Phone: $PHONE_NUMBER"
    print_status $BLUE "   🆔 Project: $PROJECT_ID"
    print_status $BLUE "   🔧 Backend PID: $API_PID"
else
    print_status $YELLOW "⚠️ No test session found. Run ./test_backend.sh first"
    PHONE_NUMBER="+919871576008"
fi

# Check if backend is running
print_status $YELLOW "🔍 Checking backend status..."
if curl -s "http://localhost:8090/healthz" > /dev/null 2>&1; then
    print_status $GREEN "✅ Backend is running"
else
    print_status $RED "❌ Backend is not running"
    print_status $YELLOW "🚀 Starting backend first..."
    ./test_backend.sh
    if [ $? -ne 0 ]; then
        print_status $RED "❌ Failed to start backend"
        exit 1
    fi
fi

# Navigate to frontend
print_status $YELLOW "📂 Navigating to frontend directory..."
cd frontend/app || {
    print_status $RED "❌ Frontend directory not found"
    exit 1
}

# Check Flutter installation
print_status $YELLOW "🔍 Checking Flutter installation..."
if ! command -v flutter &> /dev/null; then
    print_status $RED "❌ Flutter is not installed"
    exit 1
fi

# Get packages
print_status $YELLOW "📦 Getting Flutter packages..."
flutter pub get || {
    print_status $RED "❌ Failed to get Flutter packages"
    exit 1
}

print_status $GREEN "✅ Flutter setup complete"

# Display testing instructions
print_status $BLUE "📋 TESTING INSTRUCTIONS"
print_status $BLUE "======================="
print_status $YELLOW "🔐 AUTHENTICATION:"
print_status $WHITE "   1. Use phone number: $PHONE_NUMBER"
print_status $WHITE "   2. The OTP will be sent to your actual phone"
print_status $WHITE "   3. If SMS doesn't work, check backend logs:"
print_status $WHITE "      tail -f /tmp/api-gateway.log | grep OTP"

print_status $YELLOW "📁 PROJECT TESTING:"
print_status $WHITE "   1. After login, you should see the home screen"
print_status $WHITE "   2. Try creating a new project"
print_status $WHITE "   3. Test project listing and navigation"
print_status $WHITE "   4. Try editing/deleting projects"

print_status $YELLOW "🐛 DEBUGGING:"
print_status $WHITE "   • Backend logs: tail -f /tmp/api-gateway.log"
print_status $WHITE "   • Flutter logs: Check terminal output below"
print_status $WHITE "   • API endpoint: http://localhost:8090/v1"

print_status $YELLOW "🛑 TO STOP:"
print_status $WHITE "   • Press Ctrl+C in this terminal to stop Flutter"
print_status $WHITE "   • Run: kill $API_PID to stop backend"

print_status $GREEN "🚀 Starting Flutter app..."
print_status $YELLOW "   (This will open the app in your connected device/simulator)"

# Start Flutter app
flutter run

# Clean up message
print_status $YELLOW "📱 Flutter app stopped"
print_status $BLUE "Backend is still running. Use 'kill $API_PID' to stop it." 