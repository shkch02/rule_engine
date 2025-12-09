package alerter

import (
	"encoding/json"
	"log"
	"rule_engine/pkg/models"
)

// PrintAlerter는 Alerter 인터페이스의 PoC 구현체입니다. (콘솔 출력)
type PrintAlerter struct{}

func NewPrintAlerter() *PrintAlerter {
	return &PrintAlerter{}
}

// Alert는 룰 위반 정보를 콘솔에 로깅합니다.
func (pa *PrintAlerter) Alert(v models.Violation) {
	// 이벤트를 JSON으로 다시 직렬화하여 보기 좋게 출력
	eventBytes, _ := json.Marshal(v.Event)

	log.Printf(`
--------------------------------------------------
[!!!] 보안 룰 위반 감지 !!!
룰 ID:     %s
설명:      %s
원본 이벤트: %s
--------------------------------------------------
`, v.RuleID, v.Description, string(eventBytes))
}
