from utils import TestDir
from dataclasses import dataclass
from typing import Optional


@dataclass
class TestCase:
    name: str
    content_begin: Optional[str]
    content_new: Optional[str]
    content_main: Optional[str]
    final_file: Optional[str]
    trigger_conflict: bool = False


test_cases: list[TestCase] = [
    TestCase(
        "untouched",
        "will stay same",
        "will stay same",
        "will stay same",
        "will stay same",
    ),
    TestCase(
        "both_modified_conflict",
        "base content",
        "modified in new branch",
        "modified in main branch",
        "CONFLICT",
        True,
    ),
    TestCase("both_delete", "some content", None, None, None),
    TestCase(
        "deleted_in_new_unchanged_old", "some content", None, "some content", None
    ),
    TestCase(
        "deleted_in_old_unchanged_new", "some content", "some content", None, None
    ),
    TestCase(
        "deleted_in_new_modified_in_old",
        "some content",
        None,
        "modified in main",
        "CONFLICT",
        True,
    ),
    TestCase(
        "deleted_in_old_modified_in_new",
        "some content",
        "modified in new",
        None,
        "CONFLICT",
        True,
    ),
    TestCase("new_in_new", None, "new file content", None, "new file content"),
    TestCase("new_in_old", None, None, "new file content", "new file content"),
    TestCase(
        "both_modified_identical",
        "base content",
        "identical modification",
        "identical modification",
        "identical modification",
    ),
]


def apply_state(test_dir: TestDir, attr: str, base_attr: str) -> None:
    add_cmd, del_cmd = [], []
    for case in test_cases:
        base = getattr(case, base_attr)
        content = getattr(case, attr)

        if base is None and content is None:
            continue
        if base is not None and content is None:
            del_cmd.append(case.name)
        elif base != content:
            test_dir.write_file(case.name, content)
            add_cmd.append(case.name)

    if add_cmd:
        test_dir.run_command(f"add {' '.join(add_cmd)}")
    if del_cmd:
        test_dir.run_command(f"rm {' '.join(del_cmd)}")


def init_repo(test_dir: TestDir) -> None:
    test_dir.run_command("init")
    test_dir.run_command("set --set User=name")
    test_dir.run_command("set --set LogLevel=DEBUGGING")


def make_initial_commit(test_dir: TestDir) -> None:
    names = [c.name for c in test_cases if c.content_begin is not None]
    for case in test_cases:
        if case.content_begin is not None:
            test_dir.write_file(case.name, case.content_begin)
    if names:
        test_dir.run_command(f"add {' '.join(names)}")
    test_dir.run_command('commit -m "Initial commit on main"')


def verify_results(test_dir: TestDir) -> None:
    for case in test_cases:
        if case.trigger_conflict:
            test_dir.print_file(case.name)
        elif case.final_file is None:
            assert (
                case.name not in test_dir.list_dir()
            ), f"{case.name} should be deleted"
        else:
            content = test_dir.read_file(case.name)
            assert (
                content == case.final_file
            ), f"{case.name}: expected '{case.final_file}', got '{content}'"


test_dir = TestDir()
init_repo(test_dir)
make_initial_commit(test_dir)

test_dir.run_command("checkout -b feature")
apply_state(test_dir, "content_new", "content_begin")
test_dir.run_command("status")
test_dir.run_command('commit -m "Commit on feature branch"')

test_dir.run_command("checkout main")
apply_state(test_dir, "content_main", "content_begin")
test_dir.run_command("status")
test_dir.run_command('commit -m "Commit on main branch"')

test_dir.run_command("merge feature")
test_dir.run_command("status")
test_dir.write_file("both_modified_conflict", "halo")
test_dir.run_command(
    "add deleted_in_new_modified_in_old both_modified_conflict deleted_in_old_modified_in_new"
)
test_dir.run_command("restore penis")

test_dir.print_file("both_modified_identical")
test_dir.run_command("restore -staged both_modified_identical")
test_dir.run_command("restore -staged both_modified_conflict")
test_dir.print_file("both_modified_conflict")
test_dir.run_command("status")
test_dir.run_command("commit -m 'merge done' -u test")
test_dir.run_command("log", wait_time=1)

# verify_results(test_dir)
test_dir.print_file(".gvc/.log")

test_dir.cleanup()
