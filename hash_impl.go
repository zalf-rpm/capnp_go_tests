package main

import (
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"net"
	"os"
	"syscall"

	"capnproto.org/go/capnp/v3/rpc"
	"capnproto.org/go/capnp/v3/server"
	"example.com/capnp_schemas/hashes"
	"golang.org/x/net/context"
)

type listenCall struct {
	c   net.Conn
	err error
}

func tcpPipe() (t1, t2 net.Conn, err error) {
	host, err := net.LookupIP("localhost")
	if err != nil {
		return nil, nil, err
	}
	l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: host[0]})
	if err != nil {
		return nil, nil, err
	}
	ch := make(chan listenCall)
	abort := make(chan struct{})
	go func() {
		c, err := l.AcceptTCP()
		select {
		case ch <- listenCall{c, err}:
		case <-abort:
			c.Close()
		}
	}()
	laddr := l.Addr().(*net.TCPAddr)
	c2, err := net.DialTCP("tcp", nil, laddr)
	if err != nil {
		close(abort)
		l.Close()
		return nil, nil, err
	}
	lc := <-ch
	if lc.err != nil {
		c2.Close()
		l.Close()
		return nil, nil, err
	}
	return lc.c, c2, nil
}

func tcpServerConns(addr string) (ch chan listenCall, abort chan struct{}, err error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, err
	}
	ch = make(chan listenCall)
	abort = make(chan struct{})
	go func() {
		defer l.Close()
		for {
			c, err := l.Accept()
			select {
			case ch <- listenCall{c, err}:
			case <-abort:
				fmt.Println("abort closed")
				break
			}
		}
	}()
	return
}

func tcpClientConn(addr string) (c net.Conn, err error) {
	c, err = net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return
}

func chkfatal(err error) {
	if err != nil {
		panic(err)
	}
}

func socketpairPipe() (net.Conn, net.Conn) {
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	chkfatal(err)
	mkConn := func(fd int, name string) net.Conn {
		conn, err := net.FileConn(os.NewFile(uintptr(fd), name))
		chkfatal(err)
		return conn
	}
	return mkConn(fds[0], "pipe0"), mkConn(fds[1], "pipe1")
}

type hello struct{}

func (h hello) Hello(_ context.Context, call hashes.Hello_hello) error {
	res, err := call.AllocResults()
	if err != nil {
		return err
	}
	return res.SetT("Hello")
}

func (h hello) World(_ context.Context, call hashes.Hello_world) error {
	res, err := call.AllocResults()
	if err != nil {
		return err
	}
	return res.SetT("World!")
}

func serveHello(ctx context.Context, rwc io.ReadWriteCloser) error {
	// Create a new locally implemented HashFactory.
	main := hashes.Hello_ServerToClient(hello{}, nil) //&server.Policy{})

	// Listen for calls, using the HashFactory as the bootstrap interface.
	conn := rpc.NewConn(rpc.NewStreamTransport(rwc), &rpc.Options{
		ErrorReporter:   eR{},
		BootstrapClient: main.Client,
	})
	defer conn.Close()

	// Wait for connection to abort.
	select {
	case <-conn.Done():
		fmt.Print("conn.Done()")
		return nil
	case <-ctx.Done():
		fmt.Print("conn.Close()")
		return conn.Close()
	}
}

func clientHello(ctx context.Context, rwc io.ReadWriteCloser) error {
	// Create a connection that we can use to get the HashFactory.
	conn := rpc.NewConn(rpc.NewStreamTransport(rwc), &rpc.Options{ErrorReporter: eR{}}) // nil sets default options
	defer conn.Close()

	// Get the "bootstrap" interface.  This is the capability set with
	// rpc.MainInterface on the remote side.
	hello := hashes.Hello{Client: conn.Bootstrap(ctx)}

	hf, free := hello.Hello(ctx, nil)
	defer free()
	wf, free := hello.World(ctx, nil)
	hr, err := hf.Struct()
	if err != nil {
		return err
	}
	h, _ := hr.T()
	wr, err := wf.Struct()
	if err != nil {
		return err
	}
	w, _ := wr.T()
	fmt.Println(h, " ", w)
	return nil
}

// hashFactory is a local implementation of HashFactory.
type hashFactory struct{}

func (hf *hashFactory) NewSha1(_ context.Context, call hashes.HashFactory_newSha1) error {
	// Create a new locally implemented Hash capability.
	hsc := hashes.Hash_ServerToClient(&hashServer{h: sha1.New()}, &server.Policy{})

	// Notice that methods can return other interfaces.
	res, err := call.AllocResults()
	if err != nil {
		return err
	}
	return res.SetHash(hsc)
}

// hashServer is a local implementation of Hash.
type hashServer struct {
	h hash.Hash
	a int
}

func (hs *hashServer) Write(_ context.Context, call hashes.Hash_write) error {
	data, err := call.Args().Data()
	if err != nil {
		return err
	}
	hs.a += 1

	_, err = hs.h.Write(data)
	return err
}

func (hs *hashServer) Sum(_ context.Context, call hashes.Hash_sum) error {
	res, err := call.AllocResults()
	if err != nil {
		return err
	}
	fmt.Println("a:", hs.a)

	b := hs.h.Sum(nil)
	return res.SetHash(b)
}

func serveHash(ctx context.Context, rwc io.ReadWriteCloser) error {
	// Create a new locally implemented HashFactory.
	main := hashes.HashFactory_ServerToClient(&hashFactory{}, nil) //&server.Policy{})

	// Listen for calls, using the HashFactory as the bootstrap interface.
	conn := rpc.NewConn(rpc.NewStreamTransport(rwc), &rpc.Options{
		ErrorReporter:   eR{},
		BootstrapClient: main.Client,
	})
	defer conn.Close()

	// Wait for connection to abort.
	select {
	case <-conn.Done():
		fmt.Print("conn.Done()")
		return nil
	case <-ctx.Done():
		fmt.Print("conn.Close()")
		return conn.Close()
	}
}

type eR struct{}

func (er eR) ReportError(err error) {
	fmt.Println("ReportError:", err)
}

func clientHash(ctx context.Context, rwc io.ReadWriteCloser) error {
	// Create a connection that we can use to get the HashFactory.
	conn := rpc.NewConn(rpc.NewStreamTransport(rwc), &rpc.Options{ErrorReporter: eR{}}) // nil sets default options
	defer conn.Close()

	// Get the "bootstrap" interface.  This is the capability set with
	// rpc.MainInterface on the remote side.
	hf := hashes.HashFactory{Client: conn.Bootstrap(ctx)}

	// Now we can call methods on hf, and they will be sent over c.
	// The NewSha1 method does not have any parameters we can set, so we
	// pass a nil function.
	sha1f, free := hf.NewSha1(ctx, nil)
	defer free()
	sr, err := sha1f.Struct()
	if err != nil {
		fmt.Println("NewSha1 ", err)
		return err
	}
	s := sr.Hash()

	// 'NewSha1' returns a future, which allows us to pipeline calls to
	// returned values before they are actually delivered.  Here, we issue
	// calls to an as-of-yet-unresolved Sha1 instance.
	//s := sha1f.Hash()

	// s refers to a remote Hash.  Method calls are delivered in order.
	f1, free := s.Write(ctx, func(p hashes.Hash_write_Params) error {
		fmt.Println("first s.Write")
		return p.SetData([]byte("Hello, "))
	})
	defer free()
	_, err = f1.Struct()
	if err != nil {
		fmt.Println("first Write ", err)
		return err
	}
	f2, free := s.Write(ctx, func(p hashes.Hash_write_Params) error {
		fmt.Println("second s.Write")
		return p.SetData([]byte("World!"))
	})
	defer free()
	_, err = f2.Struct()
	if err != nil {
		fmt.Println("second Write ", err)
		return err
	}

	// Get the sum, waiting for the result.
	sumf, free := s.Sum(ctx, nil)
	defer free()
	result, err := sumf.Struct()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Display the result.
	sha1Val, err := result.Hash()
	if err != nil {
		return err
	}
	fmt.Printf("sha1: %x\n", sha1Val)
	return nil
}

/*
func testTCPStreamTransport(newTransport func(io.ReadWriteCloser) rpc.Transport) {
	type listenCall struct {
		c   *net.TCPConn
		err error
	}

	makePipe := func() (t1, t2 rpc.Transport, err error) {
		host, err := net.LookupIP("localhost")
		if err != nil {
			return nil, nil, err
		}
		l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: host[0]})
		if err != nil {
			return nil, nil, err
		}
		ch := make(chan listenCall)
		abort := make(chan struct{})
		go func() {
			c, err := l.AcceptTCP()
			select {
			case ch <- listenCall{c, err}:
			case <-abort:
				c.Close()
			}
		}()
		laddr := l.Addr().(*net.TCPAddr)
		c2, err := net.DialTCP("tcp", nil, laddr)
		if err != nil {
			close(abort)
			l.Close()
			return nil, nil, err
		}
		lc := <-ch
		if lc.err != nil {
			c2.Close()
			l.Close()
			return nil, nil, err
		}
		return newTransport(lc.c), newTransport(c2), nil
	}

	t.Run("ServerToClient", func(t *testing.T) {
		testTransport(t, makePipe)
	})

	t.Run("ClientToServer", func(t *testing.T) {
		testTransport(t, func() (t1, t2 rpc.Transport, err error) {
			t2, t1, err = makePipe()
			return
		})
	})
}
*/
