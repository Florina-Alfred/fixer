# Play 1: Rotate Database Passwords
#
# Task: In /etc/services/database.conf, replace all OLD_* passwords with NEW_* passwords
# Format: OLD_<SERVICE>_<NUMBERS> -> NEW_<SERVICE>_<TIMESTAMP>
# Use timestamp: 20240115
#
# Examples:
# OLD_PASS_12345 -> NEW_PASS_20240115
# OLD_REPLICA_67890 -> NEW_REPLICA_20240115
# OLD_CACHE_54321 -> NEW_CACHE_20240115
#
# Hint: Use sed -i for in-place editing
# Hint: Use sed 's/OLD_/NEW_/g' pattern
#
# Expected: 3 passwords rotated
