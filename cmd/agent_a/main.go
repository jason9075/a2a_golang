package main

import (
	"a2a/models"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func main() {
	fmt.Println("🏢 [公司差旅展示] Agent A (助理) 正在啟動...")
	time.Sleep(1 * time.Second)

	// Step 1: 與 Agent B (財務) 互動
	fmt.Println("\n=== Step 1: 與 Agent B (財務) 協調行程 ===")
	rounds := []string{
		"老闆下週一要去台北出差三天，預算一天 5,000 元，請推薦飯店。",
		"訂君悅。另外請幫忙訂週一早上 9 點從台中出發的高鐵。",
		"沒問題，直接訂票。請確認總費用。",
		"參加 Google A2A 技術研討會。",
	}

	for i, cmd := range rounds {
		fmt.Printf("\n--- 第 %d 回合 ---\n", i+1)
		sendA2AMessage("http://localhost:8080/agent/finance", cmd)
		time.Sleep(1 * time.Second)
	}

	// Step 2: 取得 Agent B 的最終報告 (SSE)
	fmt.Printf("\n--- 第 5 回合 (SSE 串流展示) ---\n")
	fmt.Println("PA: 請產出最終行程表與報帳單。")
	finalReport := streamA2AMessage("http://localhost:8080/agent/finance", "產出最終行程表與報帳單。")
	
	// Step 3: 送交 Agent C (稽核) 審核
	fmt.Println("\n=== Step 2: 送交 Agent C (稽核) 審核 ===")
	time.Sleep(1 * time.Second)
	
	fmt.Printf("PA 發送報告給稽核: %s\n", finalReport)
	
	// 這裡我們直接把 Agent B 的輸出丟給 Agent C
	// 在實際應用中，可能需要稍微整理格式，但 Agent C 的邏輯是 regex 金額，所以沒問題
	sendA2AMessage("http://localhost:8080/agent/compliance", "請審核以下報表: " + finalReport)
}

func sendA2AMessage(endpoint, text string) {
	fmt.Printf("PA -> %s: %s\n", endpoint, text)

	reqID := fmt.Sprintf("req-%d", time.Now().Unix())
	params := models.TaskSendParams{
		ID: "travel-task-123",
		Message: models.Message{
			Role: "user",
			Parts: []models.Part{
				{Text: &text},
			},
		},
	}

	rpcReq := models.JSONRPCRequest{
		JSONRPCMessage: models.JSONRPCMessage{
			JSONRPC: "2.0",
			JSONRPCMessageIdentifier: models.JSONRPCMessageIdentifier{ID: reqID},
		},
		Method: "message/send",
		Params: params,
	}

	body, _ := json.Marshal(rpcReq)
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("錯誤: %v\n", err)
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()

	var rpcResp models.JSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		fmt.Printf("Decode error: %v\n", err)
		return
	}

	// 從 Metadata 中抓取我們剛才塞的回應
	if res, ok := rpcResp.Result.(map[string]interface{}); ok {
		if meta, ok := res["metadata"].(map[string]interface{}); ok {
			fmt.Printf("RESPONSE: %v\n", meta["reply"])
		}
	}
}

// 修改後的回傳值：返回最終累積的字串，供下一步驟使用
func streamA2AMessage(endpoint, text string) string {
	reqID := "req-stream-999"
	params := models.TaskSendParams{
		ID: "travel-task-123",
		Message: models.Message{
			Role: "user",
			Parts: []models.Part{
				{Text: &text},
			},
		},
	}

	rpcReq := models.JSONRPCRequest{
		JSONRPCMessage: models.JSONRPCMessage{
			JSONRPC: "2.0",
			JSONRPCMessageIdentifier: models.JSONRPCMessageIdentifier{ID: reqID},
		},
		Method: "message/stream",
		Params: params,
	}

	body, _ := json.Marshal(rpcReq)
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("錯誤: %v\n", err)
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()

	fmt.Println(">>> 正在接收即時進度更新 (SSE)...")
	
	fullText := ""
	
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return ""
		}
		
		if strings.TrimSpace(line) == "" {
			continue
		}

		var streamResp models.SendTaskStreamingResponse
		if err := json.Unmarshal([]byte(line), &streamResp); err == nil {
			if update, ok := streamResp.Result.(map[string]interface{}); ok {
				// 處理 1: 狀態更新
				if _, ok := update["status"].(map[string]interface{}); ok {
					if update["final"] != true {
						_ = 0
					}
				}
				
				// 處理 2: 文字碎片
				if artifact, ok := update["artifact"].(map[string]interface{}); ok {
					if parts, ok := artifact["parts"].([]interface{}); ok && len(parts) > 0 {
						if part, ok := parts[0].(map[string]interface{}); ok {
							if txt, ok := part["text"].(string); ok {
								fmt.Print(txt)
								fullText += txt
							}
						}
					}
				}

				if update["final"] == true {
					fmt.Println("\n\n✅ 任務完整結束！")
					break
				}
			}
		}
	}
	return fullText
}
