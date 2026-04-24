#!/bin/sh
# Setup script for Secret Hunter lab
# Creates a realistic log environment with hidden secrets

# Create log directories
mkdir -p /var/log/app /var/log/system /var/log/security /home/admin/logs

# Create a Python script that simulates a running application with logs
cat > /usr/local/bin/log-generator.py << 'PYEOF'
#!/usr/bin/env python3
import sys
import time
import random
import os

LOG_DIR = "/var/log/app"
SECRETS = [
    "SECRET_NUM_7429",
    "API_KEY_xK9mP2nQ",
    "TOKEN_aB3cD4eF5g",
]

log_messages = [
    "User {user} logged in from {ip}",
    "Processing request for {resource}",
    "Database query completed in {time}ms",
    "Cache hit for key {key}",
    "Background job {job} started",
    "Memory usage: {mem}%",
    "CPU load: {load}",
]

users = ["alice", "bob", "charlie", "diana", "eve"]
resources = ["api/users", "api/orders", "api/products", "api/inventory"]
jobs = ["cleanup", "sync", "backup", "report"]

def generate_ip():
    return f"{random.randint(10,192)}.{random.randint(0,255)}.{random.randint(0,255)}.{random.randint(1,254)}"

def generate_log():
    msg_template = random.choice(log_messages)
    msg = msg_template.format(
        user=random.choice(users),
        ip=generate_ip(),
        resource=random.choice(resources),
        time=random.randint(5, 500),
        key=f"key_{random.randint(1000, 9999)}",
        job=random.choice(jobs),
        mem=random.randint(20, 95),
        load=random.randint(1, 100),
    )
    
    timestamp = time.strftime("%Y-%m-%d %H:%M:%S")
    level = random.choices(["INFO", "DEBUG", "WARN", "ERROR"], weights=[60, 20, 15, 5])[0]
    
    # Occasionally inject a secret
    if random.random() < 0.02:  # 2% chance
        secret = random.choice(SECRETS)
        return f"{timestamp} [{level}] {secret} - {msg}"
    
    return f"{timestamp} [{level}] {msg}"

if __name__ == "__main__":
    log_file = os.path.join(LOG_DIR, "application.log")
    with open(log_file, "a") as f:
        for _ in range(100):  # Generate 100 initial log lines
            f.write(generate_log() + "\n")
        print(f"Generated 100 log entries in {log_file}")
PYEOF

chmod +x /usr/local/bin/log-generator.py

# Run the log generator
python3 /usr/local/bin/log-generator.py

# Create additional log files with secrets hidden
cat > /var/log/system/system.log << 'EOF'
2024-01-15 08:00:00 [INFO] System startup initiated
2024-01-15 08:00:01 [INFO] Loading kernel modules
2024-01-15 08:00:02 [INFO] Network interfaces configured
2024-01-15 08:00:03 [INFO] AUTH_CODE_9182 - Authentication service started
2024-01-15 08:00:04 [INFO] Filesystem mounted
2024-01-15 08:00:05 [INFO] System ready
2024-01-15 08:15:00 [WARN] High memory usage detected
2024-01-15 08:30:00 [INFO] Scheduled backup started
2024-01-15 08:30:05 [INFO] Backup completed successfully
2024-01-15 09:00:00 [ERROR] Disk space low on /var
EOF

# Create security logs with access codes
cat > /var/log/security/access.log << 'EOF'
2024-01-15 10:00:00 ACCESS_GRANTED user=admin ip=192.168.1.100
2024-01-15 10:05:00 ACCESS_DENIED user=hacker ip=10.0.0.99
2024-01-15 10:10:00 ACCESS_GRANTED user=alice ip=192.168.1.101
2024-01-15 10:15:00 MASTER_KEY_3847 - Security audit completed
2024-01-15 10:20:00 ACCESS_GRANTED user=bob ip=192.168.1.102
2024-01-15 10:25:00 ACCESS_DENIED user=unknown ip=10.0.0.100
2024-01-15 10:30:00 ACCESS_GRANTED user=charlie ip=192.168.1.103
EOF

# Create a configuration file with embedded secrets
cat > /home/admin/.config << 'EOF'
# Admin Configuration File
# Do not share this file

database_host=db.internal.local
database_port=5432
database_name=production

# API Configuration
api_endpoint=https://api.example.com
api_version=v2

# WARNING: The following contains sensitive data
emergency_contact=555-0123
backup_code=BK_7291

# Logging settings
log_level=INFO
log_retention_days=30
EOF

chmod 600 /home/admin/.config

# Create some decoy files
echo "This is not a secret" > /var/log/app/README
echo "Random data file" > /tmp/data.txt
mkdir -p /var/cache/app
echo "cache data" > /var/cache/app/cache.db

echo "Setup complete for Secret Hunter lab"
