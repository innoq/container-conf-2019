package main

// loosely based on https://github.com/lizrice/containers-from-scratch

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)


func main() {
	if len(os.Args) == 2 {
		clone()
	} else if len(os.Args) == 3 && os.Args[1] == "clone" {
		containerizeBash()
	} else {
		panic("I need somebody, help!")
	}
}

func clone() {
	fmt.Println("Cloning process with new namespaces")

	cmd := exec.Command(os.Args[0], "clone", os.Args[1])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	proceedOrPanic(cmd.Run())
}

func containerizeBash() {

	// How much memory should be assigned
	size, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	fmt.Println("Setting up cgroup memory limit")

	cgroups := "/sys/fs/cgroup/"
	memory := filepath.Join(cgroups, "memory")
	os.Mkdir(filepath.Join(memory, "containerized_bash"), 0755)
	proceedOrPanic(ioutil.WriteFile(filepath.Join(memory, "containerized_bash/memory.limit_in_bytes"), []byte(strconv.Itoa(1024*1024*size)), 0700))

	// Notify on container exit
	proceedOrPanic(ioutil.WriteFile(filepath.Join(memory, "containerized_bash/notify_on_release"), []byte("1"), 0700))
	
	// Add process to cgroup
	proceedOrPanic(ioutil.WriteFile(filepath.Join(memory, "containerized_bash/tasks"), []byte(strconv.Itoa(os.Getpid())), 0700))
	
	cmd := exec.Command("/bin/bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Chrooting environment")

	proceedOrPanic(syscall.Sethostname([]byte("container")))
	// this path must be adapted to your local setup
	proceedOrPanic(syscall.Chroot("ContainerConf2019/examples/full"))
	proceedOrPanic(os.Chdir("/"))
	proceedOrPanic(syscall.Mount("proc", "proc", "proc", 0, ""))

	proceedOrPanic(cmd.Run())

	proceedOrPanic(syscall.Unmount("proc", 0))
}

func proceedOrPanic(err error) {
	if err != nil {
		panic(err)
	}
}