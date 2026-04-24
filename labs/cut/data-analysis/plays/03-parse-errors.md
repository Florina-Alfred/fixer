# Play 3: Parse Server Logs - Extract Error Messages
#
# Task: From /data/server-logs.txt, extract all ERROR level messages
# Save the SERVICE name and MESSAGE to /tmp/errors.txt
# Format: SERVICE:MESSAGE
#
# Hint: The file uses | as delimiter
# Hint: Fields are: TIMESTAMP|LEVEL|SERVICE|MESSAGE|USER_ID|RESPONSE_TIME
# Hint: First filter for ERROR lines, then extract fields
#
# Expected: 3 error entries
