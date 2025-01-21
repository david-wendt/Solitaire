package ioutils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Input(s string) string {
	fmt.Print(s)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n') // Reads until newline
	return strings.ToLower(strings.TrimSpace(input))
}