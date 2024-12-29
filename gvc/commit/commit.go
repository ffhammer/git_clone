package commit

// type CommitMetdata struct {
// 	ParentCommitHash string `json:"parent_commit_hash"`
// 	BranchName       string `json:"branch_name"`
// 	Author           string `json:"author"`
// 	CommitMessage    string `json:"commit_message"`
// 	Date             string `json:"date"`
// 	TreeHash         string `json:"tree_hash"`
// }

func commit(message, author string)

// Core Outline for commit
// 1. Collect Metadata
// Commit Message: Accept a commit message as input.
// Author Information: Optionally include the author's name and email.
// Timestamp: Record the current date and time.
// Parent Commit: Identify the current HEAD or parent commit for history tracking.
// 2. Generate the Snapshot  sdasd
// Traverse the Index:
// Read the index to get a list of tracked files and their hashes.
// Verify the index reflects the working directory (optional: check for staged changes).
// Generate a Tree Object:
// Create a hierarchical representation of files and directories (like a Merkle tree).
// Store each directory or file as an object with a hash.
// 3. Write the Commit Object
// Create a new commit object containing:
// A reference to the root tree object (representing the snapshot).
// Metadata: Commit message, author, timestamp, and parent commit hash(es).
// Write the commit object to storage and generate its hash.
// 4. Update the HEAD
// Update the reference (HEAD) to point to the new commit.
// Optionally, update the branch reference (e.g., refs/heads/main).
// 5. Error Handling
// Ensure a clean index and no untracked changes (optional, depending on your design).
// Handle cases like an empty index or missing parent commits gracefully.
