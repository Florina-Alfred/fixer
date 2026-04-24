# Play 5: Regional Sales Summary
#
# Task: From /sales/monthly-sales.csv, calculate total sales per region
# Sum up QUANTITY * UNIT_PRICE * (1 - DISCOUNT) for each region
# Save to /tmp/regional-summary.txt with format: REGION|TOTAL_REVENUE|SALE_COUNT
#
# Hint: Use awk associative arrays to accumulate by region
# Hint: Fields: region($3), qty($5), price($6), discount($7)
# Hint: Use END block with for loop to print results
#
# Expected: 4 regional summaries (North, South, East, West)
