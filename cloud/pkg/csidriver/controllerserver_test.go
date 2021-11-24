package csidriver

import (
	"context"
	"errors"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	. "github.com/onsi/gomega"
)

func TestCSICreateVolume(t *testing.T) {
	errBoom := errors.New("boom")
	tbl := []struct {
		name   string
		sendFn KubeEdgeSendFn
		req    *csi.CreateVolumeRequest
		res    *csi.CreateVolumeResponse
		err    string
	}{
		{
			name: "err missing volume name",
			sendFn: func(req interface{}, nodeID, volumeID, csiOp string, res interface{}, kubeEdgeEndpoint string) error {
				return nil
			},
			req: &csi.CreateVolumeRequest{},
			err: "rpc error: code = InvalidArgument desc = Name missing in request",
		},
		{
			name: "err missing volume capabilities",
			sendFn: func(req interface{}, nodeID, volumeID, csiOp string, res interface{}, kubeEdgeEndpoint string) error {
				return nil
			},
			req: &csi.CreateVolumeRequest{
				Name: "foo",
			},
			err: "rpc error: code = InvalidArgument desc = Volume Capabilities missing in request",
		},
		{
			name: "kubeedge send error",
			sendFn: func(req interface{}, nodeID, volumeID, csiOp string, res interface{}, kubeEdgeEndpoint string) error {
				return errBoom
			},
			err: errBoom.Error(),
			req: &csi.CreateVolumeRequest{
				Name:               "foo",
				VolumeCapabilities: []*csi.VolumeCapability{},
			},
		},
		{
			name: "successful response",
			sendFn: func(req interface{}, nodeID, volumeID, csiOp string, res interface{}, kubeEdgeEndpoint string) error {
				rq := req.(*csi.CreateVolumeRequest)
				rs := res.(*csi.CreateVolumeResponse)

				Expect(rq.Name).To(Equal("foo"))
				Expect(rq.CapacityRange.RequiredBytes).To(BeEquivalentTo(10))

				*rs = csi.CreateVolumeResponse{
					Volume: &csi.Volume{
						VolumeId: "1234",
					},
				}
				return nil
			},
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
					VolumeId:      "1234",
					CapacityBytes: 10,
					ContentSource: &csi.VolumeContentSource{
						Type: &csi.VolumeContentSource_Volume{
							Volume: &csi.VolumeContentSource_VolumeSource{
								VolumeId: "example",
							},
						},
					},
				},
			},
		},
	}

	for _, row := range tbl {
		RegisterTestingT(t)
		t.Run(row.name, func(t *testing.T) {
			cs := newControllerServer("foo", "foo")
			cs.sendFn = row.sendFn
			res, err := cs.CreateVolume(context.Background(), row.req)
			if row.err != "" {
				Expect(err).To(MatchError(row.err))
				Expect(res).To(BeNil())
			} else {
				Expect(err).To(BeNil())
				Expect(res).To(Equal(row.res))
			}
		})
	}
}

func TestCSIDeleteVolume(t *testing.T) {
	errBoom := errors.New("boom")
	tbl := []struct {
		name   string
		sendFn KubeEdgeSendFn
		req    *csi.DeleteVolumeRequest
		res    *csi.DeleteVolumeResponse
		err    string
	}{
		{
			name: "err missing volume id",
			sendFn: func(req interface{}, nodeID, volumeID, csiOp string, res interface{}, kubeEdgeEndpoint string) error {
				return nil
			},
			req: &csi.DeleteVolumeRequest{},
			err: "rpc error: code = InvalidArgument desc = Volume ID missing in request",
		},
		{
			name: "kubeedge send error",
			sendFn: func(req interface{}, nodeID, volumeID, csiOp string, res interface{}, kubeEdgeEndpoint string) error {
				return errBoom
			},
			err: errBoom.Error(),
			req: &csi.DeleteVolumeRequest{
				VolumeId: "foo",
			},
		},
		{
			name: "successful response",
			sendFn: func(req interface{}, nodeID, volumeID, csiOp string, res interface{}, kubeEdgeEndpoint string) error {
				rq := req.(*csi.DeleteVolumeRequest)
				rs := res.(*csi.DeleteVolumeResponse)

				Expect(rq.VolumeId).To(Equal("foo"))
				*rs = csi.DeleteVolumeResponse{}
				return nil
			},
			req: &csi.DeleteVolumeRequest{
				VolumeId: "foo",
			},
			res: &csi.DeleteVolumeResponse{},
		},
	}

	for _, row := range tbl {
		RegisterTestingT(t)
		t.Run(row.name, func(t *testing.T) {
			cs := newControllerServer("foo", "foo")
			cs.sendFn = row.sendFn
			res, err := cs.DeleteVolume(context.Background(), row.req)
			if row.err != "" {
				Expect(err).To(MatchError(row.err))
				Expect(res).To(BeNil())
			} else {
				Expect(err).To(BeNil())
				Expect(res).To(Equal(row.res))
			}
		})
	}
}

func TestCSIPublishVolume(t *testing.T) {
	errBoom := errors.New("boom")
	tbl := []struct {
		name   string
		sendFn KubeEdgeSendFn
		req    *csi.ControllerPublishVolumeRequest
		res    *csi.ControllerPublishVolumeResponse
		err    string
	}{
		{
			name: "err missing volume id",
			sendFn: func(req interface{}, nodeID, volumeID, csiOp string, res interface{}, kubeEdgeEndpoint string) error {
				return nil
			},
			req: &csi.ControllerPublishVolumeRequest{
				NodeId: "foo",
			},
			err: "rpc error: code = InvalidArgument desc = ControllerPublishVolume Volume ID must be provided",
		},
		{
			name: "err missing node id",
			sendFn: func(req interface{}, nodeID, volumeID, csiOp string, res interface{}, kubeEdgeEndpoint string) error {
				return nil
			},
			req: &csi.ControllerPublishVolumeRequest{
				VolumeId: "foo",
			},
			err: "rpc error: code = InvalidArgument desc = ControllerPublishVolume Node ID must be provided",
		},
		{
			name: "kubeedge send error",
			sendFn: func(req interface{}, nodeID, volumeID, csiOp string, res interface{}, kubeEdgeEndpoint string) error {
				return errBoom
			},
			err: errBoom.Error(),
			req: &csi.ControllerPublishVolumeRequest{
				VolumeId: "foo",
				NodeId:   "bar",
			},
		},
		{
			name: "successful response",
			sendFn: func(req interface{}, nodeID, volumeID, csiOp string, res interface{}, kubeEdgeEndpoint string) error {
				rq := req.(*csi.ControllerPublishVolumeRequest)
				rs := res.(*csi.ControllerPublishVolumeResponse)

				Expect(rq.VolumeId).To(Equal("foo"))
				Expect(rq.NodeId).To(Equal("bar"))
				*rs = csi.ControllerPublishVolumeResponse{
					PublishContext: map[string]string{
						"foo": "example",
					},
				}
				return nil
			},
			req: &csi.ControllerPublishVolumeRequest{
				VolumeId: "foo",
				NodeId:   "bar",
			},
			res: &csi.ControllerPublishVolumeResponse{
				PublishContext: map[string]string{
					"foo": "example",
				},
			},
		},
	}

	for _, row := range tbl {
		RegisterTestingT(t)
		t.Run(row.name, func(t *testing.T) {
			cs := newControllerServer("foo", "foo")
			cs.sendFn = row.sendFn
			res, err := cs.ControllerPublishVolume(context.Background(), row.req)
			if row.err != "" {
				Expect(err).To(MatchError(row.err))
				Expect(res).To(BeNil())
			} else {
				Expect(err).To(BeNil())
				Expect(res).To(Equal(row.res))
			}
		})
	}
}

func TestCSIUnpublishVolume(t *testing.T) {
	errBoom := errors.New("boom")
	tbl := []struct {
		name   string
		sendFn KubeEdgeSendFn
		req    *csi.ControllerUnpublishVolumeRequest
		res    *csi.ControllerUnpublishVolumeResponse
		err    string
	}{
		{
			name: "err missing volume id",
			sendFn: func(req interface{}, nodeID, volumeID, csiOp string, res interface{}, kubeEdgeEndpoint string) error {
				return nil
			},
			req: &csi.ControllerUnpublishVolumeRequest{
				NodeId: "foo",
			},
			err: "rpc error: code = InvalidArgument desc = ControllerUnpublishVolume Volume ID must be provided",
		},
		{
			name: "err missing node id",
			sendFn: func(req interface{}, nodeID, volumeID, csiOp string, res interface{}, kubeEdgeEndpoint string) error {
				return nil
			},
			req: &csi.ControllerUnpublishVolumeRequest{
				VolumeId: "foo",
			},
			err: "rpc error: code = InvalidArgument desc = ControllerUnpublishVolume Node ID must be provided",
		},
		{
			name: "kubeedge send error",
			sendFn: func(req interface{}, nodeID, volumeID, csiOp string, res interface{}, kubeEdgeEndpoint string) error {
				return errBoom
			},
			err: errBoom.Error(),
			req: &csi.ControllerUnpublishVolumeRequest{
				VolumeId: "foo",
				NodeId:   "bar",
			},
		},
		{
			name: "successful response",
			sendFn: func(req interface{}, nodeID, volumeID, csiOp string, res interface{}, kubeEdgeEndpoint string) error {
				rq := req.(*csi.ControllerUnpublishVolumeRequest)
				rs := res.(*csi.ControllerUnpublishVolumeResponse)

				Expect(rq.VolumeId).To(Equal("foo"))
				Expect(rq.NodeId).To(Equal("bar"))
				*rs = csi.ControllerUnpublishVolumeResponse{}
				return nil
			},
			req: &csi.ControllerUnpublishVolumeRequest{
				VolumeId: "foo",
				NodeId:   "bar",
			},
			res: &csi.ControllerUnpublishVolumeResponse{},
		},
	}

	for _, row := range tbl {
		RegisterTestingT(t)
		t.Run(row.name, func(t *testing.T) {
			cs := newControllerServer("foo", "foo")
			cs.sendFn = row.sendFn
			res, err := cs.ControllerUnpublishVolume(context.Background(), row.req)
			if row.err != "" {
				Expect(err).To(MatchError(row.err))
				Expect(res).To(BeNil())
			} else {
				Expect(err).To(BeNil())
				Expect(res).To(Equal(row.res))
			}
		})
	}
}

func TestCSIValidateVolumeCapabilities(t *testing.T) {
	tbl := []struct {
		name string
		req  *csi.ValidateVolumeCapabilitiesRequest
		res  *csi.ValidateVolumeCapabilitiesResponse
		err  string
	}{
		{
			name: "err missing volume id",

			req: &csi.ValidateVolumeCapabilitiesRequest{},
			err: "rpc error: code = InvalidArgument desc = Volume ID cannot be empty",
		},
		{
			name: "err missing node id",
			req: &csi.ValidateVolumeCapabilitiesRequest{
				VolumeId: "foo",
			},
			err: "rpc error: code = InvalidArgument desc = Volume Capabilities can not be empty",
		},
		{
			name: "err missing access type",
			req: &csi.ValidateVolumeCapabilitiesRequest{
				VolumeId: "foo",
				VolumeContext: map[string]string{
					"foo": "example",
				},
				VolumeCapabilities: []*csi.VolumeCapability{
					{},
				},
			},
			err: "rpc error: code = InvalidArgument desc = Cannot have both mount and block access type be undefined",
		},
		{
			name: "successful response",
			req: &csi.ValidateVolumeCapabilitiesRequest{
				VolumeId: "foo",
				VolumeContext: map[string]string{
					"foo": "example",
				},
				VolumeCapabilities: []*csi.VolumeCapability{
					{
						AccessType: &csi.VolumeCapability_Mount{
							Mount: &csi.VolumeCapability_MountVolume{},
						},
						AccessMode: &csi.VolumeCapability_AccessMode{
							Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY,
						},
					},
				},
			},
			res: &csi.ValidateVolumeCapabilitiesResponse{
				Confirmed: &csi.ValidateVolumeCapabilitiesResponse_Confirmed{
					VolumeContext: map[string]string{
						"foo": "example",
					},
					VolumeCapabilities: []*csi.VolumeCapability{
						{
							AccessType: &csi.VolumeCapability_Mount{
								Mount: &csi.VolumeCapability_MountVolume{},
							},
							AccessMode: &csi.VolumeCapability_AccessMode{
								Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY,
							},
						},
					},
				},
			},
		},
	}

	for _, row := range tbl {
		RegisterTestingT(t)
		t.Run(row.name, func(t *testing.T) {
			cs := newControllerServer("foo", "foo")
			res, err := cs.ValidateVolumeCapabilities(context.Background(), row.req)
			if row.err != "" {
				Expect(err).To(MatchError(row.err))
				Expect(res).To(BeNil())
			} else {
				Expect(err).To(BeNil())
				Expect(res).To(Equal(row.res))
			}
		})
	}
}
