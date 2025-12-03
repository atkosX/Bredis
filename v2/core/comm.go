package core

import "syscall"

type FDcommand struct {
    FD int
}


func (f FDcommand) Read(b []byte) (int, error) {
    return syscall.Read(f.FD,b)   
}

func (f FDcommand) Write(b []byte) (int, error) {
    return syscall.Write(f.FD, b)
}