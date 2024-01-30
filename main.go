package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type domainData struct {
	MX          bool
	SPF         bool
	DMARC       bool
	spfRecord   string
	dmarcRecord string
}

func (d domainData) String() string {
	return fmt.Sprintf(
		"hasMX:%t\nhasSPF:%t\nhasDMARC:%t\nspfRecords:%v\ndmarcRecord:%v",
		d.MX,
		d.SPF,
		d.DMARC,
		d.spfRecord,
		d.dmarcRecord,
	)

}

func checkDomain(domain string) *domainData {
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	// check for MX
	mxRecords, err := net.LookupMX(domain)

	if err != nil {
		log.Printf("Error doing mxLookUp: %v\n", err)
	}
	if len(mxRecords) > 0 {
		hasMX = true
	}

	// check for spf
	txtRecords, err := net.LookupTXT(domain)

	if err != nil {
		log.Printf("Error doing txtLookup: %v\n", err)
	}
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	// check for DMARC

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)

	if err != nil {
		log.Printf("Error doing txtLookup: %v\n", err)
	}
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}
	return &domainData{
		DMARC:       hasDMARC,
		SPF:         hasSPF,
		MX:          hasMX,
		spfRecord:   spfRecord,
		dmarcRecord: dmarcRecord,
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter domain name: ")

	for scanner.Scan() {
		fmt.Println(checkDomain(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
