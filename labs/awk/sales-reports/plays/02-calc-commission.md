# Play 2: Calculate Commissions
#
# Task: From /sales/monthly-sales.csv, calculate commission for each sale
# Commission = QUANTITY * UNIT_PRICE * (1 - DISCOUNT) * COMMISSION_RATE
# Save to /tmp/commissions.txt with format: SALES_REP|COMMISSION
#
# Hint: Use awk with -F',' to set field separator
# Hint: Fields: rep($2), qty($5), price($6), discount($7), commission_rate($8)
#
# Expected: 20 commission entries
