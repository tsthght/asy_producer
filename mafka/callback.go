package mafka

type MafkaCallBack struct{
	SuccessChan chan []interface{}
	FailureChan chan interface{}
}

func (m MafkaCallBack) OnSuccess(msgs []interface{}) {
	m.SuccessChan <- msgs
}

func (m MafkaCallBack) OnFailure(msgs []interface{}, err error) {
	m.FailureChan <- err
}
