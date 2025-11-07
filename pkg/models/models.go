package models

// Event는 동적 필드를 가진 단일 로그 이벤트를 나타냅니다.
// eBPF 로그의 JSON 구조가 syscall마다 다르므로 map을 사용합니다.
type Event map[string]interface{}

// Violation은 룰 위반 시 Alerter로 전달되는 구조체입니다.
type Violation struct {
	RuleID      string
	Description string
	Event       Event // 위반을 일으킨 원본 이벤트
}
