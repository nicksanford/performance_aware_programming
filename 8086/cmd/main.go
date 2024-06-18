package main

import (
	"flag"
	"log"
	"os"

	"github.com/nicksanford/inst"
)

func main() {
	l := log.New(os.Stderr, "", 0)
	help := flag.Bool("h", false, "help")
	simulate := flag.Bool("s", true, "simulate, defaults to true")
	flag.Parse()
	if *help {
		l.Printf("usage: %s <file>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(flag.Args()) != 1 {
		l.Printf("usage: %s <file>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	data, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		l.Fatalf(err.Error())
	}

	if *simulate {
		data, err = inst.Dasm(data)
		if err != nil {
			l.Fatal(err.Error())
		}

		lines := inst.ToLines(data)

		tokens, err := inst.Tokenize(lines)
		if err != nil {
			l.Fatal(err.Error())
		}

		result, err := inst.Simulate(tokens)
		if err != nil {
			l.Fatal(err.Error())
		}

		if _, err := os.Stdout.WriteString(result.String()); err != nil {
			l.Fatal(err.Error())
		}
	} else {
		data, err = inst.Dasm(data)
		if err != nil {
			l.Fatal(err.Error())
		}

		if _, err := os.Stdout.Write(data); err != nil {
			l.Fatal(err.Error())
		}
	}
}
