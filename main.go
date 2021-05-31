package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Pod struct {
	PodName        string `json:"PodName"`
	ContainerName  string `json:"ContainerName"`
	ContainerImage string `json:"ContainerImage"`
}

var Pods []Pod

func returnAllPods(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllPods")
	json.NewEncoder(w).Encode(Pods)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func createNewPod(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	var pod Pod
	json.Unmarshal(reqBody, &pod)

	//esm_el_pod=reqbody.podname
	//esmelcontainer=reqbody.containername
	//llllllllllllll kkkkk
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

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/pods", returnAllPods).Methods("GET")
	myRouter.HandleFunc("/pod", createNewPod).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {

	/*
		nodeList, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		for _, n := range nodeList.Items {
			fmt.Println(n.Name)
		}

		if err != nil {
			panic(err)
		}*/

	handleRequests()
}
