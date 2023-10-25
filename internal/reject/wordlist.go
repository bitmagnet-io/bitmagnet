package reject

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ContainsBanWord(title string) bool {
	file, err := os.Open("word_blacklist")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(strings.ToLower(title), strings.ToLower(scanner.Text())) {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return false
}
