package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/akamensky/argparse"
)

func getSliceFractions[Type any](slice []Type, fractionLength int) [][]Type {
	if len(slice) <= fractionLength {
		return [][]Type{
			slice,
		}
	}

	var out [][]Type

	n := 0
	for n < len(slice) {

		if (len(slice)-1)-n >= fractionLength-1 {
			out = append(out, slice[n:][:fractionLength])

			n += fractionLength
		} else {
			if n == 0 {
				out = append(out, slice)
			} else {
				out = append(out, slice[n:])
			}

			break
		}
	}

	return out
}

func formatBytes(bytes []byte) string {
	fractionedBytes := getSliceFractions(bytes, 8)

	var out []string

	for _, fb := range fractionedBytes {
		var o []string
		for _, b := range fb {
			o = append(o, fmt.Sprint(b))
		}
		out = append(out, strings.Join(o, " | "))
	}

	return strings.Join(out, "\n")
}

func readBytesFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var content []byte
	var bcache [][]byte
	for scanner.Scan() {
		b := scanner.Bytes()
		//fmt.Println(b)

		t := 0
		if len(bcache) > 1 {
			if len(bcache[len(bcache)-1]) == 0 || len(b) == 0 {
				t++
			}
		} else if len(b) == 0 {
			t++
		}

		for range t {
			content = append(content, 10)
		}

		content = append(content, b...)
		bcache = append(bcache, b)
	}
	if err := scanner.Err(); err != nil {
		return []byte{}, err
	}
	return content, nil
}

func main() {
	parser := argparse.NewParser("Bindump", "Dumps the bytes of the given file")

	inpath := parser.String("f", "file", &argparse.Options{Required: true, Help: "Path to the file you want to dump the bytes of"})
	*inpath = path.Clean(*inpath)

	parserErr := parser.Parse(os.Args)
	if parserErr != nil {
		fmt.Print(parser.Usage(parserErr))
		return
	}

	byteFcontent, byteFileErr := readBytesFile(*inpath)
	if byteFileErr != nil {
		fmt.Println(byteFileErr)
		return
	}

	formatedBytes := formatBytes(byteFcontent)

	fmt.Println("bytes:\n" + formatedBytes + "\nend")
}
