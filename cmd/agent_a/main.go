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

	rounds := []string{
		"老闆下週一要去台北出差三天，預算一天 5,000 元，請推薦飯店。",
		"訂君悅。另外請幫忙訂週一早上 9 點從台中出發的高鐵。",
		"沒問題，直接訂票。請確認總費用。",
		"參加 Google A2A 技術研討會。",
	}

	// 模擬前四回合的標準 A2A 對話
	for i, cmd := range rounds {
		fmt.Printf("\n--- 第 %d 回合 ---\n", i+1)
		sendA2AMessage(cmd)
		time.Sleep(2 * time.Second) // 留一點時間讓老闆看 Log
	}

	// 第五回合：展示 SSE 串流功能
	fmt.Printf("\n--- 第 5 回合 (SSE 串流展示) ---\n")
	fmt.Println("PA: 請產出最終行程表與報帳單。")
	streamA2AMessage("產出最終行程表與報帳單。")
}

func sendA2AMessage(text string) {
	fmt.Printf("PA 發送指令: %s\n", text)

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
	resp, err := http.Post("http://localhost:8080/a2a", "application/json", bytes.NewBuffer(body))
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
			fmt.Printf("TF 回應: %v\n", meta["reply"])
		}
	}
}

func streamA2AMessage(text string) {
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
		Method: "message/stream", // 使用串流方法
		Params: params,
	}

	body, _ := json.Marshal(rpcReq)
	resp, err := http.Post("http://localhost:8080/a2a", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("錯誤: %v\n", err)
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()

	fmt.Println(">>> 正在接收即時進度更新 (SSE)...")
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		
		if strings.TrimSpace(line) == "" {
			continue
		}

		var streamResp models.SendTaskStreamingResponse
		if err := json.Unmarshal([]byte(line), &streamResp); err == nil {
			if update, ok := streamResp.Result.(map[string]interface{}); ok {
				// 處理 1: 狀態更新 (只在狀態改變時印出，並換行)
				if _, ok := update["status"].(map[string]interface{}); ok {
					if update["final"] != true {
						// 這裡不頻繁印出狀態，以免打斷打字機
						_ = 0 // No-op to satisfy staticcheck
					}
				}
				
				// 處理 2: 文字碎片 (打字機效果)
				if artifact, ok := update["artifact"].(map[string]interface{}); ok {
					if parts, ok := artifact["parts"].([]interface{}); ok && len(parts) > 0 {
						if part, ok := parts[0].(map[string]interface{}); ok {
							if txt, ok := part["text"].(string); ok {
								fmt.Print(txt) 
							}
						}
					}
				}

				if update["final"] == true {
					fmt.Println("\n\n✅ 任務完整結束，報告已由財務專員產出！")
					break
				}
			}
		}
	}
}
