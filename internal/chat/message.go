package chat

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
	RoleSystem    Role = "system"
)

type Message struct {
	Role      Role
	Content   string
	ToolName  string             // for tool results
	ToolCalls []ToolCallRecord   // for assistant messages with tool calls
}

type ToolCallRecord struct {
	Name      string
	Arguments string // raw JSON
}
