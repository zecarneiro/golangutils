package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golangutils/pkg/logic"
	"golangutils/pkg/obj"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
)

func ReadFile(file string) (string, error) {
	body, err := os.ReadFile(ResolvePath(file))
	if err != nil {
		return "", err
	}
	return string(body), err
}

func ReadFileLineByLine(filePath string, callback func(string, error)) {
	file, err := os.Open(ResolvePath(filePath))
	if err != nil {
		callback("", err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		callback(line, nil)
	}
	if scanner.Err() != nil {
		callback("", scanner.Err())
	}
}

func WriteFile(file string, data string, isAppend bool, isCreateDir bool) error {
	return WriteFileWithEncoding(file, data, isAppend, isCreateDir, "utf-8")
}

func WriteFileWithEncoding(file string, data string, isAppend bool, isCreateDir bool, encodingName string) error {
	utfName := "utf-8"
	file = ResolvePath(file)
	encodingName = logic.Ternary(len(encodingName) > 0, encodingName, utfName)
	encodingName = logic.Ternary(encodingName == "UTF-8", utfName, encodingName)
	if isCreateDir {
		dirname := Dirname(file)
		err := CreateDirectory(dirname, true)
		if err != nil {
			return err
		}
	}
	// 1. Obter o codificador baseado no nome fornecido
	enc, err := ianaindex.MIME.Encoding(encodingName)
	if err != nil {
		return fmt.Errorf("encoding inválido: %v", err)
	}

	var fileStream *os.File
	if isAppend && FileExist(file) {
		fileStream, err = os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0o644)
	} else {
		fileStream, err = os.Create(file)
		if encodingName == utfName {
			fileStream.Write([]byte{0xEF, 0xBB, 0xBF})
		}
	}
	if err != nil {
		return err
	}
	defer fileStream.Close()

	writer := transform.NewWriter(fileStream, enc.NewEncoder())
	if _, err := fmt.Fprintln(writer, data); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return nil
}

func DeleteFile(file string) error {
	file = ResolvePath(file)
	err := os.Remove(file)
	if err != nil {
		return err
	}
	return nil
}

func IsFile(file string) bool {
	return FileExist(file) && Type(file) == File
}

func IsSymbolicLink(file string) bool {
	return FileExist(file) && Type(file) == SymbolicLink
}

func ReadJsonFile(jsonFile string, target any) error {
	data, err := ReadFile(ResolvePath(jsonFile))
	if err != nil {
		return err
	}
	return obj.StringToObject(data, &target)
}

func WriteJsonFile(jsonFile string, data any, escapeHtml bool) error {
	var dataStr string
	var err error
	if escapeHtml {
		dataStr, err = obj.ObjectToStringEscapeHtml(data)
		if err != nil {
			return err
		}
	} else {
		dataStr, err = obj.ObjectToString(data)
		if err != nil {
			return err
		}
	}
	return WriteFile(ResolvePath(jsonFile), dataStr, false, true)
}

func CopyFile(src string, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo
	src = ResolvePath(src)
	dst = ResolvePath(dst)
	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func ReadFileInByte(filename string) ([]byte, error) {
	file, err := os.Open(ResolvePath(filename))
	byteArr := []byte{}
	if err != nil {
		file.Close()
		return byteArr, err
	}
	defer file.Close()

	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		return byteArr, err
	}

	// Read the file into a byte slice
	byteArr = make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(byteArr)
	if err != nil && err != io.EOF {
		return byteArr, err
	}
	return byteArr, nil
}

func FileExtension(data string) string {
	if data == "" || data == "." || data == ".." {
		return ""
	}
	data = ResolvePath(data)
	data = Basename(data)
	ext := filepath.Ext(data)
	if ext == "." {
		return ""
	}
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}
	return ext
}

func GetFileEncoding(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Use bufio to peek at the first 1024 bytes without consuming the reader
	reader := bufio.NewReader(file)
	data, err := reader.Peek(1024)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return "", err
	}

	// DetermineEncoding detects the encoding based on the byte stream
	// and the declared Content-Type (empty string if unknown)
	_, name, _ := charset.DetermineEncoding(data, "")
	return name, nil
}
