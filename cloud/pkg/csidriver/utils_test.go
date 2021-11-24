package csidriver

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/kubeedge/beehive/pkg/core/model"
	"github.com/kubeedge/kubeedge/cloud/pkg/cloudhub/servers/udsserver"
	"github.com/kubeedge/kubeedge/common/constants"
	. "github.com/onsi/gomega"
)

func TestSendToKubeEdge(t *testing.T) {
	RegisterTestingT(t)
	addr := "unix:///tmp/udssrv-" + strconv.FormatInt(time.Now().Unix(), 10)
	srv := udsserver.NewUnixDomainSocket(addr, udsserver.DefaultBufferSize)
	go srv.StartServer()
	tbl := []struct {
		name     string
		handler  func(string) string
		req      interface{}
		nodeID   string
		volumeID string
		op       string
		res      interface{}
		err      string
	}{
		{
			name: "missing parameters",
			err:  "failed to build message resource: required parameter are not set (node id, namespace or resource type)",
			req:  &csi.CreateVolumeRequest{},
			res:  &csi.CreateVolumeResponse{},
			handler: func(s string) string {
				return ""
			},
		},
		{
			name:     "unexpected respone msg",
			req:      &csi.CreateVolumeRequest{},
			res:      &csi.CreateVolumeResponse{},
			volumeID: "foo",
			nodeID:   "bar",
			op:       constants.CSIOperationTypeCreateVolume,
			err:      "failed to unmarshal beehive msg: invalid character 'N' looking for beginning of value",
			handler: func(s string) string {
				return "NOT JSON"
			},
		},
		{
			name:     "unexpected respone operation",
			req:      &csi.CreateVolumeRequest{},
			res:      &csi.CreateVolumeResponse{},
			volumeID: "foo",
			nodeID:   "bar",
			op:       constants.CSIOperationTypeCreateVolume,
			err:      "unexpected message error: &model.Message{Header:model.MessageHeader{ID:\"\", ParentID:\"\", Timestamp:0, ResourceVersion:\"\", Sync:false, MessageType:\"\"}, Router:model.MessageRoute{Source:\"\", Destination:\"\", Group:\"\", Operation:\"error\", Resource:\"\"}, Content:\"some-content\"}",
			handler: func(s string) string {
				return `{"route":{"operation":"error"},"content":"some-content"}`
			},
		},
		{
			name:     "unexpected message content",
			req:      &csi.CreateVolumeRequest{},
			res:      &csi.CreateVolumeResponse{},
			volumeID: "foo",
			nodeID:   "bar",
			op:       constants.CSIOperationTypeCreateVolume,
			err:      "unexpected message content type: %!s(float64=1.23451234123e+11)",
			handler: func(s string) string {
				return `{"route":{"operation":"response"},"content": 123451234123}`
			},
		},
		{
			name:     "successful cycle",
			volumeID: "foo",
			nodeID:   "bar",
			op:       constants.CSIOperationTypeCreateVolume,
			req: &csi.CreateVolumeRequest{
				Name: "foo",
				CapacityRange: &csi.CapacityRange{
					RequiredBytes: 10,
					LimitBytes:    20,
				},
				VolumeContentSource: &csi.VolumeContentSource{
					Type: &csi.VolumeContentSource_Volume{
						Volume: &csi.VolumeContentSource_VolumeSource{
							VolumeId: "example",
						},
					},
				},
				VolumeCapabilities: []*csi.VolumeCapability{},
			},
			res: &csi.CreateVolumeResponse{
				Volume: &csi.Volume{
					VolumeId: "1234",
				},
			},
			handler: func(s string) string {
				// check that msg was serialized correctly
				var im model.Message
				err := json.Unmarshal([]byte(s), &im)
				Expect(err).ToNot(HaveOccurred())
				Expect(im.Content).To(Equal(`{"name":"foo","capacityRange":{"requiredBytes":"10","limitBytes":"20"},"volumeContentSource":{"volume":{"volumeId":"example"}}}`))
				Expect(im.Router.Operation).To(Equal("createvolume"))
				Expect(im.Router.Resource).To(Equal("node/bar/default/volume/foo"))

				// this basicially simulates what edged is doing
				// see edged.go
				// csi res -> json -> base64 -> model -> json string
				pbm := jsonpb.Marshaler{}
				emb, err := pbm.MarshalToString(&csi.CreateVolumeResponse{
					Volume: &csi.Volume{
						VolumeId: "1234",
					},
				})
				Expect(err).ToNot(HaveOccurred())
				m := model.NewMessage("")
				m.FillBody(emb)
				b, err := json.Marshal(m)
				Expect(err).ToNot(HaveOccurred())
				return string(b)
			},
		},
	}

	for _, row := range tbl {
		t.Run(row.name, func(t *testing.T) {
			srv.SetContextHandler(row.handler)
			// clone & reset response so we can check if the value
			// was modified as expected
			resClone := proto.Clone(row.res.(proto.Message))
			resClone.Reset()
			err := sendToKubeEdge(row.req, row.nodeID, row.volumeID, row.op, resClone, addr)
			if row.err != "" {
				Expect(err).To(MatchError(row.err))
			} else {
				Expect(err).To(BeNil())
				Expect(resClone).To(Equal(row.res))
			}
		})
	}
}
