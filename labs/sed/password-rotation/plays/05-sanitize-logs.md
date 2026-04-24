# Play 5: Sanitize Access Log
#
# Task: In /var/config/access.log, replace all passwords with [REDACTED]
# This is a security best practice - never log passwords!
#
# Examples:
# OLD_ADMIN_PASS_111 -> [REDACTED]
# sk_live_OLD1234567890abcdefghij -> [REDACTED]
# OLD_PASS_12345 -> [REDACTED]
#
# Hint: Match password patterns and replace with literal [REDACTED]
# Hint: Use word boundaries or specific patterns
#
# Expected: All passwords in log replaced with [REDACTED]
