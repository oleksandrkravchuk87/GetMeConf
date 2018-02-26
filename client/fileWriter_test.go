package main

import (
	"testing"
)



func TestWriteFile(t *testing.T) {

	 data := []byte{345, 345,467}

	type args struct {
		data     []byte
		fileName string
		outPath  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		"asfsdf", {}, false
	}

	var tst1 tests
	{
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteFile(tt.args.data, tt.args.fileName, tt.args.outPath); (err != nil) != tt.wantErr {
				t.Errorf("WriteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
