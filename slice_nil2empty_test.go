package goconvrt

import (
	"reflect"
	"testing"
)

type SimpleStruct struct {
	A int
}

type SimpleEmbedStruct struct {
	List []int
}

type SimpleEmbedStruct2 struct {
	List []int
	Num  int
}

type EmbedStructStruct struct {
	Struct SimpleEmbedStruct
	Slice  []SimpleEmbedStruct
	Map    map[string]SimpleEmbedStruct
}

type ContainInt struct {
	Num *int
}

type SliceContainInt struct {
	List []*int
}

func TestConvertNilSlice2Empty(t *testing.T) {
	tests := []struct {
		name  string
		args  interface{}
		want  interface{}
		check func(arg interface{}) bool
	}{
		{
			name: "非结构体-int",
			args: 1,
			want: 1,
		},
		{
			name: "非结构体-string",
			args: "abc",
			want: "abc",
		},
		{
			name: "非结构体-bool",
			args: true,
			want: true,
		},
		{
			name: "非结构体-nil",
			args: nil,
			want: nil,
		},
		{
			name: "slice: not empty",
			args: []int{1, 2, 3},
			want: []int{1, 2, 3},
		},
		{
			name: "slice: empty",
			args: []int{},
			want: []int{},
		},
		{
			name: "slice: nil",
			args: []int(nil),
			want: []int{},
			check: func(got interface{}) bool {
				return len(got.([]int)) == 0 && got != nil
			},
		},
		{
			name: "slice: elem is object and nil",
			args: []SimpleStruct(nil),
			want: []SimpleStruct{},
			check: func(got interface{}) bool {
				return len(got.([]SimpleStruct)) == 0 && got != nil
			},
		},
		{
			name: "slice: elem is object and not nil",
			args: []SimpleStruct{{1}},
			want: []SimpleStruct{{1}},
		},
		{
			name: "slice: elem is object embed slice and not nil",
			args: []SimpleEmbedStruct{{List: []int{1}}},
			want: []SimpleEmbedStruct{{List: []int{1}}},
		},
		{
			name: "slice: elem is object embed slice and nil",
			args: []SimpleEmbedStruct{{List: nil}},
			want: []SimpleEmbedStruct{{List: []int{}}},
			check: func(got interface{}) bool {
				return len(got.([]SimpleEmbedStruct)) == 1 && got.([]SimpleEmbedStruct)[0].List != nil && len(got.([]SimpleEmbedStruct)[0].List) == 0
			},
		},
		{
			name: "slice: elem is object embed slice and not nil",
			args: []SimpleEmbedStruct{{List: nil}, {List: []int{1, 2, 3}}},
			want: []SimpleEmbedStruct{{List: []int{}}, {List: []int{1, 2, 3}}},
			check: func(got interface{}) bool {
				data := got.([]SimpleEmbedStruct)
				return len(data) == 2 && data[0].List != nil && len(data[0].List) == 0 && data[1].List[0] == 1 &&
					data[1].List[1] == 2 && data[1].List[2] == 3

			},
		},
		{
			name: "struct: elem is slice and not nil",
			args: SimpleEmbedStruct{List: []int{1, 2, 3}},
			want: SimpleEmbedStruct{List: []int{1, 2, 3}},
		},
		{
			name: "struct: elem is slice and nil",
			args: SimpleEmbedStruct{List: nil},
			check: func(got interface{}) bool {
				data := got.(SimpleEmbedStruct)
				return data.List != nil && len(data.List) == 0
			},
		},
		{
			name: "struct:  slice is nil and other not",
			args: SimpleEmbedStruct2{List: nil, Num: 1},
			check: func(got interface{}) bool {
				data := got.(SimpleEmbedStruct2)
				return data.List != nil && len(data.List) == 0 && data.Num == 1
			},
		},
		{
			name: "map: simple",
			args: map[string]interface{}{"a": 1},
			check: func(got interface{}) bool {
				data := got.(map[string]interface{})
				return data["a"] == 1
			},
		},
		{
			name: "map: contain nil slice",
			args: map[string]interface{}{"a": []SimpleStruct(nil)},
			check: func(got interface{}) bool {
				data := got.(map[string]interface{})
				return data["a"].([]SimpleStruct) != nil && len(data["a"].([]SimpleStruct)) == 0
			},
		},
		{
			name: "pointer: contain nil slice",
			args: &EmbedStructStruct{},
			check: func(got interface{}) bool {
				data := got.(*EmbedStructStruct)
				return data.Slice != nil && len(data.Slice) == 0
			},
		},
		{
			name: "nil int not nil",
			args: &one,
			check: func(got interface{}) bool {
				return *(got.(*int)) == 1
			},
		},
		{
			name: "nil int",
			args: nilOne,
			check: func(got interface{}) bool {
				return got.(*int) == nil
			},
		},
		{
			name: "struct contain int ptr",
			args: ContainInt{&one},
			check: func(got interface{}) bool {
				return *got.(ContainInt).Num == 1
			},
		},
		{
			name: "struct contain nil int ptr",
			args: ContainInt{nil},
			check: func(got interface{}) bool {
				return got.(ContainInt).Num == nil
			},
		},
		{
			name: "slice contain nil int ptr",
			args: SliceContainInt{List: []*int{nilOne}},
			check: func(got interface{}) bool {
				return got.(SliceContainInt).List[0] == nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertNilSlice2Empty(tt.args)
			equal := false
			if tt.check != nil {
				equal = tt.check(got)
			} else {
				equal = reflect.DeepEqual(got, tt.want)
			}
			if !equal {
				t.Errorf("ConvertNilSlice2Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewZeroStruct(t *testing.T) {
	type args struct {
		typ reflect.Type
	}
	tests := []struct {
		name  string
		args  args
		want  reflect.Value
		check func(got reflect.Value) bool
	}{
		{
			name: "nil struct",
			args: args{typ: reflect.TypeOf(SimpleStruct{})},
			want: reflect.ValueOf(SimpleStruct{}),
		},
		{
			name: "embed struct",
			args: args{typ: reflect.TypeOf(SimpleEmbedStruct{})},
			want: reflect.ValueOf(SimpleEmbedStruct{}),
			check: func(got reflect.Value) bool {
				list := got.Interface().(SimpleEmbedStruct).List
				return list != nil && len(list) == 0
			},
		},
		{
			name: "embed struct struct",
			args: args{typ: reflect.TypeOf(EmbedStructStruct{})},
			want: reflect.ValueOf(EmbedStructStruct{}),
			check: func(got reflect.Value) bool {
				data := got.Interface().(EmbedStructStruct)
				if data.Struct.List == nil || len(data.Struct.List) != 0 {
					return false
				}
				if data.Slice == nil || len(data.Slice) != 0 {
					return false
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newZeroStruct(tt.args.typ)
			pass := false
			if tt.check != nil {
				pass = tt.check(got)
			} else {
				pass = reflect.DeepEqual(got.Interface(), tt.want.Interface())
			}
			if !pass {
				t.Errorf("NewZeroStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}

var one = 1
var nilOne = (*int)(nil)
