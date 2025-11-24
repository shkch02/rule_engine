package input

import (
	"context"
	"encoding/json"
	"log"
	"rule_engine/pkg/models"
	"time"

	"github.com/segmentio/kafka-go"
)

// 카프카 컨슘을 위한 구조체
type KafkaSource struct {
	reader *kafka.Reader
}

func NewKafkaSource(brokers []string, topic, groupID string) *KafkaSource {
	dialer := &kafka.Dialer{
		Timeout:   kafka.DefaultDialer.Timeout,
		DualStack: true,
		KeepAlive: 30 * time.Second,
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		Dialer:         dialer,      //타임아웃 방지1
		CommitInterval: time.Second, //1초마다 커밋하도록 변경(기존엔 매 매세지마다 커밋함)
	})
	return &KafkaSource{reader: r}
}

func (ks *KafkaSource) Stream(ctx context.Context) (<-chan models.Event, error) {
	eventCh := make(chan models.Event, 1000) //버퍼 천개 할당

	go func() { //고루틴으로 비동기 처리
		defer close(eventCh)
		defer ks.reader.Close()

		for {
			//1. 컨텍스트 취소 확인
			if ctx.Err() != nil {
				return
			}

			//2. 카프카에서 메시지 읽기
			m, err := ks.reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() == nil {
					log.Printf("카프카 메시지 읽기 오류 : %v", err)
				}
				return
			}

			//3. 메시지 파싱
			var event models.Event
			if err := json.Unmarshal(m.Value, &event); err != nil {
				log.Printf("경고 : 카프카 메시지 파싱 실패 : %v", err)
			} else {
				//4.파싱 성공시 채널로 이벤트 전송
				select {
				case eventCh <- event:
				case <-ctx.Done():
					return
				}
			}

			//5. 메시지 커밋 주요 병목 지점 -> 메 메시지마다 커밋 안하고 1초마다 자동 커밋
			//if err := ks.reader.CommitMessages(ctx, m); err != nil {
			//	log.Printf("카프카 메시지 커밋 실패 : %v", err)
			//}
		}
	}()
	return eventCh, nil
}
