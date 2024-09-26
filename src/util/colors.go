package util

import (
	"fmt"
)

func Colors() {
	// Print the colors in rows of 6
	for i := 0; i < 256; i++ {
		// Set the color using ANSI escape codes
		fmt.Printf("\033[48;5;%dm%4d\033[0m ", i, i)

		// Print a newline after every 6 colors for better readability
		if (i+1)%6 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}
