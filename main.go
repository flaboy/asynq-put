package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hibiken/asynq"
)

func main() {
	// 定义命令行参数
	queueName := flag.String("q", "default", "Queue name")
	redisPort := flag.Int("p", 6379, "Redis port")
	flag.Parse()

	// 检查必需的参数
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Usage: asynq-put <topic> <payload> [-q queue_name] [-p redis_port]")
		os.Exit(1)
	}

	topic := args[0]
	payload := args[1]

	// 创建 Redis 连接选项
	redisOpt := asynq.RedisClientOpt{
		Addr: fmt.Sprintf("localhost:%d", *redisPort),
	}

	// 创建 Asynq 客户端
	client := asynq.NewClient(redisOpt)
	defer client.Close()

	// 创建任务
	task, err := createTask(topic, payload)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}

	// 入队任务
	info, err := client.Enqueue(task, asynq.Queue(*queueName))
	if err != nil {
		log.Fatalf("Failed to enqueue task: %v", err)
	}

	fmt.Printf("Task enqueued successfully. ID: %s, Queue: %s\n", info.ID, info.Queue)
}

func createTask(topic, payload string) (*asynq.Task, error) {
	// 尝试将 payload 解析为 JSON
	var payloadMap map[string]interface{}
	err := json.Unmarshal([]byte(payload), &payloadMap)
	if err != nil {
		// 如果不是有效的 JSON，则将 payload 作为字符串处理
		return asynq.NewTask(topic, []byte(payload)), nil
	}

	// 如果是有效的 JSON，将其转换回 JSON 字节
	payloadBytes, err := json.Marshal(payloadMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	return asynq.NewTask(topic, payloadBytes), nil
}
