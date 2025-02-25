import os
import tempfile
import subprocess
from pathlib import Path
from utils import TestDir, run_command

# Initialize test environment
test_dir = TestDir()

# Run GVC commands inside the test directory
run_command("init", cwd=test_dir.path)
run_command("set --set User=name", cwd=test_dir.path)
run_command("set --set LogLevel=DEBUGGING", cwd=test_dir.path)

# Step 1: On main branch, create and commit the initial file.
test_dir.write_file("a.py", "original content")
run_command("add a.py", cwd=test_dir.path)
run_command('commit -m "first commit"', cwd=test_dir.path)

# Step 2: Create a new branch 'new' from main.
run_command("checkout -b new", cwd=test_dir.path)

# On branch 'new', modify the file to create a conflicting change.
test_dir.write_file("a.py", "changed in new branch")
run_command("add a.py", cwd=test_dir.path)
run_command('commit -m "commit on new branch"', cwd=test_dir.path)

# Step 3: Switch back to main.
run_command("checkout main", cwd=test_dir.path)

# On main, modify the file differently to cause a conflict.
test_dir.write_file("a.py", "changed in main branch")
run_command("add a.py", cwd=test_dir.path)
run_command('commit -m "commit on main branch"', cwd=test_dir.path)

# Step 4: Merge branch 'new' into main (should trigger a merge conflict).
run_command("merge new", cwd=test_dir.path)
test_dir.print_file("a.py")
run_command("status", cwd=test_dir.path)

# Step 5: Print the merge log file to inspect conflict markers.
test_dir.print_file(".gvc/.log")

# Clean up after test
test_dir.cleanup()
