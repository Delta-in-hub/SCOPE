package redis

import (
	"context"
	"encoding/json"
	"scope/internal/models"

	"github.com/redis/go-redis/v9"
)

// NodeStore 实现基于Redis的节点存储
type NodeStore struct {
	client *redis.Client
}

// NewNodeStore 创建一个新的Redis节点存储
func NewNodeStore(client *redis.Client) *NodeStore {
	return &NodeStore{
		client: client,
	}
}

func (n *NodeStore) UpdateNode(ctx context.Context, node models.NodeInfo) error {
	k := node.ID
	v, err := json.Marshal(node)
	if err != nil {
		return err
	}
	return n.client.HSet(ctx, "nodes", k, v).Err()
}

func (n *NodeStore) GetNode(ctx context.Context, id string) (models.NodeInfo, error) {
	v, err := n.client.HGet(ctx, "nodes", id).Result()
	if err != nil {
		return models.NodeInfo{}, err
	}
	var node models.NodeInfo
	err = json.Unmarshal([]byte(v), &node)
	if err != nil {
		return models.NodeInfo{}, err
	}
	return node, nil
}

func (n *NodeStore) DeleteNode(ctx context.Context, id string) error {
	return n.client.HDel(ctx, "nodes", id).Err()
}

func (n *NodeStore) ListNodes(ctx context.Context) ([]models.NodeInfo, error) {
	v, err := n.client.HGetAll(ctx, "nodes").Result()
	if err != nil {
		return nil, err
	}
	var nodes []models.NodeInfo
	for _, v := range v {
		var node models.NodeInfo
		err = json.Unmarshal([]byte(v), &node)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}
