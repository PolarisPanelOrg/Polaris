package xray

import (
	"testing"

	"github.com/cuijianzhuang/polaris/internal/model"
)

func TestRenderRulesIncludesSniffedProtocol(t *testing.T) {
	rules, _ := renderRules([]model.RouteRule{{Match: model.RouteMatch{Protocol: []string{"bittorrent"}}, Outbound: "block"}}, nil)
	if len(rules) != 2 {
		t.Fatalf("rules = %#v", rules)
	}
	rule := rules[1].(map[string]any)
	protocol, ok := rule["protocol"].([]string)
	if !ok || len(protocol) != 1 || protocol[0] != "bittorrent" || rule["outboundTag"] != "block" {
		t.Fatalf("protocol route = %#v", rule)
	}
}
