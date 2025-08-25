package bootstrap

import (
	"auth-service/constants"
	loggerI "auth-service/domain/service/logger"
	queueS "auth-service/domain/service/queue"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

type queueClient struct {
	client  *asynq.Client
	retry   int
	timeout time.Duration
	log     loggerI.Log
}

func NewQueueClient(env *Env, log loggerI.Log) *queueClient {
	opt := asynq.RedisClientOpt{
		Addr:     env.QUEUE.Addr,
		DB:       env.QUEUE.DB,
		Password: env.QUEUE.Password,
		Network:  env.QUEUE.Network,
	}
	client := asynq.NewClient(opt)

	if client.Ping() != nil {
		log.Fatal("Failed to connect to the queue server: " + client.Ping().Error())
	}

	return &queueClient{
		client:  client,
		retry:   5,
		timeout: 10 * time.Minute,
		log:     log,
	}
}

func (qc *queueClient) EnqueueMail(payload queueS.Payload) (string, error) {
	return qc.EnqueueAnyTask(constants.QUEUE_MAIL, payload)
}

func (qc *queueClient) EnqueueAnyTask(taskType constants.QueueType, payload queueS.Payload) (string, error) {
	task, err := qc.NewTask(string(taskType), payload)
	if err != nil {
		return "", err
	}
	i, err := qc.client.Enqueue(task)
	if err != nil {
		return "", err
	}
	return i.ID, nil
}

func (qc *queueClient) NewTask(typeTask string, payload queueS.Payload, opts ...asynq.Option) (*asynq.Task, error) {
	defaultOpts := []asynq.Option{
		asynq.MaxRetry(qc.retry),
		asynq.Timeout(qc.timeout),
	}
	opts = append(defaultOpts, opts...)
	pl, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(typeTask, pl, opts...), nil
}

func (qc *queueClient) Close() {
	qc.client.Close()
}

func (qc *queueClient) Ping() error {
	return qc.client.Ping()
}
