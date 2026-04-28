#!/usr/bin/env python3
"""Validate Cut Data Analysis lab."""

import sys

errors = 0

# Check emails.txt exists
if open("/tmp/labs/emails.txt").read():
    print("✓ Task 1: Emails extracted")
else:
    print("✗ Task 1: No emails")
    errors += 1

if errors == 0:
    print("✓ All tasks passed!")
    sys.exit(0)
else:
    print(f"✗ {errors} task(s) failed")
    sys.exit(1)