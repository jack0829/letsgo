package kafka

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"os"
	"testing"
	"time"
)

var (
	hosts  []string
	topic  string
	writer *kafka.Writer
	reader *kafka.Reader
)

func TestMain(m *testing.M) {

	ctx := context.Background()

	hosts = []string{
		"127.0.0.1:9092",
	}

	topic = "DemoTopic"

	groupID := "demo-group"

	mechanism := plain.Mechanism{
		Username: "user",
		Password: "client_PAs5W0rd",
	}

	dialer := &kafka.Dialer{
		ClientID:      "Jack-MBP-Test",
		Timeout:       10 * time.Second,
		SASLMechanism: mechanism,
	}

	tr := &kafka.Transport{
		SASL: mechanism,
	}

	conn, err := dialer.DialContext(ctx, "tcp", hosts[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer conn.Close()

	writer = &kafka.Writer{
		Addr:      kafka.TCP(hosts...),
		Topic:     topic,
		Balancer:  &kafka.LeastBytes{},
		Transport: tr,
	}
	defer writer.Close()

	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     hosts,
		Topic:       topic,
		Dialer:      dialer,
		GroupID:     groupID,
		GroupTopics: []string{topic},
	})
	defer reader.Close()

	m.Run()
}

func TestKafka(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	go testWrite(ctx, time.Second*2)

	testRead(ctx)
}

func testWrite(ctx context.Context, delay time.Duration) {

	fmt.Fprintln(os.Stdout, "写入开始")
	defer fmt.Fprintln(os.Stdout, "写入结束")

Loop:
	for {

		select {
		case <-ctx.Done():
			break Loop

		case t := <-time.After(delay):

			val, _ := jsoniter.Marshal(t)

			m := kafka.Message{
				Value: val,
			}

			if err := writer.WriteMessages(ctx, m); err != nil {
				fmt.Fprintln(os.Stderr, "写入错误", err)
				continue Loop
			}

			fmt.Fprintln(os.Stdout, "写入", t)
		}

	}
}

func testRead(ctx context.Context) {

	fmt.Fprintln(os.Stdout, "读取开始")
	defer fmt.Fprintln(os.Stdout, "读取结束")

	var t time.Time
Loop:
	for {

		select {
		case <-ctx.Done():
			break Loop

		default:

			m, err := reader.ReadMessage(ctx)
			if err != nil {
				fmt.Fprintln(os.Stderr, "读取错误", err)
				continue Loop
			}

			if err = jsoniter.Unmarshal(m.Value, &t); err != nil {
				fmt.Fprintln(os.Stderr, "解析错误", err)
				continue Loop
			}

			fmt.Fprintln(os.Stdout, "读取", t)
		}

	}
}
