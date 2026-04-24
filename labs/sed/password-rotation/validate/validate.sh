#!/bin/sh
# Validation script for Password Rotation lab

ERRORS=0

# Check Play 1: Database passwords
if grep -q "NEW_PASS_20240115" /etc/services/database.conf && \
   grep -q "NEW_REPLICA_20240115" /etc/services/database.conf && \
   grep -q "NEW_CACHE_20240115" /etc/services/database.conf && \
   ! grep -q "OLD_PASS_12345" /etc/services/database.conf && \
   ! grep -q "OLD_REPLICA_67890" /etc/services/database.conf && \
   ! grep -q "OLD_CACHE_54321" /etc/services/database.conf; then
    echo "✓ Play 1: Database passwords rotated"
else
    echo "✗ Play 1: Database passwords not properly rotated"
    ERRORS=$((ERRORS + 1))
fi

# Check Play 2: API keys
if grep -q "sk_live_NEW1234567890abcdefghij" /app/config/api-keys.conf && \
   grep -q "AKIANEW1234567890A" /app/config/api-keys.conf && \
   grep -q "newsecretkey1234567890abcdefghij" /app/config/api-keys.conf && \
   grep -q "SG.newapikey1234567890abcdef" /app/config/api-keys.conf && \
   grep -q "ghp_newtoken1234567890abcdefghijklmno" /app/config/api-keys.conf && \
   ! grep -q "OLD" /app/config/api-keys.conf && \
   ! grep -q "old" /app/config/api-keys.conf; then
    echo "✓ Play 2: API keys rotated"
else
    echo "✗ Play 2: API keys not properly rotated"
    ERRORS=$((ERRORS + 1))
fi

# Check Play 3: Nginx auth passwords
if grep -q "NEW_ADMIN_PASS_111" /etc/services/nginx-auth.conf && \
   grep -q "NEW_JOHN_PASS_222" /etc/services/nginx-auth.conf && \
   grep -q "NEW_JANE_PASS_333" /etc/services/nginx-auth.conf && \
   ! grep -q "OLD_ADMIN_PASS_111" /etc/services/nginx-auth.conf && \
   ! grep -q "OLD_JOHN_PASS_222" /etc/services/nginx-auth.conf && \
   ! grep -q "OLD_JANE_PASS_333" /etc/services/nginx-auth.conf; then
    echo "✓ Play 3: Nginx auth passwords rotated"
else
    echo "✗ Play 3: Nginx auth passwords not properly rotated"
    ERRORS=$((ERRORS + 1))
fi

# Check Play 4: Application secrets
if grep -q "DB_NEW_SECRET_2024" /var/config/app-secrets.yml && \
   grep -q "ENC_NEW_KEY_12345" /var/config/app-secrets.yml && \
   grep -q "JWT_NEW_SECRET_KEY_67890" /var/config/app-secrets.yml && \
   grep -q "SESSION_NEW_TOKEN_ABCDEFGH" /var/config/app-secrets.yml && \
   grep -q "MASTER_NEW_KEY_99999" /var/config/app-secrets.yml && \
   ! grep -q "OLD" /var/config/app-secrets.yml; then
    echo "✓ Play 4: Application secrets rotated"
else
    echo "✗ Play 4: Application secrets not properly rotated"
    ERRORS=$((ERRORS + 1))
fi

# Check Play 5: Log sanitization
if ! grep -q "OLD_ADMIN_PASS_111" /var/config/access.log && \
   ! grep -q "sk_live_OLD1234567890abcdefghij" /var/config/access.log && \
   ! grep -q "OLD_PASS_12345" /var/config/access.log && \
   ! grep -q "ghp_oldtoken" /var/config/access.log && \
   ! grep -q "OLD_CACHE_54321" /var/config/access.log && \
   grep -q "\[REDACTED\]" /var/config/access.log; then
    echo "✓ Play 5: Access log sanitized"
else
    echo "✗ Play 5: Access log not properly sanitized"
    ERRORS=$((ERRORS + 1))
fi

echo ""
if [ $ERRORS -eq 0 ]; then
    echo "🎉 All plays completed successfully!"
    echo "All passwords have been rotated and logs sanitized."
    exit 0
else
    echo "❌ $ERRORS play(s) failed. Review the errors above."
    exit 1
fi
