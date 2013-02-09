/*
	Provides necessary data to Munin to enable monitoring usage of Huge pages,
	if configured.

	This would probably be better implemented with a shell script than a
	very primitive Go app...
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type HugePageInfo struct {
	hptotal, hprsvd, hpfree, hpsurp, hpsize uint64
}

func main() {

	var hpi *HugePageInfo = new(HugePageInfo)

	if len(os.Args) >= 2 && os.Args[1] == "config" {
		outputCfg()
	} else {
		getData(hpi)
		printValues(hpi)
	}
}

func printValues(hpi *HugePageInfo) {
	fmt.Printf("hptotal.value %d\n", (*hpi).hpsize*(*hpi).hptotal)
	fmt.Printf("hpfree.value %d\n", (*hpi).hpsize*(*hpi).hpfree)
	fmt.Printf("hprsvd.value %d\n", (*hpi).hpsize*(*hpi).hprsvd)
	fmt.Printf("hpsurp.value %d\n", (*hpi).hpsize*(*hpi).hpsurp)
}

func outputCfg() {
	fmt.Print(
		`graph_title Hugepage usage
graph_args --base 1024 -l 0
graph_vlabel Pages
graph_category system
graph_info Utilisation of Huge pages
hptotal.draw AREA
hptotal.label Total huge pages
hptotal.type GAUGE
hpfree.draw AREA
hpfree.label Free pages
hpfree.type GAUGE
hprsvd.draw LINE2
hprsvd.label Reserved pages
hprsvd.type GAUGE
hpsurp.draw LINE2
hpsurp.label Surplus pages
hpsurp.type GAUGE
`)

}

func getData(hpi *HugePageInfo) {

	fd, err := os.Open("/proc/meminfo")
	if err != nil {
		fmt.Println("error opening file: ", err)
		os.Exit(1)
	}

	r := bufio.NewReader(fd)

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
			fmt.Println("Finished")
		}
		reg := regexp.MustCompile(`^(.+):\s*(\d*)`)
		if str := reg.FindStringSubmatch(line); str != nil {
			var i uint64

			i, err := strconv.ParseUint(str[2], 10, 64)
			if err != nil {
				panic(fmt.Sprintf("Failed to parse %s", str[1]))
			}

			switch str[1] {
			case `HugePages_Total`:
				(*hpi).hptotal = i
			case `HugePages_Rsvd`:
				(*hpi).hprsvd = i
			case `HugePages_Free`:
				(*hpi).hpfree = i
			case `HugePages_Surp`:
				(*hpi).hpsurp = i
			case `Hugepagesize`:
				(*hpi).hpsize = 1024 * i
			}
		}
	}
}
