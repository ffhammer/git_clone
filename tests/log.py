from utils import TestDir
import time

# Initialize test environment
test_dir = TestDir()

# Run GVC commands inside the test directory
test_dir.run_command("init")
test_dir.run_command("set --set User=name")
test_dir.run_command("set --set LogLevel=DEBUGGING")

# Create and track a file
test_dir.write_file("a.py", "some bs")
test_dir.write_file("b.py", "some bs")
test_dir.run_command("add a.py b.py")
test_dir.run_command('commit -m "first commit"')

# Create a new branch and modisfy the file
test_dir.write_file("a.py", "changed")
test_dir.write_file("c.py", "new")
test_dir.run_command("diff")
test_dir.run_command("add a.py c.py")
test_dir.run_command("rm b.py")
test_dir.run_command('commit -m "save changes to a"')


test_dir.run_command("log")
test_dir.run_command("log")
test_dir.cleanup()
