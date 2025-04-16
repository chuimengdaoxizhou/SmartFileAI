package models

import (
	"encoding/json"
	"strings"
)

type UserHistoryMessage struct {
	UserID         string           `json:"userid"`
	CreateTime     int64            `json:"createtime"`
	HistoryMessage []HistoryMessage `json:"historymessage"`
}

type ManageHistoryMessage struct {
	// 使用UserID+CreateTime作为key
	Data map[string]*UserHistoryMessage
}

// ContentBlock 表示消息中的一个内容块（可以是文本、工具使用、工具结果等）
type ContentBlock struct {
	Type    string          `json:"type"`              // 内容块的类型（例如 "text", "tool_use", "tool_result"）
	Text    string          `json:"text,omitempty"`    // 如果类型是 "text"，这是文本内容
	ID      string          `json:"id,omitempty"`      // 内容块的唯一 ID
	Name    string          `json:"name,omitempty"`    // 工具的名称（如果是工具调用）
	Input   json.RawMessage `json:"input,omitempty"`   // 工具调用的输入参数（如果是工具调用）
	Content interface{}     `json:"content,omitempty"` // 其他类型的内容（例如工具结果的返回值）
}

// HistoryMessage 实现了 llm.Message 接口，用于存储消息
type HistoryMessage struct {
	Role    string         `json:"role"`    // 消息的角色（如 user, assistant）
	Content []ContentBlock `json:"content"` // 消息内容，包含多个内容块
}

// GetRole 返回消息的角色（例如 "user" 或 "assistant"）
func (m *HistoryMessage) GetRole() string {
	return m.Role
}

// GetContent 返回消息的所有文本内容（将多个内容块中的文本合并成一个字符串）
func (m *HistoryMessage) GetContent() string {
	var content string
	// 遍历所有内容块，提取文本类型的内容
	for _, block := range m.Content {
		if block.Type == "text" {
			content += block.Text + " " // 合并文本内容
		}
	}
	return strings.TrimSpace(content) // 去除多余的空格
}

// IsToolResponse 判断消息是否是工具的响应（类型为 "tool_result"）
func (m *HistoryMessage) IsToolResponse() bool {
	// 遍历所有内容块，查找是否有工具结果类型的内容块
	for _, block := range m.Content {
		if block.Type == "tool_result" {
			return true
		}
	}
	return false
}

// GetUsage 返回消息的使用情况（这里不跟踪使用情况，因此返回 (0, 0)）
func (m *HistoryMessage) GetUsage() (int, int) {
	return 0, 0 // 当前历史记录不跟踪使用情况
}

// HistoryToolCall 实现了 llm.ToolCall 接口，用于存储工具调用的记录
type HistoryToolCall struct {
	id   string          // 工具调用的 ID
	name string          // 工具的名称
	args json.RawMessage // 工具调用的参数
}

// GetID 返回工具调用的 ID
func (t *HistoryToolCall) GetID() string {
	return t.id
}

// GetName 返回工具调用的名称
func (t *HistoryToolCall) GetName() string {
	return t.name
}

// GetArguments 返回工具调用的参数（将原始 JSON 解析为 map）
func (t *HistoryToolCall) GetArguments() map[string]interface{} {
	var args map[string]interface{}
	// 解析工具调用的参数
	if err := json.Unmarshal(t.args, &args); err != nil {
		return make(map[string]interface{}) // 如果解析失败，返回一个空的 map
	}
	return args
}
