package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain, hasMx, hasSPF, sprRecord, hasDMARC, dmarcRecord")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error: could not read from input: %v\n",err)
	}

}

func checkDomain(domain string) {
	var hasMx, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	mxRecord, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error: could not lookup MX record for %s: %v\n", domain, err)
	}
	if len(mxRecord) > 0 {
		hasMx = true
	}

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		hasSPF = false
	} else {
		for _, txt := range txtRecords {
			if strings.HasPrefix(txt, "v=spf1") {
				hasSPF = true
				spfRecord = txt
				break
			}
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		hasDMARC = false
	} else {
		for _, txt := range dmarcRecords {
			if strings.HasPrefix(txt, "v=DMARC1") {
				hasDMARC = true
				dmarcRecord = txt
				break
			}
		}
	}
	fmt.Printf("%s, %t, %t, %s, %t, %s\n", domain, hasMx, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}