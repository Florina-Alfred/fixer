#!/usr/bin/env python3
"""Setup script for Awk Sales Reports lab."""

import os

os.makedirs("/tmp/labs", exist_ok=True)

# Sales data (pipe-delimited)
sales = """date|product|amount|region|status
2024-01-15|Widget|100.00|North|SHIPPED
2024-01-16|Gadget|200.00|South|SHIPPED
2024-01-17|Tool|150.00|East|PENDING
2024-01-18|Widget|100.00|West|SHIPPED
"""
with open("/tmp/labs/sales.csv", "w") as f:
    f.write(sales)

print("Setup complete")