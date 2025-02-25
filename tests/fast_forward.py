from utils import TestDir, run_command

# Initialize test environment
test_dir = TestDir()

# Run GVC commands inside the test directory
run_command("init", cwd=test_dir.path)
run_command("set --set User=name", cwd=test_dir.path)
run_command("set --set LogLevel=DEBUGGING", cwd=test_dir.path)

# Create and track a file
test_dir.write_file("a.py", "some bs")
run_command("status", cwd=test_dir.path)
run_command("add a.py", cwd=test_dir.path)
run_command('commit -m "first commit"', cwd=test_dir.path)

# Create a new branch and modify the file
run_command("checkout -b new", cwd=test_dir.path)
test_dir.write_file("a.py", "changed")
run_command("add a.py", cwd=test_dir.path)
run_command('commit -m "save changes to a"', cwd=test_dir.path)

# Checkout main and merge the new branch
run_command("checkout main", cwd=test_dir.path)
run_command("merge new", cwd=test_dir.path)
test_dir.print_file(".gvc/.log")


# Clean up after test
test_dir.cleanup()
