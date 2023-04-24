package main

import (
	"bufio"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var ipv4RegExp = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

func main() {
	logsGzPath := os.Args[1] // path to directory containing `.gz` log files.
	ipPath := os.Args[2]     // path to text file containing IPs

	targetIPs, err := extractTargetIPs(ipPath)
	if err != nil {
		log.Fatalf("failed to extract target IPs: %v", err)
	}

	if len(targetIPs) == 0 {
		log.Fatalf("no targey IP addresses found in %s", ipPath)
	} else {
		log.Printf("found %d target IP address(es) in %s\n", len(targetIPs), ipPath)
	}

	entries, err := os.ReadDir(logsGzPath)
	if err != nil {
		log.Fatal(err)
	}
	if len(entries) == 0 {
		log.Fatalf("no files found in %s", logsGzPath)
	}
	var totalFoundIPs, totalGz int
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".gz") {
			totalGz++
			var foundIPs int
			gzf, err := os.Open(filepath.Join(logsGzPath, e.Name()))
			if err != nil {
				log.Fatal(err)
			}
			defer gzf.Close()
			gz, err := gzip.NewReader(gzf)
			if err != nil {
				log.Fatal(err)
			}
			s := bufio.NewScanner(gz)
			s.Split(bufio.ScanLines)
			for s.Scan() {
				text := s.Text()
				iips := ipv4RegExp.FindAllString(text, -1)
				for _, p := range iips {
					if _, ok := targetIPs[p]; ok {
						foundIPs++
					}
				}
			}
			log.Printf("found %d IPs in %s\n", foundIPs, e.Name())
			totalFoundIPs += foundIPs
		}
	}
	if totalGz == 0 {
		log.Fatalf("no .gz files found in %s", logsGzPath)
	}
	log.Printf("found total of %d IPs across %d .gz files\n", totalFoundIPs, totalGz)
}

func extractTargetIPs(ipFile string) (map[string]struct{}, error) {
	ipf, err := os.Open(ipFile)
	if err != nil {
		return nil, err
	}
	defer ipf.Close()
	content, err := io.ReadAll(ipf)
	if err != nil {
		return nil, err
	}
	ips := ipv4RegExp.FindAllString(string(content), -1)
	ipsMap := make(map[string]struct{})
	for _, ip := range ips {
		ipsMap[ip] = struct{}{}
	}
	return ipsMap, nil
}
