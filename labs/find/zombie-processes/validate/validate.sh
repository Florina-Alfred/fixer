#!/bin/sh
# Validation script for Zombie Processes lab

ERRORS=0

# Check Play 1: Config files
if [ -f /tmp/config-files.txt ]; then
    CONFIG_COUNT=$(wc -l < /tmp/config-files.txt)
    HAS_YML=$(grep -c "\.yml" /tmp/config-files.txt || echo 0)
    HAS_CONF=$(grep -c "\.conf" /tmp/config-files.txt || echo 0)
    
    if [ "$HAS_YML" -gt 0 ] && [ "$HAS_CONF" -gt 0 ]; then
        echo "✓ Play 1: Found $CONFIG_COUNT config files (yml and conf)"
    else
        echo "✗ Play 1: Missing some config types. YML=$HAS_YML, CONF=$HAS_CONF"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "✗ Play 1: /tmp/config-files.txt not found"
    ERRORS=$((ERRORS + 1))
fi

# Check Play 2: Hidden files
if [ -f /tmp/hidden-files.txt ]; then
    HIDDEN_COUNT=$(wc -l < /tmp/hidden-files.txt)
    
    HAS_BACKUP=$(grep -c "\.backup-creds" /tmp/hidden-files.txt || echo 0)
    HAS_EMERGENCY=$(grep -c "\.emergency-keys" /tmp/hidden-files.txt || echo 0)
    HAS_INTERNAL=$(grep -c "\.internal" /tmp/hidden-files.txt || echo 0)
    
    if [ "$HAS_BACKUP" -gt 0 ] && [ "$HAS_EMERGENCY" -gt 0 ] && [ "$HAS_INTERNAL" -gt 0 ]; then
        echo "✓ Play 2: Found $HIDDEN_COUNT hidden files/directories"
    else
        echo "✗ Play 2: Missing some hidden files. backup=$HAS_BACKUP, emergency=$HAS_EMERGENCY, internal=$HAS_INTERNAL"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "✗ Play 2: /tmp/hidden-files.txt not found"
    ERRORS=$((ERRORS + 1))
fi

# Check Play 3: Small files
if [ -f /tmp/small-files.txt ]; then
    SMALL_COUNT=$(wc -l < /tmp/small-files.txt)
    
    if [ "$SMALL_COUNT" -ge 3 ]; then
        echo "✓ Play 3: Found $SMALL_COUNT small files"
    else
        echo "✗ Play 3: Expected at least 3 small files, found $SMALL_COUNT"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "✗ Play 3: /tmp/small-files.txt not found"
    ERRORS=$((ERRORS + 1))
fi

# Check Play 4: Restricted files
if [ -f /tmp/restricted-files.txt ]; then
    RESTRICTED_COUNT=$(wc -l < /tmp/restricted-files.txt)
    HAS_API_KEYS=$(grep -c "api_keys.conf" /tmp/restricted-files.txt || echo 0)
    
    if [ "$HAS_API_KEYS" -gt 0 ]; then
        echo "✓ Play 4: Found $RESTRICTED_COUNT restricted files (including api_keys.conf)"
    else
        echo "✗ Play 4: api_keys.conf not found in restricted files"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "✗ Play 4: /tmp/restricted-files.txt not found"
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
