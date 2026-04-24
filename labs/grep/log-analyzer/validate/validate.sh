#!/bin/bash

cd /app || exit 1

echo "=== Log Analyzer Validation ==="
echo ""

PASSED=0
FAILED=0

# Check if user created analysis files
if [ -f "analysis.txt" ]; then
    echo "✓ Created analysis.txt"
    ((PASSED++))
else
    echo "✗ analysis.txt not found"
    ((FAILED++))
fi

# Check for error count extraction
if [ -f "errors.txt" ]; then
    ERROR_COUNT=$(wc -l < errors.txt)
    if [ "$ERROR_COUNT" -eq 7 ]; then
        echo "✓ Extracted all 7 error lines"
        ((PASSED++))
    else
        echo "✗ Wrong number of error lines (expected 7, got $ERROR_COUNT)"
        ((FAILED++))
    fi
else
    echo "✗ errors.txt not found"
    ((FAILED++))
fi

# Check for 4xx status codes
if [ -f "client_errors.txt" ]; then
    ERROR_COUNT=$(wc -l < client_errors.txt)
    if [ "$ERROR_COUNT" -eq 6 ]; then
        echo "✓ Extracted all 6 client error requests (4xx)"
        ((PASSED++))
    else
        echo "✗ Wrong number of client errors (expected 6, got $ERROR_COUNT)"
        ((FAILED++))
    fi
else
    echo "✗ client_errors.txt not found"
    ((FAILED++))
fi

# Check for suspicious IP
if [ -f "suspicious.txt" ]; then
    if grep -q "203.0.113.42" suspicious.txt; then
        echo "✓ Identified suspicious IP 203.0.113.42"
        ((PASSED++))
    else
        echo "✗ Suspicious IP 203.0.113.42 not found in suspicious.txt"
        ((FAILED++))
    fi
else
    echo "✗ suspicious.txt not found"
    ((FAILED++))
fi

# Check for failed login attempts
if [ -f "failed_logins.txt" ]; then
    LOGIN_COUNT=$(wc -l < failed_logins.txt)
    if [ "$LOGIN_COUNT" -ge 2 ]; then
        echo "✓ Found failed login attempts"
        ((PASSED++))
    else
        echo "✗ Not enough failed login attempts found"
        ((FAILED++))
    fi
else
    echo "✗ failed_logins.txt not found"
    ((FAILED++))
fi

echo ""
echo "Results: $PASSED passed, $FAILED failed"

if [ $FAILED -eq 0 ]; then
    echo "✓ All validation checks passed!"
    exit 0
else
    echo "✗ Some checks failed"
    exit 1
fi
