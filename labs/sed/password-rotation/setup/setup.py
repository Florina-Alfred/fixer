#!/usr/bin/env python3
"""Setup script for Sed Password Rotation lab."""

import os

os.makedirs("/tmp/labs", exist_ok=True)

# database.conf
database_conf = """db_master = OLD_PASS_12345
db_replica = OLD_REPLICA_67890
db_cache = OLD_CACHE_54321
"""
with open("/tmp/labs/database.conf", "w") as f:
    f.write(database_conf)

# api-keys.conf
api_conf = """api_key = old_api_key_12345
api_secret = old_secret_67890
"""
with open("/tmp/labs/api-keys.conf", "w") as f:
    f.write(api_conf)

# nginx-auth.conf
nginx_conf = """admin = OLD_ADMIN_111
john = OLD_JOHN_222
jane = OLD_JANE_333
"""
with open("/tmp/labs/nginx-auth.conf", "w") as f:
    f.write(nginx_conf)

print("Setup complete")