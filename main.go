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
	"math/rand"
  	"time"
	"strconv"
	"strings"
	pretty "github.com/inancgumus/prettyslice"
)

func main() {
	log.Println("Starting Server ....")
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/hello", hello).Methods("GET")
	r.HandleFunc("/pods", pods).Methods("GET")
	r.HandleFunc("/config", config_handeler).Methods("GET")
	r.HandleFunc("/kill", kill_handeler).Methods("GET")
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

	// configration 
	namespace := "default"
	interval := "40"
	
	os.Setenv("chaos_namespace", namespace)
	os.Setenv("chaos_interval", interval)

	// loop killer
	go loopkiller()

	log.Println("Web Server is Starting on 7070 Port ....")
	log.Fatal(http.ListenAndServe(":7070", nil))
}

func pods(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	namespace := query.Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	log.Printf("Received a request for pods in %s namespace \n", namespace)
	
	fmt.Println("################## List of POD #############")
	w.Write([]byte(fmt.Sprintf("\n################## List of POD #############\n")))

	pods_arr := getPods(namespace)
	for _, pod := range pods_arr {
		w.Write([]byte(fmt.Sprintf("Pod: %s\n", pod)))
	}
	w.Write([]byte(fmt.Sprintf(" ########################################### \n")))
	pretty.Show("Pods :", pods_arr)

	fmt.Println("\n #############################################\n")
	log.Printf("")

}

func getPods(namespace string) []string{
	client, err := newClient("")
	if err != nil {
	log.Fatal(err)
	}

	pods, err := client.CoreV1().Pods(namespace).List(context.Background(),meta_v1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	pods_arr := []string{}

	for _, pod := range pods.Items {
		pods_arr=append(pods_arr,pod.Name);
	}
	return pods_arr;
}

func config_handeler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received a request for config_handeler")
	
	query := r.URL.Query()
	namespace := query.Get("namespace")
	if namespace == "" {
		namespace = "default"
	}else
	{
		fmt.Fprintf(w, "Make sure that you deploy the clusterrolebinding service-admin to %s namespace\n", namespace)
	}
	os.Setenv("chaos_namespace", namespace)

	interval := query.Get("interval")
	if interval == "" {
		interval = "40"
	}
	os.Setenv("chaos_interval", interval)
	fmt.Fprintf(w, "namespace = %s\n", namespace)
	fmt.Fprintf(w, "interval = %s\n", interval)
	
	log.Printf("namespace = %s\n", namespace)
	log.Printf("interval = %s\n", interval)	
}

func get_random(min int,max int) int{
	rand.Seed(time.Now().UnixNano())
	n := min + rand.Intn(max-min+1)
	return n;
}

func kill_handeler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request for kill_handeler")
	
	namespace := os.Getenv("chaos_namespace")
	pods_arr := getPods(namespace)
	
	log.Printf("Pods before killing")
	pretty.Show("Pods :", pods_arr)
	
	if(len(pods_arr)==0){
		log.Printf("No pods to Kill")
	}else{
		random := get_random(0,len(pods_arr))
		random_pod := pods_arr[random]
		
		log.Printf("Picking a random Pod to kill \n")
		excute_kill(namespace,random_pod)
		time.Sleep(15 * time.Second)

		pods_arr = getPods(namespace)
		log.Printf("Pods after 15 Secoonds ")
		pretty.Show("Pods :", pods_arr)
	}
}

func excute_kill(namespace string,podname string){
	log.Printf("!!!! Chaos Monkey is killing  %s  !!!!! \n", podname)
	
	client, err := newClient("")
	if err != nil {
		log.Fatal(err)
	}
	err2 := client.CoreV1().Pods(namespace).Delete(context.Background(),podname, meta_v1.DeleteOptions{})
	if err2 != nil {
		log.Fatal(err2)
	}
}

func loopkiller(){
	log.Printf("Start Loop killer \n")

	for {
		log.Printf("####################### Chaos Monkey is Playing ####################### \n")
		random_pod := ""
		namespace := ""
		interval := ""
		wait := 0
		for{
			namespace = os.Getenv("chaos_namespace")
			pods_arr := getPods(namespace)
			log.Printf("Pods in %s Before Killing  \n",namespace)
			if len(pods_arr)>2 {
				pretty.Show("Pods :", pods_arr)
			}
			interval = os.Getenv("chaos_interval")
			i, err := strconv.Atoi(interval)
			wait=i;
			if err != nil {
				log.Fatal(err)
			}

			if len(pods_arr)==0 {
				log.Printf("No pods to Kill")
			}else{
				log.Printf("Choosing a Random pod to kill")
				random := get_random(0,(len(pods_arr)-1))
				random_pod = pods_arr[random]
				log.Printf("Random Pod %s ",random_pod)

				if len(pods_arr)<2 {
					log.Printf("no Pods to kill .....\n")
				}
				// to avoid killing him self
				myself := strings.Contains(random_pod, "chaos-monkey")
				if myself {
					log.Printf("I will not kill myself ! ")
					log.Printf("Deploy some test pods to kill !")
					break
				}else{
					go excute_kill(namespace,random_pod)
				}
			}
			log.Printf("Sleeping for %s seconds \n", interval)
			time.Sleep(time.Duration(wait) * time.Second)
		}
		log.Printf("Sleeping for %s seconds \n", interval)
		time.Sleep(time.Duration(wait) * time.Second)
	}	
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
	fmt.Fprintf(w, "Home Page ! \n")
	fmt.Fprintf(w, "Try /pods  : to see pods \n")
	fmt.Fprintf(w, "Try /kill  : to kill a pod \n")
	fmt.Fprintf(w, "Try /config  : to config the chaos Monkey \n")
}