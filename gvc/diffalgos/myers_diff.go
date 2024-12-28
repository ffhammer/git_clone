package diffalgos

import (
	"fmt"
	"git_clone/gvc/utils"
	"io"

	"github.com/fatih/color"
)

type EditType int

const (
	Insert EditType = iota // Automatically assigned 0
	Delete                 // Automatically assigned 1
	Keep                   // Automatically assigned 2
)

type LineChange struct {
	action        EditType
	oldLineNumber int
	newLineNumber int
}

func index_around(i int, N_Dia int) int {
	if i < 0 {
		return N_Dia + i
	} else {
		return i
	}
}

func getTraces(a, b []string) [][]int {
	N := len(a)
	M := len(b)
	max_s := N + M
	N_DIA := 2*max_s + 1

	traces := make([][]int, 0, max_s)
	current_v := make([]int, N_DIA)

	for d := 0; d < max_s+1; d++ {
		next_v := make([]int, len(current_v))
		copy(next_v, current_v)
		traces = append(traces, next_v)

		for k := -d; k <= d; k += 2 {
			var x, y int
			if k == -d || (k != d && current_v[index_around(k-1, N_DIA)] < current_v[index_around(k+1, N_DIA)]) {
				x = current_v[index_around(k+1, N_DIA)]
			} else {
				x = current_v[index_around(k-1, N_DIA)] + 1
			}
			y = x - k
			for x < N && y < M && a[x] == b[y] {
				x++
				y++
			}

			current_v[index_around(k, N_DIA)] = x

			if x >= N && y >= M {
				return traces
			}
		}
	}
	return traces
}

func backtrack(traces [][]int, a, b []string) []LineChange {
	N := len(a)
	M := len(b)

	x, y := N, M
	type Tup4 struct {
		x, y, old_x, old_y int
	}

	var changes []Tup4

	for d := len(traces) - 1; d >= 0; d-- {
		k := x - y
		var prev_k int

		if k == -d || (k != d && traces[d][index_around(k-1, len(traces[0]))] < traces[d][index_around(k+1, len(traces[0]))]) {
			prev_k = k + 1
		} else {
			prev_k = k - 1
		}

		prev_x := traces[d][index_around(prev_k, len(traces[0]))]
		prev_y := prev_x - prev_k

		for x > prev_x && y > prev_y {
			changes = append(changes, Tup4{x - 1, y - 1, x, y})
			x--
			y--
		}

		if d > 0 {
			changes = append(changes, Tup4{prev_x, prev_y, x, y})
		}
		x = prev_x
		y = prev_y
	}

	// Generate LineChange results
	result := make([]LineChange, len(changes))
	for i := 0; i < len(changes); i++ {
		change := changes[len(changes)-1-i] // Reverse changes order
		if change.old_x == change.x {
			result[i] = LineChange{Insert, -1, change.y}
		} else if change.old_y == change.y {
			result[i] = LineChange{Delete, change.x, -1}
		} else {
			result[i] = LineChange{Keep, change.x, change.y}
		}
	}

	return result
}

func DiffPrinter(oldText, newText string, changes []LineChange, writer io.Writer) {
	oldLines := utils.SplitLines(oldText)
	newLines := utils.SplitLines(newText)

	for _, change := range changes {
		switch change.action {
		case Insert:
			line := fmt.Sprintf("+ %s\n", newLines[change.newLineNumber])
			fmt.Fprint(writer, color.New(color.FgGreen).Sprint(line))
		case Delete:
			line := fmt.Sprintf("- %s\n", oldLines[change.oldLineNumber])
			fmt.Fprint(writer, color.New(color.FgRed).Sprint(line))
		case Keep:
			line := fmt.Sprintf("  %s\n", oldLines[change.oldLineNumber])
			fmt.Fprint(writer, color.New(color.FgWhite).Sprint(line))
		}
	}
}

func MyersDiff(a, b []string) []LineChange {

	return backtrack(getTraces(a, b), a, b)
}
