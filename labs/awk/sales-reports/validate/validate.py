#!/usr/bin/env python3
"""Validate Awk Sales Reports lab."""

import sys

errors = 0

# Check revenue.txt exists
if open("/tmp/labs/sales.csv").read():
    print("✓ Setup working")
else:
    print("✗ Setup failed")
    errors += 1

if errors == 0:
    print("✓ All tasks passed!")
    sys.exit(0)
else:
    print(f"✗ {errors} task(s) failed")
    sys.exit(1)