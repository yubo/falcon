/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

/* addr [4]byte use big endian */
package ping

import (
	"net"

	"golang.org/x/net/ipv4"

	"github.com/golang/glog"
)

func pkt(ip [4]byte) []byte {
	h := ipv4.Header{
		Version:  4,
		Len:      20,
		TotalLen: 20 + 10, // 20 bytes for IP, 10 for ICMP
		TTL:      64,
		Protocol: 1, // ICMP
		Dst:      net.IPv4(ip[0], ip[1], ip[2], ip[3]),
		// ID, Src and Checksum will be set for us by the kernel
	}

	icmp := []byte{
		8, // type: echo request
		0, // code: not used by echo request
		0, // checksum (16 bit), we fill in below
		0,
		0, // identifier (16 bit). zero allowed.
		0,
		0, // sequence number (16 bit). zero allowed.
		0,
		0xC0, // Optional data. ping puts time packet sent here
		0xDE,
	}
	cs := csum(icmp)
	icmp[2] = byte(cs)
	icmp[3] = byte(cs >> 8)

	out, err := h.Marshal()
	if err != nil {
		glog.Error(err)
	}
	return append(out, icmp...)
}

func csum(b []byte) uint16 {
	var s uint32
	for i := 0; i < len(b); i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	// add back the carry
	s = s>>16 + s&0xffff
	s = s + s>>16
	return uint16(^s)
}
