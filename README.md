# k8s-secret-admin

This program tries to fix a problem when deploying new software to kubernetes clusters: generating secure
passwords and other keys automatically, while ensuring they aren't changed and co-existing with static
values in the same `Secret` that can be updated.

Originally built to be used in a one-time job after a Helm update, it may also be used as an init container
or something alike.


# Usage

Normally you want to use this inside a pod with a service account allowed to read/create/update secrets.
In-Cluster authentication is used automatically, if it fails, it looks at the `apiserver` and `kubeconfig`
flags and tries to get a connection to the cluster with these.

You have to give the `name` and `namespace` of the secret you want to administrate and probably want to give
some actions to do.


## Actions

You can use as much actions as you want, but you will catch some undefined behaviour if you mutate the same
key with different actions.


### static

Just sets a given key to a given value. If the key already exists in the secret, the value is updated since
the admin knows best when values have to change.

`k8s-secret-admin --name secret --namespace kube-public --static k8s=great`


### password

Will never change the existing value in a secret but only add the key. Generates a new secure password with
the given length.

`k8s-secret-admin --name secret --namespace kube-public --password root_password=32`


### bytes

Will never change the existing valie in a secret but only add the key. Generates a new secure random byte
sequence with the given length. This was built to generate a AES key for Yesods session cookie, but may come
handy in a lot of use cases.

`k8s-secret-admin --name secret --namespace kube-public --bytes session_key=128`


# Maturity

Hacked together mostly in a train and finished after arriving home. The available features are simple enough
to probably work most of the time but the interface isn't stable yet. Use exactly pinned versions.


# See also

If you have to generate a lot of random strings and stuff for your secrets, you may want an operator for that:  
https://github.com/mittwald/kubernetes-secret-generator


# License

MIT license, so no problem to use commercially. But please give changes upstream and support minorities in IT.




Trans rights! Black lives matter!
