package main

import (
	"a2a/models"
	"a2a/server"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	// 1. 定義 Agent B 的身分卡 (AgentCard)
	financeAgentCard := models.AgentCard{
		Name:        "FinanceTravelExpert",
		Description: stringPtr("專門處理公司差旅預算與訂票的財務助理"),
		Version:     "1.0.0",
		URL:         "http://localhost:8080/a2a",
		Capabilities: models.AgentCapabilities{
			Streaming: boolPtr(true),
		},
		Skills: []models.AgentSkill{
			{ID: "travel-booking", Name: "差旅訂票", Description: stringPtr("處理飯店與高鐵訂位")},
			{ID: "budget-check", Name: "預算審核", Description: stringPtr("確保開支符合公司政策")},
		},
	}

	// 2. 實作任務處理邏輯 (支援即時串流更新)
	handler := func(task *models.Task, msg *models.Message, update server.TaskUpdateFunc) (*models.Task, error) {
		text := ""
		if len(msg.Parts) > 0 && msg.Parts[0].Text != nil {
			text = *msg.Parts[0].Text
		}

		fmt.Printf("[Agent B] 收到指令: %s\n", text)

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
			// 模擬打字機效果的串流輸出 (修正亂碼：按 Rune 迭代)
			report := "【最終行程報告】\n- 飯店：君悅飯店 (3晚)\n- 交通：高鐵台中-台北來回\n- 事由：A2A技術研討會\n- 總預算：$15,500\n✅ 報帳單已產出並歸檔。"
			
			for _, charRune := range report {
				char := string(charRune)
				// 推送 Artifact 片段
				update(models.TaskArtifactUpdateEvent{
					ID: task.ID,
					Artifact: models.Artifact{
						Parts: []models.Part{
							{Text: &char},
						},
					},
					Final: boolPtr(false),
				})
				time.Sleep(40 * time.Millisecond) 
			}
			
			responseText = "報告已完成。"
			responseState = models.TaskStateCompleted
		default:
			responseText = "收到您的訊息，正在處理中..."
		}

		fmt.Printf("[Agent B] 處理完畢\n")
		
		task.Status.State = responseState
		if task.Metadata == nil {
			task.Metadata = make(map[string]interface{})
		}
		task.Metadata["reply"] = responseText

		return task, nil
	}

	// 3. 建立並啟動 A2A Server
	srv := server.NewA2AServer(financeAgentCard, handler)
	
	// 設定 Server 參數 (這部分的 API 在範例 server.go 中是私有的，我們直接在 Start 前設定好)
	// 註：在 a2a-samples 的實作中，我們需要確保 server 結構體有對外暴露設定埠口的方法
	// 這裡假設我們直接使用預設行為
	
	fmt.Println("🚀 Agent B (財務專員) 已啟動，等待任務中... (Port 8080)")
	
	// 為了 Demo 方便，我們手動啟動 HTTP
	http.Handle("/a2a", srv)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func stringPtr(s string) *string { return &s }
func boolPtr(b bool) *bool     { return &b }
