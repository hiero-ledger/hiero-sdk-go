package hiero

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// SPDX-License-Identifier: Apache-2.0

type HookExtensionPoint int32

const (
	ACCOUNT_ALLOWANCE_HOOK HookExtensionPoint = iota
)

// The details of a hook's creation.
type HookCreationDetails struct {
	extensionPoint HookExtensionPoint
	hookId         int64
	lambdaEvmHook  LambdaEvmHook
	adminKey       Key
}

// NewHookCreationDetails creates a new HookCreationDetails
func NewHookCreationDetails() *HookCreationDetails {
	return &HookCreationDetails{}
}

// GetExtensionPoint returns the hook extension point
func (hcd HookCreationDetails) GetExtensionPoint() HookExtensionPoint {
	return hcd.extensionPoint
}

// SetExtensionPoint sets the hook extension point
func (hcd *HookCreationDetails) SetExtensionPoint(extensionPoint HookExtensionPoint) *HookCreationDetails {
	hcd.extensionPoint = extensionPoint
	return hcd
}

// GetHookId returns the hook ID
func (hcd HookCreationDetails) GetHookId() int64 {
	return hcd.hookId
}

// SetHookId sets the hook ID
func (hcd *HookCreationDetails) SetHookId(hookId int64) *HookCreationDetails {
	hcd.hookId = hookId
	return hcd
}

// GetLambdaEvmHook returns the lambda EVM hook
func (hcd HookCreationDetails) GetLambdaEvmHook() LambdaEvmHook {
	return hcd.lambdaEvmHook
}

// SetLambdaEvmHook sets the lambda EVM hook
func (hcd *HookCreationDetails) SetLambdaEvmHook(lambdaEvmHook LambdaEvmHook) *HookCreationDetails {
	hcd.lambdaEvmHook = lambdaEvmHook
	return hcd
}

// GetAdminKey returns the admin key
func (hcd HookCreationDetails) GetAdminKey() Key {
	return hcd.adminKey
}

// SetAdminKey sets the admin key
func (hcd *HookCreationDetails) SetAdminKey(adminKey Key) *HookCreationDetails {
	hcd.adminKey = adminKey
	return hcd
}

func hookCreationDetailsFromProtobuf(pb *services.HookCreationDetails) HookCreationDetails {
	var key Key
	if pb.GetAdminKey() != nil {
		key, _ = _KeyFromProtobuf(pb.GetAdminKey())
	}
	return HookCreationDetails{
		extensionPoint: HookExtensionPoint(pb.GetExtensionPoint()),
		hookId:         pb.GetHookId(),
		lambdaEvmHook:  lambdaEvmHookFromProtobuf(pb.GetLambdaEvmHook()),
		adminKey:       key,
	}
}

func (hcd HookCreationDetails) toProtobuf() *services.HookCreationDetails {
	var adminKey *services.Key
	if hcd.adminKey != nil {
		adminKey = hcd.adminKey._ToProtoKey()
	}

	protoBody := &services.HookCreationDetails{
		ExtensionPoint: services.HookExtensionPoint(hcd.extensionPoint),
		HookId:         hcd.hookId,
		Hook: &services.HookCreationDetails_LambdaEvmHook{
			LambdaEvmHook: hcd.lambdaEvmHook.toProtobuf(),
		},
		AdminKey: adminKey,
	}

	return protoBody
}
