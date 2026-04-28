#!/usr/bin/env python3
"""Setup script for CSV Extractor lab."""

import os

os.makedirs("/tmp/labs", exist_ok=True)

# Employees CSV
employees = """id,name,email,department,salary,phone
1,John Smith,john.smith@company.com,Engineering,75000,555-0101
2,Jane Doe,jane.doe@company.com,Marketing,65000,555-0102
3,Bob Johnson,bob.j@company.com,Engineering,80000,555-0103
4,Alice Brown,alice.b@company.com,Sales,70000,555-0104
5,Charlie Wilson,charlie.w@company.com,Marketing,62000,555-0105
6,Diana Lee,diana.l@company.com,Engineering,85000,555-0106
7,Edward Davis,edward.d@company.com,Sales,68000,555-0107
8,Fiona Garcia,fiona.g@company.com,HR,60000,555-0108
9,George Martinez,george.m@company.com,Engineering,90000,555-0109
10,Helen Anderson,helen.a@company.com,HR,58000,555-0110
"""
with open("/tmp/labs/employees.csv", "w") as f:
    f.write(employees)

# Customers CSV
customers = """customer_id,first_name,last_name,email,city,country,purchase_amount
101,Michael,Johnson,michael.j@email.com,New York,USA,150.99
102,Sarah,Williams,sarah.w@email.com,London,UK,230.50
103,David,Brown,david.b@email.com,Toronto,Canada,89.99
104,Emma,Davis,emma.d@email.com,Paris,France,320.00
105,James,Wilson,james.w@email.com,Sydney,Australia,175.25
106,Olivia,Taylor,olivia.t@email.com,Berlin,Germany,210.75
107,Robert,Thomas,robert.t@email.com,Tokyo,Japan,145.50
108,Sophia,Jackson,sophia.j@email.com,Madrid,Spain,280.00
"""
with open("/tmp/labs/customers.csv", "w") as f:
    f.write(customers)

print("Setup complete")