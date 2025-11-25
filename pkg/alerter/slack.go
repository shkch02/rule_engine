package alerter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rule_engine/pkg/models"
	"time"
)

// SlackAlerter는 향후 구현될 Slack 알림 모듈입니다.
type SlackAlerter struct {
	WebhookURL string
}

func NewSlackAlerter(webhookURL string) *SlackAlerter {
	return &SlackAlerter{WebhookURL: webhookURL}
}

func (sa *SlackAlerter) sendSlackNotification(message string) error {
	payload := map[string]string{"text": message}
	jsonPayload, err := json.Marshal(payload)

	req, _ := http.NewRequest("POST", sa.WebhookURL, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("슬랙 알림 실패, 상태 코드: %d", resp.StatusCode)
	}
	return nil
}

// Alert는 SlackAlerter가 Alerter 인터페이스를 만족시키기 위한 스텁입니다.
func (sa *SlackAlerter) Alert(v models.Violation) {
	eventBytes, _ := json.MarshalIndent(v.Event, "", "  ")
	msg := fmt.Sprintf(`
[!!!] 보안 룰 위반 감지
룰 ID:     %s
설명:      %s
원본 이벤트: %s
`, v.RuleID, v.Description, string(eventBytes))

	go func(msg string) {
		err := sa.sendSlackNotification(msg)
		if err != nil {
			log.Printf("슬랙 알림 전송 실패: %v", err)
		}
	}(msg)

}
