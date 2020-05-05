package mafkasync

import (
	"fmt"
	"testing"
)

type Message struct {
	Ts int64
	Msg string
}

func (m Message) GetKey() int64 {
	return m.Ts
}

func Printf(data Keyer) {
	message := data.(Message)
	fmt.Printf("%v, %v\n", message.Ts, message.Msg)
}

func TestList( t *testing.T) {
	l := NewMapList()
	for i:=int64(0); i< 10; i++ {
		l.Push(Message{i, "hello"})
	}
	//l.Walk(Printf)
	fmt.Printf("size = %d\n", l.Size())
	l.RemoveBefore(8)
	l.Walk(Printf)
	fmt.Printf("size = %d\n", l.Size())
}