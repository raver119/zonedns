package plugin

import (
	"context"
	"github.com/miekg/dns"
)

type ZonedNS struct {
}

func (z ZonedNS) ServeDNS(ctx context.Context, writer dns.ResponseWriter, msg *dns.Msg) (int, error) {
	panic("implement me")
}

func (z ZonedNS) Name() string {
	return "zonedns"
}
