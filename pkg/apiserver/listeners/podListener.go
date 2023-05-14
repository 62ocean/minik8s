package listeners

import (
	"encoding/json"
	"fmt"
	storagepb2 "github.com/coreos/etcd/storage/storagepb"
	"k8s/pkg/global"
	"k8s/pkg/util/msgQueue/publisher"
	log "log"
)

/*-----------------Pod Etcd Listener---------------*/

type PodListener struct {
	publisher *publisher.Publisher
}

func NewPodListener() *PodListener {
	newPublisher, _ := publisher.NewPublisher(global.MQHost)
	listener := PodListener{
		publisher: newPublisher,
	}
	return &listener
}

/*-----------------Pod Etcd Handler-----------------*/

// OnSet apiserver设置了对该资源的监听时回调
func (p PodListener) OnSet(kv storagepb2.KeyValue) {
	log.Printf("ETCD: set watcher of key " + string(kv.Key) + "\n")
	return
}

// OnModify etcd中对应资源被修改时回调
func (p PodListener) OnModify(kv storagepb2.KeyValue) {
	log.Printf("ETCD: modify kye:" + string(kv.Key) + " value:" + string(kv.Value) + "\n")
	jsonMsg, err := json.Marshal(kv)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(jsonMsg))
	//fmt.Println("lalalala")
	//fmt.Println(p)
	//fmt.Println(p.publisher)
	err = p.publisher.Publish("pods", jsonMsg, "PUT")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//fmt.Println("lalalala")
	return
}

// OnDelete etcd中对应资源被删除时回调
func (p PodListener) OnDelete(kv storagepb2.KeyValue) {
	log.Printf("ETCD: delete kye:" + string(kv.Key) + "\n")
	jsonMsg, err := json.Marshal(kv)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(jsonMsg))
	err = p.publisher.Publish("pods", jsonMsg, "DEL")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}
