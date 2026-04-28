#!/usr/bin/env python3
"""Setup script for Cut Data Analysis lab."""

import os

os.makedirs("/tmp/labs", exist_ok=True)

# Employees (pipe-delimited)
employees = """id|name|email|department|salary|phone
1|John|john@test.com|Engineering|75000|555-0101
2|Jane|jane@test.com|Marketing|65000|555-0102
3|Bob|bob@test.com|Engineering|80000|555-0103
4|Alice|alice@test.com|Sales|70000|555-0104
"""
with open("/tmp/labs/employees.csv", "w") as f:
    f.write(employees)

# Server logs
server_logs = """2024-01-15 08:00:00 [INFO] Server started
2024-01-15 08:10:00 [ERROR] auth-service: connection failed
2024-01-15 08:15:00 [ERROR] payment: timeout
2024-01-15 08:20:00 [ERROR] database: connection lost
2024-01-15 08:30:00 [INFO] Server running
"""
with open("/tmp/labs/server_logs.txt", "w") as f:
    f.write(server_logs)

# Network data
network = """ip,bytes,status
192.168.1.1,1024,active
10.0.0.1,2048,active
192.168.1.2,512,error
"""
with open("/tmp/labs/network.csv", "w") as f:
    f.write(network)

# Sales data
sales = """date,product,region,amount,customer
2024-01-15,Widget,North,100.00,Acme
2024-01-16,Gadget,South,200.00,Beta
"""
with open("/tmp/labs/sales.csv", "w") as f:
    f.write(sales)

print("Setup complete")