package csidriver

import (
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"golang.org/x/net/context"
	"k8s.io/klog/v2"
)

type nodeServer struct {
	nodeID string
}

func newNodeServer(nodeID string) *nodeServer {
	return &nodeServer{
		nodeID: nodeID,
	}
}

const (
	DriverName = "csidriver.kubeedge.com"

	// users must label their nodes in the following way
	//
	TopologyEdgeKey = "topology." + DriverName + "/edge"

	WellKnownTopologyKey = "topology.kubernetes.io/zone"
)

var ErrNotImpl = fmt.Errorf("not implemented")

func (d *nodeServer) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	klog.V(4).Infof("NodeExpandVolume: called with args %+v", *req)
	return nil, ErrNotImpl
}

func (d *nodeServer) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	klog.V(4).Infof("NodeGetCapabilities: called with args %+v", *req)
	return nil, ErrNotImpl
}

func (d *nodeServer) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	klog.V(4).Infof("NodeGetInfo: called with args %+v", *req)

	return &csi.NodeGetInfoResponse{
		NodeId: fmt.Sprintf("csidriver-kubeedge-%s", d.nodeID),
		// no topology constraints by this driver
		// TODO: add topology constraints to allow group edge nodes by csidriver
		AccessibleTopology: nil,
	}, nil

}

func (d *nodeServer) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	klog.V(4).Infof("NodeGetVolumeStats: called with args %+v", *req)
	return nil, ErrNotImpl
}

func (d *nodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	klog.V(4).Infof("NodePublishVolume: called with args %+v", *req)
	return nil, ErrNotImpl
}

func (d *nodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	klog.V(4).Infof("NodeUnpublishVolume: called with args %+v", *req)
	return nil, ErrNotImpl
}

func (d *nodeServer) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	klog.V(4).Infof("NodeStageVolume: called with args %+v", *req)
	return nil, ErrNotImpl
}

func (d *nodeServer) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	klog.V(4).Infof("NodeUnstageVolume: called with args %+v", *req)
	return nil, ErrNotImpl
}
