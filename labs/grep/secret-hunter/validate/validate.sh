#!/bin/sh
# Validation script for Secret Hunter lab
# Validates that the learner has correctly extracted all secrets

ERRORS=0

# Check Play 1: AUTH_CODE extraction
if [ -f /tmp/auth-code.txt ]; then
    AUTH_CODE=$(cat /tmp/auth-code.txt | tr -d '[:space:]')
    if [ "$AUTH_CODE" = "AUTH_CODE_9182" ]; then
        echo "✓ Play 1: AUTH_CODE correctly extracted"
    else
        echo "✗ Play 1: AUTH_CODE incorrect. Expected: AUTH_CODE_9182, Got: $AUTH_CODE"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "✗ Play 1: /tmp/auth-code.txt not found"
    ERRORS=$((ERRORS + 1))
fi

# Check Play 2: All secrets extraction
if [ -f /tmp/secrets_found.txt ]; then
    SECRET_COUNT=$(wc -l < /tmp/secrets_found.txt)
    
    # Check for each expected secret
    HAS_SECRET_NUM=$(grep -c "SECRET_NUM_7429" /tmp/secrets_found.txt || echo 0)
    HAS_API_KEY=$(grep -c "API_KEY_xK9mP2nQ" /tmp/secrets_found.txt || echo 0)
    HAS_TOKEN=$(grep -c "TOKEN_aB3cD4eF5g" /tmp/secrets_found.txt || echo 0)
    
    if [ "$HAS_SECRET_NUM" -gt 0 ] && [ "$HAS_API_KEY" -gt 0 ] && [ "$HAS_TOKEN" -gt 0 ]; then
        echo "✓ Play 2: All secrets correctly extracted ($SECRET_COUNT found)"
    else
        echo "✗ Play 2: Missing some secrets. Found: SECRET_NUM=$HAS_SECRET_NUM, API_KEY=$HAS_API_KEY, TOKEN=$HAS_TOKEN"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "✗ Play 2: /tmp/secrets_found.txt not found"
    ERRORS=$((ERRORS + 1))
fi

# Check Play 3: MASTER_KEY extraction
if [ -f /tmp/master-key.txt ]; then
    if grep -q "MASTER_KEY_3847" /tmp/master-key.txt; then
        echo "✓ Play 3: MASTER_KEY correctly found"
    else
        echo "✗ Play 3: MASTER_KEY incorrect or missing"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "✗ Play 3: /tmp/master-key.txt not found"
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
