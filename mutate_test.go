package main

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/slok/kubewebhook/v2/pkg/log"
	kwhmodel "github.com/slok/kubewebhook/v2/pkg/model"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_swapPodMutator(t *testing.T) {
	type args struct {
		cfg             *ImageSwapConfig
		ctx             context.Context
		admissionReview *kwhmodel.AdmissionReview
		obj             metav1.Object
		logger          log.Logger
	}
	tests := []struct {
		name    string
		args    args
		want    metav1.Object
		wantErr bool
	}{
		{
			name: "basic test",
			args: args{
				cfg:             initConfig("./tests/imageswap-test.yaml"),
				ctx:             context.TODO(),
				admissionReview: &kwhmodel.AdmissionReview{},
				obj: &corev1.Pod{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "test",
								Image: "gcr.io/library/testing:latest",
							},
						},
					},
				},
				logger: log.Noop,
			},
			want: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "test",
							Image: "registry.testing.com/library/testing:latest",
						},
					},
				},
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		fmt.Println(tt.args.cfg)

		t.Run(tt.name, func(t *testing.T) {
			got, err := swapPodMutator(tt.args.cfg, tt.args.ctx, tt.args.admissionReview, tt.args.obj, tt.args.logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("swapPodMutator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.MutatedObject, tt.want) {
				t.Errorf("swapPodMutator() = %v, want %v", got.MutatedObject, tt.want)
			}
		})
	}
}
