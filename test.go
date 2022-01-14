package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"example.com/capnp_schemas/books"

	"capnproto.org/go/capnp/v3"
	"example.com/greetings"
	"rsc.io/quote"
)

func main() {

	args := os.Args[1:]
	switch args[0] {

	case "send":
		// Make a brand new empty message.  A Message allocates Cap'n Proto structs.
		msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
		if err != nil {
			panic(err)
		}

		// Create a new Book struct.  Every message must have a root struct.
		book, err := books.NewRootBook(seg)
		if err != nil {
			panic(err)
		}
		book.SetTitle("War and Peace")
		book.SetPageCount(1440)

		// Write the message to stdout.
		err = capnp.NewEncoder(os.Stdout).Encode(msg)
		if err != nil {
			panic(err)
		}
	case "recv":
		// Read the message from stdin.
		msg, err := capnp.NewDecoder(os.Stdin).Decode()
		if err != nil {
			panic(err)
		}

		// Extract the root struct from the message.
		book, err := books.ReadRootBook(msg)
		if err != nil {
			panic(err)
		}

		// Access fields from the struct.
		title, err := book.Title()
		if err != nil {
			panic(err)
		}
		pageCount := book.PageCount()
		fmt.Printf("%q has %d pages\n", title, pageCount)
	case "hash":
		ctx := context.Background()
		c1, c2, err := tcpPipe()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		go serveHash(ctx, c1)
		err = clientHash(ctx, c2)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "hello":
		ctx := context.Background()
		c1, c2, _ := tcpPipe()
		go serveHello(ctx, c1)
		err := clientHello(ctx, c2)
		if err != nil {
			panic(err)
		}
	case "chello":
		ctx := context.Background()
		c, err := tcpClientConn("localhost:8000")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		//err = clientHello(ctx, c)
		err = clientHash(ctx, c)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "shello":
		ctx := context.Background()
		ch, abort, err := tcpServerConns(":8000")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for {
			select {
			case lc := <-ch:
				if lc.err != nil {
					close(abort)
				}
				//go serveHello(ctx, lc.c)
				go serveHash(ctx, lc.c)
			}
		}
	default:
		// Set properties of the predefined Logger, including
		// the log entry prefix and a flag to disable printing
		// the time, source file, and line number.
		log.SetPrefix("greetings: ")
		log.SetFlags(0)

		// Request a greeting message.
		message, err := greetings.Hello("")
		// If an error was returned, print it to the console and
		// exit the program.
		if err != nil {
			log.Fatal(err)
		}

		// If no error was returned, print the returned message
		// to the console.
		fmt.Println(message)
		fmt.Println(quote.Go())
	}
}
