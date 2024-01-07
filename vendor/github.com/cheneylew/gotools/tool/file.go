package tool

import (
	"io/ioutil"
	"os"
	"net/http"
	"fmt"
	"io"
	"bufio"
	"encoding/json"
)

func ReadFile(path string) ([]byte, error) {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	return ioutil.ReadAll(fi)
}

func ReadFileString(path string) (string, error) {
	bytes, err := ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func WriteFile(filename string, content string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(content))
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

func WriteFileWithInterface(filename string, v interface{}) error {
	res, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return WriteFile(filename, string(res))
}


func AppendFile(filename string, content string) error {
	f, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		return err
	}

	return nil
}

func AppendFileWithInterface(filename string, v interface{}) error {
	res, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return AppendFile(filename, string(res))
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
func WriteLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

// Exist returns true if a file or directory exists.
func FileExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func DownloadFileToBytes(url string) ([]byte, error) {
	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return bodyBytes, nil
	}

	return nil, fmt.Errorf("download file! http status not ok")
}

// DownloadFile will download a url and store it in local filepath.
// It writes to the destination file as it downloads it, without
// loading the entire file into memory.
func DownloadFileToPath(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}