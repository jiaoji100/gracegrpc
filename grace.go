/**
 * @Author: jiaoji1@staff.weibo.com
 * @Date: 2020/5/28 2:19 下午
 * @Description: 平滑重启 gRPC server
 */
package gracegrpc

import (
	"os"
	"fmt"
	"log"
	"net"
	"os/exec"
	"syscall"
	"os/signal"

	"google.golang.org/grpc"
)

const listenerFilename = "LISTENER-FILENAME"

var (
	didInherit = os.Getenv(listenerFilename) != ""
	ppid       = os.Getppid()
)
type app struct {
	listener net.Listener
	server   *grpc.Server
	addr     string
}

func newApp(server *grpc.Server, addr string) (a *app, err error) {
	a = &app{
		server: server,
		addr:   addr,
	}

	// 创建或继承listener
	a.listener, err = inheritOrCreateListener(a.addr)
	if err != nil {
		return nil, fmt.Errorf("create or import failed,err:", err)
	}

	return a, nil
}

func Serve(server *grpc.Server, addr string) error {

	a, err := newApp(server, addr)
	if err != nil {
		return err
	}

	return a.run()
}

func (a *app) run() (err error) {

	//serve
	go func() {
		err := a.server.Serve(a.listener)
		if err != nil {
			log.Printf("gRPC server start failed,err:%v\n", err)
			panic(err)
		}
	}()
	// Close the parent if we inherited and it wasn't init that started us.
	if didInherit && ppid != 1 {
		if err := syscall.Kill(ppid, syscall.SIGTERM); err != nil {
			return fmt.Errorf("failed to close parent: %s", err)
		}
	}

	a.waitForSignals()
	return nil
}

// 继承或创建listener
func inheritOrCreateListener(addr string) (net.Listener, error) {
	if didInherit {
		return inheritListener()
	}

	return createListener(addr)
}

func createListener(addr string) (net.Listener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return ln, nil
}

func inheritListener() (net.Listener, error) {
	// 从环境变量获取listener信息
	listenerFilename := os.Getenv(listenerFilename)
	if listenerFilename == "" {
		return nil, fmt.Errorf("unable to find LISTENER environment variable")
	}

	// 根据环境变量中的文件名和描述符，创建一个新的文件
	listenerFile := os.NewFile(uintptr(3), listenerFilename)
	if listenerFile == nil {
		return nil, fmt.Errorf("unable to create listener file : %s", listenerFilename)
	}
	defer listenerFile.Close()

	// 根据创建的文件，创建listener
	ln, err := net.FileListener(listenerFile)
	if err != nil {
		return nil, err
	}

	return ln, nil
}

func (a *app) waitForSignals() {
	signalCh := make(chan os.Signal, 10)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	for {
		s := <-signalCh
		log.Printf("%v signal received.\n", s)
		switch s {
		case syscall.SIGUSR2:
			child, err := forkChild(a.listener)
			if err != nil {
				log.Printf("Unable to fork child: %v.\n", err)
				continue
			}
			log.Printf("Forked child %v\n", child.Pid)

		case syscall.SIGINT, syscall.SIGTERM:
			signal.Stop(signalCh)
			a.server.GracefulStop()
			log.Printf("Receive quit signal and quit %v \n", os.Getpid())
			return
		}
	}
}

var originalWD, _ = os.Getwd()

func forkChild(ln net.Listener) (*os.Process, error) {
	// 获取当前进程的listener的文件描述符
	lnFile, err := getListenerFile(ln)
	if err != nil {
		return nil, err
	}
	defer lnFile.Close()

	//当前进程的listener的文件描述符名字通过环境变量传递给子进程
	environment := append(os.Environ(), listenerFilename+"="+lnFile.Name())

	argv0, err := exec.LookPath(os.Args[0])
	if err != nil {
		return nil, err
	}

	// 标准输入、标准输出、标准错误输出、当前进程的listener的文件描述符，4个文件传递给进程
	files := []*os.File{os.Stdin, os.Stdout, os.Stderr, lnFile}

	// 启动子进程
	child, err := os.StartProcess(argv0, os.Args, &os.ProcAttr{
		Dir:   originalWD,
		Env:   environment,
		Files: files,
		Sys:   &syscall.SysProcAttr{},
	})

	return child, err
}

func getListenerFile(ln net.Listener) (*os.File, error) {
	tcpListener, ok := ln.(*net.TCPListener)
	if !ok {
		return nil, fmt.Errorf("unsupported listener: %T", ln)
	}
	return tcpListener.File()
}
