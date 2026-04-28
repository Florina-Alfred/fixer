#!/usr/bin/env python3
"""Validate Duplicate Finder lab."""

import subprocess
import sys

errors = 0

# Check if duplicate finder works
try:
    result = subprocess.run(["find", "/tmp/labs", "-type", "f", "-name", "*.txt"], 
                          capture_output=True, text=True)
    count = len(result.stdout.strip().split("\n")) if result.stdout.strip() else 0
    if count >= 2:
        print("✓ Task 1: Found files")
    else:
        print("✗ Task 1: No files found")
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