package alerter

import "rule_engine/pkg/models"

// Alerter는 룰 위반 알림을 보내는 모든 모듈(print, slack 등)의 인터페이스입니다.
type Alerter interface {
	Alert(violation models.Violation)
}
