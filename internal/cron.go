package internal

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

type Cron struct {
	interval    time.Duration
	mu          sync.Mutex
	syncedFiles map[string]struct{}
}

func NewCron(interval time.Duration) *Cron {
	return &Cron{
		interval:    interval,
		syncedFiles: make(map[string]struct{}),
	}
}

func (c *Cron) Start(ctx context.Context, origin, dest string) {
	log.Println("sync successfully started")
	fmt.Println("sync successfully started")

	ticker := time.NewTicker(c.interval)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("sync successfully stopped")
				return
			case <-ticker.C:

				originFiles, err := ioutil.ReadDir(origin)
				if err != nil {
					log.Fatal(err)
				}
				for _, file := range originFiles {
					f, err := os.ReadFile(fmt.Sprintf("%s/%s", origin, file.Name()))
					if err != nil {
						log.Fatal(err)
					}

					buffer, err := c.copyFile(file.Name(), f, dest) // произвели копирование существующих файлов
					if err != nil {
						log.Fatal(err)
					}

					err = ioutil.WriteFile(fmt.Sprintf("%s/%s", dest, file.Name()), buffer, 0777)
					if err != nil {
						log.Fatal(err)
					}

					c.mu.Lock()
					c.syncedFiles[file.Name()] = struct{}{} // регистрируем новый синхронизируемый файлов
					c.mu.Unlock()
					log.Printf("file %s synced. size: %d bytes", file.Name(), file.Size())
				}

				_, err = c.clean(c.syncedFiles, originFiles, dest)
				if err != nil {
					log.Fatalf("delete synced files error: %v", err)
				}

			}
		}
	}()
}

func (c *Cron) copyFile(fileName string, originFile []byte, dest string) ([]byte, error) {
	pipeReader, pipeWriter := io.Pipe()

	go func() {
		defer pipeWriter.Close()

		writer := bufio.NewWriter(pipeWriter)
		writer.Write(originFile)
		writer.Flush()
	}()

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(pipeReader)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *Cron) clean(syncedFiles map[string]struct{}, originFiles []fs.FileInfo, dir string) (int, error) {
	var cleaned int

	for k, _ := range syncedFiles {
		var counter int
		for _, f := range originFiles {
			if k == f.Name() {
				counter++
			}
		}
		if counter == 0 {
			err := os.Remove(fmt.Sprintf("%s/%s", dir, k))
			if err != nil {
				return cleaned, err
			}
			cleaned++
			log.Printf("file %s deleted", k)
		}
	}
	return cleaned, nil
}
