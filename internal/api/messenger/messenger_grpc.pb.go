// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: api/proto/messenger.proto

package messenger

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MessengerClient is the client API for Messenger service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessengerClient interface {
	UserChatsWithOtherUsers(ctx context.Context, in *FeedChatRequest, opts ...grpc.CallOption) (*FeedChat, error)
	SendMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*MsgID, error)
	MessageFromChat(ctx context.Context, in *FeedMessageRequest, opts ...grpc.CallOption) (*FeedMessage, error)
	UpdateMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*empty.Empty, error)
	DeleteMessage(ctx context.Context, in *MsgID, opts ...grpc.CallOption) (*empty.Empty, error)
	GetMessage(ctx context.Context, in *MsgID, opts ...grpc.CallOption) (*Message, error)
}

type messengerClient struct {
	cc grpc.ClientConnInterface
}

func NewMessengerClient(cc grpc.ClientConnInterface) MessengerClient {
	return &messengerClient{cc}
}

func (c *messengerClient) UserChatsWithOtherUsers(ctx context.Context, in *FeedChatRequest, opts ...grpc.CallOption) (*FeedChat, error) {
	out := new(FeedChat)
	err := c.cc.Invoke(ctx, "/messenger.Messenger/UserChatsWithOtherUsers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messengerClient) SendMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*MsgID, error) {
	out := new(MsgID)
	err := c.cc.Invoke(ctx, "/messenger.Messenger/SendMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messengerClient) MessageFromChat(ctx context.Context, in *FeedMessageRequest, opts ...grpc.CallOption) (*FeedMessage, error) {
	out := new(FeedMessage)
	err := c.cc.Invoke(ctx, "/messenger.Messenger/MessageFromChat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messengerClient) UpdateMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/messenger.Messenger/UpdateMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messengerClient) DeleteMessage(ctx context.Context, in *MsgID, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/messenger.Messenger/DeleteMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messengerClient) GetMessage(ctx context.Context, in *MsgID, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/messenger.Messenger/GetMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessengerServer is the server API for Messenger service.
// All implementations must embed UnimplementedMessengerServer
// for forward compatibility
type MessengerServer interface {
	UserChatsWithOtherUsers(context.Context, *FeedChatRequest) (*FeedChat, error)
	SendMessage(context.Context, *Message) (*MsgID, error)
	MessageFromChat(context.Context, *FeedMessageRequest) (*FeedMessage, error)
	UpdateMessage(context.Context, *Message) (*empty.Empty, error)
	DeleteMessage(context.Context, *MsgID) (*empty.Empty, error)
	GetMessage(context.Context, *MsgID) (*Message, error)
	mustEmbedUnimplementedMessengerServer()
}

// UnimplementedMessengerServer must be embedded to have forward compatible implementations.
type UnimplementedMessengerServer struct {
}

func (UnimplementedMessengerServer) UserChatsWithOtherUsers(context.Context, *FeedChatRequest) (*FeedChat, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserChatsWithOtherUsers not implemented")
}
func (UnimplementedMessengerServer) SendMessage(context.Context, *Message) (*MsgID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedMessengerServer) MessageFromChat(context.Context, *FeedMessageRequest) (*FeedMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MessageFromChat not implemented")
}
func (UnimplementedMessengerServer) UpdateMessage(context.Context, *Message) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMessage not implemented")
}
func (UnimplementedMessengerServer) DeleteMessage(context.Context, *MsgID) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMessage not implemented")
}
func (UnimplementedMessengerServer) GetMessage(context.Context, *MsgID) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMessage not implemented")
}
func (UnimplementedMessengerServer) mustEmbedUnimplementedMessengerServer() {}

// UnsafeMessengerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessengerServer will
// result in compilation errors.
type UnsafeMessengerServer interface {
	mustEmbedUnimplementedMessengerServer()
}

func RegisterMessengerServer(s grpc.ServiceRegistrar, srv MessengerServer) {
	s.RegisterService(&Messenger_ServiceDesc, srv)
}

func _Messenger_UserChatsWithOtherUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FeedChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessengerServer).UserChatsWithOtherUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messenger.Messenger/UserChatsWithOtherUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessengerServer).UserChatsWithOtherUsers(ctx, req.(*FeedChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Messenger_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessengerServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messenger.Messenger/SendMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessengerServer).SendMessage(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

func _Messenger_MessageFromChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FeedMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessengerServer).MessageFromChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messenger.Messenger/MessageFromChat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessengerServer).MessageFromChat(ctx, req.(*FeedMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Messenger_UpdateMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessengerServer).UpdateMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messenger.Messenger/UpdateMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessengerServer).UpdateMessage(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

func _Messenger_DeleteMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessengerServer).DeleteMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messenger.Messenger/DeleteMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessengerServer).DeleteMessage(ctx, req.(*MsgID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Messenger_GetMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessengerServer).GetMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messenger.Messenger/GetMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessengerServer).GetMessage(ctx, req.(*MsgID))
	}
	return interceptor(ctx, in, info, handler)
}

// Messenger_ServiceDesc is the grpc.ServiceDesc for Messenger service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Messenger_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "messenger.Messenger",
	HandlerType: (*MessengerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserChatsWithOtherUsers",
			Handler:    _Messenger_UserChatsWithOtherUsers_Handler,
		},
		{
			MethodName: "SendMessage",
			Handler:    _Messenger_SendMessage_Handler,
		},
		{
			MethodName: "MessageFromChat",
			Handler:    _Messenger_MessageFromChat_Handler,
		},
		{
			MethodName: "UpdateMessage",
			Handler:    _Messenger_UpdateMessage_Handler,
		},
		{
			MethodName: "DeleteMessage",
			Handler:    _Messenger_DeleteMessage_Handler,
		},
		{
			MethodName: "GetMessage",
			Handler:    _Messenger_GetMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/messenger.proto",
}
