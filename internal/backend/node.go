package backend

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

func NodePingChecker(wg *sync.WaitGroup, handler *Handler) {
	defer wg.Done()
	nodestore := handler.nodeHandler.nodeService.nodeStore
	timeticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-timeticker.C:
			nodes, err := nodestore.ListNodes(context.Background())
			if err != nil {
				log.Println("Error listing nodes:", err)
				continue
			}
			for _, node := range nodes {
				// Ping node
				success := false
				for _, ip := range node.IPs {
					if success {
						break
					}
					url := "http://" + ip + ":18090" + "/ping"
					// Send ping request
					ts := time.Now()
					resp, err := http.Get(url)
					if err != nil {
						log.Println("Error pinging node:", err)
						success = false
						continue
					}
					data := map[string]string{}
					if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
						log.Println("Error decoding response:", err)
						success = false
						continue
					}
					resp.Body.Close()
					if data["timestamp"] == "" {
						log.Println("Invalid response from node:", resp.Status)
						success = false
						continue
					}
					success = true
					// Update last seen time
					node.LastSeen = time.Now()
					node.Status = "online"
					node.Latency = time.Since(ts)
					err = nodestore.UpdateNode(context.Background(), node)
					if err != nil {
						log.Println("Error updating node:", err)
						continue
					}
				}
				if !success {
					node.Status = "offline"
					nodestore.UpdateNode(context.Background(), node)
				}
			}
		default:
			time.Sleep(5 * time.Second)
			continue
		}
	}
}
