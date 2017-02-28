package client

import (
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

var (
	Debug      = false
	protoDebug = &jsonpb.Marshaler{Indent: "\t"}
)

func debugProto(action string, proto proto.Message) {
	if Debug && proto != nil {
		debug, err := protoDebug.MarshalToString(proto)
		if err != nil {
			return
		}
		fmt.Printf("%s - %s\n", action, debug)
	}
}
