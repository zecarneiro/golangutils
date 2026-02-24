package file

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golangutils/pkg/enums"
	"golangutils/pkg/models"
	"golangutils/pkg/obj"
	"golangutils/pkg/str"

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

func WriteFile(config models.FileWriterConfig) error {
	utf8Encoding := "utf-8"
	config.File = ResolvePath(config.File)
	if str.IsEmpty(config.EncodingName) || config.EncodingName == "UTF-8" {
		config.EncodingName = utf8Encoding
	}
	if config.IsCreateDir {
		dirname := Dirname(config.File)
		err := CreateDirectory(dirname, true)
		if err != nil {
			return err
		}
	}
	// 1. Obter o codificador baseado no nome fornecido
	enc, err := ianaindex.MIME.Encoding(config.EncodingName)
	if err != nil {
		return fmt.Errorf("encoding inválido: %v", err)
	}

	var fileStream *os.File
	if config.IsAppend && FileExist(config.File) {
		fileStream, err = os.OpenFile(config.File, os.O_APPEND|os.O_WRONLY, 0o644)
	} else {
		fileStream, err = os.Create(config.File)
		if config.WithUtf8BOM && config.EncodingName == utf8Encoding {
			fileStream.Write([]byte{0xEF, 0xBB, 0xBF})
		}
	}
	if err != nil {
		return err
	}
	defer fileStream.Close()

	var writer io.Writer = fileStream
	if config.EncodingName != utf8Encoding {
		tWriter := transform.NewWriter(fileStream, enc.NewEncoder())
		defer tWriter.Close()
		writer = tWriter
	}
	if _, err := fmt.Fprintln(writer, config.Data); err != nil {
		return err
	}
	return nil
}

func DeleteFile(file string) error {
	file = ResolvePath(file)
	if IsFile(file) {
		if err := os.Remove(file); err != nil {
			return err
		}
	}
	return nil
}

func IsFile(file string) bool {
	return FileExist(file) && Type(file) == enums.File
}

func IsSymbolicLink(file string) bool {
	return FileExist(file) && Type(file) == enums.SymbolicLink
}

func ReadJsonFile[T any](jsonFile string) (T, error) {
	var data T
	jsonFile = ResolvePath(jsonFile)
	file, err := os.Open(jsonFile)
	if err != nil {
		return data, fmt.Errorf("failed to open the file: %w", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	// "Spy" the first 3 bytes (size of the BOM UTF-8)
	bom, err := reader.Peek(3)
	if err == nil && bytes.Equal(bom, []byte("\xef\xbb\xbf")) {
		// If the BOM exists, we will delete, by move the reader in 3 bytes
		reader.Discard(3)
	}
	// Now we put the reader to the decoder
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err != nil {
		return data, fmt.Errorf("failed on decode json: %w", err)
	}

	return data, nil
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
	fileConfig := models.FileWriterConfig{
		File:        ResolvePath(jsonFile),
		Data:        dataStr,
		IsAppend:    false,
		IsCreateDir: true,
		WithUtf8BOM: false,
	}
	return WriteFile(fileConfig)
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
	if str.IsEmpty(data) || data == "." || data == ".." {
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

func FileName(data string) string {
	basename := ResolvePath(Basename(data))
	return strings.TrimSuffix(basename, filepath.Ext(basename))
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
