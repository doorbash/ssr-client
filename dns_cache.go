package main

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

var ErrNoRecords = errors.New("no records")

type DNSRecord struct {
	ips  []net.IP
	time time.Time
}

type DNSCache struct {
	records map[string]DNSRecord
	mutex   sync.Mutex
}

func (d *DNSCache) Put(domain string, ips []net.IP) {
	if len(ips) == 0 {
		return
	}
	d.mutex.Lock()
	defer d.mutex.Unlock()
	_, ok := d.records[domain]
	if ok {
		return
	}
	d.records[domain] = DNSRecord{
		ips:  ips,
		time: time.Now(),
	}
}

func (d *DNSCache) Get(domain string) ([]net.IP, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	record, ok := d.records[domain]
	if ok {
		return record.ips, nil
	}
	return nil, ErrNoRecords
}

func (d *DNSCache) cleanUp() {
	log.Printf("DNSCache: Starting DNS Cache CleanUp")
	r := 0
	for range time.Tick(time.Minute) {
		d.mutex.Lock()
		for k, v := range d.records {
			if time.Since(v.time).Minutes() > 10 {
				delete(d.records, k)
				r++
			}
		}
		d.mutex.Unlock()
		if r > 0 {
			log.Printf("DNSCache: %d dns records has been removed", r)
		}
		r = 0
	}
}

func NewDNSCache() *DNSCache {
	cache := &DNSCache{
		records: make(map[string]DNSRecord),
	}
	go cache.cleanUp()
	return cache
}
