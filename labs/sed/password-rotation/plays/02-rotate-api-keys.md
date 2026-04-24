# Play 2: Rotate API Keys
#
# Task: In /app/config/api-keys.conf, replace all OLD* keys with NEW* keys
# Format: OLD -> NEW prefix
#
# Examples:
# sk_live_OLD1234567890abcdefghij -> sk_live_NEW1234567890abcdefghij
# AKIAOLD1234567890A -> AKIANEW1234567890A
# oldsecretkey1234567890abcdefghij -> newsecretkey1234567890abcdefghij
# SG.oldapikey1234567890abcdef.ghijklmnop -> SG.newapikey1234567890abcdef.ghijklmnop
# ghp_oldtoken1234567890abcdefghijklmno -> ghp_newtoken1234567890abcdefghijklmno
#
# Hint: Use sed with multiple -e flags or combine patterns
# Hint: Be careful with case sensitivity (OLD vs old)
#
# Expected: 5 API keys rotated
