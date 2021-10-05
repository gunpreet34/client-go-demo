package main

import (
	"bytes"
	"context"
	"firstClientGoProject/stateFulSets"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"time"
)

var config *rest.Config
var clientSet *kubernetes.Clientset

func main() {
	var err error
	config, err = getConfig()
	// Get clientSet (to use kubectl commands using GoLang)
	clientSet, err = getClientSet(config)
	if err != nil {
		panic(err.Error())
	}

	stateFulSets.CreateBusyBoxStateFulSet(clientSet)
	time.Sleep(time.Second * 10)
	podList, err := getPodList(clientSet, apiv1.NamespaceDefault)
	if err!=nil {
		fmt.Printf("couldn't find pod: 'busybox' in namespace: %v, err: %v", apiv1.NamespaceDefault, err.Error())
	} else {
		executeCommandInPod(podList, "busybox-0", []string{"/bin/sh", "-c", "sleep 240" , "ls -la"})

	}
	stateFulSets.DeleteBusyBoxStateFulSet(clientSet)

}

func getConfig() (*rest.Config, error){
	var kubeConfig *string

	// Parsing kubeConfig
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeconfig", filepath.Join(home , ".kube", "config-file-name"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeConfig = flag.String("kubeconfig", "", "/path-to-your-config-file/config-file-name")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func getClientSet(config *rest.Config) (*kubernetes.Clientset, error){

	// create the client set
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientSet, err
}

func getPodList(clientSet *kubernetes.Clientset, nameSpace string) (*apiv1.PodList, error) {

	pods, err := clientSet.CoreV1().Pods(nameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return pods, nil
}

func getPodFromList(pods *apiv1.PodList, podName string) (apiv1.Pod, string) {
	for _, pod := range pods.Items {
		if pod.GetName() == podName {
			return pod, ""
		}
	}
	return apiv1.Pod{}, "No Pod found"
}

func executeCommandInPod(podList *apiv1.PodList, podName string, command []string) {
	pod, message := getPodFromList(podList, podName)
	// Error handling
	if message != "" {
		fmt.Println(message)
		return
	}

	result, resErr, Err := executeRemoteCommand(&pod, command)

	// Error handling
	if Err != nil {
		fmt.Printf("error while executing command in pod: %v", Err.Error())
	} else if resErr != "" {
		fmt.Println(resErr)
	} else {
		fmt.Println(result)
	}
}

func executeRemoteCommand(pod *apiv1.Pod, command []string) (string, string, error) {


	request := clientSet.RESTClient().
		Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec")

	scheme := runtime.NewScheme()
	if err:= apiv1.AddToScheme(scheme); err != nil {
		return "", "", err
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	request.VersionedParams(&apiv1.PodExecOptions{
		Command: command,
		Stdin:   false,
		Stdout:  true,
		Stderr:  true,
		TTY:     false,
	}, parameterCodec)


	exec, err := remotecommand.NewSPDYExecutor(config, "POST", request.URL())
	if err != nil {
		return "", "", err
	}

	var buf, errBuf bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &buf,
		Stderr: &errBuf,
		Tty: false,
	})

	if err != nil {
		return "", "", errors.Wrapf(err, "Failed executing command %s on %v/%v", command, pod.Namespace, pod.Name)
	}

	return buf.String(), errBuf.String(), nil
}