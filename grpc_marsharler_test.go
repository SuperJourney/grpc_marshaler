package marshaler

import (
	"errors"
	"testing"

	demo "github.com/SuperJourney/grpc_marshaler/example"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

var grpcErrWithDetail, _ = gstatus.New(codes.InvalidArgument, "这个是一个错误").WithDetails(&demo.ErrMsg{
	BusinessCode: 100,
	BusinessMsg:  "测试错误",
})

func TestMarshaler_MarshalWrapper(t *testing.T) {
	type args struct {
		respBody []interface{}
	}

	tests := []struct {
		name      string
		m         *Marshaler
		args      args
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "test response return succ case",
			m:    &Marshaler{},
			args: args{
				respBody: []interface{}{
					&demo.GetResponse{
						Name: "小红",
					},
					nil,
				},
			},

			want:      `{"Body":"CgblsI/nuqI=","Err":null}`,
			assertion: assert.NoError,
		},

		{
			name: "test response return grpc err case",
			m:    &Marshaler{},
			args: args{
				respBody: []interface{}{
					nil,
					gstatus.New(codes.InvalidArgument, "这个是一个错误").Err(),
				},
			},
			want:      `{"Body":null,"Err":{"code":3,"message":"这个是一个错误"}}`,
			assertion: assert.NoError,
		},

		{
			name: "test response return grpc with detail err case",
			m:    &Marshaler{},
			args: args{
				respBody: []interface{}{
					nil,
					grpcErrWithDetail.Err(),
				},
			},
			want:      `{"Body":null,"Err":{"code":3,"message":"这个是一个错误","details":[{"type_url":"type.googleapis.com/demo.ErrMsg","value":"CGQSDOa1i+ivlemUmeivrw=="}]}}`,
			assertion: assert.NoError,
		},

		{
			name: "test response return custom err case",
			m:    &Marshaler{},
			args: args{
				respBody: []interface{}{
					nil,
					errors.New("这个是一个错误"),
				},
			},
			want:      `{"Body":null,"Err":{"code":2,"message":"这个是一个错误"}}`,
			assertion: assert.NoError,
		},
		{
			name: "test response return nil case",
			m:    &Marshaler{},
			args: args{
				respBody: []interface{}{
					(*demo.GetResponse)(nil),
					errors.New("这个是一个错误"),
				},
			},
			want:      `{"Body":null,"Err":{"code":2,"message":"这个是一个错误"}}`,
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.MarshalWrapper(tt.args.respBody...)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

func TestMarshaler_UnMarshalWrapper(t *testing.T) {
	var err error
	type args struct {
		strings string
		resp    any
	}
	tests := []struct {
		name      string
		m         *Marshaler
		args      args
		want      []interface{}
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "test return succ case",
			m:    &Marshaler{},
			args: args{
				strings: `{"Body":"CgblsI/nuqI=","Err":null}`,
				resp:    []interface{}{&demo.GetResponse{}, err},
			},
			want: []interface{}{&demo.GetResponse{
				Name: "小红",
			}, nil},
			assertion: assert.NoError,
		},

		{
			name: "test response return grpc err case",
			m:    &Marshaler{},
			args: args{
				strings: `{"Body":null,"Err":{"code":3,"message":"这个是一个错误"}}`,
				resp:    []interface{}{&demo.GetResponse{}, err},
			},
			want:      []interface{}{nil, gstatus.New(codes.InvalidArgument, "这个是一个错误").Err()},
			assertion: assert.NoError,
		},

		{
			name: "test response return grpc with detail err case",
			m:    &Marshaler{},
			args: args{
				strings: `{"Body":null,"Err":{"code":3,"message":"这个是一个错误","details":[{"type_url":"type.googleapis.com/demo.ErrMsg","value":"CGQSDOa1i+ivlemUmeivrw=="}]}}`,
				resp:    []interface{}{&demo.GetResponse{}, err},
			},
			want:      []interface{}{nil, grpcErrWithDetail.Err()},
			assertion: assert.NoError,
		},
		{
			name: "test response return grpc with detail err case",
			m:    &Marshaler{},
			args: args{
				strings: `{"Body":null,"Err":{"code":2,"message":"这个是一个错误"}}`,
				resp:    []interface{}{&demo.GetResponse{}, err},
			},
			want:      []interface{}{nil, gstatus.New(codes.Unknown, "这个是一个错误").Err()},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.UnMarshalWrapper([]byte(tt.args.strings), tt.args.resp)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
