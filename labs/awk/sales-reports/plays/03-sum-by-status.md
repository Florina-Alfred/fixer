# Play 3: Sum Orders by Status
#
# Task: From /sales/orders.log, calculate total value for each order status
# Total = QUANTITY * PRICE
# Save to /tmp/orders-by-status.txt with format: STATUS|TOTAL_VALUE
#
# Hint: Use awk associative arrays to accumulate totals
# Hint: Fields: order_id($1), customer($2), product($3), qty($4), price($5), status($6)
# Hint: Use END block to print accumulated totals
#
# Expected: 3 status totals (SHIPPED, PENDING, CANCELLED)
