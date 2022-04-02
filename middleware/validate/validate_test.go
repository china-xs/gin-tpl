package validate

import (
	"reflect"
	"testing"
)

func TestGetValidateKey(t *testing.T) {
	type Want struct {
		Key    string
		Params map[string]interface{}
	}
	tests := []struct {
		name string
		args string
		want Want
	}{
		{
			name: "float.const",
			args: "invalid CreateUserRequest.Price: value must equal 0.99",
			want: Want{
				Key: "CreateUserRequest.Price.const",
				Params: map[string]interface{}{
					"const": "0.99",
				},
			},
		},
		{
			name: "float.lt",
			args: "invalid CreateUserRequest.Price: value must be less than 0.99",
			want: Want{
				"CreateUserRequest.Price.lt",
				map[string]interface{}{
					"lt": "0.99",
				},
			},
		},
		{
			name: "float.lte",
			args: "invalid CreateUserRequest.Price: value must be less than or equal to 0.99",
			want: Want{
				"CreateUserRequest.Price.lte",
				map[string]interface{}{
					"lte": "0.99",
				},
			},
		},
		{
			name: "float.gt",
			args: "invalid CreateUserRequest.Price: value must be greater than 0.99",
			want: Want{
				"CreateUserRequest.Price.gt",
				map[string]interface{}{
					"gt": "0.99",
				},
			},
		},
		{
			name: "float.gte",
			args: "invalid CreateUserRequest.Price: value must be greater than or equal to 0.99",
			want: Want{
				"CreateUserRequest.Price.gte",
				map[string]interface{}{
					"gte": "0.99",
				},
			},
		}, {
			name: "float.in",
			args: "invalid CreateUserRequest.Price: value must be in list [0.99 1 1.1]",
			want: Want{
				"CreateUserRequest.Price.in",
				map[string]interface{}{
					"in": "[0.99 1 1.1]",
				},
			},
		}, {
			name: "float.not_in",
			args: "invalid CreateUserRequest.Price: value must not be in list [0.99 1 1.1]",
			want: Want{
				"CreateUserRequest.Price.not_in",
				map[string]interface{}{
					"not_in": "[0.99 1 1.1]",
				},
			},
		},
		{
			name: "bool.const",
			args: "invalid CreateUserRequest.IsTrue: value must equal true",
			want: Want{
				"CreateUserRequest.IsTrue.const",
				map[string]interface{}{
					"const": "true",
				},
			},
		}, {
			name: "string.const",
			args: "invalid CreateUserRequest.Name: value must equal 闪灵",
			want: Want{
				"CreateUserRequest.Name.const",
				map[string]interface{}{
					"const": "闪灵",
				},
			},
		}, {
			name: "string.len",
			args: "invalid CreateUserRequest.Name: value length must be 2 runes",
			want: Want{
				"CreateUserRequest.Name.len",
				map[string]interface{}{
					"len": "2",
				},
			},
		}, {
			name: "string.min_len",
			args: "invalid CreateUserRequest.Name: value length must be at least 2 runes",
			want: Want{
				"CreateUserRequest.Name.min_len",
				map[string]interface{}{
					"min_len": "2",
				},
			},
		}, {
			name: "string.max_len",
			args: "invalid CreateUserRequest.Name: value length must be at most 2 runes",
			want: Want{
				"CreateUserRequest.Name.max_len",
				map[string]interface{}{
					"max_len": "2",
				},
			},
		}, {
			name: "string.min_len&max_len",
			args: "invalid CreateUserRequest.Name: value length must be between 2 and 10 runes, inclusive",
			want: Want{
				"CreateUserRequest.Name.between",
				map[string]interface{}{
					"min_len": "2",
					"max_len": "10",
				},
			},
		}, {
			name: "string.len_bytes",
			args: "invalid CreateUserRequest.Name: value length must be 2 bytes",
			want: Want{
				"CreateUserRequest.Name.len_bytes",
				map[string]interface{}{
					"len_bytes": "2",
				},
			},
		}, {
			name: "string.min_bytes",
			args: "invalid CreateUserRequest.Name: value length must be at least 2 bytes",
			want: Want{
				"CreateUserRequest.Name.min_bytes",
				map[string]interface{}{
					"min_bytes": "2",
				},
			},
		}, {
			name: "string.max_bytes",
			args: "invalid CreateUserRequest.Name: value length must be at most 2 bytes",
			want: Want{
				"CreateUserRequest.Name.max_bytes",
				map[string]interface{}{
					"max_bytes": "2",
				},
			},
		}, {
			name: "string.min_bytes&max_bytes",
			args: "invalid CreateUserRequest.Name: value length must be between 2 and 10 bytes, inclusive",
			want: Want{
				"CreateUserRequest.Name.between",
				map[string]interface{}{
					"min_bytes": "2",
					"max_bytes": "10",
				},
			},
		}, {
			name: "string.pattern",
			args: "invalid CreateUserRequest.Name: value does not match regex pattern \"^[^[0-9]A-Za-z]+( [^[0-9]A-Za-z]+)*$\"",
			want: Want{
				"CreateUserRequest.Name.pattern",
				map[string]interface{}{},
			},
		}, {
			name: "string.prefix",
			args: "invalid CreateUserRequest.Name: value does not have prefix \"YJF_\"",
			want: Want{
				"CreateUserRequest.Name.prefix",
				map[string]interface{}{},
			},
		}, {
			name: "string.suffix",
			args: "invalid CreateUserRequest.Name: value does not have suffix \"_YJF\"",
			want: Want{
				"CreateUserRequest.Name.suffix",
				map[string]interface{}{},
			},
		}, {
			name: "string.contains",
			args: "invalid CreateUserRequest.Name: value does not contain substring \"YJF\"",
			want: Want{
				"CreateUserRequest.Name.contains",
				map[string]interface{}{},
			},
		}, {
			name: "strings.not_contains",
			args: "invalid CreateUserRequest.Name: value contains substring \"YJF\"",
			want: Want{
				"CreateUserRequest.Name.not_contains",
				map[string]interface{}{},
			},
		}, {
			name: "string.in",
			args: "invalid CreateUserRequest.Name: value must be in list [莫佳品 小米 华为]",
			want: Want{
				"CreateUserRequest.Name.in",
				map[string]interface{}{
					"in": "[莫佳品 小米 华为]",
				},
			},
		}, {
			name: "string.not_in",
			args: "invalid CreateUserRequest.Name: value must not be in list [莫佳品 小米 华为]",
			want: Want{
				"CreateUserRequest.Name.not_in",
				map[string]interface{}{
					"not_in": "[莫佳品 小米 华为]",
				},
			},
		}, {
			name: "string.email",
			args: "invalid CreateUserRequest.Email: value must be a valid email address | caused by: mail: missing '@' or angle-addr",
			want: Want{
				"CreateUserRequest.Email.email",
				map[string]interface{}{},
			},
		}, {
			name: "string.hostname",
			args: "invalid CreateUserRequest.Hostname: value must be a valid hostname | caused by: hostname part must be non-empty and cannot exceed 63 characters",
			want: Want{
				"CreateUserRequest.Hostname.hostname",
				map[string]interface{}{},
			},
		}, {
			name: "stirng.ip",
			args: "invalid CreateUserRequest.Ip: value must be a valid IP address",
			want: Want{
				"CreateUserRequest.Ip.ip",
				map[string]interface{}{},
			},
		}, {
			name: "string.ipv4",
			args: "invalid CreateUserRequest.Ip: value must be a valid IPv4 address",
			want: Want{
				"CreateUserRequest.Ip.ipv4",
				map[string]interface{}{},
			},
		}, {
			name: "string.ipv6",
			args: "invalid CreateUserRequest.Ip: value must be a valid IPv6 address",
			want: Want{
				"CreateUserRequest.Ip.ipv6",
				map[string]interface{}{},
			},
		}, {
			name: "string.uri",
			args: "invalid CreateUserRequest.Uri: value must be absolute",
			want: Want{
				"CreateUserRequest.Uri.uri",
				map[string]interface{}{},
			},
		}, {
			name: "string.address",
			args: "invalid CreateUserRequest.Addr: value must be a valid hostname, or ip address",
			want: Want{
				"CreateUserRequest.Addr.address",
				map[string]interface{}{},
			},
		}, {
			name: "string.uuid",
			args: "invalid CreateUserRequest.No: value must be a valid UUID | caused by: invalid uuid format",
			want: Want{
				"CreateUserRequest.No.uuid",
				map[string]interface{}{},
			},
		},
		{
			name: "repeated.between",
			args: "invalid CreateUserRequest.Ids: value must contain between 1 and 10 items, inclusive",
			want: Want{
				"CreateUserRequest.Ids.repeated.between",
				map[string]interface{}{"min_items": "1", "max_items": "10"},
			},
		}, {
			name: "required",
			args: "invalid CreateUserRequest.Super: value is required",
			want: Want{
				"CreateUserRequest.Super.required",
				map[string]interface{}{},
			},
		}, {
			name: "enum",
			args: "invalid CreateUserRequest.UserType: value must be one of the defined enum values",
			want: Want{
				"CreateUserRequest.UserType.enum",
				map[string]interface{}{},
			},
		}, {
			name: "enum-in",
			args: "invalid CreateUserRequest.UserType: value must be in list [1 2]",
			want: Want{
				"CreateUserRequest.UserType.in",
				map[string]interface{}{"in": "[1 2]"},
			},
		}, {
			name: "嵌套验证",
			args: "invalid CreateUserRequest.Super: EMBEDDED MESSAGE FAILED VALIDATION | caused by: invalid Superman.Role: value is required",
			want: Want{
				"CreateUserRequest.Super.Superman.Role.required",
				map[string]interface{}{},
			},
		}, {
			name: "嵌套&范围",
			args: "invalid CreateUserRequest.Super: embedded message failed validation | caused by: invalid Superman.Role: embedded message failed validation | caused by: invalid Role.Id: value must be inside range [2, 10)",
			want: Want{
				"CreateUserRequest.Super.Superman.Role.Role.Id.between",
				map[string]interface{}{
					"gte": "2",
					"lt":  "10",
				},
			},
		}, {
			name: "嵌套&字符串长度",
			args: "invalid CreateUserRequest.Super: embedded message failed validation | caused by: invalid Superman.Role: embedded message failed validation | caused by: invalid Role.RoleName: value length must be at most 5 bytes",
			want: Want{
				"CreateUserRequest.Super.Superman.Role.Role.RoleName.max_bytes",
				map[string]interface{}{"max_bytes": "5"},
			},
		}, {
			name: "repeated.between",
			args: "invalid CreateUserRequest.Ids: value must contain between 1 and 10 items, inclusive",
			want: Want{
				"CreateUserRequest.Ids.repeated.between",
				map[string]interface{}{"max_items": "10", "min_items": "1"},
			},
		}, {
			name: "repeated.unique",
			args: "invalid CreateUserRequest.Ids[4]: repeated value must contain unique items",
			want: Want{
				"CreateUserRequest.Ids.repeated.unique",
				map[string]interface{}{"key": "4"},
			},
		},
	}
	var res Want
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, p := getValidateKey(tt.args)
			res = Want{
				Key:    key,
				Params: p,
			}
			if !reflect.DeepEqual(res, tt.want) {
				t.Errorf("getValidateKey() = %v, want %v", res, tt.want)
			}
		})
	}

	//repeatedBetween := "invalid CreateUserRequest.Ids: value must contain between 1 and 10 items, inclusive"
	//key,pMap := getValidateKey(repeatedBetween)
	//
	//fmt.Printf("key:%v,map:%v\n",key,pMap)

}
