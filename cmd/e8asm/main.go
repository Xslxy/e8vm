package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"e8vm.io/e8vm/arch8"
	"e8vm.io/e8vm/asm8"
	"e8vm.io/e8vm/dasm8"
	"e8vm.io/e8vm/lex8"
)

var (
	doDasm      = flag.Bool("d", false, "do dump")
	ncycle      = flag.Int("n", 100000, "max cycles to execute")
	memSize     = flag.Int("m", 0, "memory size; 0 for full 4GB")
	printStatus = flag.Bool("s", false, "print status after execution")
	randSeed    = flag.Int64("seed", 0, "random seed, 0 for using the time")
)

func run(bs []byte) (int, error) {
	// create a single core machine
	m := arch8.NewMachine(&arch8.Config{
		MemSize:  uint32(*memSize),
		RandSeed: *randSeed,
	})
	if err := m.LoadImageBytes(bs); err != nil {
		return 0, err
	}

	ret, exp := m.Run(*ncycle)
	if *printStatus {
		m.PrintCoreStatus()
	}

	if exp == nil {
		return ret, nil
	}
	return ret, exp
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		log.Fatal("need exactly one input file\n")
	}

	fname := args[0]
	var bs []byte
	f, e := os.Open(fname)
	if e != nil {
		log.Fatalf("open: %s", e)
	}

	var es []*lex8.Error
	if strings.HasSuffix(fname, "_bare.s") {
		bs, es = asm8.BuildBareFunc(fname, f)
	} else {
		bs, es = asm8.BuildSingleFile(fname, f)
	}

	if len(es) > 0 {
		for _, e := range es {
			fmt.Println(e)
		}
		os.Exit(-1)
		return
	}

	if *doDasm {
		lines := dasm8.Dasm(bs, arch8.InitPC)
		for _, line := range lines {
			fmt.Println(line)
		}
	} else {
		n, e := run(bs)
		fmt.Printf("(%d cycles)\n", n)
		if e != nil {
			if !arch8.IsHalt(e) {
				fmt.Println(e)
			}
		} else {
			fmt.Println("(end of time)")
		}
	}
}
