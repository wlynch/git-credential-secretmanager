package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	version = flag.String("version", "", "Secret Manager version (e.g. projects/my-project/secrets/my-secret/versions/latest)")
)

func main() {
	flag.Parse()

	v := *version
	if v == "" {
		v = os.Getenv("GIT_K8S_SECRET")
	}
	if v == "" {
		fmt.Fprintf(os.Stderr, "%s: cannot determine Secret Manager version, --version or ${GIT_SECRET_MANAGER_VERSION} not specified\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	if os.Args[len(os.Args)-1] != "get" {
		return
	}

	fmt.Fprintln(os.Stderr, "fetching", v)

	/*
		config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
		if err != nil {
			panic(err.Error())
		}
	*/
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// Read in Git credential config - https://git-scm.com/docs/git-credential#IOFMT
	cred, err := read(os.Stdin)
	if err != nil {
		log.Fatal("error reading stdin: %v", err)
	}

	ctx := context.Background()
	result, err := clientset.CoreV1().Secrets("default").Get(ctx, v, v1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Write secret back out to Git credential.
	cred.username = string(result.Data["username"])
	cred.password = string(result.Data["password"])
	cred.write(os.Stdout)
}

type credential struct {
	protocol string
	host     string
	path     string
	username string
	password string
	url      string
}

func read(r io.Reader) (credential, error) {
	var c credential
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s := strings.SplitN(scanner.Text(), "=", 2)
		switch s[0] {
		case "protocol":
			c.protocol = s[1]
		case "host":
			c.host = s[1]
		case "path":
			c.path = s[1]
		case "username":
			c.username = s[1]
		case "password":
			c.password = s[1]
		case "url":
			c.url = s[1]
		}
	}
	return c, scanner.Err()
}

func (c credential) write(w io.Writer) {
	printIfSet(w, "protocol", c.protocol)
	printIfSet(w, "host", c.host)
	printIfSet(w, "path", c.path)
	printIfSet(w, "username", c.username)
	printIfSet(w, "password", c.password)
	printIfSet(w, "url", c.url)
}

func printIfSet(w io.Writer, k, v string) {
	if v != "" {
		fmt.Fprintf(w, "%s=%s\n", k, v)
	}
}
