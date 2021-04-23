package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

const MaxLength = 50

var (
	types = []string{"chore", "docs", "feat", "fix", "refactor", "revert", "style", "test"}
	scopes = []string{"cli", "config", "proxy", "redis", "server"}
)

func main() {
	if len(os.Args) != 2 {
		fail("One argument is expected")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fail("Can't open file")
	}
	defer func() {
		if error := file.Close(); error != nil {
			fmt.Printf("Error closing file: %s\n", error)
		}
	}()

	reader := bufio.NewReader(file)
	s, err := reader.ReadString('\n')
	if err != nil {
		fail("Can't read file")
	}

	r := regexp.MustCompile(`^(\w*)(\((\w*)\))?: .*`)
	match := r.FindStringSubmatch(s)
	if len(match) == 0 {
		fail("Commit message should start with `type(scope?): subject`")
	}
	if len(s) > MaxLength {
		fail(fmt.Sprintf("First line of commit message is loger then %v characters", MaxLength))
	}

	if !contains(types, match[1]) {
		fail(fmt.Sprintf("Unknown type `%v`.\n Available types: %v", match[1], types))
	}

	if match[3] != "" && !contains(scopes, match[3]) {
		fail(fmt.Sprintf("Unknown scope `%v`.\n Available scopes: %v", match[3], scopes))
	}
}

func fail(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
