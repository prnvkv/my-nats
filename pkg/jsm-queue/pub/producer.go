package pub

import (
	"context"
	"encoding/json"
	"errors"

	fmt "github.com/sirupsen/logrus"

	"github.com/nats-io/nats.go"
)

type Message struct {
	Data               []byte `json:"data"`
	QueueMessageObject interface{}
}

func (m *Message) Ack(ctx context.Context) error {
	groupError := "QUEUE_MESSAGE_ACK"
	if m.QueueMessageObject == nil {
		// do nothing
		err := errors.New("no queue message object found to acknowledge")
		return err
	}

	msgObj := m.QueueMessageObject.(*nats.Msg)
	err := msgObj.Ack(nats.Context(ctx))
	if err != nil {
		fmt.Error("Error occurred %q with group %s in ACK", err.Error(), groupError)
		return err
	}

	return nil
}

type Queue interface {
	Publish(msg *Message) error
	Subscribe(ch chan *Message) error
	Unsubscribe() error
}

type queue struct {
	streamName          string
	subject             string
	jetStreamClient     nats.JetStreamContext
	stream              *nats.StreamInfo
	subscription        *nats.Subscription
	subscriptionChannel chan *nats.Msg
	closeSubscription   chan bool
}

func (q *queue) initalize() error {
	groupErr := "QUEUE_INITIALIZE"

	// check if stream already exists
	stream, err := q.jetStreamClient.StreamInfo(q.streamName)
	if err != nil {
		fmt.Errorf("Error occurred %q with %s", err.Error(), groupErr)
	}
	err = nil
	if stream != nil {
		fmt.Info("Found existing Stream")
		q.stream = stream
	} else {
		stream, err := q.jetStreamClient.AddStream(&nats.StreamConfig{
			Name:     q.streamName,
			Subjects: []string{q.subject},
		})
		if err != nil {
			fmt.Errorf("Error creating queue %q with group %s", err.Error(), groupErr)
			return err
		}
		fmt.Infof("Stream created %v...", q.stream.Config.Name)
		q.stream = stream
	}

	subjectPresent := false
	for _, subject := range q.stream.Config.Subjects {
		if subject == q.subject {
			subjectPresent = true
			fmt.Info("Subject found in the current stream")
			break
		}
	}
	if !subjectPresent {
		fmt.Info("Subject not present in the stream")
		q.stream.Config.Subjects = append(q.stream.Config.Subjects, q.subject)
		q.stream, err = q.jetStreamClient.UpdateStream(&q.stream.Config)
		if err != nil {
			fmt.Errorf("Error occurred %q with group %s", err.Error(), groupErr)
			return err
		}
	}
	return nil
}

func (q *queue) Publish(msg *Message) error {
	groupError := "QUEUE_PUBLISH"
	data := msg.Data
	byteData, err := json.Marshal(data)
	if err != nil {
		fmt.WithError(err).Error(groupError)
		return err
	}

	_, err = q.jetStreamClient.Publish(q.subject, byteData)
	if err != nil {
		fmt.WithError(err).Error(groupError)
		return err
	}
	return nil
}

func (q *queue) Subscribe(ch chan *Message) error {
	groupError := "QUEUE_SUBSCRIBE"

	if ch == nil {
		err := errors.New("subscription channel cannot be nil")
		return err
	}

	channelCapacity := cap(ch)
	if channelCapacity == 0 {
		q.subscriptionChannel = make(chan *nats.Msg)
	} else {
		q.subscriptionChannel = make(chan *nats.Msg, channelCapacity)
	}

	var err error
	q.subscription, err = q.jetStreamClient.Subscribe(
		q.subject,
		func(msg *nats.Msg) {

			// forward to ch
			data := &Message{
				Data:               msg.Data,
				QueueMessageObject: msg,
			}
			ch <- data
		},
		nats.ManualAck(),
		nats.Durable("NEW"))
	if err != nil {
		fmt.Errorf("Failed to subscribe to stream for subject %s", q.subject)
		fmt.WithError(err).Error(groupError)
		return err
	}

	return nil
}

func (q *queue) Unsubscribe() error {
	groupError := "QUEUE_UNSUBSCRIBE"

	if q.subscription == nil {
		// nothing to unsubscribe
		return nil
	}
	// fire the event for closing subscription
	q.closeSubscription <- true
	close(q.closeSubscription)
	close(q.subscriptionChannel)

	err := q.subscription.Unsubscribe()
	if err != nil {
		fmt.WithError(err).Error(groupError)
		return err
	}

	return nil
}

// NewQueue initializes the queue client
func NewQueue(
	streamName string,
	subject string,
	jetStreamClient nats.JetStreamContext,
) (Queue, error) {
	s := &queue{
		streamName:      streamName,
		subject:         subject,
		jetStreamClient: jetStreamClient,
	}

	fmt.Infoln("Queue recieved values", streamName, subject, jetStreamClient)

	// initialize
	err := s.initalize()
	return s, err
}
