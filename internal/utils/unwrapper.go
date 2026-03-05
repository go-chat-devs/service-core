package utils

import "google.golang.org/protobuf/types/known/wrapperspb"

func UnwrapStringValue(wrapper *wrapperspb.StringValue) *string {
	if wrapper == nil {
		return nil
	}
	return &wrapper.Value
}
