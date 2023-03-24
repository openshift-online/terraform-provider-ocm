// Code generated by MockGen. DO NOT EDIT.
// Source: provider/cluster_resource.go

// Package provider is a generated GoMock package.
package provider

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	diag "github.com/hashicorp/terraform-plugin-framework/diag"
	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

// MockClusterClient is a mock of ClusterClient interface.
type MockClusterClient struct {
	ctrl     *gomock.Controller
	recorder *MockClusterClientMockRecorder
}

// MockClusterClientMockRecorder is the mock recorder for MockClusterClient.
type MockClusterClientMockRecorder struct {
	mock *MockClusterClient
}

// NewMockClusterClient creates a new mock instance.
func NewMockClusterClient(ctrl *gomock.Controller) *MockClusterClient {
	mock := &MockClusterClient{ctrl: ctrl}
	mock.recorder = &MockClusterClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClusterClient) EXPECT() *MockClusterClientMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockClusterClient) Create(ctx context.Context, cluster *v1.Cluster) (*v1.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, cluster)
	ret0, _ := ret[0].(*v1.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockClusterClientMockRecorder) Create(ctx, cluster interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockClusterClient)(nil).Create), ctx, cluster)
}

// Delete mocks base method.
func (m *MockClusterClient) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockClusterClientMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockClusterClient)(nil).Delete), ctx, id)
}

// Get mocks base method.
func (m *MockClusterClient) Get(ctx context.Context, id string) (*v1.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*v1.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockClusterClientMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockClusterClient)(nil).Get), ctx, id)
}

// PollReady mocks base method.
func (m *MockClusterClient) PollReady(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PollReady", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// PollReady indicates an expected call of PollReady.
func (mr *MockClusterClientMockRecorder) PollReady(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PollReady", reflect.TypeOf((*MockClusterClient)(nil).PollReady), ctx, id)
}

// PollRemoved mocks base method.
func (m *MockClusterClient) PollRemoved(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PollRemoved", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// PollRemoved indicates an expected call of PollRemoved.
func (mr *MockClusterClientMockRecorder) PollRemoved(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PollRemoved", reflect.TypeOf((*MockClusterClient)(nil).PollRemoved), ctx, id)
}

// Update mocks base method.
func (m *MockClusterClient) Update(ctx context.Context, id string, cluster *v1.Cluster) (*v1.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, cluster)
	ret0, _ := ret[0].(*v1.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockClusterClientMockRecorder) Update(ctx, id, cluster interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockClusterClient)(nil).Update), ctx, id, cluster)
}

// MockClusterResourceUtils is a mock of ClusterResourceUtils interface.
type MockClusterResourceUtils struct {
	ctrl     *gomock.Controller
	recorder *MockClusterResourceUtilsMockRecorder
}

// MockClusterResourceUtilsMockRecorder is the mock recorder for MockClusterResourceUtils.
type MockClusterResourceUtilsMockRecorder struct {
	mock *MockClusterResourceUtils
}

// NewMockClusterResourceUtils creates a new mock instance.
func NewMockClusterResourceUtils(ctrl *gomock.Controller) *MockClusterResourceUtils {
	mock := &MockClusterResourceUtils{ctrl: ctrl}
	mock.recorder = &MockClusterResourceUtilsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClusterResourceUtils) EXPECT() *MockClusterResourceUtilsMockRecorder {
	return m.recorder
}

// createClusterObject mocks base method.
func (m *MockClusterResourceUtils) createClusterObject(ctx context.Context, state *ClusterState, diags diag.Diagnostics) (*v1.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "createClusterObject", ctx, state, diags)
	ret0, _ := ret[0].(*v1.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// createClusterObject indicates an expected call of createClusterObject.
func (mr *MockClusterResourceUtilsMockRecorder) createClusterObject(ctx, state, diags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "createClusterObject", reflect.TypeOf((*MockClusterResourceUtils)(nil).createClusterObject), ctx, state, diags)
}

// populateClusterState mocks base method.
func (m *MockClusterResourceUtils) populateClusterState(object *v1.Cluster, state *ClusterState) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "populateClusterState", object, state)
}

// populateClusterState indicates an expected call of populateClusterState.
func (mr *MockClusterResourceUtilsMockRecorder) populateClusterState(object, state interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "populateClusterState", reflect.TypeOf((*MockClusterResourceUtils)(nil).populateClusterState), object, state)
}

// MockDiagnostics is a mock of Diagnostics interface.
type MockDiagnostics struct {
	ctrl     *gomock.Controller
	recorder *MockDiagnosticsMockRecorder
}

// MockDiagnosticsMockRecorder is the mock recorder for MockDiagnostics.
type MockDiagnosticsMockRecorder struct {
	mock *MockDiagnostics
}

// NewMockDiagnostics creates a new mock instance.
func NewMockDiagnostics(ctrl *gomock.Controller) *MockDiagnostics {
	mock := &MockDiagnostics{ctrl: ctrl}
	mock.recorder = &MockDiagnosticsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDiagnostics) EXPECT() *MockDiagnosticsMockRecorder {
	return m.recorder
}

// AddError mocks base method.
func (m *MockDiagnostics) AddError(title, description string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddError", title, description)
}

// AddError indicates an expected call of AddError.
func (mr *MockDiagnosticsMockRecorder) AddError(title, description interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddError", reflect.TypeOf((*MockDiagnostics)(nil).AddError), title, description)
}