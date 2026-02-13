package agents

import (
	"a2a/models"
	"a2a/server"
	"fmt"
	"strings"
	"time"
)

// FinanceAgent (Agent B)
func NewFinanceAgent() *server.Server {
	card := models.AgentCard{
		Name:        "FinanceTravelExpert",
		Description: models.StringPtr("專門處理公司差旅預算與訂票的財務助理"),
		Version:     "1.0.0",
		URL:         "http://localhost:8080/agent/finance",
		Capabilities: models.AgentCapabilities{
			Streaming: models.BoolPtr(true),
		},
		Skills: []models.AgentSkill{
			{ID: "travel-booking", Name: "差旅訂票", Description: models.StringPtr("處理飯店與高鐵訂位")},
			{ID: "budget-check", Name: "預算審核", Description: models.StringPtr("確保開支符合公司政策")},
		},
	}

	handler := func(task *models.Task, msg *models.Message, update server.TaskUpdateFunc) (*models.Task, error) {
		text := ""
		if len(msg.Parts) > 0 && msg.Parts[0].Text != nil {
			text = *msg.Parts[0].Text
		}

		fmt.Printf("[Agent B (Finance)] 收到指令: %s\n", text)

		responseState := models.TaskStateWorking
		var responseText string

		switch {
		case strings.Contains(text, "下週一"):
			responseText = "【第一回合】已為您找到兩間符合政策的飯店：1. 君悅 ($4,800) 2. 寒舍艾美 ($5,000)。請問要訂哪一間？"
		case strings.Contains(text, "君悅"):
			responseText = "【第二回合】君悅飯店已保留。關於高鐵，週一 09:10 有班次 ($700)，是否直接訂購？"
		case strings.Contains(text, "直接訂票"):
			responseText = "【第三回合】機票與飯店已確認，總計 $15,500。請問此行出差事由為何？財務部報支需要。"
		case strings.Contains(text, "研討會"):
			responseText = "【第四回合】收到。我現在開始為您準備完整的行程摘要與報帳草案，請稍候..."
			responseState = models.TaskStateCompleted
		case strings.Contains(text, "產出"):
			// 模擬打字機效果的串流輸出
			report := "【最終行程報告】\n- 飯店：君悅飯店 (3晚)\n- 交通：高鐵台中-台北來回\n- 事由：A2A技術研討會\n- 總預算：$15,500\n✅ 報帳單已產出並歸檔。"
			
			for _, charRune := range report {
				char := string(charRune)
				update(models.TaskArtifactUpdateEvent{
					ID: task.ID,
					Artifact: models.Artifact{
						Parts: []models.Part{
							{Text: &char},
						},
					},
					Final: models.BoolPtr(false),
				})
				time.Sleep(20 * time.Millisecond) // Slightly faster for demo
			}
			
			responseText = report // Return full report as final result
			responseState = models.TaskStateCompleted
		default:
			responseText = "收到您的訊息，正在處理中..."
		}

		task.Status.State = responseState
		if task.Metadata == nil {
			task.Metadata = make(map[string]interface{})
		}
		task.Metadata["reply"] = responseText

		return task, nil
	}

	return server.NewA2AServer(card, handler)
}
