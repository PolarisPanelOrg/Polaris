package singbox

import (
	"testing"

	"github.com/cuijianzhuang/polaris/internal/model"
)

func TestRenderRulesIncludesSniffedProtocol(t *testing.T) {
	rules, _, _ := renderRules([]model.RouteRule{{Match: model.RouteMatch{Protocol: []string{"bittorrent"}}, Outbound: "block"}})
	if len(rules) != 1 {
		t.Fatalf("rules = %#v", rules)
	}
	rule := rules[0].(map[string]any)
	protocol, ok := rule["protocol"].([]string)
	if !ok || len(protocol) != 1 || protocol[0] != "bittorrent" || rule["action"] != "reject" {
		t.Fatalf("protocol route = %#v", rule)
	}
}
