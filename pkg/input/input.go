package input

import (
	"context"
	"rule-engine-poc/pkg/models"
)

// Source는 모든 이벤트 입력 소스(Kafka, File 등)가 구현할 인터페이스입니다.
// PoC (main.go)에서는 아직 사용되지 않습니다.
type Source interface {
	// Stream은 이벤트를 읽어 채널로 전송합니다.
	Stream(ctx context.Context) (<-chan models.Event, error)
}
