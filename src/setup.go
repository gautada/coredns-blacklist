package blacklist

import (
//	"os"
//	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
)

var log = clog.NewWithPlugin("blacklist")

func init() { plugin.Register("blacklist", setup) }

func periodicBlacklistUpdate(b *Blacklist) chan bool {
	parseChan := make(chan bool)

	if b.options.reload == 0 {
		return parseChan
	}

	go func() {
		ticker := time.NewTicker(b.options.reload)
		for {
			select {
			case <-parseChan:
				return
			case <-ticker.C:
				b.readLists()
			}
		}
	}()
	return parseChan
}

func setup(c *caddy.Controller) error {
	h, err := blacklistParse(c)
	if err != nil {
		return plugin.Error("blacklist", err)
	}

	parseChan := periodicBlacklistUpdate(&h)

	c.OnStartup(func() error {
		h.readLists()
		return nil
	})

	c.OnShutdown(func() error {
		close(parseChan)
		return nil
	})

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		h.Next = next
		return h
	})

	return nil
}

func blacklistParse(c *caddy.Controller) (Blacklist, error) {
	// config := dnsserver.GetConfig(c)

	h := Blacklist{
		ListFile: &ListFile{
			// blacklist:    "/etc/coredns/blacklist",
                        // path_w:    "/etc/coredns/whitelist",
			// path_c:    "/etc/coredns/cache",
			options: newOptions(),
		},
	}

	inline := []string{}
	i := 0
	for c.Next() {
		if i > 0 {
			return h, plugin.ErrOnce
		}
		i++

//		args := c.RemainingArgs()

/*		if len(args) >= 1 {
			// h.path = args[0]
			args = args[1:]
			if !filepath.IsAbs(h..options.blacklist) && config.Root != "" {
				h.path = filepath.Join(config.Root, h.path)
			}
			s, err := os.Stat(h.path)
			if err != nil {
				if os.IsNotExist(err) {
					log.Warningf("File does not exist: %s", h.path)
				} else {
					return h, c.Errf("unable to access hosts file '%s': %v", h.path, err)
				}
			}
			if s != nil && s.IsDir() {
				log.Warningf("Hosts file %q is a directory", h.path)
			}
		}

		origins := make([]string, len(c.ServerBlockKeys))
		copy(origins, c.ServerBlockKeys)
		if len(args) > 0 {
			origins = args
		}

		for i := range origins {
			origins[i] = plugin.Host(origins[i]).Normalize()
		}
		h.Origins = origins
*/
		for c.NextBlock() {
			switch c.Val() {
			case "fallthrough":
				h.Fall.SetZonesFromArgs(c.RemainingArgs())
			case "ipv4":
				remaining := c.RemainingArgs()
				h.options.ipv4 = remaining[0]
			case "ipv6":
				remaining := c.RemainingArgs()
                                h.options.ipv6 = remaining[0]
			case "blacklist":
                                remaining := c.RemainingArgs()
                                h.options.blacklist = remaining[0]
			case "ttl":
				remaining := c.RemainingArgs()
				if len(remaining) < 1 {
					return h, c.Errf("ttl needs a time in second")
				}
				ttl, err := strconv.Atoi(remaining[0])
				if err != nil {
					return h, c.Errf("ttl needs a number of second")
				}
				if ttl <= 0 || ttl > 65535 {
					return h, c.Errf("ttl provided is invalid")
				}
				h.options.ttl = uint32(ttl)
			case "reload":
				remaining := c.RemainingArgs()
				if len(remaining) != 1 {
					return h, c.Errf("reload needs a duration (zero seconds to disable)")
				}
				reload, err := time.ParseDuration(remaining[0])
				if err != nil {
					return h, c.Errf("invalid duration for reload '%s'", remaining[0])
				}
				if reload < 0 {
					return h, c.Errf("invalid negative duration for reload '%s'", remaining[0])
				}
				h.options.reload = reload
			default:
				if len(h.Fall.Zones) == 0 {
					line := strings.Join(append([]string{c.Val()}, c.RemainingArgs()...), " ")
					inline = append(inline, line)
					continue
				}
				return h, c.Errf("unknown property '%s'", c.Val())
			}
		}
	}

	return h, nil
}
