package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"encoding/json"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("ENter the domain\n")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: could not read from input: %v\n", err)
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
	result := map[string]interface{}{
		"domain":      domain,
		"hasMx":       hasMx,
		"hasSPF":      hasSPF,
		"spfRecord":   spfRecord,
		"hasDMARC":    hasDMARC,
		"dmarcRecord": dmarcRecord,
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error: could not marshal result to JSON for domain %s: %v\n", domain, err)
		return
	}

	fmt.Println(string(jsonResult))
}