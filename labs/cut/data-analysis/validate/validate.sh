#!/bin/sh
# Validation script for Data Analysis lab

ERRORS=0

# Check Play 1: Email extraction
if [ -f /tmp/emails.txt ]; then
    EMAIL_COUNT=$(wc -l < /tmp/emails.txt)
    HAS_AT=$(grep -c "@" /tmp/emails.txt || echo 0)
    
    if [ "$EMAIL_COUNT" -eq 14 ] && [ "$HAS_AT" -eq 14 ]; then
        echo "✓ Play 1: Extracted 14 email addresses"
    else
        echo "✗ Play 1: Expected 14 emails, found $EMAIL_COUNT (with @: $HAS_AT)"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "✗ Play 1: /tmp/emails.txt not found"
    ERRORS=$((ERRORS + 1))
fi

# Check Play 2: Department and salary
if [ -f /tmp/dept-salary.txt ]; then
    DS_COUNT=$(wc -l < /tmp/dept-salary.txt)
    HAS_COLON=$(grep -c ":" /tmp/dept-salary.txt || echo 0)
    
    if [ "$DS_COUNT" -eq 14 ] && [ "$HAS_COLON" -eq 14 ]; then
        echo "✓ Play 2: Extracted 14 department:salary entries"
    else
        echo "✗ Play 2: Expected 14 entries, found $DS_COUNT (with colon: $HAS_COLON)"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "✗ Play 2: /tmp/dept-salary.txt not found"
    ERRORS=$((ERRORS + 1))
fi

# Check Play 3: Error messages
if [ -f /tmp/errors.txt ]; then
    ERROR_COUNT=$(wc -l < /tmp/errors.txt)
    
    HAS_AUTH=$(grep -c "auth-service" /tmp/errors.txt || echo 0)
    HAS_PAYMENT=$(grep -c "payment" /tmp/errors.txt || echo 0)
    HAS_DATABASE=$(grep -c "database" /tmp/errors.txt || echo 0)
    
    if [ "$ERROR_COUNT" -eq 3 ] && [ "$HAS_AUTH" -gt 0 ] && [ "$HAS_PAYMENT" -gt 0 ] && [ "$HAS_DATABASE" -gt 0 ]; then
        echo "✓ Play 3: Extracted 3 error messages"
    else
        echo "✗ Play 3: Expected 3 errors, found $ERROR_COUNT"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "✗ Play 3: /tmp/errors.txt not found"
    ERRORS=$((ERRORS + 1))
fi

# Check Play 4: Network traffic
if [ -f /tmp/network-traffic.txt ]; then
    NET_COUNT=$(wc -l < /tmp/network-traffic.txt)
    
    if [ "$NET_COUNT" -eq 10 ]; then
        echo "✓ Play 4: Extracted 10 network traffic entries"
    else
        echo "✗ Play 4: Expected 10 entries, found $NET_COUNT"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "✗ Play 4: /tmp/network-traffic.txt not found"
    ERRORS=$((ERRORS + 1))
fi

# Check Play 5: Sales by region
if [ -f /tmp/sales-by-region.txt ]; then
    SALES_COUNT=$(wc -l < /tmp/sales-by-region.txt)
    
    if [ "$SALES_COUNT" -eq 10 ]; then
        echo "✓ Play 5: Extracted 10 sales entries"
    else
        echo "✗ Play 5: Expected 10 entries, found $SALES_COUNT"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "✗ Play 5: /tmp/sales-by-region.txt not found"
    ERRORS=$((ERRORS + 1))
fi

echo ""
if [ $ERRORS -eq 0 ]; then
    echo "🎉 All plays completed successfully!"
    exit 0
else
    echo "❌ $ERRORS play(s) failed. Review the errors above."
    exit 1
fi
