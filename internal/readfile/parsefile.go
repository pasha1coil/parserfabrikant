package readfile

import (
	"bufio"
	"os"
)

// ReadFile открытие файла с данными и его чтение
func ReadFile(path string, ch chan string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ch <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
