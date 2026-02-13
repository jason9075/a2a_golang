package main

import (
	"a2a/internal/agents"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 1. Initialize Agents
	financeAgent := agents.NewFinanceAgent()
	complianceAgent := agents.NewComplianceAgent()

	// 2. Register Routes (Single Port, Multiple Paths)
	http.Handle("/agent/finance", financeAgent)
	http.Handle("/agent/compliance", complianceAgent)

	// 3. Start Server
	port := ":8080"
	fmt.Printf("🚀 A2A Server Cluster Started on %s\n", port)
	fmt.Println("   - Agent B (Finance):    http://localhost:8080/agent/finance")
	fmt.Println("   - Agent C (Compliance): http://localhost:8080/agent/compliance")
	
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
