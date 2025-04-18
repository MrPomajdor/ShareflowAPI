package filesystem

import (
	"io"
	"mime/multipart"
	"os"

	"github.com/MrPomajdor/ShareFlowAPI/internal/entity"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	os.IsNotExist(err)
	return err == nil
}

func WriteFile(file multipart.File, file_struct *entity.FileNode, path string) error {
	// TODO : Implement file encryption

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0o664)
	defer file.Close()
	defer f.Close()
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err_f := io.Copy(f, file); err_f != nil {
		return err
	}
	return nil
}

func RemoveFile(path string) error {
	return os.Remove(path)
}

func CreateDirectory(path string) error {
	return os.Mkdir(path, 0o775)
}
