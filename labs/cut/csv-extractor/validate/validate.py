#!/usr/bin/env python3
"""Validate CSV Extractor lab."""

import subprocess
import sys

errors = 0

# Check employee_names.txt
try:
    result = subprocess.run(["wc", "-l", "/tmp/labs/employee_names.txt"], 
                        capture_output=True, text=True)
    count = int(result.stdout.strip().split()[0]) if result.returncode == 0 else 0
    if count >= 1:
        print(f"✓ Task 1: Found {count} names")
    else:
        print("✗ Task 1: No names")
        errors += 1
except Exception as e:
    print(f"✗ Task 1: {e}")
    errors += 1

print()
if errors == 0:
    print("✓ All tasks passed!")
    sys.exit(0)
else:
    print(f"✗ {errors} task(s) failed")
    sys.exit(1)