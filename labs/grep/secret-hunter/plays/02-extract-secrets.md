# Play 2: Extract All Secrets from Application Logs
#
# Task: Find all secret patterns in /var/log/app/application.log
# Secrets follow these patterns:
# - SECRET_NUM_XXXX (4-digit number)
# - API_KEY_xXXXXXXX (alphanumeric)
# - TOKEN_xXXXXXXX (alphanumeric)
#
# Save all found secrets to /tmp/secrets_found.txt (one per line)
#
# Hint: Use grep -o to extract only the matching parts
# Hint: Use regex patterns like [A-Z0-9_]+
#
# Expected: At least 3 secrets found
