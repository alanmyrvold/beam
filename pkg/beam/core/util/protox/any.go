package protox

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	protobuf "github.com/golang/protobuf/ptypes/any"
	protobufw "github.com/golang/protobuf/ptypes/wrappers"
)

const (
	bytesValueTypeURL = "type.googleapis.com/google.protobuf.BytesValue"
)

// Unpack decodes a proto.
func Unpack(data *protobuf.Any, url string, ret proto.Message) error {
	if data.TypeUrl != url {
		return fmt.Errorf("Bad type: %v, want %v", data.TypeUrl, url)
	}
	return proto.Unmarshal(data.Value, ret)
}

// UnpackProto decodes a BytesValue wrapped proto.
func UnpackProto(data *protobuf.Any, ret proto.Message) error {
	buf, err := UnpackBytes(data)
	if err != nil {
		return err
	}
	return proto.Unmarshal(buf, ret)
}

// PackProto encodes a proto message and wraps into a BytesValue.
func PackProto(in proto.Message) (*protobuf.Any, error) {
	b, err := proto.Marshal(in)
	if err != nil {
		return nil, err
	}

	return PackBytes(b)
}

// UnpackBase64Proto decodes a BytesValue + base64 wrapped proto.
func UnpackBase64Proto(data *protobuf.Any, ret proto.Message) error {
	buf, err := UnpackBytes(data)
	if err != nil {
		return err
	}
	return DecodeBase64(string(buf), ret)
}

// PackBase64Proto encodes a proto message into a base64-encoded BytesValue message.
func PackBase64Proto(in proto.Message) (*protobuf.Any, error) {
	data, err := EncodeBase64(in)
	if err != nil {
		return nil, err
	}
	return PackBytes([]byte(data))
}

// UnpackBytes removes the BytesValue wrapper.
func UnpackBytes(data *protobuf.Any) ([]byte, error) {
	if data.TypeUrl != bytesValueTypeURL {
		return nil, fmt.Errorf("Bad type: %v, want %v", data.TypeUrl, bytesValueTypeURL)
	}

	var buf protobufw.BytesValue
	if err := proto.Unmarshal(data.Value, &buf); err != nil {
		return nil, fmt.Errorf("BytesValue unmarshal failed: %v", err)
	}
	return buf.Value, nil
}

// PackBytes applies a BytesValue wrapper to the supplied bytes.
func PackBytes(in []byte) (*protobuf.Any, error) {
	var buf protobufw.BytesValue
	buf.Value = in
	b, err := proto.Marshal(&buf)
	if err != nil {
		return nil, err
	}

	var out protobuf.Any
	out.TypeUrl = bytesValueTypeURL
	out.Value = b
	return &out, nil
}
