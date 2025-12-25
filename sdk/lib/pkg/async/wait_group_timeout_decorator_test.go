package task

//
//import (
//	"sync"
//	"testing"
//	"time"
//)
//
//func TestTimeoutWaitGroupDecorator_WaitOrTimeout(t *testing.T) {
//	type fields struct {
//		WaitGroup *sync.WaitGroup
//	}
//	type args struct {
//		timeout time.Duration
//	}
//	tests := []struct {
//		name       string
//		fields     fields
//		args       args
//		preAction  func(wd *sync.WaitGroup)
//		postAction func(wd *sync.WaitGroup) // runs in a go routine
//		want       bool
//	}{
//		{
//			name: "If no Add() was invoked, then ends with success returning false",
//			args: args{timeout: time.Second},
//			want: false,
//		},
//		{
//			name: "If Add() was invoked, then ends with timeout returning true",
//			args: args{timeout: time.Second},
//			preAction: func(wd *sync.WaitGroup) {
//				wd.Add(1)
//			},
//			want: true,
//		},
//		{
//			name: "If Add() was invoked, but Done() is called before timeout, then ends with success returning false",
//			args: args{timeout: time.Second},
//			preAction: func(wd *sync.WaitGroup) {
//				wd.Add(1)
//			},
//			postAction: func(wd *sync.WaitGroup) {
//				time.Sleep(500 * time.Millisecond)
//				wd.Done()
//			},
//			want: false,
//		},
//		{
//			name: "If Add() was invoked, but Done() is called after timeout, then ends with success returning false",
//			args: args{timeout: time.Second},
//			preAction: func(wd *sync.WaitGroup) {
//				wd.Add(1)
//			},
//			postAction: func(wd *sync.WaitGroup) {
//				time.Sleep(1500 * time.Millisecond)
//				wd.Done()
//			},
//			want: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			w := &TimeoutWaitGroupDecorator{}
//			if tt.preAction != nil {
//				tt.preAction(&w.WaitGroup)
//			}
//
//			if tt.postAction != nil {
//				go tt.postAction(&w.WaitGroup)
//			}
//
//			if got := w.WaitOrTimeout(tt.args.timeout); got != tt.want {
//				t.Errorf("WaitOrTimeout() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
