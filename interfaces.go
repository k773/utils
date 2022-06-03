package utils

import (
	"bytes"
	"os"
	"time"
)

type PacketReadWriter interface {
	ReadPacket() (*bytes.Buffer, error)
	WritePacket(data *bytes.Buffer) error
}

type FileInfoWithFullPath struct {
	os.FileInfo
	FullPath string
}

// FileInfo implements os.FileInfo.
// It may be needed because Go's builtin os.FileInfo is an interface and os.fileStat is a non-exported platform-dependent structure
type FileInfo struct {
	Name_    string
	Size_    int64
	Mode_    os.FileMode
	ModTime_ time.Time
	IsDir_   bool
}

func (f *FileInfo) Name() string {
	return f.Name_
}
func (f *FileInfo) Size() int64 {
	return f.Size_
}
func (f *FileInfo) Mode() os.FileMode {
	return f.Mode_
}
func (f *FileInfo) ModTime() time.Time {
	return f.ModTime_
}
func (f *FileInfo) IsDir() bool {
	return f.IsDir_
}
func (f *FileInfo) Sys() any {
	return nil
}
