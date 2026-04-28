#!/usr/bin/env python3
"""Validate Grep Secret Hunter lab."""

import os
import sys

errors = 0

# Check Task 1: AUTH_CODE
auth_code_file = "/tmp/auth-code.txt"
if os.path.exists(auth_code_file):
    with open(auth_code_file) as f:
        content = f.read().strip()
    if "AUTH_CODE_9182" in content:
        print("✓ Task 1: AUTH_CODE found")
    else:
        print(f"✗ Task 1: AUTH_CODE not found (got: {content})")
        errors += 1
else:
    print(f"✗ Task 1: {auth_code_file} not found")
    errors += 1

# Check Task 2: Secrets
secrets_file = "/tmp/secrets_found.txt"
if os.path.exists(secrets_file):
    with open(secrets_file) as f:
        content = f.read()
    has_secret_num = "SECRET_NUM_7429" in content
    has_api_key = "API_KEY_xK9mP2nQ" in content
    has_token = "TOKEN_aB3cD4eF5g" in content
    if has_secret_num and has_api_key and has_token:
        print("✓ Task 2: All secrets found")
    else:
        print(f"✗ Task 2: Missing secrets (SECRET_NUM:{has_secret_num}, API_KEY:{has_api_key}, TOKEN:{has_token})")
        errors += 1
else:
    print(f"✗ Task 2: {secrets_file} not found")
    errors += 1

# Check Task 3: MASTER_KEY
master_file = "/tmp/master-key.txt"
if os.path.exists(master_file):
    with open(master_file) as f:
        content = f.read()
    if "MASTER_KEY_3847" in content:
        print("✓ Task 3: MASTER_KEY found")
    else:
        print(f"✗ Task 3: MASTER_KEY not found")
        errors += 1
else:
    print(f"✗ Task 3: {master_file} not found")
    errors += 1

print()
if errors == 0:
    print("✓ All tasks passed!")
    sys.exit(0)
else:
    print(f"✗ {errors} task(s) failed")
    sys.exit(1)