package marshaler

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	gstatus "google.golang.org/grpc/status"
)

type Marshaler struct {
}

func NewMarshaler() *Marshaler {
	return &Marshaler{}
}

type protoBody struct {
	Body []byte
	Err  *spb.Status
}

func (m *Marshaler) MarshalWrapper(rets ...interface{}) ([]byte, error) {
	if len(rets) != 2 {
		return nil, fmt.Errorf("MarshalWrapper err len %d", len(rets))
	}
	var bodyStr []byte
	var bodyErr *spb.Status
	if rets[0] != nil && !reflect.ValueOf(rets[0]).IsZero() {
		respbody := rets[0].(proto.Message)
		var err error
		bodyStr, err = proto.Marshal(respbody)
		if err != nil {
			return nil, err
		}
	}

	if rets[1] != nil {
		resperr := rets[1].(error)
		// 不管是错误都需要转换成 grpc err ， 才能方便存储
		s, _ := gstatus.FromError(resperr)
		bodyErr = s.Proto()
	}

	protoBody := &protoBody{
		Body: bodyStr,
		Err:  bodyErr,
	}

	return json.Marshal(protoBody)
}

func (m *Marshaler) UnMarshalWrapper(strings []byte, resp any) ([]interface{}, error) {
	protoBody := &protoBody{}
	err := json.Unmarshal(strings, protoBody)
	if err != nil {
		return nil, err
	}

	v := resp.([]interface{})
	var bodyF proto.Message
	if protoBody.Body != nil {
		var ok bool
		bodyF, ok = v[0].(proto.Message)
		if !ok {
			return nil, fmt.Errorf("UnMarshalWrapper err type %T", v[0])
		}
		if err := proto.Unmarshal(protoBody.Body, bodyF); err != nil {
			return nil, err
		}
	} else {
		rv := reflect.ValueOf(v[0]).Elem()
		if rv.CanSet() {
			rv.Set(reflect.Zero(rv.Type()))
			bodyF = rv.Addr().Interface().(proto.Message)
		} else {
			bodyF = v[0].(proto.Message)
		}
	}
	var respErr error
	if protoBody.Err != nil {
		respErr = gstatus.FromProto(protoBody.Err).Err()
	} else {
		respErr = nil
	}

	return []interface{}{bodyF, respErr}, nil
}
