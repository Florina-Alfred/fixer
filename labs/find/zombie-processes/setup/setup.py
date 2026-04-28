#!/usr/bin/env python3
"""Setup script for Find Zombie Processes lab."""

import os

os.makedirs("/tmp/labs/config", exist_ok=True)
os.makedirs("/tmp/labs/.hidden", exist_ok=True)

# Config files
open("/tmp/labs/config/app.yml", "w").write("db_host: localhost\n")
open("/tmp/labs/config/database.yml", "w").write("cache: redis\n")

# Hidden files  
open("/tmp/labs/.hidden/backup-creds", "w").write("secret123\n")
open("/tmp/labs/.hidden/api-keys", "w").write("key456\n")

# Small files
open("/tmp/labs/tiny1.txt", "w").write("x\n")
open("/tmp/labs/tiny2.txt", "w").write("y\n")

# Large file
with open("/tmp/labs/big.dat", "wb") as f:
    f.write(b"x" * 10240)

print("Setup complete")