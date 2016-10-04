package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func stringsToGrdIteration(line []string) grdIteration {
	numline := make([]float64, len(line))
	hasZeroGrd := false
	for i, val := range line {
		if n, err := strconv.ParseFloat(val, 64); err == nil {
			numline[i] = n
			if n == 0 {
				hasZeroGrd = true
			}
		}
	}
	return grdIteration{hasZeroGrd, numline}
}

type grdIteration struct {
	HasZeroGradient bool
	Gradients       []float64
}

// GrdFile contains information about the gradient file
type GrdFile struct {
	Header    []string
	Gradients []grdIteration
}

func main() {
	start := time.Now()
	// usr, _ := user.Current()
	// dir := usr.HomeDir
	// fpath := filepath.Join(dir, "tiny01.grd")
	grdFile := new(GrdFile)
	var gradients []grdIteration
	fpath := filepath.Join("fixtures", "tiny01.grd")
	fmt.Println(fpath)
	file, _ := os.Open(fpath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, "TABLE") {
			continue
		}
		if strings.Contains(text, "ITERATION") {
			line := strings.Fields(text)
			grdFile.Header = line
		} else {
			line := strings.Fields(text)
			grd := stringsToGrdIteration(line)
			gradients = append(gradients, grd)
		}
	}
	grdFile.Gradients = gradients
	elapsed := time.Since(start)
	fmt.Println(grdFile)
	log.Printf("read and parsing took %s", elapsed)
}
