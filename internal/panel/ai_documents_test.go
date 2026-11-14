package panel

import (
	"strings"
	"testing"

	"github.com/cuijianzhuang/polaris/internal/store"
)

func TestAIKnowledgeDocumentsRespectPortalAudience(t *testing.T) {
	s := newTestServer(t)
	for _, doc := range []store.AIKnowledgeDocument{
		{ID: "kb-portal", Title: "Portal connectivity", Body: "portal unique-network-help", Audience: "portal", Enabled: true, Revision: 1},
		{ID: "kb-admin", Title: "Admin secret", Body: "admin unique-network-help", Audience: "admin", Enabled: true, Revision: 1},
		{ID: "kb-disabled", Title: "Disabled", Body: "unique-network-help disabled", Audience: "portal", Enabled: false, Revision: 1},
	} {
		if err := s.st.DB.Create(&doc).Error; err != nil {
			t.Fatal(err)
		}
	}
	portal := s.aiKnowledgeDocumentsFor("unique-network-help", true)
	if len(portal) != 1 || portal[0].Title != "Portal connectivity" || strings.Contains(portal[0].Excerpt, "admin") {
		t.Fatalf("portal docs = %#v", portal)
	}
	admin := s.aiKnowledgeDocumentsFor("unique-network-help", false)
	if len(admin) != 2 {
		t.Fatalf("admin docs = %#v", admin)
	}
}
