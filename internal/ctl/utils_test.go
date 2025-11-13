package ctl

import (
	"errors"
	"io/fs"
	"testing"
	"time"

	"github.com/etoneja/go-keeper/internal/ctl/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock файловой информации
type mockFileInfo struct {
	size int64
	name string
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return m.size }
func (m mockFileInfo) Mode() fs.FileMode  { return 0644 }
func (m mockFileInfo) ModTime() time.Time { return time.Now() }
func (m mockFileInfo) IsDir() bool        { return false }
func (m mockFileInfo) Sys() interface{}   { return nil }

func TestFileChecker_CheckFileSize(t *testing.T) {
	t.Run("valid file size", func(t *testing.T) {
		checker := &FileChecker{
			statFunc: func(name string) (fs.FileInfo, error) {
				return mockFileInfo{
					size: constants.MaxFileSize - 100,
					name: "testfile.txt",
				}, nil
			},
		}

		err := checker.CheckFileSize("testfile.txt")
		assert.NoError(t, err)
	})

	t.Run("file too large", func(t *testing.T) {
		checker := &FileChecker{
			statFunc: func(name string) (fs.FileInfo, error) {
				return mockFileInfo{
					size: constants.MaxFileSize + 100,
					name: "largefile.bin",
				}, nil
			},
		}

		err := checker.CheckFileSize("largefile.bin")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "file too large")
		assert.Contains(t, err.Error(), "bytes (max:")
	})

	t.Run("file not found", func(t *testing.T) {
		checker := &FileChecker{
			statFunc: func(name string) (fs.FileInfo, error) {
				return nil, fs.ErrNotExist
			},
		}

		err := checker.CheckFileSize("nonexistent.txt")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get file info")
		assert.ErrorIs(t, err, fs.ErrNotExist)
	})

	t.Run("exactly max size", func(t *testing.T) {
		checker := &FileChecker{
			statFunc: func(name string) (fs.FileInfo, error) {
				return mockFileInfo{
					size: constants.MaxFileSize,
					name: "maxfile.dat",
				}, nil
			},
		}

		err := checker.CheckFileSize("maxfile.dat")
		assert.NoError(t, err)
	})

	t.Run("permission denied", func(t *testing.T) {
		checker := &FileChecker{
			statFunc: func(name string) (fs.FileInfo, error) {
				return nil, errors.New("permission denied")
			},
		}

		err := checker.CheckFileSize("/root/protected.txt")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get file info")
		assert.Contains(t, err.Error(), "permission denied")
	})

	t.Run("zero size file", func(t *testing.T) {
		checker := &FileChecker{
			statFunc: func(name string) (fs.FileInfo, error) {
				return mockFileInfo{
					size: 0,
					name: "empty.txt",
				}, nil
			},
		}

		err := checker.CheckFileSize("empty.txt")
		assert.NoError(t, err)
	})

	t.Run("directory instead of file", func(t *testing.T) {
		checker := &FileChecker{
			statFunc: func(name string) (fs.FileInfo, error) {
				return mockFileInfo{
					size: 4096,
					name: "somedir",
				}, nil
			},
		}

		err := checker.CheckFileSize("somedir")
		assert.NoError(t, err) // Размер директории тоже можно проверить
	})
}

func TestNewFileChecker(t *testing.T) {
	checker := NewFileChecker()

	// Проверяем что используется реальный os.Stat
	// Создаем временный файл для проверки
	tmpFile := t.TempDir() + "/testfile"

	// Должна вернуться ошибка "file not found"
	err := checker.CheckFileSize(tmpFile)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get file info")
}
