package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"math/big"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// get my args.
	flag.Parse()
	myargs := flag.Args()

	// make the ps args with "u" option.
	// set the "u" option means show memory usage.
	psargs := "u" + strings.Join(myargs, "")

	// prepare the ps command.
	cmd := exec.Command("ps", psargs)
	// for get raw stdout
	var rawout bytes.Buffer
	cmd.Stdout = &rawout

	// execute the ps command.
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// stdout buffer for human-readable VSZ and RSS field
	// as like "2.6Gi"
	var hrout bytes.Buffer

	// initialize binary prefix unit.
	initbpu()

	// get header line and keep VSZ and RSS field position.
	headerline, err := rawout.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	hrout.WriteString(headerline)
	vszpos, rsspos := getVszAndRssPos(headerline)

	// get list line and format VSZ and RSS filed to human-readable
	for {
		listline, err := rawout.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
			break
		}
		if listline == headerline {
			hrout.WriteString(listline)
			continue
		}
		listfields := strings.Fields(listline)

		vszrawbyte, rssrawbyte := listfields[vszpos], listfields[rsspos]
		vszfmtbyte, rssfmtbyte := fmtByte(vszrawbyte), fmtByte(rssrawbyte)

		vszsplitline := strings.SplitN(listline, vszrawbyte, 2)
		vszfmtline := strings.Join(vszsplitline, vszfmtbyte)

		rsssplitline := strings.SplitN(vszfmtline, rssrawbyte, 2)
		vssandrssfmtline := strings.Join(rsssplitline, rssfmtbyte)

		hrout.WriteString(vssandrssfmtline)
	}

	// show human-readable stdout.
	hrout.WriteTo(os.Stdout)

}

// get header line and keep VSZ and RSS field position.
func getVszAndRssPos(headerline string) (vszpos int, rsspos int) {
	vszpos, rsspos = -1, -1
	headerfields := strings.Fields(headerline)
	for index, name := range headerfields {
		switch {
		case name == "VSZ":
			vszpos = index
		case name == "RSS":
			rsspos = index
		}
	}
	return vszpos, rsspos
}

// convert raw byte to easy-to-read format.
// style: digit 3 char ##0.# + unit 2 char [BKMGTPEZY]i
// ex: in  16          byte, out " 16Bi"
// ex: in 512 * 1024^1 byte, out "512Ki"
// ex: in   2 * 1024^3 byte, out "2.0Gi"
func fmtByte(rawbyte string) string {
	rawbi, _ := new(big.Int).SetString(rawbyte, 10)
	var fmtint, fmtdec, fmtsuf string
	switch {
	case bpu.KiB.Cmp(rawbi) == 1: // v < KiB then Bi
		fmtint, fmtdec, fmtsuf = rawbyte, "0", "Bi"
	case bpu.MiB.Cmp(rawbi) == 1: // v < MiB then Ki
		div, mod := new(big.Int).DivMod(rawbi, bpu.KiB, new(big.Int))
		fmtint, fmtdec, fmtsuf = div.String(), mod.String()[0:1], "Ki"
	case bpu.GiB.Cmp(rawbi) == 1: // v < GiB then Mi
		div, mod := new(big.Int).DivMod(rawbi, bpu.MiB, new(big.Int))
		fmtint, fmtdec, fmtsuf = div.String(), mod.String()[0:1], "Mi"
	case bpu.TiB.Cmp(rawbi) == 1: // v < TiB then Gi
		div, mod := new(big.Int).DivMod(rawbi, bpu.GiB, new(big.Int))
		fmtint, fmtdec, fmtsuf = div.String(), mod.String()[0:1], "Gi"
	case bpu.PiB.Cmp(rawbi) == 1: // v < PiB then Ti
		div, mod := new(big.Int).DivMod(rawbi, bpu.TiB, new(big.Int))
		fmtint, fmtdec, fmtsuf = div.String(), mod.String()[0:1], "Ti"
	case bpu.EiB.Cmp(rawbi) == 1: // v < EiB then Pi
		div, mod := new(big.Int).DivMod(rawbi, bpu.PiB, new(big.Int))
		fmtint, fmtdec, fmtsuf = div.String(), mod.String()[0:1], "Pi"
	case bpu.ZiB.Cmp(rawbi) == 1: // v < ZiB then Ei
		div, mod := new(big.Int).DivMod(rawbi, bpu.EiB, new(big.Int))
		fmtint, fmtdec, fmtsuf = div.String(), mod.String()[0:1], "Ei"
	case bpu.YiB.Cmp(rawbi) == 1: // v < YiB then Zi
		div, mod := new(big.Int).DivMod(rawbi, bpu.ZiB, new(big.Int))
		fmtint, fmtdec, fmtsuf = div.String(), mod.String()[0:1], "Zi"
	case bpu.YiB.Cmp(rawbi) < 1: // v >= YiB then Yi
		div, mod := new(big.Int).DivMod(rawbi, bpu.YiB, new(big.Int))
		fmtint, fmtdec, fmtsuf = div.String(), mod.String()[0:1], "Yi"
	default:
		fmtint, fmtdec, fmtsuf = rawbyte, "0", ""
	}
	// if integer part has 1 length, return "#.#si" (ex: "2.6Gi")
	if len(fmtint) == 1 {
		return fmtint + "." + fmtdec[0:1] + fmtsuf
	}
	// else, return "##0si" (ex: "16Mi", "256Bi")
	return fmtint + fmtsuf
}

func initbpu() {
	bpu = new(BinaryPrefixUnit)

	base := big.NewInt(1024)

	bpu.KiB = new(big.Int).Exp(base, big.NewInt(1), nil)
	bpu.MiB = new(big.Int).Exp(base, big.NewInt(2), nil)
	bpu.GiB = new(big.Int).Exp(base, big.NewInt(3), nil)
	bpu.TiB = new(big.Int).Exp(base, big.NewInt(4), nil)
	bpu.PiB = new(big.Int).Exp(base, big.NewInt(5), nil)
	bpu.EiB = new(big.Int).Exp(base, big.NewInt(6), nil)
	bpu.ZiB = new(big.Int).Exp(base, big.NewInt(7), nil)
	bpu.YiB = new(big.Int).Exp(base, big.NewInt(8), nil)
}

var bpu *(BinaryPrefixUnit)

// BinaryPrefixUnit ...
type BinaryPrefixUnit struct {
	// BiB big.Int 1024^0=1
	KiB *(big.Int) // 1024^1
	MiB *(big.Int) // 1024^2
	GiB *(big.Int) // 1024^3
	TiB *(big.Int) // 1024^4
	PiB *(big.Int) // 1024^5
	EiB *(big.Int) // 1024^6
	ZiB *(big.Int) // 1024^7
	YiB *(big.Int) // 1024^8
}
