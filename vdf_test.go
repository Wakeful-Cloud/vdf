package vdf

import (
	"os"
	"reflect"
	"testing"
)

func TestReadVdf(t *testing.T) {
	type args struct {
		data []byte
	}

	//Read
	bytes, err := os.ReadFile("./example/read-test.vdf")

	if err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		args    args
		want    Map
		wantErr bool
	}{
		{
			name: "read example",
			args: args{
				data: bytes,
			},
			want: Map{
				"key1": "value1",
				"key2": uint32(3),
				"key3": Map{
					"key4": "value2",
					"key5": "value3",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadVdf(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadVdf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadVdf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nextMap(t *testing.T) {
	type args struct {
		buffer *buffer
	}

	//Read
	bytes, err := os.ReadFile("./example/map-test.vdf")

	if err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		args    args
		want    Map
		wantErr bool
	}{
		{
			name: "get next map",
			args: args{
				buffer: &buffer{
					Data:     bytes,
					Position: 0,
				},
			},
			want: Map{
				"key1": "value1",
				"key2": uint32(3),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nextMap(tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("nextMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nextMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nextMapItem(t *testing.T) {
	type args struct {
		buffer *buffer
	}

	//Read
	mapBytes, err := os.ReadFile("./example/map-only-test.vdf")

	if err != nil {
		panic(err)
	}

	stringBytes, err := os.ReadFile("./example/string-test.vdf")

	if err != nil {
		panic(err)
	}

	numberBytes, err := os.ReadFile("./example/number-test.vdf")

	if err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		args    args
		want    mapItem
		wantErr bool
	}{
		{
			name: "get next map item",
			args: args{
				buffer: &buffer{
					Data:     mapBytes,
					Position: 0,
				},
			},
			want: mapItem{
				Type: vdfMapStart,
				Name: "key1",
				Value: Map{
					"key2": "value1",
					"key3": "value2",
				},
			},
		},
		{
			name: "get next string item",
			args: args{
				buffer: &buffer{
					Data:     stringBytes,
					Position: 0,
				},
			},
			want: mapItem{
				Type:  vdfString,
				Name:  "key",
				Value: "value",
			},
		},
		{
			name: "get next number item",
			args: args{
				buffer: &buffer{
					Data:     numberBytes,
					Position: 0,
				},
			},
			want: mapItem{
				Type:  vdfNumber,
				Name:  "key",
				Value: uint32(3),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nextMapItem(tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("nextMapItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nextMapItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nextStringZero(t *testing.T) {
	type args struct {
		buffer *buffer
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "get next string zero for a normal buffer",
			args: args{
				buffer: &buffer{
					Data:     []byte("Hello,\x00 world!"),
					Position: 0,
				},
			},
			want: "Hello,",
		},
		{
			name: "get next string zero for an overflowed buffer",
			args: args{
				buffer: &buffer{
					Data:     []byte("Hello, world!"),
					Position: 13,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nextStringZero(tt.args.buffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("nextStringZero() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("nextStringZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteVdf(t *testing.T) {
	type args struct {
		vdfMap Map
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "write example",
			args: args{
				vdfMap: Map{
					"key1": "value1",
					"key2": uint32(3),
					"key3": Map{
						"key4": "value2",
						"key5": "value3",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := WriteVdf(tt.args.vdfMap)
			if err != nil {
				t.Errorf("WriteVdf() error = %v", err)
				return
			}
			got, err := ReadVdf(bytes)
			if err != nil {
				t.Errorf("WriteVdf() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.args.vdfMap) {
				t.Errorf("WriteVdf() = %v", got)
			}
		})
	}
}

func Test_addMap(t *testing.T) {
	type args struct {
		vdfMap Map
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "add a map",
			args: args{
				vdfMap: Map{
					"key1": "value1",
					"key2": uint32(3),
					"key3": Map{
						"key4": "value2",
						"key5": "value3",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := addMap(tt.args.vdfMap)
			if err != nil {
				t.Errorf("addMap() error = %v", err)
				return
			}
			got, err := nextMap(&buffer{
				Data:     bytes,
				Position: 0,
			})
			if err != nil {
				t.Errorf("addMap() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.args.vdfMap) {
				t.Errorf("addMap() = %v", got)
			}
		})
	}
}

func Test_addString(t *testing.T) {
	type args struct {
		value string
	}

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "add a string",
			args: args{
				value: "Hello, world!",
			},
			want: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := addString(tt.args.value)
			if err != nil {
				t.Errorf("addString() error = %v", err)
				return
			}
			got, err := nextStringZero(&buffer{
				Data:     bytes,
				Position: 0,
			})
			if err != nil {
				t.Errorf("addString() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.args.value) {
				t.Errorf("addString() = %v", got)
			}
		})
	}
}

func Test_addKT(t *testing.T) {
	type args struct {
		Type byte
		key  string
	}

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "add a string key-type pair",
			args: args{
				Type: vdfString,
				key:  "key1",
			},
			want: []byte{1, 107, 101, 121, 49, 0},
		},
		{
			name: "add a number key-type pair",
			args: args{
				Type: vdfNumber,
				key:  "key2",
			},
			want: []byte{2, 107, 101, 121, 50, 0},
		},
		{
			name: "add a map key-type pair",
			args: args{
				Type: vdfMapStart,
				key:  "key3",
			},
			want: []byte{0, 107, 101, 121, 51, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := addKT(tt.args.Type, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("addKT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addKT() = %v, want %v", got, tt.want)
			}
		})
	}
}
