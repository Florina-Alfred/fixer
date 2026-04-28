#!/usr/bin/env python3
"""Validate Log Analyzer lab."""

import os
import subprocess
import sys

errors = 0

# Check Task 1: errors.txt
try:
    result = subprocess.run(["grep", "-c", "ERROR", "/tmp/labs/error.log"], 
                          capture_output=True, text=True)
    count = int(result.stdout.strip()) if result.returncode == 0 else 0
    if count >= 7:
        print(f"✓ Task 1: Found {count} error lines")
    else:
        print(f"✗ Task 1: Expected 7, got {count}")
        errors += 1
except Exception as e:
    print(f"✗ Task 1: {e}")
    errors += 1

# Check Task 2: client_errors.txt
try:
    result = subprocess.run(["grep", "-c", " 4[0-9][0-9] ", "/tmp/labs/access.log"], 
                          capture_output=True, text=True)
    count = int(result.stdout.strip()) if result.returncode == 0 else 0
    if count >= 6:
        print(f"✓ Task 2: Found {count} client errors")
    else:
        print(f"✗ Task 2: Expected 6, got {count}")
        errors += 1
except Exception as e:
    print(f"✗ Task 2: {e}")
    errors += 1

# Check Task 3: suspicious.txt
try:
    result = subprocess.run(["grep", "203.0.113.42", "/tmp/labs/suspicious.txt"], 
                          capture_output=True, text=True)
    if result.returncode == 0:
        print("✓ Task 3: Suspicious IP found")
    else:
        print("✗ Task 3: IP not found")
        errors += 1
except Exception as e:
    print(f"✗ Task 3: {e}")
    errors += 1

print()
if errors == 0:
    print("✓ All tasks passed!")
    sys.exit(0)
else:
    print(f"✗ {errors} task(s) failed")
    sys.exit(1)