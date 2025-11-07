package alerter

import (
	"log"
	"rule-engine-poc/pkg/models"
)

// SlackAlerter는 향후 구현될 Slack 알림 모듈입니다.
type SlackAlerter struct {
	WebhookURL string
}

// Alert는 SlackAlerter가 Alerter 인터페이스를 만족시키기 위한 스텁입니다.
func (sa *SlackAlerter) Alert(v models.Violation) {
	// TODO: Slack Webhook 전송 로직 구현
	log.Printf("알림 (스텁): Slack으로 전송 -> %s", v.RuleID)
}
