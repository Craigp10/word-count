package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	configValue configStruct
	filePaths   []string
)

type configStruct struct {
	Count      bool
	Lines      bool
	Characters bool
	Words      bool
}

type countStruct struct {
	Count      int64
	Lines      int32
	Characters int64
	Words      int32
}

var wcCmd = &cobra.Command{
	Use:   "wc [files...]",
	Short: "My word counter CLI Application",
	Long:  "This is a sample word counter CLI application using Cobra as the CLI reader.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var c countStruct
		for _, arg := range args {
			// Perform actions with the variadic strings here
			// filePaths = append(filePaths, arg)
			processFile(&c, &configValue, arg)
		}
	},
}

func processFile(c *countStruct, cfg *configStruct, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if cfg.Count {
		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}
		fileSize := fileInfo.Size()
		c.Count = fileSize
	}

	scanner := bufio.NewScanner(file)
	fmt.Println(scanner.Text())
	for scanner.Scan() {
		if cfg.Lines {
			c.Lines++
		}
		words := strings.Fields(scanner.Text())
		if cfg.Words {
			count := len(words)
			c.Words += int32(count)
		}
		if cfg.Characters && !cfg.Count {
			// Count characters as all characters not counting spaces,
			c.Characters += int64(len(strings.Join(words, " ")))
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Printf("%d %d %d %s\n", c.Lines, c.Words, func() int64 {
		if cfg.Count {
			return c.Count
		} else {
			return c.Characters
		}
	}(), filePath)

	return nil
}

func main() {
	wcCmd.Flags().BoolVarP(&configValue.Count, "count", "c", false, "The number of bytes in each input file is written to the standard output. This will cancel out any usage of the -m option.")
	wcCmd.Flags().BoolVarP(&configValue.Lines, "lines", "l", false, "The number of lines in each input file is written to the standard output.")
	wcCmd.Flags().BoolVarP(&configValue.Characters, "characters", "m", false, "The number of characters in each input file is written to the standard output. If the current locale does not support multibyte characters, this is equivalent to the -c option.")
	wcCmd.Flags().BoolVarP(&configValue.Words, "words", "w", false, "The number of words in each input file is written to the standard output.")
	// wcCmd.Flags().StringSliceVarP(&filePaths, "files", "f", []string{}, "Input file paths")

	if err := wcCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
