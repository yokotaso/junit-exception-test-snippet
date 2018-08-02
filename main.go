package main

import (
	"strings"
	"io/ioutil"
	"log"
	"bufio"
	"os"
	"flag"
)

const sampleCode = `
import java.lang.RuntimeException;
import org.junit.Test;

public class SampleTest {
    @Test( expected = RuntimeException.class )
    public void testException() {
    }
}
`

func main() {
	var example = flag.Bool("example", false, "display replacement of sample")
	var inputfile = flag.String("input", "", "path of input soy file")
	var outputfile = flag.String("output", "", "path of output soy file")
	flag.Parse()

	if ! *example {
		if *inputfile == "" {
			log.Fatal("flag -input required.")
		}

		if *outputfile == "" {
			log.Fatal("flag -output required.")
		}
	}

	reader := getReader(*example, *inputfile)
	writer := getWriter(*example, *outputfile)
	defer writer.Flush()

	if *example {
		writer.Write([]byte("SampleCode:\n" + sampleCode))
		writer.Write([]byte("=======================\n"))
		writer.Write([]byte("ModifiedCode:\n"))
	}

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text() + "\n"

		ok := hasTestAnnotation(text)
		if !ok {
			writer.Write([]byte(text))
			continue
		}
		if ok, snippet := ParseTestAnnotation([]byte(text)); ok {
			snippet.Write(writer)
		} else {
			writer.Write([]byte(text))
			continue
		}
	}
	bufio.NewWriter(writer).Flush()
}

func hasTestAnnotation(text string) bool {
	return strings.Contains(text, "@Test")
}

func getReader(isExample bool, inputfile string) *strings.Reader {
	if isExample {
		return strings.NewReader(sampleCode)
	} else {
		bytes, err := ioutil.ReadFile(inputfile)
		if err != nil {
			log.Fatal(err)
		}
		return strings.NewReader(string(bytes))
	}
}

func getWriter(isExample bool, outputfile string) *bufio.Writer {
	if isExample {
		return bufio.NewWriter(os.Stdout)
	} else {
		f, err := os.Create(outputfile)
		if err != nil {
			log.Fatal(f)
		}
		return bufio.NewWriter(f)
	}
}
