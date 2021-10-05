### Creating a stateful set using Client-go to do the following functions:

- Spin up a pod using stateful set
- Display pod list for a namespace
- Get Pod object from list
- Execute remote command and get the output inside a pod (Similar to kubectl 'exec')

**The main method resides inside controller directory - kubectlClient.go**

>Some util methods are there:

```
Method: getConfig()
Description: To get the k8 cluster config object based on the '.kube/config' file from the cluster   
args: None
return variables: (*rest.Config, error) => k8 rest config object 
```

```
Method: getClientSet()
Description: To get the k8 cluster clientSet object - Analogous to kubectl   
args: config *rest.Config => k8 config object fetched using getConfig() method
return variables: (*kubernetes.Clientset, error) => k8 clientSet object 
```

```
Method: getPodList()
Description: To fetch the Pod List for the given namespace from k8 cluster   
args: (clientSet *kubernetes.Clientset, nameSpace string) => clientSet is to use kubectl commands
return variables: (*k8s.io/api/core/v1.PodList, error) => Pod list object 
```

```
Method: getPodFromList()
Description: To fetch a Pod from the List with the given name from given Pod List   
args: (podList *k8s.io/api/core/v1.PodList, podName string) => podList to fetch the pod object for given podName 
return variables: (*k8s.io/api/core/v1.Pod, string) => Pod object, error string 
```

```
Method: executeCommandInPod()
Description: To execute a command in the Pod
args: (podList *k8s.io/api/core/v1.PodList, podName string, command []string) => podList to fetch the pod object for given podName, command arg to provide multiple commands
return variables: None 
```

```
Method: executeRemoteCommand()
Description: Execute remote command utility method
args: (pod *k8s.io/api/core/v1.Pod, command []string) => podList to fetch the pod object for given podName, command arg to provide multiple commands
return variables: None
```



**The stateFulSets package has go file busybox which helps in creation and deletion of the statefulset(and pods)**

>The methods are:

```
Method: CreateBusyBoxStateFulSet()
Description: Create a statefulset to spin up busybox pod
args: (clientSet *kubernetes.Clientset) => clientSet is to use kubectl commands
return variables: None 
```

```
Method: DeleteBusyBoxStateFulSet()
Description: Deletes statefulset to spin up busybox pod
args: (clientSet *kubernetes.Clientset) => clientSet is to use kubectl commands
return variables: None 
```