package main

import (
	"net/http/httptest"
	"reflect"
	"testing"

	utiltesting "k8s.io/client-go/util/testing"
)

var (
	respone = `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"OpenEBS__iops","instance":"172.17.0.2:9500","job":"cluster_uuid_9aba2480-a180-41ca-b5cb-f4a099376a16_openebs-volumes","kubernetes_pod_name":"pvc-4fa13b09-6242-11e8-a310-1458d00e6b83-ctrl-745784bb48-z9pl8","openebs_pv":"pvc-4fa13b09-6242-11e8-a310-1458d00e6b83"},"value":[1528354477.902,"0"]}]}}`
)

func TestGetValues(t *testing.T) {
	cases := map[string]struct {
		fakeHandler utiltesting.FakeHandler
		queryType   string
		channel     string
	}{
		"When getting data for OpenEBS_read_iops:": {
			fakeHandler: utiltesting.FakeHandler{
				StatusCode:   200,
				ResponseBody: string(respone),
				T:            t,
			},
			queryType: "OpenEBS_read_iops",
			channel:   "read",
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(&tt.fakeHandler)
			go getValues(server.URL, tt.queryType)
			if tt.channel == "read" {
				read := <-readch
				t.Logf("%v", read)
				//t.Errorf("%#v", read)
			} else if tt.channel == "write" {
				write := <-writech
				t.Logf("%v", write)
				//t.Errorf("%#v", write)
			}

		})
	}
}

func Test_getPVs(t *testing.T) {
	tests := []struct {
		name string
		want map[string]pvdata
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPVs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPVs() = %v, want %v", got, tt.want)
			}
		})
	}
}
