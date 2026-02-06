// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.
package ip

import (
	"bufio"
	"chatwiki/internal/app/chatwiki/define"
	"encoding/binary"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/zhimaAi/go_tools/logs"
)

type IPRange struct {
	Start uint32 // start IP
	End   uint32 // end IP
}

var instance *ChinaIPChecker

type ChinaIPChecker struct {
	ranges []IPRange
}

func NewChinaIPChecker(filename string) (*ChinaIPChecker, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ranges []IPRange
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// parse format: apnic|CN|ipv4|223.202.0.0|131072|20100713|allocated
		parts := strings.Split(line, "|")
		if len(parts) < 5 || parts[1] != "CN" || parts[2] != "ipv4" {
			continue
		}

		startIP := parts[3]
		count, err := strconv.ParseUint(parts[4], 10, 32)
		if err != nil {
			continue
		}

		// format to number
		start := ipToUint32(net.ParseIP(startIP).To4())
		if start == 0 {
			continue
		}

		// format end IP: start + count - 1
		end := start + uint32(count) - 1

		ranges = append(ranges, IPRange{
			Start: start,
			End:   end,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// sort by start IP
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].Start < ranges[j].Start
	})

	return &ChinaIPChecker{ranges: ranges}, nil
}

func (c *ChinaIPChecker) IsChinaIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	ip = ip.To4()
	if ip == nil {
		//IPv6 not supported
		return false
	}
	ipVal := ipToUint32(ip)
	//search
	idx := sort.Search(len(c.ranges), func(i int) bool {
		return c.ranges[i].Start >= ipVal
	})
	// check found position
	if idx < len(c.ranges) && c.ranges[idx].Start <= ipVal && ipVal <= c.ranges[idx].End {
		return true
	}
	// check previous range
	if idx > 0 {
		prev := c.ranges[idx-1]
		if prev.Start <= ipVal && ipVal <= prev.End {
			return true
		}
	}
	return false
}

func ipToUint32(ip net.IP) uint32 {
	if len(ip) != 4 {
		return 0
	}
	return binary.BigEndian.Uint32(ip)
}

func uint32ToIP(n uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, n)
	return ip.String()
}

func init() {
	var err error
	instance, err = NewChinaIPChecker(define.AppRoot + "/ip/apnic.txt")
	if err != nil {
		logs.Error("NewChinaIPChecker failed: %v", err)
	}
}

func IsChinaIP(ipStr string) bool {
	if instance == nil {
		return false
	}
	return instance.IsChinaIP(ipStr)
}
