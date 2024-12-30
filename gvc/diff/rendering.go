package diff

import (
	"errors"
	"fmt"
	"git_clone/gvc/config"
	"git_clone/gvc/diffalgos"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

// controls sepration and also when a diff began
const nLinesBuffer = 3

var (
	greenColor = color.New(color.FgGreen)
	redColor   = color.New(color.FgRed)
	blueColor  = color.RGB(53, 175, 183)
)

func getContext(lines []string, chunkStart int, fileEnding string) string {
	if fileEnding != "" {
		fmt.Print("diff context specific to file endings is not implemented yet\n")
	}

	for i := chunkStart; i >= 0; i-- {
		if len(lines[i]) > 0 && !unicode.IsSpace(rune(lines[i][0])) {
			return lines[i]
		}
	}
	return ""
}

func chunkLinesToString(chunk []diffalgos.LineChange, a, b []string, builder *strings.Builder) {
	for _, change := range chunk {
		switch change.Action {
		case diffalgos.Insert:
			builder.WriteString(greenColor.Sprint(fmt.Sprintf("+%s\n", b[change.NewLineNumber])))
		case diffalgos.Delete:
			builder.WriteString(redColor.Sprint(fmt.Sprintf("-%s\n", a[change.OldLineNumber])))
		case diffalgos.Keep:
			builder.WriteString(fmt.Sprintf(" %s\n", a[change.OldLineNumber]))
		}
	}
}

func generateBody(a, b []string, output *strings.Builder) {
	if len(a) == 0 {
		output.WriteString(greenColor.Sprint("+" + strings.Join(b, "\n+ ") + "\n"))
	}
	if len(b) == 0 {
		output.WriteString(redColor.Sprint("-" + strings.Join(a, "\n- ") + "\n"))
	}

	diffs := diffalgos.MyersDiff(a, b)

	for startIndex := 0; startIndex < len(diffs); {
		if diffs[startIndex].Action == diffalgos.Keep {
			startIndex++
			continue
		}

		linesBlank := 0
		endIndex := startIndex + 1

		lastValidOldLineNumber := max(1, diffs[startIndex].OldLineNumber)
		lastValidNewLineNumber := max(1, diffs[startIndex].NewLineNumber)
		firstValidOldLineNumber := max(0, diffs[startIndex].OldLineNumber)
		firstValidNewLineNumber := max(0, diffs[startIndex].NewLineNumber)

		for endIndex < len(diffs) && linesBlank < nLinesBuffer {
			if diffs[endIndex].Action == diffalgos.Keep {
				linesBlank++
			} else {
				linesBlank = 0
			}
			lastValidOldLineNumber = max(lastValidOldLineNumber, diffs[endIndex].OldLineNumber)
			lastValidNewLineNumber = max(lastValidNewLineNumber, diffs[endIndex].NewLineNumber)
			firstValidOldLineNumber = min(firstValidOldLineNumber, max(0, diffs[startIndex].OldLineNumber))
			firstValidNewLineNumber = min(firstValidNewLineNumber, max(0, diffs[startIndex].OldLineNumber))
			endIndex++
		}
		endIndex--

		blankStartIndex := max(0, startIndex-nLinesBuffer)
		for i := startIndex - 1; i <= 0 && (-i) < blankStartIndex; i-- {
			firstValidOldLineNumber = min(firstValidOldLineNumber, max(0, diffs[startIndex].OldLineNumber))
			firstValidNewLineNumber = min(firstValidNewLineNumber, max(0, diffs[startIndex].OldLineNumber))
		}

		output.WriteString(blueColor.Sprint(fmt.Sprintf("@@ -%d,%d +%d,%d @@",
			firstValidOldLineNumber+1, lastValidOldLineNumber-firstValidOldLineNumber+1,
			firstValidNewLineNumber+1, lastValidNewLineNumber-firstValidNewLineNumber+1,
		)))

		output.WriteString("\n")
		chunkLinesToString(diffs[blankStartIndex:endIndex], a, b, output)
		output.WriteString("\n")

		startIndex = endIndex + 1
	}

}

func GenerateFileDiff(oldHash, oldRelPath, newHash, newRelPath string, oldFile, newFile []string) (string, error) {

	if oldHash == config.DOES_NOT_EXIST_HASH && newHash == config.DOES_NOT_EXIST_HASH {
		return "", errors.New("encountered logic assertion in generateFileDiff, both hashes empty")
	}

	var outputs strings.Builder

	outputs.WriteString(fmt.Sprintf("diff --git a/%s b/%s\n", oldRelPath, newRelPath))
	if oldHash == config.DOES_NOT_EXIST_HASH {
		outputs.WriteString("new file\n")
	}
	if newHash == config.DOES_NOT_EXIST_HASH {
		outputs.WriteString("deleted file\n")
	}

	outputs.WriteString(fmt.Sprintf("index %s..%s\n", oldHash, newHash))

	if oldHash == config.DOES_NOT_EXIST_HASH {
		outputs.WriteString("--- dev/null\n")
	} else {
		outputs.WriteString(fmt.Sprintf("--- a/%s\n", oldRelPath))
	}
	if newHash == config.DOES_NOT_EXIST_HASH {
		outputs.WriteString("+++ dev/null\n")
	} else {
		outputs.WriteString(fmt.Sprintf("+++ b/%s\n", oldRelPath))
	}

	generateBody(oldFile, newFile, &outputs)

	return outputs.String(), nil
}
