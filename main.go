package main

import (
	"context"
	"github.com/mkideal/cli"
	"github.com/thanhpk/randstr"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

type argT struct {
	Name       string `cli:"*name" usage:"Name of the secret to administrate"`
	Namespace  string `cli:"*namespace" usage:"Namespace where the secret is stored"`
	APIServer  string `cli:"apiserver" usage:"Kubernetes apiserver to talk to"`
	Kubeconfig string `cli:"kubeconfig" usage:"Kubeconfig to use for connection"`

	Help bool `cli:"h,help" usage:"show help"`

	Static              map[string]string `cli:"static,s" usage:"Static values to set, entries will be updated. key=value"`
	GeneratePassword    map[string]int    `cli:"password,p" usage:"Generate passwords and never change them. key=length"`
	GenerateRandomBytes map[string]int    `cli:"bytes,b" usage:"Generate random bytes and never change them, for example session cookie key. key=length"`
}

func (argv *argT) AutoHelp() bool {
	return argv.Help
}

func mutateSecret(args *argT, secret *v1.Secret) {
	for k, v := range args.Static {
		secret.Data[k] = []byte(v)
	}

	for k, v := range args.GeneratePassword {
		if _, ok := secret.Data[k]; !ok {
			secret.Data[k] = []byte(randstr.String(v))
		}
	}

	for k, v := range args.GenerateRandomBytes {
		if _, ok := secret.Data[k]; !ok {
			secret.Data[k] = []byte(randstr.Bytes(v))
		}
	}
}

func getConfig(args *argT) (*rest.Config, error) {
	if config, inClusterError := rest.InClusterConfig(); inClusterError != nil {
		if config, kubeconfigError := clientcmd.BuildConfigFromFlags(args.APIServer, args.Kubeconfig); kubeconfigError != nil {
			return nil, inClusterError
		} else {
			return config, nil
		}
	} else {
		return config, nil
	}
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		args := ctx.Argv().(*argT)
		createNew := false

		if config, err := getConfig(args); err != nil {
			log.Fatal(err)
		} else {
			if cs, err := kubernetes.NewForConfig(config); err != nil {
				log.Fatal(err)
			} else {
				secretsClient := cs.CoreV1().Secrets(args.Namespace)
				var secret *v1.Secret

				if existing, err := secretsClient.Get(context.Background(), args.Name, metav1.GetOptions{}); err != nil {
					secret = &v1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Name: args.Name,
						},
						Data: map[string][]byte{},
					}

					createNew = true
				} else {
					secret = existing
				}

				mutateSecret(args, secret)

				if createNew {
					if _, err := secretsClient.Create(context.Background(), secret, metav1.CreateOptions{}); err != nil {
						log.Fatal(err)
					}
				} else {
					if _, err := secretsClient.Update(context.Background(), secret, metav1.UpdateOptions{}); err != nil {
						log.Fatal(err)
					}
				}
			}
		}

		return nil
	})
}
