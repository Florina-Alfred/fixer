# Play 4: Rotate Application Secrets
#
# Task: In /var/config/app-secrets.yml, replace all OLD_* secrets with NEW_* secrets
# Format: OLD -> NEW prefix, keep the rest
#
# Examples:
# DB_OLD_SECRET_2024 -> DB_NEW_SECRET_2024
# ENC_OLD_KEY_12345 -> ENC_NEW_KEY_12345
# JWT_OLD_SECRET_KEY_67890 -> JWT_NEW_SECRET_KEY_67890
# SESSION_OLD_TOKEN_ABCDEFGH -> SESSION_NEW_TOKEN_ABCDEFGH
# MASTER_OLD_KEY_99999 -> MASTER_NEW_KEY_99999
#
# Hint: Use sed with -i flag for in-place editing
# Hint: Pattern: s/OLD_/NEW_/g
#
# Expected: 5 secrets rotated
