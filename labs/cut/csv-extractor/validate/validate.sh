#!/bin/bash

cd /app || exit 1

echo "=== CSV Extractor Validation ==="
echo ""

PASSED=0
FAILED=0

# Check for employee names extraction
if [ -f "employee_names.txt" ]; then
    NAME_COUNT=$(wc -l < employee_names.txt)
    if [ "$NAME_COUNT" -eq 10 ]; then
        echo "✓ Extracted all 10 employee names"
        ((PASSED++))
    else
        echo "✗ Wrong number of names (expected 10, got $NAME_COUNT)"
        ((FAILED++))
    fi
else
    echo "✗ employee_names.txt not found"
    ((FAILED++))
fi

# Check for email extraction
if [ -f "emails.txt" ]; then
    EMAIL_COUNT=$(wc -l < emails.txt)
    if [ "$EMAIL_COUNT" -eq 10 ]; then
        echo "✓ Extracted all 10 employee emails"
        ((PASSED++))
    else
        echo "✗ Wrong number of emails (expected 10, got $EMAIL_COUNT)"
        ((FAILED++))
    fi
else
    echo "✗ emails.txt not found"
    ((FAILED++))
fi

# Check for department extraction
if [ -f "departments.txt" ]; then
    DEPT_COUNT=$(wc -l < departments.txt)
    if [ "$DEPT_COUNT" -eq 10 ]; then
        echo "✓ Extracted all 10 departments"
        ((PASSED++))
    else
        echo "✗ Wrong number of departments (expected 10, got $DEPT_COUNT)"
        ((FAILED++))
    fi
else
    echo "✗ departments.txt not found"
    ((FAILED++))
fi

# Check for customer names (first and last)
if [ -f "customer_names.txt" ]; then
    CUSTOMER_COUNT=$(wc -l < customer_names.txt)
    if [ "$CUSTOMER_COUNT" -eq 8 ]; then
        echo "✓ Extracted all 8 customer names"
        ((PASSED++))
    else
        echo "✗ Wrong number of customers (expected 8, got $CUSTOMER_COUNT)"
        ((FAILED++))
    fi
else
    echo "✗ customer_names.txt not found"
    ((FAILED++))
fi

# Check for phone numbers
if [ -f "phones.txt" ]; then
    PHONE_COUNT=$(wc -l < phones.txt)
    if [ "$PHONE_COUNT" -eq 10 ]; then
        echo "✓ Extracted all 10 phone numbers"
        ((PASSED++))
    else
        echo "✗ Wrong number of phones (expected 10, got $PHONE_COUNT)"
        ((FAILED++))
    fi
else
    echo "✗ phones.txt not found"
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
