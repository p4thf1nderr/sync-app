package internal

import (
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCopyFile(t *testing.T) {

	cron := NewCron(time.Millisecond)

	req := require.New(t)

	input := []byte("test")

	output, err := cron.copyFile("test", input, "test")

	req.NoError(err)
	req.Equal(input, output)
}

func TestClean(t *testing.T) {

	req := require.New(t)

	cron := NewCron(time.Millisecond)

	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Error(err)
	}

	input := []byte("test")

	t.Run("successfull clean", func(t *testing.T) {
		if err := ioutil.WriteFile(filepath.Join(dir, "test.txt"), input, 0666); err != nil {
			t.Error(err)
		}
		syncedFiles := map[string]struct{}{
			"test.txt": {},
		}

		cleaned, err := cron.clean(syncedFiles, []fs.FileInfo{}, dir)

		req.Equal(1, cleaned)
		req.NoError(err)
		syncedFiles = nil
	})

	t.Run("nothing to clean", func(t *testing.T) {
		if err := ioutil.WriteFile(filepath.Join(dir, "test.txt"), input, 0666); err != nil {
			t.Error(err)
		}
		syncedFiles := map[string]struct{}{}

		cleaned, err := cron.clean(syncedFiles, []fs.FileInfo{}, dir)

		req.Equal(0, cleaned)
		req.NoError(err)

	})
}

func BenchmarkCopy(b *testing.B) {
	cron := NewCron(time.Millisecond)

	for i := 0; i < b.N; i++ {

		input := []byte("test")

		cron.copyFile("test", input, "test")
	}
}
