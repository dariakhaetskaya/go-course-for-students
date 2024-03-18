package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Options struct {
	From        string
	To          string
	Offset      int
	Limit       int64
	BlockSize   int64
	Conversions []string
}

func ParseFlags() (*Options, error) {
	var opts Options
	var conversions string

	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.IntVar(&opts.Offset, "offset", 0, "offset in bytes to start copying from")
	flag.Int64Var(&opts.Limit, "limit", 0, "max number of bytes to copy")
	flag.Int64Var(&opts.BlockSize, "block-size", 1024, "block size in bytes for reading and writing")
	flag.StringVar(&conversions, "conv", "", "conversions to apply: upper_case,lower_case,trim_spaces")
	flag.Parse()

	if strings.Contains(conversions, "upper_case") && strings.Contains(conversions, "lower_case") {
		return nil, fmt.Errorf("both upper_case and lower_case specified")
	}

	opts.Conversions = strings.Split(conversions, ",")

	if opts.Offset < 0 {
		return nil, errors.New("offset must be positive integer")
	}

	return &opts, nil
}

func transformText(input string, conversions []string) (string, error) {
	for _, conv := range conversions {
		switch conv {
		case "upper_case":
			input = strings.ToUpper(input)
		case "lower_case":
			input = strings.ToLower(input)
		case "trim_spaces":
			input = strings.TrimSpace(input)
		}
	}
	return input, nil
}

func Copy(opts *Options) error {
	var input io.Reader = os.Stdin
	if opts.From != "" {
		file, err := os.Open(opts.From)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error opening source file:", err)
			return err
		}
		defer file.Close()
		input = file
	}

	var output io.Writer = os.Stdout
	if opts.To != "" {
		file, err := os.OpenFile(opts.To, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating destination file:", err)
			return err
		}
		defer file.Close()
		output = file
	}

	bufferedReader := bufio.NewReader(input)
	bufferedWriter := bufio.NewWriter(output)
	defer bufferedWriter.Flush()

	discarded, err := bufferedReader.Discard(opts.Offset)
	if discarded < opts.Offset && err != nil {
		fmt.Fprintln(os.Stderr, "Offset size is bigger than file size:", err)
		return err
	}

	var readBytes int64

	for {
		bytesToRead := opts.BlockSize
		if opts.Limit > 0 {
			if readBytes >= opts.Limit {
				break
			}
			if opts.Limit-readBytes < opts.BlockSize {
				bytesToRead = opts.Limit - readBytes
			}
		}
		_, err = io.CopyN(bufferedWriter, bufferedReader, bytesToRead)

		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "can not read:", err)
			os.Exit(1)
		}
		readBytes += bytesToRead
	}
	return nil
}

func main() {
	opts, err := ParseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}

	err = Copy(opts)
	if err != nil {
		os.Exit(1)
	}

}
