package main

import (
	
    "fmt"
    "log"
	"os"
    "net/http"
	"github.com/gorilla/mux"
	"gopkg.in/natefinch/lumberjack.v2"
    "context"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	log.Println("Starting Server ....")
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/hello", hello).Methods("GET")
	r.HandleFunc("/pods", pods).Methods("GET")
	http.Handle("/", r)

	// Configure Logging
	LOG_FILE_LOCATION := os.Getenv("LOG_FILE_LOCATION")
	if LOG_FILE_LOCATION != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   LOG_FILE_LOCATION,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		})
	}
	log.Println("Web Server is Starting on 7070 Port ....")
	log.Fatal(http.ListenAndServe(":7070", nil))
}

func pods(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	namespace := query.Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	log.Printf("Received request for pods in %s\n", namespace)
	
	fmt.Println("################## List of POD #############")
	w.Write([]byte(fmt.Sprintf("################## List of POD #############")))
	
	contextName := ""
	if len(os.Args) >= 2 {
		contextName = os.Args[1]
	}

	client, err := newClient(contextName)
	if err != nil {
		log.Fatal(err)
	}

	pods, err := client.CoreV1().Pods(namespace).List(context.Background(),meta_v1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, pod := range pods.Items {
		fmt.Println("Pod: "+ pod.Name)
		w.Write([]byte(fmt.Sprintf("Pod: , %s\n", pod.Name)))

	}
	fmt.Println("########################################")
	log.Printf("")

}

func hello(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "Guest"
	}
	log.Printf("Received request for %s\n", name)
	w.Write([]byte(fmt.Sprintf("Hello, %s\n", name)))
}

func newClient(contextName string) (kubernetes.Interface, error) {
	configOverrides := &clientcmd.ConfigOverrides{CurrentContext: contextName}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides).ClientConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page !")
}
