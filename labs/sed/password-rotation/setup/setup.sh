#!/bin/sh
# Setup script for Password Rotation lab
# Creates realistic config files with passwords to rotate

mkdir -p /etc/services /var/config /app/config

# Create database configuration files
cat > /etc/services/database.conf << 'EOF'
# Database Configuration
# Last updated: 2024-01-15

[primary]
host=db-primary.internal.local
port=5432
database=production
username=app_user
password=OLD_PASS_12345
ssl_mode=require

[replica]
host=db-replica.internal.local
port=5432
database=production
username=replica_user
password=OLD_REPLICA_67890
ssl_mode=require

[cache]
host=redis.internal.local
port=6379
password=OLD_CACHE_54321
EOF

# Create API configuration
cat > /app/config/api-keys.conf << 'EOF'
# API Keys and Secrets
# WARNING: Rotate these regularly!

STRIPE_SECRET_KEY=sk_live_OLD1234567890abcdefghij
STRIPE_PUBLISHABLE_KEY=pk_live_OLD0987654321zyxwvuts

AWS_ACCESS_KEY_ID=AKIAOLD1234567890A
AWS_SECRET_ACCESS_KEY=oldsecretkey1234567890abcdefghij

SENDGRID_API_KEY=SG.oldapikey1234567890abcdef.ghijklmnop

GITHUB_TOKEN=ghp_oldtoken1234567890abcdefghijklmno
EOF

# Create nginx configuration with auth
cat > /etc/services/nginx-auth.conf << 'EOF'
# Nginx Authentication Config
# For admin panel access

auth_basic "Admin Panel";
auth_basic_user_file /etc/nginx/.htpasswd;

# htpasswd entries (base64 encoded)
# admin:$apr1$OLDHASH1234$xyz
# john:$apr1$OLDHASH5678$abc
# jane:$apr1$OLDHASH9012$def

# Backup password file
# Format: username:password
admin:OLD_ADMIN_PASS_111
john:OLD_JOHN_PASS_222
jane:OLD_JANE_PASS_333
EOF

# Create application secrets file
cat > /var/config/app-secrets.yml << 'EOF'
# Application Secrets
# DO NOT COMMIT TO VERSION CONTROL

database:
  password: DB_OLD_SECRET_2024
  encryption_key: ENC_OLD_KEY_12345

jwt:
  secret: JWT_OLD_SECRET_KEY_67890
  expiry_hours: 24

session:
  secret: SESSION_OLD_TOKEN_ABCDEFGH
  timeout_minutes: 30

encryption:
  master_key: MASTER_OLD_KEY_99999
EOF

# Create a log file with some passwords accidentally logged (security issue!)
cat > /var/config/access.log << 'EOF'
2024-01-15 10:00:00 [INFO] User admin logged in with password OLD_ADMIN_PASS_111
2024-01-15 10:05:00 [INFO] API call with key sk_live_OLD1234567890abcdefghij
2024-01-15 10:10:00 [WARN] Failed login for user john with OLD_JOHN_PASS_222
2024-01-15 10:15:00 [INFO] Database connection with OLD_PASS_12345
2024-01-15 10:20:00 [ERROR] Invalid token: ghp_oldtoken1234567890abcdefghijklmno
2024-01-15 10:25:00 [INFO] Cache connection with OLD_CACHE_54321
EOF

# Create a template file for reference
cat > /var/config/password-template.txt << 'EOF'
# Password Rotation Template
# Replace OLD_ with NEW_ followed by timestamp

Format: NEW_<SERVICE>_<TIMESTAMP>
Example: NEW_DB_20240115

Services to rotate:
- database (primary, replica, cache)
- api_keys (stripe, aws, sendgrid, github)
- app_secrets (database, jwt, session, encryption)
- nginx_auth (admin, john, jane)
EOF

echo "Setup complete for Password Rotation lab"
echo "Configuration files created with old passwords"
echo "Task: Rotate all passwords to NEW_<SERVICE>_<TIMESTAMP> format"
