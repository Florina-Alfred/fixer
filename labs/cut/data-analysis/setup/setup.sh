#!/bin/sh
# Setup script for Data Analysis lab
# Creates realistic CSV and data files for cut exercises

mkdir -p /data /reports /exports

# Create a realistic employee database
cat > /data/employees.csv << 'EOF'
employee_id,first_name,last_name,email,department,salary,hire_date,manager_id
E001,John,Doe,john.doe@company.com,Engineering,75000,2020-03-15,E100
E002,Jane,Smith,jane.smith@company.com,Marketing,68000,2019-07-22,E101
E003,Bob,Johnson,bob.johnson@company.com,Sales,72000,2021-01-10,E102
E004,Alice,Brown,alice.brown@company.com,Engineering,82000,2018-11-05,E100
E005,Charlie,Wilson,charlie.wilson@company.com,HR,65000,2020-09-18,E103
E006,Diana,Miller,diana.miller@company.com,Sales,70000,2019-04-30,E102
E007,Eve,Davis,eve.davis@company.com,Engineering,78000,2021-06-12,E100
E008,Frank,Garcia,frank.garcia@company.com,Marketing,66000,2020-02-28,E101
E009,Grace,Martinez,grace.martinez@company.com,HR,67000,2019-12-01,E103
E010,Henry,Anderson,henry.anderson@company.com,Sales,74000,2018-08-20,E102
E100,Michael,Chen,michael.chen@company.com,Engineering,95000,2015-01-10,
E101,Sarah,Lee,sarah.lee@company.com,Marketing,92000,2015-06-15,
E102,David,Kim,david.kim@company.com,Sales,90000,2014-09-01,
E103,Emily,Taylor,emily.taylor@company.com,HR,88000,2016-03-22,
EOF

# Create a server log in a custom format
# Format: TIMESTAMP|LEVEL|SERVICE|MESSAGE|USER_ID|RESPONSE_TIME
cat > /data/server-logs.txt << 'EOF'
2024-01-15T10:00:00|INFO|auth-service|User login successful|U1001|45ms
2024-01-15T10:00:05|DEBUG|api-gateway|Request received|U1002|12ms
2024-01-15T10:00:10|ERROR|database|Connection timeout|U1003|5000ms
2024-01-15T10:00:15|WARN|cache|Cache miss rate high|U1001|230ms
2024-01-15T10:00:20|INFO|auth-service|Token refreshed|U1004|67ms
2024-01-15T10:00:25|ERROR|payment|Transaction failed|U1005|3200ms
2024-01-15T10:00:30|INFO|api-gateway|Health check passed|SYSTEM|5ms
2024-01-15T10:00:35|DEBUG|database|Query executed|U1002|156ms
2024-01-15T10:00:40|ERROR|auth-service|Invalid credentials|U1006|23ms
2024-01-15T10:00:45|INFO|cache|Cache cleared|ADMIN|89ms
2024-01-15T10:00:50|WARN|api-gateway|Rate limit approaching|U1007|34ms
2024-01-15T10:00:55|INFO|payment|Refund processed|U1008|445ms
2024-01-15T10:01:00|ERROR|database|Deadlock detected|U1009|1200ms
2024-01-15T10:01:05|INFO|auth-service|User logout|U1001|12ms
2024-01-15T10:01:10|DEBUG|cache|Cache hit|U1010|3ms
EOF

# Create a network traffic log
# Format: SRC_IP:PORT -> DST_IP:PORT | PROTOCOL | BYTES | STATUS
cat > /data/network.log << 'EOF'
192.168.1.100:45123 -> 10.0.0.5:443 | HTTPS | 15234 | OK
192.168.1.101:52341 -> 10.0.0.5:80 | HTTP | 2341 | OK
10.0.0.99:12345 -> 10.0.0.5:22 | SSH | 567 | BLOCKED
192.168.1.102:33456 -> 10.0.0.10:3306 | MySQL | 8923 | OK
172.16.0.50:44567 -> 10.0.0.5:443 | HTTPS | 45123 | OK
192.168.1.100:45124 -> 10.0.0.5:443 | HTTPS | 12456 | TIMEOUT
10.0.0.100:55678 -> 10.0.0.5:8080 | HTTP | 3456 | OK
192.168.1.103:23456 -> 10.0.0.15:5432 | PostgreSQL | 6789 | OK
172.16.0.51:34567 -> 10.0.0.5:443 | HTTPS | 23456 | OK
192.168.1.104:45678 -> 10.0.0.5:80 | HTTP | 1234 | BLOCKED
EOF

# Create a sales report
# Format: DATE|REGION|PRODUCT|QUANTITY|UNIT_PRICE|TOTAL|SALES_REP
cat > /data/sales-report.csv << 'EOF'
DATE,REGION,PRODUCT,QUANTITY,UNIT_PRICE,TOTAL,SALES_REP
2024-01-01,North,Widget A,150,25.00,3750.00,John Smith
2024-01-01,South,Widget B,200,30.00,6000.00,Jane Doe
2024-01-01,East,Widget A,175,25.00,4375.00,Bob Johnson
2024-01-02,North,Widget C,100,45.00,4500.00,Alice Brown
2024-01-02,South,Widget A,225,25.00,5625.00,Charlie Wilson
2024-01-02,West,Widget B,180,30.00,5400.00,Diana Miller
2024-01-03,East,Widget C,125,45.00,5625.00,Eve Davis
2024-01-03,North,Widget B,190,30.00,5700.00,Frank Garcia
2024-01-03,South,Widget C,160,45.00,7200.00,Grace Martinez
2024-01-04,West,Widget A,210,25.00,5250.00,Henry Anderson
EOF

# Create an IP address list with metadata
# Format: IP_ADDRESS|COUNTRY|CITY|ISP|TYPE
cat > /data/ip-geolocation.csv << 'EOF'
IP_ADDRESS,COUNTRY,CITY,ISP,TYPE
192.168.1.1,USA,New York,Comcast,Residential
10.0.0.5,USA,Chicago,Verizon,Data Center
172.16.0.1,Germany,Berlin,Deutsche Telekom,Enterprise
192.168.1.100,USA,Los Angeles,AT&T,Residential
10.0.0.10,UK,London,BT Group,Data Center
172.16.0.50,France,Paris,Orange,SMB
192.168.1.101,USA,Miami,Comcast,Residential
10.0.0.15,Japan,Tokyo,NTT,Enterprise
172.16.0.51,Canada,Toronto,Rogers,SMB
192.168.1.102,Australia,Sydney,Telstra,Residential
EOF

echo "Setup complete for Data Analysis lab"
echo "Data files created in /data/"
ls -la /data/
