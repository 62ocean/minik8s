package listeners

import (
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/msgQueue/publisher"
	"log"
)

/*-----------------Function Etcd Listener---------------*/

type FunctionListener struct {
	publisher *publisher.Publisher
}

func NewfunctionListener() *FunctionListener {
	newPublisher, _ := publisher.NewPublisher(global.MQHost)
	listener := FunctionListener{
		publisher: newPublisher,
	}
	return &listener
}

/*-----------------Function Etcd Handler-----------------*/

// OnSet apiserver设置了对该资源的监听时回调
func (p FunctionListener) OnSet(kv mvccpb.KeyValue) {
	log.Printf("ETCD: set watcher of key " + string(kv.Key) + "\n")
	return
}

// OnCreate etcd中对应资资源被创建时回调
func (p FunctionListener) OnCreate(kv mvccpb.KeyValue) {
	log.Printf("ETCD: create key:" + string(kv.Key) + " value:" + string(kv.Value) + "\n")
	jsonMsg := publisher.ConstructPublishMsg(kv, kv, object.CREATE)
	err := p.publisher.Publish("functions", jsonMsg, "CREATE")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}

// OnModify etcd中对应资源被修改时回调
func (p FunctionListener) OnModify(kv mvccpb.KeyValue, prevkv mvccpb.KeyValue) {
	log.Printf("ETCD: modified new kye:" + string(kv.Key) + " value:" + string(kv.Value) + "\n")
	jsonMsg := publisher.ConstructPublishMsg(kv, prevkv, object.UPDATE)
	err := p.publisher.Publish("functions", jsonMsg, "PUT")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}

// OnDelete etcd中对应资源被删除时回调
func (p FunctionListener) OnDelete(kv mvccpb.KeyValue, prevkv mvccpb.KeyValue) {
	log.Printf("ETCD: delete kye:" + string(prevkv.Key) + "\n")
	jsonMsg := publisher.ConstructPublishMsg(kv, prevkv, object.DELETE)
	err := p.publisher.Publish("functions", jsonMsg, "DEL")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}
