#!/bin/sh
# Validation script for Sales Reports lab

ERRORS=0

# Check Play 1: Revenue calculation
if [ -f /tmp/revenue.txt ]; then
    REV_COUNT=$(wc -l < /tmp/revenue.txt)
    
    if [ "$REV_COUNT" -eq 20 ]; then
        echo "✓ Play 1: Calculated 20 revenue entries"
    else
        echo "✗ Play 1: Expected 20 entries, found $REV_COUNT"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "✗ Play 1: /tmp/revenue.txt not found"
    ERRORS=$((ERRORS + 1))
fi

# Check Play 2: Commission calculation
if [ -f /tmp/commissions.txt ]; then
    COM_COUNT=$(wc -l < /tmp/commissions.txt)
    
    if [ "$COM_COUNT" -eq 20 ]; then
        echo "✓ Play 2: Calculated 20 commission entries"
    else
        echo "✗ Play 2: Expected 20 entries, found $COM_COUNT"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "✗ Play 2: /tmp/commissions.txt not found"
    ERRORS=$((ERRORS + 1))
fi

# Check Play 3: Orders by status
if [ -f /tmp/orders-by-status.txt ]; then
    STATUS_COUNT=$(wc -l < /tmp/orders-by-status.txt)
    HAS_SHIPPED=$(grep -c "SHIPPED" /tmp/orders-by-status.txt || echo 0)
    HAS_PENDING=$(grep -c "PENDING" /tmp/orders-by-status.txt || echo 0)
    HAS_CANCELLED=$(grep -c "CANCELLED" /tmp/orders-by-status.txt || echo 0)
    
    if [ "$STATUS_COUNT" -eq 3 ] && [ "$HAS_SHIPPED" -gt 0 ] && [ "$HAS_PENDING" -gt 0 ] && [ "$HAS_CANCELLED" -gt 0 ]; then
        echo "✓ Play 3: Calculated totals for all 3 order statuses"
    else
        echo "✗ Play 3: Expected 3 status totals, found $STATUS_COUNT"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "✗ Play 3: /tmp/orders-by-status.txt not found"
    ERRORS=$((ERRORS + 1))
fi

# Check Play 4: Inventory reorder
if [ -f /tmp/reorder-needed.txt ]; then
    REORDER_COUNT=$(wc -l < /tmp/reorder-needed.txt)
    
    if [ "$REORDER_COUNT" -ge 1 ]; then
        echo "✓ Play 4: Found $REORDER_COUNT products needing reorder"
    else
        echo "✗ Play 4: Expected at least 1 product needing reorder, found $REORDER_COUNT"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "✗ Play 4: /tmp/reorder-needed.txt not found"
    ERRORS=$((ERRORS + 1))
fi

# Check Play 5: Regional summary
if [ -f /tmp/regional-summary.txt ]; then
    REGION_COUNT=$(wc -l < /tmp/regional-summary.txt)
    HAS_NORTH=$(grep -c "North" /tmp/regional-summary.txt || echo 0)
    HAS_SOUTH=$(grep -c "South" /tmp/regional-summary.txt || echo 0)
    HAS_EAST=$(grep -c "East" /tmp/regional-summary.txt || echo 0)
    HAS_WEST=$(grep -c "West" /tmp/regional-summary.txt || echo 0)
    
    if [ "$REGION_COUNT" -eq 4 ] && [ "$HAS_NORTH" -gt 0 ] && [ "$HAS_SOUTH" -gt 0 ] && [ "$HAS_EAST" -gt 0 ] && [ "$HAS_WEST" -gt 0 ]; then
        echo "✓ Play 5: Calculated summaries for all 4 regions"
    else
        echo "✗ Play 5: Expected 4 regional summaries, found $REGION_COUNT"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "✗ Play 5: /tmp/regional-summary.txt not found"
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
