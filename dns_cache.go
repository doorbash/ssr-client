package main

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

var (
	ErrNoRecords = errors.New("no records")
)

type DNSCache struct {
	records map[string]DNSRecord
	mutex   sync.Mutex
}

type DNSRecord struct {
	addresses []net.IP
	time      time.Time
}

func (dns *DNSCache) Put(domain string, addresses []net.IP) {
	if len(addresses) > 0 {
		dns.mutex.Lock()
		defer dns.mutex.Unlock()
		_, contains := dns.records[domain]
		if contains {
			return
		}
		dns.records[domain] = DNSRecord{
			addresses: addresses,
			time:      time.Now(),
		}
	}
}

func (dns *DNSCache) Get(domain string) ([]net.IP, error) {
	dns.mutex.Lock()
	defer dns.mutex.Unlock()
	_, contains := dns.records[domain]
	if contains {
		return dns.records[domain].addresses, nil
	}
	return nil, ErrNoRecords
}

func (dns *DNSCache) CleanUp() {
	log.Printf("DNSCache: Starting DNS Cache CleanUp")
	records := 0
	for range time.Tick(time.Minute) {
		dns.mutex.Lock()
		for k, v := range dns.records {
			if time.Since(v.time).Minutes() > 10 {
				delete(dns.records, k)
				records++
			}
		}
		dns.mutex.Unlock()
		if records > 0 {
			log.Printf("DNSCache: %d dns records has been removed", records)
		}
		records = 0
	}
}

func NewDNSCache() *DNSCache {
	cache := &DNSCache{
		records: make(map[string]DNSRecord),
	}
	go cache.CleanUp()
	return cache
}
