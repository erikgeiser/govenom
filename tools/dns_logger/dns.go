package main

import (
	"net"

	"github.com/miekg/dns"
)

type logHandler interface {
	Handle(string) error
}

type dnsHandler struct {
	logHandler
}

func newDNSHandler(handler logHandler) *dnsHandler {
	return &dnsHandler{handler}
}

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}

	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		domain := msg.Question[0].Name
		err := h.Handle(domain)
		if err != nil {
			logError(err)
		}

		msg.Authoritative = true
		msg.Answer = append(msg.Answer, &dns.A{
			Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 1},
			A:   net.ParseIP("127.0.0.1"),
		})
	}
	w.WriteMsg(&msg)
}
