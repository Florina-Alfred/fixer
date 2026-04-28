#!/usr/bin/env python3
"""Setup script for Duplicate Finder lab."""

import os
import stat

os.makedirs("/tmp/labs", exist_ok=True)

# Duplicate content
open("/tmp/labs/file1.txt", "w").write("same content\n")
open("/tmp/labs/file2.txt", "w").write("same content\n")
open("/tmp/labs/unique1.txt", "w").write("unique 1\n")

# Empty dir
os.makedirs("/tmp/labs/empty_dir", exist_ok=True)

# Empty file
open("/tmp/labs/empty_file.txt", "w").write("")

# World-writable (insecure)
open("/tmp/labs/insecure.conf", "w").write("")
os.chmod("/tmp/labs/insecure.conf", 0o777)

# SUID binary
open("/tmp/labs/special.bin", "w").write("")
os.chmod("/tmp/labs/special.bin", 0o4755)

print("Setup complete")