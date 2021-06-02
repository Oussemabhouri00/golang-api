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
	"k8s.io/client-go/rest"
)

type Pod struct {
	PodName        string `json:"PodName"`
	ContainerName  string `json:"ContainerName"`
	ContainerImage string `json:"ContainerImage"`
	PodNamespace   string `json:"Namespace"`
}

type Namespace struct {
	NsName string `json:"NsName"`
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
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	newPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: pod.PodName,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: pod.ContainerName, Image: pod.ContainerImage, Command: []string{"sleep", "100000"}},
			},
		},
	}

	createdPod, err := clientset.CoreV1().Pods(pod.PodNamespace).Create(context.Background(), newPod, metav1.CreateOptions{})

	fmt.Println(createdPod)
	fmt.Fprintf(w, "%+v", string(pod.PodName))
}

func createNewNamespace(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var namespace Pod
	json.Unmarshal(reqBody, &namespace)
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	newNamespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "My-namespacee",
		},
	}

	createdNs, err := clientset.CoreV1().Namespaces().Create(context.Background(), newNamespace, metav1.CreateOptions{})
	fmt.Println(createdNs)

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
	//myRouter.HandleFunc("/pod", createNewPod).Methods("POST")
	myRouter.HandleFunc("/namespace", createNewNamespace).Methods(("POST"))
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
	Pods = []Pod{
		Pod{PodName: "Hello", ContainerName: "Pod ContainerName", ContainerImage: "Pod ContainerImage"},
	}

	handleRequests()
}
