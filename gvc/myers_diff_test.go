package gvc

import (
	"bytes"
	"fmt"
	"testing"
)

func TestMyersDiffExamples(t *testing.T) {
	// Define test cases
	examples := []struct {
		name string
		a, b string
	}{
		{
			name: "Simple Example",
			a:    "line1\nline2\nline3",
			b:    "line1\nline2 modified\nline3\nline4",
		},
		{
			name: "Insertion Only",
			a:    "line1\nline2",
			b:    "line1\nline2\nline3\nline4",
		},
		{
			name: "Deletion Only",
			a:    "line1\nline2\nline3\nline4",
			b:    "line1\nline2",
		},
		{
			name: "No Changes",
			a:    "line1\nline2\nline3",
			b:    "line1\nline2\nline3",
		},
	}

	fmt.Print("welcome to test \n")
	for _, example := range examples {
		t.Run(example.name, func(t *testing.T) {
			// Compute the diff
			changes := MyersDiff(SplitLines(example.a), SplitLines(example.b))

			// Print the diff to stdout
			fmt.Printf("\n=== %s ===\n\n", example.name)
			fmt.Printf("Input A:\n%s\n\n", example.a)
			fmt.Printf("Input B:\n%s\n", example.b)

			// fmt.Println("\nDiff Output:")
			// DiffPrinter(changes, os.Stdout)
			fmt.Print("\n=== changes ===\n\n")
			for _, change := range changes {
				fmt.Printf("%d - %d - %d\n", change.action, change.oldLineNumber, change.newLineNumber)
			}
			// Optionally capture the diff into a buffer for testing
			var buf bytes.Buffer
			DiffPrinter(example.a, example.b, changes, &buf)
			fmt.Println("\nCaptured Output:")
			fmt.Println(buf.String())
		})
	}
}
