import os
import subprocess

EXECUTABLE_FILE = "/Users/felix/Desktop/gvc"
SAVE_TO = "command_infos.md"


def run_command(command: str) -> str:
    result = subprocess.run(
        [EXECUTABLE_FILE] + command.split(),
        cwd=os.getcwd(),
        text=True,
        capture_output=True,
    )
    if result.stdout.strip():
        return result.stdout.strip()
    if result.stderr.strip():
        return result.stderr.strip()
    return ""


commands = [
    "init",
    "add",
    "rm",
    "commit",
    "status",
    "restore",
    "diff",
    "log",
    "branch",
    "checkout",
    "merge",
    "set",
]

with open(SAVE_TO, "w") as f:
    for cmd in commands:
        help_text = run_command(f"{cmd} -h")
        f.write(f"### {cmd}\n\n")
        f.write(help_text + "\n\n")
