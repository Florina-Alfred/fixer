#!/bin/bash

cd /app || exit 1

echo "=== Duplicate Finder Validation ==="
echo ""

PASSED=0
FAILED=0

# Check for duplicate files by size
if [ -f "duplicates.txt" ]; then
    DUP_COUNT=$(wc -l < duplicates.txt)
    if [ "$DUP_COUNT" -ge 2 ]; then
        echo "✓ Found duplicate files"
        ((PASSED++))
    else
        echo "✗ Not enough duplicates found"
        ((FAILED++))
    fi
else
    echo "✗ duplicates.txt not found"
    ((FAILED++))
fi

# Check for empty directories
if [ -f "empty_dirs.txt" ]; then
    DIR_COUNT=$(wc -l < empty_dirs.txt)
    if [ "$DIR_COUNT" -ge 2 ]; then
        echo "✓ Found empty directories"
        ((PASSED++))
    else
        echo "✗ Not enough empty directories found (expected at least 2)"
        ((FAILED++))
    fi
else
    echo "✗ empty_dirs.txt not found"
    ((FAILED++))
fi

# Check for empty files
if [ -f "empty_files.txt" ]; then
    FILE_COUNT=$(wc -l < empty_files.txt)
    if [ "$FILE_COUNT" -ge 2 ]; then
        echo "✓ Found empty files"
        ((PASSED++))
    else
        echo "✗ Not enough empty files found (expected at least 2)"
        ((FAILED++))
    fi
else
    echo "✗ empty_files.txt not found"
    ((FAILED++))
fi

# Check for world-writable files
if [ -f "world_writable.txt" ]; then
    if grep -q "insecure.conf" world_writable.txt; then
        echo "✓ Found world-writable files"
        ((PASSED++))
    else
        echo "✗ World-writable file insecure.conf not found"
        ((FAILED++))
    fi
else
    echo "✗ world_writable.txt not found"
    ((FAILED++))
fi

# Check for SUID files
if [ -f "suid_files.txt" ]; then
    if grep -q "special.bin" suid_files.txt; then
        echo "✓ Found SUID files"
        ((PASSED++))
    else
        echo "✗ SUID file special.bin not found"
        ((FAILED++))
    fi
else
    echo "✗ suid_files.txt not found"
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
