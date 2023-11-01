package testread

import (
	"os"
	"reflect"
	"testing"
	"testingParser/internal/readfile"
)

func TestReadFile(t *testing.T) {
	// Создаем временный файл для тестирования
	tmpFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	// удаляем файл по завершению
	defer os.Remove(tmpFile.Name())
	// закрываем по завершению
	defer tmpFile.Close()

	// Записываем тестовые данные в файл
	testData := "test 1\ntest 2\ntest 3"
	tmpFile.WriteString(testData)

	// Создаем канал для передачи данных
	ch := make(chan string, 3)

	// Запускаем функцию ReadFile в горутине
	go func() {
		err := readfile.ReadFile(tmpFile.Name(), ch)
		if err != nil {
			t.Fatalf("ReadFile error: %v", err)
		}
		close(ch)
	}()

	// Собираем данные из канала
	var lines []string
	for line := range ch {
		lines = append(lines, line)
	}

	// Ожидаем, что считанные строки соответствуют тестовым данным
	expectedLines := []string{"test 1", "test 2", "test 3"}
	if !reflect.DeepEqual(lines, expectedLines) {
		t.Errorf("ReadFile did not produce the expected output. Got: %v, Expected: %v", lines, expectedLines)
	}
}
