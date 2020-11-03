// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file is a modified version of net/hosts.go from the golang repo

package blacklist

import (
	"bufio"
//	"bytes"
	"io"
//	"net"
	"os"
//	"strings"
	"sync"
	"time"
//        "sort"

//	"github.com/coredns/coredns/plugin"
//	"github.com/stretchr/testify/assert"
)

type options struct {
	// the default ipv4 address for a provided host
	ipv4 string

	// the default ipv6 address for a provided host
        ipv6 string

	// The TTL of the record we generate
	ttl uint32

	// The time between two reload of the configuration
	reload time.Duration

	// Explicit match domains
	explicit bool

	// Path to blacklist file
	blacklist string

	// Path to whitelist file
	// whitelist string

	// Path to cache file
	// cache string
}

func newOptions() *options {
	return &options{
		ipv4: 		"1.0.1.0",
                ipv6: 		"0:0:0:0:0:ffff:100:100",
		blacklist:	"/etc/coredns/blacklist",
                explicit:	true,
        	// whitelist: 	"/etc/coredns/whitelist",
		// cache:		"/etc/coredns/cache",
		ttl:         	3600,
		reload:      	time.Duration(86400 * time.Second),
	}
}

// Hostsfile contains known host entries.
type ListFile struct {
	sync.RWMutex

        // lists
        list map[string]bool

	// mtime and size are only read and modified by a single goroutine
	mtime time.Time
	size  int64

	options *options
}

// readHosts determines if the cached data needs to be updated based on the size and modification time of the hostsfile.
func (lf *ListFile) readLists() {
	lf.list, lf.mtime, lf.size = lf.readList(lf.options.blacklist)
	// sort.Strings(lf.blacklist)
}

func (lf ListFile) readList(path string) (map[string]bool, time.Time, int64) {
	file, err := os.Open(path)
	log.Debug("Parsing %s file", path)
        if err != nil {
                // We already log a warning if the file doesn't exist or can't be opened on setup. No need to return the error here.
                return nil, time.Now(), 0
        }
        defer file.Close()

        stat, err := file.Stat()
        // lf.RLock()
        // size := lf.size
        // lf.RUnlock()

        // if err == nil && lf.mtime.Equal(stat.ModTime()) && size == stat.Size() {
        //        return
        // }
	entries := lf.parseList(file)
        log.Debugf("Parsed blacklist file into %d entries", len(entries))

	return entries, stat.ModTime(), stat.Size()
}

// @to-do: strip the any white space and skip empty strings
func (lf ListFile) parseList(r io.Reader) map[string]bool {
	entries := make(map[string]bool)
	reader := bufio.NewReader(r)
        for {
        	line, _, err := reader.ReadLine()
                l := string(line)
		entries[l] = true
		//lst = append(lst, l)
                // log.Debug("LIST LINE: ")
                // log.Debug(l)
                if err == io.EOF {
                	break
                 }
         }
	return entries
}

func (lf ListFile) reverseTokens(tokens []string) []string {
        rtokens := make([]string, len(tokens))
        for i, j := 0, len(tokens)-1; i <= j; i, j = i+1, j-1 {
                rtokens[i], rtokens[j] = tokens[j], tokens[i]
        }
        return rtokens
}

func (lf ListFile) lookupDomain(fqdn string) bool {
	// list := []string{"org.gautier", "foo.string", "com.example", "org.gautier.www"}
	// sort.Strings(list)

	// tokens := strings.Split(fqdn, ".")
        // rtokens := lf.reverseTokens(tokens)
	// rfqdn :=strings.Join(rtokens, ".")
        // qtokens := make([]string, 0)
	// results := make([]string, 0)
        // for t := range rtokens {
        //         qtokens = append(qtokens, rtokens[t])
        //        if 1 < len(qtokens) { // make sure we hace tld.domain
        //                query := strings.Join(qtokens, ".")
	//		result := sort.SearchStrings(list, query)
	//		if result < len(list) {
	//			results = append(results, list[result])
	//		}
        //	}
	// }
	// if lf.options.explicit {
	//	eresult := sort.SearchStrings(results, rfqdn)
	//	if eresult < len(results) {
	//		return true
	//	}
	// } else {
	//	if 0 < len(results) {
	//		return true
	//	}
	// }
	// return false
	// result := sort.SearchStrings(lf.blacklist, fqdn)
        // log.Debug("Search result: ", result)
        // log.Debug("Blacklist count: ", len(lf.blacklist))
	// log.Debug("Found ", lf.blacklist[result])
        // log.Debug("Query: ", fqdn)
	/*
        if result < len(lf.blacklist) {
        	return true
        } else {
        	return false
	}
	*/
	return lf.list[fqdn]
}

