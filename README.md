Since I had some time during my gap year and historically sucked at anything git-related,  
I thought: *why not create a git clone in Go?*  
Thus, **gvc** (Go Version Control) was born.

Having finished it did not help — I still feel like I suck at git (:

Also, if you're proficient in Go: for your own safety, please refrain from looking at the code.

I completed the following functionality:

### init

gvc init
initialize a new .gvc in the working directory.
Fails if already exists

### add

gvc add [options] <path>
  <path> can be:
    - a direct path to a single file
    - a directory path (adds all non-ignored files in that dir)
    - a glob pattern (adds all matching non-ignored files)
Options:
  -f    Force adding files even if they are ignored in .gvcignore

### rm

gvc rm [options] <path>...
Removes file(s) from the working tree and/or the index.

Options:
  --cached   Only remove from index, keep file in working dir
  -f         Force removal, even if file is staged or ignored
  -r         Allow recursive removal of directories

### commit

gvc commit -m <message> -u <user>
Commits the staged changes with the provided commit message and user.
If no user is provided the user set by 'gvc set' is used
Fails if any positional arguments are supplied.

### status

gvc status
Show the current repository status, similar to 'git status'.
Displays the current branch, staged changes, modifications, untracked files, and merge state.

### restore

gvc restore [options] <path>
Restores files to their state in HEAD.

Options:
  --staged     Restore the index (remove from staging)
  --worktree   Restore the file content in the working directory

If neither --staged nor --worktree is given, neither are restored.

### diff

gvc diff [options] [args]
Show changes between commits, index, and working tree.

Modes:
  --cached       Diff index vs HEAD (staged changes)
  --no-index     Diff two files directly (bypasses index)
  <commit>       Diff commit vs working tree
  <c1> <c2>      Diff two commits
  (default)      Diff working tree vs index

Note: --cached and --no-index are mutually exclusive.

### log

gvc log [options]
Display the commit history.

Options:
  --patch         Show full diff for each commit
  --since <date>  Only show commits after this date (YYYY-MM-DD)
  --until <date>  Only show commits before this date (YYYY-MM-DD)
  --author <name> Filter by author
  --grep <msg>    Filter by message substring

### branch

gvc branch [options] [<branch-name>...]
Create, list, or delete branches.

Options:
  -d        Delete the specified branch(es)
  --help    Show this message

Usage:
  gvc branch           List all branches
  gvc branch <name>    Create a new branch
  gvc branch -d <name> Delete branch

### checkout

gvc checkout [options] <branch>
Switch branches or restore working tree.

Options:
  -b        Create the branch before switching to it

### merge

gvc merge [options] <branch>
Merge the specified branch into the current branch.

Options:
  -u <user>   Specify commit user (overrides config setting)

### set

gvc set [--list | --set key=value]
View or update GVC settings.

Options:
  --list           List current settings
  --set key=value  Set a setting value (e.g. --set User=felix)


Also includes a `.gvcignore` file for ignoring specific files, globs, and folders.
