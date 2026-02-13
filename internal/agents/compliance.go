package agents

import (
	"a2a/models"
	"a2a/server"
	"fmt"
	"strings"
	"time"
)

// ComplianceAgent (Agent C)
func NewComplianceAgent() *server.Server {
	card := models.AgentCard{
		Name:        "ComplianceOfficer",
		Description: models.StringPtr("稽核專員，負責審查最終報表是否合規"),
		Version:     "1.0.0",
		URL:         "http://localhost:8080/agent/compliance",
		Capabilities: models.AgentCapabilities{
			Streaming: models.BoolPtr(false),
		},
		Skills: []models.AgentSkill{
			{ID: "audit-report", Name: "報表稽核", Description: models.StringPtr("審查報支金額")},
		},
	}

	handler := func(task *models.Task, msg *models.Message, update server.TaskUpdateFunc) (*models.Task, error) {
		text := ""
		if len(msg.Parts) > 0 && msg.Parts[0].Text != nil {
			text = *msg.Parts[0].Text
		}

		fmt.Printf("[Agent C (Compliance)] 收到指令: %s\n", text)
		
		// 模擬稽核邏輯
		time.Sleep(1 * time.Second) // 模擬審查時間
		
		var responseText string
		if strings.Contains(text, "$15,500") || strings.Contains(text, "15,500") {
			responseText = "✅ [核准] 總金額 $15,500 符合部門預算 ($20,000)。核准代碼: COMP-2026-OK"
		} else if strings.Contains(text, "$30,000") || strings.Contains(text, "30,000") {
			responseText = "❌ [退回] 總金額 $30,000 超出預算上限 ($20,000)。請重新檢視行程。"
		} else {
			responseText = "⚠️ [無法判定] 報表中未發現明確金額，請補充資訊。"
		}

		task.Status.State = models.TaskStateCompleted
		if task.Metadata == nil {
			task.Metadata = make(map[string]interface{})
		}
		task.Metadata["reply"] = responseText

		return task, nil
	}

	return server.NewA2AServer(card, handler)
}
