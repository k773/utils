package typedBuffer

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/k773/utils"
	"io"
	"os"
	"time"
)

type TypedBuffer struct {
	Buf *bytes.Buffer
}

func WrapBuffer(buf *bytes.Buffer) *TypedBuffer {
	return &TypedBuffer{Buf: buf}
}

/*
	File part
*/

// WriteFilePart
// Pass n<=1 to read full file
func (t *TypedBuffer) WriteFilePart(f io.ReaderAt, off, n int64) error {
	t.WriteInt64(off)
	t.WriteInt64(n)

	var reader io.Reader = &OffsetReader{ReaderAt: f, Offset: off}
	if n > 0 {
		reader = io.LimitReader(reader, n)
	}

	wrote, e := t.Buf.ReadFrom(reader)
	if e == nil && wrote != n {
		e = io.EOF
	}
	if e != nil {
		// If an error has occurred, change file part length to an actually wrote length && write an error after the file part
		var raw = t.Buf.Bytes()
		binary.LittleEndian.PutUint64(raw[len(raw)-int(wrote)-8:], uint64(wrote))
	}
	t.WriteError(e)
	return e
}

// ReadFilePart does not copy an underlying buffer; no extra allocations are made.
// The returned slice is valid only until the next read/write call
func (t *TypedBuffer) ReadFilePart() (dst []byte, offset, n int64, e error) {
	offset = t.ReadInt64()
	n = t.ReadInt64()
	dst = t.Buf.Next(int(n))
	e = t.ReadError()
	return
}

/*
	File info
*/

// List of full paths with file info

func (t *TypedBuffer) WriteListOfFileInfoWithFullPath(info []*utils.FileInfoWithFullPath) {
	t.WriteInt64(int64(len(info)))
	for _, i := range info {
		t.WriteFileInfoWithFullPathS(i)
	}
}

func (t *TypedBuffer) ReadListOfFileInfoWithFullPath() (info []*utils.FileInfoWithFullPath) {
	info = make([]*utils.FileInfoWithFullPath, t.ReadInt64())
	for i := range info {
		info[i] = t.ReadFileInfoWithFullPathS()
	}
	return
}

// Full file path + file info

func (t *TypedBuffer) WriteFileInfoWithFullPathS(info *utils.FileInfoWithFullPath) {
	t.WriteFileInfoWithFullPath(info.FullPath, info.FileInfo)
}

func (t *TypedBuffer) ReadFileInfoWithFullPathS() (info *utils.FileInfoWithFullPath) {
	info = new(utils.FileInfoWithFullPath)
	info.FullPath, info.FileInfo = t.ReadFileInfoWithFullPath()
	return info
}

func (t *TypedBuffer) WriteFileInfoWithFullPath(fullPath string, fi os.FileInfo) {
	t.WriteString(fullPath)
	t.WriteFileInfo(fi)
}

func (t *TypedBuffer) ReadFileInfoWithFullPath() (fullPath string, fi os.FileInfo) {
	fullPath = t.ReadString()
	fi = t.ReadFileInfo()
	return
}

// File info

func (t *TypedBuffer) WriteFileInfo(fi os.FileInfo) {
	t.WriteString(fi.Name())
	t.WriteInt64(fi.Size())
	t.WriteUint32(uint32(fi.Mode()))
	t.WriteTime(fi.ModTime())
	t.WriteBool(fi.IsDir())
}

func (t *TypedBuffer) ReadFileInfo() os.FileInfo {
	return &utils.FileInfo{
		Name_:    t.ReadString(),
		Size_:    t.ReadInt64(),
		Mode_:    os.FileMode(t.ReadUint32()),
		ModTime_: t.ReadTime(),
		IsDir_:   t.ReadBool(),
	}
}

/*
	Error
*/

func (t *TypedBuffer) WriteError(e error) {
	t.WriteBool(e != nil)
	if e != nil {
		t.WriteString(e.Error())
	}
}

// ReadError will return error containing the same text as the original error.
// However, it is not the same object as the original, and should be treated accordingly.
// Only io.EOF and io.ErrClosedPipe errors are decoded into their originals.
func (t *TypedBuffer) ReadError() error {
	isErr := t.ReadBool()
	if isErr {
		switch errMsg := t.ReadString(); errMsg {
		case io.EOF.Error():
			return io.EOF
		case io.ErrClosedPipe.Error():
			return io.ErrClosedPipe
		default:
			return errors.New(errMsg)
		}
	}
	return nil
}

/*
	String
*/

func (t *TypedBuffer) WriteString(s string) {
	t.WriteInt64(int64(len(s)))
	t.Buf.Write([]byte(s))
}

func (t *TypedBuffer) ReadString() string {
	return string(t.Buf.Next(int(t.ReadInt64())))
}

/*
	Bool
*/

func (t *TypedBuffer) WriteBool(i bool) {
	t.Buf.WriteByte(utils.If(i, byte(1), byte(0)))
}

func (t *TypedBuffer) ReadBool() bool {
	b, e := t.Buf.ReadByte()
	if e != nil {
		panic(e)
	}
	return b == 1
}

/*
	Numbers
*/

// time.Time

func (t *TypedBuffer) WriteTime(i time.Time) {
	WriteNumber(t.Buf, i.UnixNano())
}

func (t *TypedBuffer) ReadTime() time.Time {
	return time.Unix(0, ReadNumber[int64](t.Buf))
}

// time.Duration

func (t *TypedBuffer) WriteDuration(i time.Duration) {
	WriteNumber(t.Buf, uint64(i))
}

func (t *TypedBuffer) ReadDuration() time.Duration {
	return time.Duration(ReadNumber[uint64](t.Buf))
}

// Uint64

func (t *TypedBuffer) WriteUint64(i uint64) {
	WriteNumber(t.Buf, i)
}

func (t *TypedBuffer) ReadUInt64() uint64 {
	return ReadNumber[uint64](t.Buf)
}

// Uint32

func (t *TypedBuffer) WriteUint32(i uint32) {
	WriteNumber(t.Buf, i)
}

func (t *TypedBuffer) ReadUint32() uint32 {
	return ReadNumber[uint32](t.Buf)
}

// Uint16

func (t *TypedBuffer) WriteUint16(i uint16) {
	WriteNumber(t.Buf, i)
}

func (t *TypedBuffer) ReadUint16() uint16 {
	return ReadNumber[uint16](t.Buf)
}

// Uint8

func (t *TypedBuffer) WriteUint8(i uint8) {
	WriteNumber(t.Buf, i)
}

func (t *TypedBuffer) ReadUint8() uint8 {
	return ReadNumber[uint8](t.Buf)
}

// int64

func (t *TypedBuffer) WriteInt64(i int64) {
	WriteNumber(t.Buf, i)
}

func (t *TypedBuffer) ReadInt64() int64 {
	return ReadNumber[int64](t.Buf)
}

// int32

func (t *TypedBuffer) WriteInt32(i int32) {
	WriteNumber(t.Buf, i)
}

func (t *TypedBuffer) ReadInt32() int32 {
	return ReadNumber[int32](t.Buf)
}

// int16

func (t *TypedBuffer) WriteInt16(i int16) {
	WriteNumber(t.Buf, i)
}

func (t *TypedBuffer) ReadInt16() int16 {
	return ReadNumber[int16](t.Buf)
}

// int8

func (t *TypedBuffer) WriteInt8(i int8) {
	WriteNumber(t.Buf, i)
}

func (t *TypedBuffer) ReadInt8() int8 {
	return ReadNumber[int8](t.Buf)
}

/*
	Generics
*/

func WriteNumber[T utils.Ints | utils.Uints](b *bytes.Buffer, i T) {
	_ = binary.Write(b, binary.LittleEndian, i)
}

func ReadNumber[T utils.Ints | utils.Uints](b *bytes.Buffer) T {
	var dst = new(T)
	_ = binary.Read(b, binary.LittleEndian, dst)
	return *dst
}
