package files

import (
	"os"
)

type TmpDir struct {
	Path string
}

func NewTmpDir(root, pattern string) (t *TmpDir, e error) {
	t = new(TmpDir)
	_ = os.MkdirAll(root, 0600)
	t.Path, e = os.MkdirTemp(root, pattern)
	return
}

func (t *TmpDir) Done() {
	_ = os.RemoveAll(t.Path)
}
