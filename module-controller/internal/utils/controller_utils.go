package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/sofastack/sofa-serverless/internal/constants/label"
)

func AddFinalizer(meta *metav1.ObjectMeta, finalizer string) bool {
	if meta.Finalizers == nil {
		meta.Finalizers = []string{}
	}
	for _, s := range meta.Finalizers {
		if s == finalizer {
			return false
		}
	}
	meta.Finalizers = append(meta.Finalizers, finalizer)
	return true
}

func RemoveFinalizer(meta *metav1.ObjectMeta, needle string) bool {
	finalizers := make([]string, 0)
	found := false
	if meta.Finalizers != nil {
		for _, finalizer := range meta.Finalizers {
			if finalizer != needle {
				finalizers = append(finalizers, finalizer)
			} else {
				found = true
			}
		}
		meta.Finalizers = finalizers
	}

	return found
}

func HasFinalizer(meta *metav1.ObjectMeta, needle string) bool {
	if meta.Finalizers != nil {
		for _, finalizer := range meta.Finalizers {
			if finalizer == needle {
				return true
			}
		}
	}

	return false
}

func Key(req ctrl.Request) string {
	return fmt.Sprintf("%s/%s", req.Namespace, req.Name)
}

func GetNextReconcileTime(currentTime time.Time) time.Duration {
	timeDuration := time.Now().Sub(currentTime)
	var nextDuration time.Duration
	if timeDuration.Hours() > 1 {
		nextDuration = time.Minute * 10
	} else if timeDuration.Minutes() > 30 {
		nextDuration = time.Minute * 5
	} else if timeDuration.Minutes() > 10 {
		nextDuration = time.Minute * 1
	} else {
		nextDuration = time.Second * 10
	}
	return nextDuration
}

func GetModuleCountFromPod(pod *corev1.Pod) (count int) {
	for k, _ := range pod.Labels {
		if strings.HasPrefix(k, label.ModuleNameLabel) {
			count += 1
		}
	}
	return count
}

func Error(err error, msg string, keysAndValues ...interface{}) error {
	log.Log.Error(err, msg, keysAndValues...)
	return err
}

func GetModuleInstanceCount(pod corev1.Pod) int {
	if pod.Labels[label.ModuleInstanceCount] == "" {
		return 0
	}
	count, err := strconv.Atoi(pod.Labels[label.ModuleInstanceCount])
	if err != nil {
		return 0
	}
	return count
}
