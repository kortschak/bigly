package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	arg "github.com/alexflint/go-arg"
	"github.com/biogo/hts/sam"
	"github.com/brentp/bigly"
	"github.com/brentp/faidx"
)

type cliarg struct {
	bigly.Options
	Reference string `arg:"-r,help:optional path to reference fasta."`
	BamPath   string `arg:"positional,required"`
	Region    string `arg:"positional,required"`
}

func (c cliarg) Version() string {
	return "bigly 0.2.0"
}

func main() {
	cli := &cliarg{}
	cli.Options.MinBaseQuality = 10
	cli.Options.MinMappingQuality = 5
	cli.Options.MinClipLength = 15
	arg.MustParse(cli)
	if cli.ExcludeFlag == 0 {
		cli.ExcludeFlag = uint16(sam.Unmapped | sam.QCFail | sam.Duplicate)
	}
	/*
		f, err := os.Create("bigly.cpu.pprof")
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	*/

	stdout := bufio.NewWriter(os.Stdout)
	defer stdout.Flush()

	chromse := strings.Split(cli.Region, ":")
	se := strings.Split(chromse[1], "-")

	start, err := strconv.Atoi(se[0])
	if err != nil {
		log.Fatal(err)
	}
	end, err := strconv.Atoi(se[1])
	if err != nil {
		log.Fatal(err)
	}

	var ref *faidx.Faidx
	if cli.Reference != "" {
		var err error
		ref, err = faidx.New(cli.Reference)
		if err != nil {
			log.Fatal(err)
		}
	}

	it := bigly.Up(cli.BamPath, cli.Options, bigly.Position{chromse[0], start - 1, end}, ref)
	for it.Next() {
		p := it.Pile()
		fmt.Fprintln(stdout, p.TabString(cli.Options))
	}
	if err := it.Error(); err != nil {
		log.Fatal(err)
	}

}
