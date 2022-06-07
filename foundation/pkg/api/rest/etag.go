package rest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cespare/xxhash"
	"github.com/emicklei/go-restful/v3"
	"github.com/golang/protobuf/proto"
	"github.com/oxtoacart/bpool"
)

type ETagResponder struct {
	responseBufferPool *bpool.BufferPool
}

func NewETagResponder(bufferPoolSize int) *ETagResponder {
	return &ETagResponder{responseBufferPool: bpool.NewBufferPool(bufferPoolSize)}
}

func (eTagResponder *ETagResponder) RespondGetJSON(
	req *restful.Request, resp *restful.Response, data interface{},
) error {
	buf := eTagResponder.responseBufferPool.Get()
	defer eTagResponder.responseBufferPool.Put(buf)

	encoder := json.NewEncoder(buf)
	err := encoder.Encode(data)
	if err != nil {
		return err
	}

	h := xxhash.Sum64(buf.Bytes())
	eTagStr := fmt.Sprintf(`W/"%d-%x"`, buf.Len(), h)

	inm := req.Request.Header.Get("If-None-Match")
	if inm == eTagStr {
		resp.WriteHeader(http.StatusNotModified)
		return nil
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("ETag", eTagStr)
	resp.WriteHeader(http.StatusOK)
	resp.Write(buf.Bytes())
	return nil
}

func (eTagResponder *ETagResponder) RespondGetProtoMessage(
	req *restful.Request, resp *restful.Response, protoMsg proto.Message,
) error {
	dataBytes, err := proto.Marshal(protoMsg)
	if err != nil {
		panic(err)
	}

	messageType := proto.MessageName(protoMsg)

	h := xxhash.Sum64(dataBytes)
	eTagStr := fmt.Sprintf(`W/"%d-%x"`, len(dataBytes), h)

	inm := req.Request.Header.Get("If-None-Match")
	if inm == eTagStr {
		resp.WriteHeader(http.StatusNotModified)
		return nil
	}

	encoding := req.QueryParameter("__pbenc")
	if encoding == "base64" {
		resp.Header().Set("Content-Type", fmt.Sprintf(
			"application/x-protobuf-base64; messageType=%q", messageType))
		resp.Header().Set("ETag", eTagStr)
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(base64.StdEncoding.EncodeToString(dataBytes)))
		return nil
	}

	resp.Header().Set("Content-Type", fmt.Sprintf(
		"application/x-protobuf; messageType=%q", messageType))
	resp.Header().Set("ETag", eTagStr)
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(dataBytes))
	return nil
}

func (eTagResponder *ETagResponder) RespondGetOctetStream(
	req *restful.Request, resp *restful.Response, data []byte,
) error {
	h := xxhash.Sum64(data)
	eTagStr := fmt.Sprintf(`W/"%d-%x"`, len(data), h)

	inm := req.Request.Header.Get("If-None-Match")
	if inm == eTagStr {
		resp.WriteHeader(http.StatusNotModified)
		return nil
	}

	resp.Header().Set("Content-Type", "application/octet-stream")
	resp.Header().Set("ETag", eTagStr)
	resp.WriteHeader(http.StatusOK)
	//TODO: use buffer pool for the encoding process
	resp.Write([]byte(base64.StdEncoding.EncodeToString(data)))
	return nil
}
