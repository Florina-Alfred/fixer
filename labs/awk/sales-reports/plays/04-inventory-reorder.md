# Play 4: Inventory Analysis
#
# Task: From /sales/inventory.csv, find products that need reordering
# A product needs reordering if QTY_ON_HAND <= REORDER_POINT
# Save to /tmp/reorder-needed.txt with format: PRODUCT|WAREHOUSE|QTY_ON_HAND|REORDER_POINT
#
# Hint: Use awk pattern/action with comparison
# Hint: Fields: product($1), warehouse($2), qty_on_hand($3), reorder_point($5)
#
# Expected: Multiple products needing reorder
