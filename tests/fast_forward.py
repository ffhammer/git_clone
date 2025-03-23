from utils import TestDir

# Initialize test environment
test_dir = TestDir()

# Run GVC commands inside the test directory
test_dir.run_command("init")
test_dir.run_command("set --set User=name")
test_dir.run_command("set --set LogLevel=DEBUGGING")

# Create and track a file
test_dir.write_file("a.py", "some bs")
test_dir.write_file("b.py", "some bs")
test_dir.run_command("status")
test_dir.run_command("add a.py b.py")
test_dir.run_command('commit -m "first commit"')

# Create a new branch and modify the file
test_dir.run_command("checkout -b new")
test_dir.write_file("a.py", "changed")
test_dir.write_file("c.py", "new")
test_dir.run_command("add a.py c.py")
test_dir.run_command("rm b.py")
test_dir.run_command('commit -m "save changes to a"')

# Checkout main and merge the new branch
test_dir.run_command("checkout main")
test_dir.list_specific_dir(".gvc/refs")
print(test_dir.read_file(".gvc/refs/main"))
print(test_dir.read_file(".gvc/HEAD"))

test_dir.run_command("merge new")
test_dir.print_file(".gvc/.log")
test_dir.print_files()


# Clean up after test
test_dir.cleanup()
