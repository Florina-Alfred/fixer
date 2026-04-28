#!/usr/bin/env python3
"""Setup script for Grep Secret Hunter lab."""

import os

os.makedirs("/tmp/labs", exist_ok=True)

# Write system.log with AUTH_CODE
system_log = """2024-01-15 08:00:03 [INFO] AUTH_CODE_9182 - Authentication service started
2024-01-15 08:00:04 [INFO] Filesystem mounted
2024-01-15 08:00:05 [INFO] System ready
"""
with open("/tmp/labs/system.log", "w") as f:
    f.write(system_log)

# Write app.log with secrets
app_log = """2024-01-15 08:10:00 [DEBUG] SECRET_NUM_7429 loaded for processing
2024-01-15 08:15:00 [INFO] Database connected
2024-01-15 08:20:00 [WARN] API_KEY_xK9mP2nQ expiring soon
2024-01-15 08:25:00 [INFO] Token refresh: TOKEN_aB3cD4eF5g
"""
with open("/tmp/labs/app.log", "w") as f:
    f.write(app_log)

# Write security.log with MASTER_KEY
security_log = """2024-01-15 10:15:00 MASTER_KEY_3847 - Security audit completed
"""
with open("/tmp/labs/security.log", "w") as f:
    f.write(security_log)

print("Setup complete")