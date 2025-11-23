package input

import (
	"context"
	"fmt"
	"rule_engine/pkg/models"

	"github.com/segmentio/kafka-go"
)

//카프카 컨슘을 위한 구조체
type KafkaSource struct {
	reader *kafka.Reader
}

func NewKafkaSource(brokers []string, topic, groupID string) *KafkaSource {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	return &KafkaSource{reader: r}
}

func (ks *KafkaSource) Stream(ctx context.Context) (<-chan models.Event, error) {
	eventCh := make(chan models.Event)
	
	go func() { //고루틴으로 비동기 처리
		defer close(eventCh)
		defer ks.reader.Close()

		for {
				//1. 컨텍스트 취소 확인
				if ctx.err() != nil{return}

				//2. 카프카에서 메시지 읽기
				m, err != ks.reader.FetchMessage(ctx)
				if err != nil {
					if ctx.Err() == nil {
						log.printf("카프카 메시지 읽기 오류 : %v",err)
					}
					return
				}

				//3. 메시지 파싱
				var event models.Event
				if err := json.Unmarshal(m.Value,&event); err != nil {
					log.Printf("경고 : 카프카 메시지 파싱 실패 : %v",err)
				} else {
						//4.파싱 성공시 채널로 이벤트 전송
						select {
						case eventCh <- event:
						case <-ctx.Done():
							return
					}
				}

				//5. 메시지 커밋
				if err := ks.reader.CommitMessages(ctx,m); err != nil {
					log.Printf("카프카 메시지 커밋 실패 : %v",err)
				}
		} 
	}()
	return eventCh,nil
}