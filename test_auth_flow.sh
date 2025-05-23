#!/bin/bash

# Test script for Vaani authentication flow
echo "Testing Vaani Authentication Flow"
echo "================================="

PHONE_NUMBER="+911234567890"
BASE_URL="http://localhost:8090/v1"

echo "1. Testing send-otp endpoint..."
SEND_OTP_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/send-otp" \
  -H "Content-Type: application/json" \
  -d "{\"phoneNumber\": \"$PHONE_NUMBER\"}")

echo "Send OTP Response: $SEND_OTP_RESPONSE"

# Check if send-otp was successful
if echo "$SEND_OTP_RESPONSE" | grep -q '"success":true'; then
  echo "✅ Send OTP successful"
else
  echo "❌ Send OTP failed"
  exit 1
fi

echo ""
echo "2. Testing verify-otp endpoint with invalid OTP..."
VERIFY_OTP_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/verify-otp" \
  -H "Content-Type: application/json" \
  -d "{\"phoneNumber\": \"$PHONE_NUMBER\", \"otp\": \"000000\"}")

echo "Verify OTP Response (invalid): $VERIFY_OTP_RESPONSE"

# Check if verify-otp correctly rejects invalid OTP
if echo "$VERIFY_OTP_RESPONSE" | grep -q '"success":false'; then
  echo "✅ Invalid OTP correctly rejected"
else
  echo "❌ Invalid OTP should be rejected"
fi

echo ""
echo "3. Testing health endpoint..."
HEALTH_RESPONSE=$(curl -s "$BASE_URL/../healthz")
echo "Health Response: $HEALTH_RESPONSE"

if [ "$HEALTH_RESPONSE" = "OK" ]; then
  echo "✅ Health check successful"
else
  echo "❌ Health check failed"
fi

echo ""
echo "4. Testing language detection endpoint..."
LANG_RESPONSE=$(curl -s -X POST "$BASE_URL/detect-language" \
  -H "Content-Type: application/json" \
  -d "{\"text\": \"Hello world\"}")

echo "Language Detection Response: $LANG_RESPONSE"

if echo "$LANG_RESPONSE" | grep -q '"language"'; then
  echo "✅ Language detection successful"
else
  echo "❌ Language detection failed"
fi

echo ""
echo "Authentication flow test completed!"
echo "Note: To complete the auth flow, use the OTP from the server logs" 