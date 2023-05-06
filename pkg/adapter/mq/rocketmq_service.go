package mq

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/golang/glog"
	"os"
	"strings"
	"sync"
	"time"
)

var doOnce sync.Once

type RocketmqService struct {
	Producer rocketmq.Producer
	Consumer rocketmq.PushConsumer
}

type IMqService interface {
	ProduceMsg(topic string, data interface{}) error
	ProduceMsgWithTag(topic, tag string, data interface{}) error
	ConsumeMsg(topic string)
	ConsumeMsgWithTag(topic, tag string)
	Shutdown()
}

func NewRocketmqService(endpoints []string, group string) *RocketmqService {
	var rockermqService *RocketmqService
	doOnce.Do(func() {
		rockermqService = &RocketmqService{
			Producer: NewProducer(endpoints, group),
			Consumer: NewPushConsumer(endpoints, group),
		}
	})
	return rockermqService
}

func NewProducer(endpoints []string, group string) rocketmq.Producer {
	p, _ := rocketmq.NewProducer(
		producer.WithNameServer(endpoints),
		producer.WithGroupName(group),
		producer.WithRetry(5),
		producer.WithSendMsgTimeout(3*time.Second),
		producer.WithQueueSelector(producer.NewHashQueueSelector()),
	)
	err := p.Start()
	if err != nil {
		glog.Errorf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	return p
}
func NewPushConsumer(endpoints []string, group string) rocketmq.PushConsumer {
	//c, _ := rocketmq.NewPushConsumer(
	//	consumer.WithGroupName(group),
	//	consumer.WithNameServer(endpoints),
	//)
	c, err:= rocketmq.NewPushConsumer(
		consumer.WithGroupName(group),
		consumer.WithNameServer(endpoints),
	)
	if err != nil {
		hlog.Error("err is"+err.Error())
	}
	return c
}

func (r *RocketmqService) Shutdown() {
	if r.Producer != nil {
		err := r.Producer.Shutdown()
		if err != nil {
			glog.Errorf("shutdown producer error: %s", err.Error())
		}
	}
}

func (r *RocketmqService) ProduceMsg(topic string, data interface{}) error {
	msg := primitive.NewMessage(topic, Marshal(data))
	res, err := r.Producer.SendSync(context.Background(), msg)
	if err != nil {
		glog.Errorf("Send message error: %s\n", err)
		return err
	}

	glog.V(5).Infof("Send message success: result=%s\n", res.String())
	return nil
}

func (r *RocketmqService) ProduceMsgWithShardingKey(topic string, data interface{}, shardingKey string) error {
	msg := primitive.NewMessage(topic, Marshal(data))
	msg.WithShardingKey(shardingKey)
	res, err := r.Producer.SendSync(context.Background(), msg)
	if err != nil {
		glog.Errorf("Send message error: %s\n", err)
		return err
	}

	glog.V(5).Infof("Send message success: result=%s\n", res.String())
	return nil
}
func (r *RocketmqService) ProduceMsgAndKeyWithShardingKey(topic, tag string, data interface{}, shardingKey string) error {
	msg := primitive.NewMessage(topic, Marshal(data))
	msg.WithShardingKey(shardingKey)
	msg.WithKeys([]string{tag})
	res, err := r.Producer.SendSync(context.Background(), msg)
	if err != nil {
		glog.Errorf("Send message error: %s\n", err)
		return err
	}

	glog.V(5).Infof("Send message success: result=%s\n", res.String())
	return nil
}

func (r *RocketmqService) ProduceMsgWithTag(topic, tag string, data interface{}) error {
	msg := primitive.NewMessage(topic, Marshal(data))
	msg.WithTag(tag)
	res, err := r.Producer.SendSync(context.Background(), msg)
	if err != nil {
		glog.Errorf("Send message error: %s\n", err)
		return err
	}

	glog.V(5).Infof("Send message success: result=%s\n", res.String())
	return nil
}

func (r *RocketmqService) ConsumeMsgWithTag(topic string, tags []string) error {
	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: strings.Join(tags, " || "),
	}
	err := r.Consumer.Subscribe(topic, selector, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		glog.V(5).Infof("subscribe callback: %v \n", msgs)
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		glog.Errorf(err.Error())
		return err
	}
	err = r.Consumer.Start()
	if err != nil {
		glog.Errorf("consumer start failed %v", err.Error())
		return err
	}
	return nil
}
func (r *RocketmqService) ConsumeMsg(topic string) error {
	selector := consumer.MessageSelector{}
	err := r.Consumer.Subscribe(topic, selector, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		glog.V(5).Infof("subscribe callback: %v \n", msgs)
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		glog.Errorf(err.Error())
	}
	err = r.Consumer.Start()
	if err != nil {
		glog.Errorf("start producer error: %s", err.Error())
		return err
	}
	return nil
}

func (r *RocketmqService) ConsumeMsgDefine(topic string, f func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)) error {
	selector := consumer.MessageSelector{}
	err := r.Consumer.Subscribe(topic, selector, f)
	if err != nil {
		glog.Errorf(err.Error())
	}
	err = r.Consumer.Start()
	if err != nil {
		glog.Errorf("start producer error: %s", err.Error())
		return err
	}
	return nil
}

func Marshal(data interface{}) []byte {
	marshal, err := sonic.Marshal(data)
	if err != nil {
		glog.Errorf("json marshal failed")
		return []byte{}
	}
	return marshal
}
