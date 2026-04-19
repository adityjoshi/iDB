package core

import "syscall"

type FileDescriptor struct {
	FD int
}

func (f *FileDescriptor) Write(b []byte) (int, error) {
	return syscall.Write(f.FD, b)
}

func (f *FileDescriptor) Read(b []byte) (int, error) {
	return syscall.Read(f.FD, b)
}
