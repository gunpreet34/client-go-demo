package stateFulSets

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)


func CreateBusyBoxStateFulSet(clientSet *kubernetes.Clientset) {
	stateFulSetClient := clientSet.AppsV1().StatefulSets(apiv1.NamespaceDefault)
	// Defining values to be passed as pointers to client go api calls

	var numOfReplicas = int32(1)

	var statefulSet = &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "busybox",
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &numOfReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "busybox",
				},
			},
			ServiceName: "busybox",
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "busybox",
					},
				},
				Spec: apiv1.PodSpec{
					InitContainers: []apiv1.Container{},
					Containers: []apiv1.Container{
						{
							Name:            "busybox",
							Image:           "busybox",
							ImagePullPolicy: "IfNotPresent",
						},
					},
				},
			},
		},

	}
	// Deploy statefulset
	fmt.Println("Creating statefulset...")
	result, err := stateFulSetClient.Create(context.TODO(), statefulSet, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created statefulset %q.\n", result.GetObjectMeta().GetName())
}

func DeleteBusyBoxStateFulSet(clientSet *kubernetes.Clientset) {
	stateFulSetClient := clientSet.AppsV1().StatefulSets(apiv1.NamespaceDefault)
	err := stateFulSetClient.Delete(context.TODO(), "busybox", metav1.DeleteOptions{})
	if err != nil {
		panic(err.Error())
	}
}

