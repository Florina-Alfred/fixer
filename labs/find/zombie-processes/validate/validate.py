#!/usr/bin/env python3
"""Validate Find Zombie Processes lab."""

import subprocess
import sys

errors = 0

# Task 1: config files
try:
    result = subprocess.run(["find", "/tmp/labs", "-name", "*.yml"], capture_output=True, text=True)
    if ".yml" in result.stdout:
        print("✓ Task 1: Config files found")
    else:
        print("✗ Task 1: No config files")
        errors += 1
except Exception as e:
    print(f"✗ Task 1: {e}")
    errors += 1

# Task 2: hidden files
try:
    result = subprocess.run(["find", "/tmp/labs", "-name", ".*"], capture_output=True, text=True)
    if result.stdout.strip():
        print("✓ Task 2: Hidden files found")
    else:
        print("✗ Task 2: No hidden files")
        errors += 1
except Exception as e:
    print(f"✗ Task 2: {e}")
    errors += 1

print()
if errors == 0:
    print("✓ All tasks passed!")
    sys.exit(0)
else:
    print(f"✗ {errors} task(s) failed")
    sys.exit(1)