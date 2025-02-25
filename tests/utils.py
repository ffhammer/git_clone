import os
import tempfile
import subprocess
from pathlib import Path

EXECUTABLE_FILE = "/Users/felix/Desktop/gvc"

# ANSI color codes
COLOR_RESET = "\033[0m"
COLOR_BLUE = "\033[1;34m"  # Bold Blue for Commands
COLOR_GREEN = "\033[1;32m"  # Green for stdout (Success)
COLOR_RED = "\033[1;31m"  # Red for stderr (Errors)
COLOR_YELLOW = "\033[1;33m"  # Yellow for Warnings


def run_command(command, cwd=None):
    """Run a shell command and print formatted output."""
    print(f"{COLOR_BLUE}> {command}{COLOR_RESET}")  # Command in blue

    result = subprocess.run(
        [EXECUTABLE_FILE] + command.split(), cwd=cwd, text=True, capture_output=True
    )

    if result.stdout.strip():
        print(f"{result.stdout.strip()}{COLOR_RESET}")  # Success output
    if result.stderr.strip():
        print(f"{result.stderr.strip()}{COLOR_RESET}")  # Error output


class TestDir:
    """Creates a temporary directory for testing GVC commands."""

    def __init__(self):
        self.dir = tempfile.TemporaryDirectory()
        self.path = Path(self.dir.name)

    def write_file(self, relative_path, content):
        """Write content to a file in the test directory."""
        file_path = self.path / relative_path
        file_path.parent.mkdir(parents=True, exist_ok=True)
        file_path.write_text(content)
        print(f"{COLOR_GREEN}Created: {file_path}{COLOR_RESET}")

    def delete_file(self, relative_path):
        """Delete a file in the test directory."""
        file_path = self.path / relative_path
        if file_path.exists():
            file_path.unlink()
            print(f"{COLOR_GREEN}Deleted: {file_path}{COLOR_RESET}")
        else:
            print(
                f"{COLOR_YELLOW}Warning: File {relative_path} does not exist.{COLOR_RESET}"
            )

    def list_files(self):
        """List files in the test directory."""
        return [
            str(p.relative_to(self.path)) for p in self.path.glob("**/*") if p.is_file()
        ]

    def cleanup(self):
        """Clean up the temporary directory."""
        self.dir.cleanup()

    def print_file(self, relative_path):
        """Print the contents of a file."""
        file_path = self.path / relative_path
        if not file_path.exists():
            print(
                f"{COLOR_YELLOW}Warning: File {relative_path} does not exist.{COLOR_RESET}"
            )
            return

        print(f"{COLOR_GREEN}Content of file {relative_path}:{COLOR_RESET}")
        print(f"{'#'*20}\n{file_path.read_text()}\n{'#'*20}")
