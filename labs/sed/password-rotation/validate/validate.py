#!/usr/bin/env python3
"""Validate Sed Password Rotation lab."""

import sys

errors = 0

# Check database.conf has OLD first (setup)
if "OLD_PASS" in open("/tmp/labs/database.conf").read():
    print("✓ Task: Setup complete (OLD_PASS found)")
    print("  Use 'sed -i s/OLD_PASS/NEW_PASS/g /tmp/labs/database.conf'")
else:
    print("✗ Setup not found")
    errors += 1

if errors == 0:
    print("✓ Ready for sed practice!")
    sys.exit(0)
else:
    print(f"✗ {errors} issue(s)")
    sys.exit(1)