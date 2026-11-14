package panel

import (
	"errors"
	"net/http"
	"strings"

	"github.com/cuijianzhuang/polaris/internal/store"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func validKnowledgeDocument(d *store.AIKnowledgeDocument) error {
	d.Title, d.Body, d.Audience = strings.TrimSpace(d.Title), strings.TrimSpace(d.Body), strings.TrimSpace(d.Audience)
	if d.Title == "" || d.Body == "" || len([]rune(d.Title)) > 160 || len([]rune(d.Body)) > 12000 || (d.Audience != "portal" && d.Audience != "admin") {
		return errors.New("知识库文档无效：标题最多 160 字，正文最多 12000 字，受众须为 portal 或 admin")
	}
	return nil
}
func (s *Server) listAIKnowledgeDocuments(c *gin.Context) {
	if !aiOwner(c) {
		return
	}
	var docs []store.AIKnowledgeDocument
	if err := s.st.DB.Order("updated_at DESC, id DESC").Find(&docs).Error; err != nil {
		serverError(c, err)
		return
	}
	c.JSON(http.StatusOK, docs)
}
func (s *Server) createAIKnowledgeDocument(c *gin.Context) {
	if !aiOwner(c) {
		return
	}
	var doc store.AIKnowledgeDocument
	if !aiReadJSON(c, &doc) {
		return
	}
	if err := validKnowledgeDocument(&doc); err != nil {
		badRequest(c, err)
		return
	}
	doc.ID, doc.Revision = store.NewID(), 1
	if err := s.st.DB.Create(&doc).Error; err != nil {
		serverError(c, err)
		return
	}
	s.auditCtx(c, "ai.knowledge_create", doc.ID, doc.Audience)
	c.JSON(http.StatusCreated, doc)
}
func (s *Server) updateAIKnowledgeDocument(c *gin.Context) {
	if !aiOwner(c) {
		return
	}
	var req struct {
		Title, Body, Audience string
		Enabled               bool
		Revision              int
	}
	if !aiReadJSON(c, &req) {
		return
	}
	doc := store.AIKnowledgeDocument{Title: req.Title, Body: req.Body, Audience: req.Audience, Enabled: req.Enabled}
	if err := validKnowledgeDocument(&doc); err != nil {
		badRequest(c, err)
		return
	}
	res := s.st.DB.Model(&store.AIKnowledgeDocument{}).Where("id = ? AND revision = ?", c.Param("id"), req.Revision).Updates(map[string]any{"title": doc.Title, "body": doc.Body, "audience": doc.Audience, "enabled": doc.Enabled, "revision": gorm.Expr("revision + 1")})
	if res.Error != nil {
		serverError(c, res.Error)
		return
	}
	if res.RowsAffected == 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "知识库文档已被修改或不存在"})
		return
	}
	var out store.AIKnowledgeDocument
	s.st.DB.First(&out, "id = ?", c.Param("id"))
	s.auditCtx(c, "ai.knowledge_update", out.ID, out.Audience)
	c.JSON(http.StatusOK, out)
}
func (s *Server) deleteAIKnowledgeDocument(c *gin.Context) {
	if !aiOwner(c) {
		return
	}
	res := s.st.DB.Delete(&store.AIKnowledgeDocument{}, "id = ?", c.Param("id"))
	if res.Error != nil {
		serverError(c, res.Error)
		return
	}
	if res.RowsAffected == 0 {
		c.Status(http.StatusNotFound)
		return
	}
	s.auditCtx(c, "ai.knowledge_delete", c.Param("id"), "")
	c.Status(http.StatusNoContent)
}
