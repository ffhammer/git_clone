import os
import tempfile
import subprocess
from pathlib import Path
from utils import TestDir
from dataclasses import dataclass
from typing import Optional


@dataclass
class TestCase:
    name: str  # will corespond to file name
    content_begin: Optional[str]  # if None -> none existing
    content_new: Optional[str]
    content_main: Optional[str]
    final_file: Optional[str]
    trigger_conflict: bool = False


test_cases: list[TestCase] = [
    TestCase(
        name="untouched",
        content_begin="will stay same",
        content_new="will stay same",
        content_main="will stay same",
        final_file="will stay same",
    ),
    TestCase(
        name="both_modified_conflict",
        content_begin="base content",
        content_new="modified in new branch",
        content_main="modified in main branch",
        final_file="CONFLICT",
        trigger_conflict=True,
    ),
    TestCase(
        name="both_delete",
        content_begin="some content",
        content_new=None,
        content_main=None,
        final_file=None,
    ),
    TestCase(
        name="deleted_in_new_unchanged_old",
        content_begin="some content",
        content_new=None,
        content_main="some content",
        final_file=None,
    ),
    TestCase(
        name="deleted_in_old_unchanged_new",
        content_begin="some content",
        content_new="some content",
        content_main=None,
        final_file=None,
    ),
    TestCase(
        name="deleted_in_new_modified_in_old",
        content_begin="some content",
        content_new=None,
        content_main="modified in main",
        final_file="CONFLICT",
        trigger_conflict=True,
    ),
    TestCase(
        name="deleted_in_old_modified_in_new",
        content_begin="some content",
        content_new="modified in new",
        content_main=None,
        final_file="CONFLICT",
        trigger_conflict=True,
    ),
    TestCase(
        name="new_in_new",
        content_begin=None,
        content_new="new file content",
        content_main=None,
        final_file="new file content",
    ),
    TestCase(
        name="new_in_old",
        content_begin=None,
        content_new=None,
        content_main="new file content",
        final_file="new file content",
    ),
    TestCase(
        name="both_modified_identical",
        content_begin="base content",
        content_new="identical modification",
        content_main="identical modification",
        final_file="identical modification",
    ),
]


# Initialize test environment
test_dir = TestDir()

test_dir.run_command("init")
test_dir.run_command("set --set User=name")
test_dir.run_command("set --set LogLevel=DEBUGGING")

# initial commit
add_command = "add"
for case in test_cases:

    if case.content_begin is None:
        continue

    test_dir.write_file(case.name, case.content_begin)
    add_command = f"{add_command} {case.name}"

test_dir.run_command(add_command)
test_dir.run_command('commit -m "Initial commit on main"')


# feature
test_dir.run_command("checkout -b feature")

add_command = "add"
del_command = "rm"
for case in test_cases:

    # both none -> skip
    if case.content_begin is None and case.content_new is None:
        continue

    # delet
    if case.content_begin is not None and case.content_new is None:
        del_command = f"{del_command} {case.name}"
        continue

    # add or modify
    if case.content_begin != case.content_new:
        test_dir.write_file(case.name, case.content_new)
        add_command = f"{add_command} {case.name}"

test_dir.run_command(add_command)
test_dir.run_command(del_command)
test_dir.run_command("status")
test_dir.run_command('commit -m "Commit on feature branch"')


# Step 3: Switch back to main. and change
test_dir.run_command("checkout main")

add_command = "add"
del_command = "rm"
for case in test_cases:

    # both none -> skip
    if case.content_begin is None and case.content_main is None:
        continue

    # delet
    if case.content_begin is not None and case.content_main is None:
        del_command = f"{del_command} {case.name}"
        continue

    # add or modify
    if case.content_begin != case.content_main:
        test_dir.write_file(case.name, case.content_main)
        add_command = f"{add_command} {case.name}"

test_dir.run_command(add_command)
test_dir.run_command(del_command)

test_dir.run_command("status")
test_dir.run_command('commit -m "Commit on main branch"')

# Step 4: Merge branch 'feature' into main (this should trigger a three-way merge conflict).
test_dir.run_command("merge feature")
test_dir.run_command("status")

for case in test_cases:

    if not case.trigger_conflict:
        continue

    test_dir.print_file(case.name)


# test_dir.print_file(".gvc/.log")

# Clean up after test.
test_dir.cleanup()
