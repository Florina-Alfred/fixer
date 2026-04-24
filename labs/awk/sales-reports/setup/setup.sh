#!/bin/sh
# Setup script for Sales Reports lab
# Creates realistic sales data for awk exercises

mkdir -p /sales /reports /exports

# Create monthly sales data
# Format: DATE|SALES_REP|REGION|PRODUCT|QUANTITY|UNIT_PRICE|DISCOUNT|COMMISSION_RATE
cat > /sales/monthly-sales.csv << 'EOF'
DATE,SALES_REP,REGION,PRODUCT,QUANTITY,UNIT_PRICE,DISCOUNT,COMMISSION_RATE
2024-01-05,Alice Johnson,North,Widget A,50,25.00,0.10,0.05
2024-01-08,Bob Smith,South,Widget B,75,30.00,0.05,0.05
2024-01-12,Carol White,East,Widget C,30,45.00,0.15,0.07
2024-01-15,David Brown,West,Widget A,60,25.00,0.00,0.05
2024-01-18,Eve Davis,North,Widget B,45,30.00,0.10,0.05
2024-01-22,Frank Miller,South,Widget C,40,45.00,0.00,0.07
2024-01-25,Grace Lee,East,Widget A,55,25.00,0.05,0.05
2024-01-28,Henry Wilson,West,Widget B,70,30.00,0.10,0.05
2024-02-02,Alice Johnson,North,Widget C,35,45.00,0.00,0.07
2024-02-05,Bob Smith,South,Widget A,65,25.00,0.10,0.05
2024-02-08,Carol White,East,Widget B,80,30.00,0.05,0.05
2024-02-12,David Brown,West,Widget C,25,45.00,0.20,0.07
2024-02-15,Eve Davis,North,Widget A,70,25.00,0.00,0.05
2024-02-18,Frank Miller,South,Widget B,50,30.00,0.10,0.05
2024-02-22,Grace Lee,East,Widget C,45,45.00,0.00,0.07
2024-02-25,Henry Wilson,West,Widget A,85,25.00,0.05,0.05
2024-02-28,Alice Johnson,North,Widget B,60,30.00,0.15,0.05
2024-03-02,Bob Smith,South,Widget C,55,45.00,0.10,0.07
2024-03-05,Carol White,East,Widget A,40,25.00,0.00,0.05
2024-03-08,David Brown,West,Widget B,90,30.00,0.05,0.05
EOF

# Create customer orders
# Format: ORDER_ID|CUSTOMER|PRODUCT|QUANTITY|PRICE|STATUS|DATE
cat > /sales/orders.log << 'EOF'
ORD001|Acme Corp|Widget A|100|25.00|SHIPPED|2024-01-10
ORD002|TechStart Inc|Widget B|50|30.00|PENDING|2024-01-12
ORD003|Global Solutions|Widget C|75|45.00|SHIPPED|2024-01-15
ORD004|StartupXYZ|Widget A|25|25.00|CANCELLED|2024-01-18
ORD005|MegaCorp|Widget B|200|30.00|SHIPPED|2024-01-20
ORD006|SmallBiz LLC|Widget C|30|45.00|PENDING|2024-01-22
ORD007|Enterprise Co|Widget A|150|25.00|SHIPPED|2024-01-25
ORD008|Local Shop|Widget B|20|30.00|SHIPPED|2024-01-28
ORD009|BigRetail|Widget C|100|45.00|PENDING|2024-02-01
ORD010|Medium Inc|Widget A|80|25.00|SHIPPED|2024-02-05
ORD011|TinyCo|Widget B|15|30.00|CANCELLED|2024-02-08
ORD012|HugeCorp|Widget C|120|45.00|SHIPPED|2024-02-10
ORD013|Startup2|Widget A|35|25.00|PENDING|2024-02-12
ORD014|MegaRetail|Widget B|180|30.00|SHIPPED|2024-02-15
ORD015|LocalBiz|Widget C|40|45.00|SHIPPED|2024-02-18
EOF

# Create inventory data
# Format: PRODUCT|WAREHOUSE|QTY_ON_HAND|QTY_ORDERED|REORDER_POINT|UNIT_COST
cat > /sales/inventory.csv << 'EOF'
PRODUCT,WAREHOUSE,QTY_ON_HAND,QTY_ORDERED,REORDER_POINT,UNIT_COST
Widget A,North-Warehouse,500,200,100,15.00
Widget A,South-Warehouse,350,150,100,15.00
Widget A,East-Warehouse,420,100,100,15.00
Widget A,West-Warehouse,280,250,100,15.00
Widget B,North-Warehouse,600,100,150,18.00
Widget B,South-Warehouse,450,200,150,18.00
Widget B,East-Warehouse,380,150,150,18.00
Widget B,West-Warehouse,520,0,150,18.00
Widget C,North-Warehouse,200,100,50,28.00
Widget C,South-Warehouse,150,150,50,28.00
Widget C,East-Warehouse,180,50,50,28.00
Widget C,West-Warehouse,220,100,50,28.00
EOF

# Create employee commission data
# Format: EMPLOYEE_ID|NAME|REGION|TOTAL_SALES|COMMISSION_RATE|BONUS
cat > /sales/employees.csv << 'EOF'
EMPLOYEE_ID,NAME,REGION,TOTAL_SALES,COMMISSION_RATE,BONUS
E001,Alice Johnson,North,45000.00,0.05,2500.00
E002,Bob Smith,South,52000.00,0.05,3000.00
E003,Carol White,East,38000.00,0.07,1500.00
E004,David Brown,West,48000.00,0.05,2000.00
E005,Eve Davis,North,41000.00,0.05,1800.00
E006,Frank Miller,South,36000.00,0.07,1200.00
E007,Grace Lee,East,43000.00,0.05,2200.00
E008,Henry Wilson,West,55000.00,0.05,3500.00
EOF

# Create a product catalog
# Format: PRODUCT_ID|NAME|CATEGORY|BASE_PRICE|COST|PROFIT_MARGIN
cat > /sales/products.csv << 'EOF'
PRODUCT_ID,NAME,CATEGORY,BASE_PRICE,COST,PROFIT_MARGIN
WGT-A,Widget A,Standard,25.00,15.00,0.40
WGT-B,Widget B,Standard,30.00,18.00,0.40
WGT-C,Widget C,Premium,45.00,28.00,0.38
WGT-D,Widget D,Premium,60.00,38.00,0.37
WGT-E,Widget E,Deluxe,80.00,48.00,0.40
EOF

echo "Setup complete for Sales Reports lab"
echo "Sales data created in /sales/"
ls -la /sales/
