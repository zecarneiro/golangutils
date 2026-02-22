package net

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync/atomic"
	"time"

	"golangutils/pkg/conv"
	"golangutils/pkg/logger"
	"golangutils/pkg/str"
)

// ProgressWriter ajuda a calcular e exibir a percentagem
type progressWriter struct {
	Total      int64
	Downloaded int64
	lastWrite  atomic.Int64 // unix nano timestamp
}

type downloader struct {
	Url      string
	Filepath string
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.Downloaded += int64(n)
	pw.lastWrite.Store(time.Now().UnixNano())
	if pw.Total > 0 {
		// Calcula a percentagem baseada no total
		percent := float64(pw.Downloaded) / float64(pw.Total) * 100
		fmt.Printf("\rProgress: %.2f%% | Total size: %.2f MB", percent, conv.BytesToMB(pw.Total))
	} else {
		// Se o servidor não informar o tamanho, mostra apenas o que já baixou
		fmt.Printf("\rDownloading: %.2f MB (Unknown total size)", conv.BytesToMB(pw.Downloaded))
	}
	return n, nil
}

func (d *downloader) resolveFilepath() (string, error) {
	if d.Filepath != "" {
		return d.Filepath, nil
	}

	parsed, err := url.Parse(d.Url)
	if err != nil {
		return "", err
	}

	filename := path.Base(parsed.Path)

	// Caso a URL não tenha um nome de ficheiro válido
	if filename == "." || filename == "/" || str.IsEmpty(filename) || filename == ".." {
		filename = "downloaded_file"
	}

	return filename, nil
}

func Download(urlData string, output string) error {
	downloader := downloader{
		Url:      urlData,
		Filepath: output,
	}
	filepath, err := downloader.resolveFilepath()
	logger.Info("Downloading file: " + filepath)
	if err != nil {
		return err
	}
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second, // timeout de conexão
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", downloader.Url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad status: %s", resp.Status)
	}
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	pw := &progressWriter{
		Total: resp.ContentLength,
	}
	pw.lastWrite.Store(time.Now().UnixNano())

	// Watchdog de inatividade: cancela se não receber dados em 20s
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				last := time.Unix(0, pw.lastWrite.Load())
				if time.Since(last) > 20*time.Second {
					cancel()
					return
				}
			case <-ctx.Done():
				// Download terminou ou foi cancelado
				return
			}
		}
	}()

	writer := io.MultiWriter(out, pw)
	_, err = io.Copy(writer, resp.Body)
	fmt.Println()
	if err != nil {
		if err == context.Canceled {
			return fmt.Errorf("Download cancelled due to inactivity")
		}
		return err
	}
	return nil
}
