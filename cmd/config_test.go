package cmd

import (
	"reflect"
	"testing"

	"github.com/rokeller/spartan/server"
)

func Test_getConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    *server.Config
		wantErr bool
	}{
		{
			name: "Default",
			want: &server.Config{
				Server: server.ServerConfig{Port: 8080, StaticContentDir: "/content"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("getConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
