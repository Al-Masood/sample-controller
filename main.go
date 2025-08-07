package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	klient "github.com/al-masood/sample-controller/pkg/generated/clientset/versioned"
)

func main() {
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "")
	}

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		log.Printf("Building config from flag, %s", err.Error())
	}

	klientset, err := klient.NewForConfig(config)
	if err != nil {
		log.Printf(err.Error())
	}

	fmt.Println(klientset)

	foos, err := klientset.SampleV1alpha1().Foos("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf(err.Error())
	}

	fmt.Printf("length of klusters is %d\n", len(foos.Items))
}
