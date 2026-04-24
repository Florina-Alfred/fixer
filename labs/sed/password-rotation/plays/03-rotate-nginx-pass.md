# Play 3: Rotate Nginx Auth Passwords
#
# Task: In /etc/services/nginx-auth.conf, replace all OLD_* passwords
# Format: OLD_<USER>_PASS_<NUMBERS> -> NEW_<USER>_PASS_<NUMBERS>
#
# Examples:
# OLD_ADMIN_PASS_111 -> NEW_ADMIN_PASS_111
# OLD_JOHN_PASS_222 -> NEW_JOHN_PASS_222
# OLD_JANE_PASS_333 -> NEW_JANE_PASS_333
#
# Hint: Use sed to match the full pattern
# Hint: Use backreferences: s/OLD_\([A-Z]*\)_PASS_\([0-9]*\)/NEW_\1_PASS_\2/
#
# Expected: 3 passwords rotated
