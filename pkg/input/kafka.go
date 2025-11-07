package input

import (
	"context"
	"fmt"
	"rule_engine/pkg/models"
)

// KafkaSource는 향후 구현될 Kafka Consumer입니다.
type KafkaSource struct {
	// ... (예: kafka.Reader, topic, brokers)
}

// Stream은 KafkaSource가 Source 인터페이스를 만족시키기 위한 스텁입니다.
func (ks *KafkaSource) Stream(ctx context.Context) (<-chan models.Event, error) {
	// TODO: Kafka Consumer 로직 구현
	return nil, fmt.Errorf("KafkaSource.Stream()은 아직 구현되지 않았습니다")
}
