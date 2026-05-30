package commonmodel

type KafkaGroupNotice struct {
	TraceID        string         `json:"trace_id"`
	GoalUserID     []int64        `json:"goal_user_id"`
	SessionID      int64          `json:"group_id,string"`
	SendTimeSecond int64          `json:"send_time_second"`
	Data           map[string]any `json:"data"`
	MessageType    MessageType    `json:"message_type"`
	MessageCode    MessageCode    `json:"message_code"`
}

type KafkaNewMessageNotice struct {
	TraceID        string      `json:"trace_id"`
	GoalUserID     []int64     `json:"goal_user_id"`
	SessionID      int64       `json:"group_id,string"`
	SendTimeSecond int64       `json:"send_time_second"`
	MessageType    MessageType `json:"message_type"`
	MessageCode    MessageCode `json:"message_code"`
}
type KafkaSystemMessage struct {
	TraceID        string         `json:"trace_id"`
	GoalUserID     []int64        `json:"goal_user_id"`
	Data           map[string]any `json:"data"`
	SendTimeSecond int64          `json:"send_time_second"`
	MessageType    MessageType    `json:"message_type"`
	MessageCode    MessageCode    `json:"message_code"`
}
