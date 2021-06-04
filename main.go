package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Deployment struct {
	DeploymentNamespace string `json:"DeploymentNamespace"`
	Replicas            int    `json:"Replicas"`
}

type Pod struct {
	PodName        string `json:"PodName"`
	ContainerName  string `json:"ContainerName"`
	ContainerImage string `json:"ContainerImage"`
}

type Namespace struct {
	NsName string `json:"NsName"`
}

var Pods []Pod

func returnAllPods(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllPods")
	json.NewEncoder(w).Encode(Pods)
}

//------------------------------------------------------------------------------------------------------------------------

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

//------------------------------------------------------------------------------------------------------------------------

func createNewdeployment(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var deploy Deployment
	json.Unmarshal(reqBody, &deploy)
	/*rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err)
	}*/
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	//clientset := kubernetes.NewForConfigOrDie(config)
	deploymentRes := schema.GroupVersionResource{Group: "kubedb.com", Version: "v1alpha1", Resource: "postgreses"}
	deployment := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "kubedb.com/v1alpha1",
			"kind":       "Postgres",
			"metadata": map[string]interface{}{
				"name":      "postgres-deployment",
				"namespace": deploy.DeploymentNamespace,
			},
			"spec": map[string]interface{}{
				"version":       "11.1",
				"replicas":      deploy.Replicas,
				"standbyMode":   "Hot",
				"streamingMode": "asynchronous",
				"leaderElection": map[string]interface{}{
					"leaseDurationSeconds": 5,
					"renewDeadlineSeconds": 3,
					"retryPeriodSeconds":   2,
				},
				"storageType": "Durable",
				"storage": map[string]interface{}{
					"storageClassName": "nas",
					"accessModes":      []string{"ReadWriteMany"},
					"resources": map[string]interface{}{
						"requests": map[string]interface{}{
							"storage": "1Gi",
						},
					},
				},
				"terminationPolicy": "WipeOut",
			},
		},
	}
	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := client.Resource(deploymentRes).Namespace(deploy.DeploymentNamespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetName())
}

//------------------------------------------------------------------------------------------------------------------------

func createNewPod(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	var pod Pod
	json.Unmarshal(reqBody, &pod)
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err)
	}
	clientset := kubernetes.NewForConfigOrDie(config)
	newPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: pod.PodName,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "nginx", Image: "nginx:latest", Command: []string{"sleep", "100000"}},
			},
		},
	}
	createdPod, err := clientset.CoreV1().Pods("pg-namespace").Create(context.Background(), newPod, metav1.CreateOptions{})
	fmt.Println(createdPod)
	Pods = []Pod{
		Pod{PodName: "Hello", ContainerName: "Pod ContainerName", ContainerImage: "Pod ContainerImage"},
	}
	fmt.Fprintf(w, "%+v", string(pod.PodName))
}

//------------------------------------------------------------------------------------------------------------------------

func createNewNamespace(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var namespace Namespace
	json.Unmarshal(reqBody, &namespace)
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	config, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err)
	}
	clientset := kubernetes.NewForConfigOrDie(config)
	if err != nil {
		panic(err.Error())
	}

	newNamespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace.NsName,
		},
	}

	createdNs, err := clientset.CoreV1().Namespaces().Create(context.Background(), newNamespace, metav1.CreateOptions{})
	fmt.Println(createdNs)
	fmt.Fprintf(w, "%+v", string(namespace.NsName))

}

//------------------------------------------------------------------------------------------------------------------------

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/pods", returnAllPods).Methods("GET")
	//myRouter.HandleFunc("/pod", createNewPod).Methods("POST")
	myRouter.HandleFunc("/namespace", createNewNamespace).Methods(("POST"))
	myRouter.HandleFunc("/postgres", createNewdeployment).Methods(("POST"))

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
func main() {

	/*	rules := clientcmd.NewDefaultClientConfigLoadingRules()
		kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
		config, err := kubeconfig.ClientConfig()
		if err != nil {
			panic(err)
		}
		clientset := kubernetes.NewForConfigOrDie(config)

		nodeList, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		for _, n := range nodeList.Items {
			fmt.Println(n.Name)
		}
		if err != nil {
			panic(err)
		}
	*/handleRequests()
}
