# Getting started with the OpenShift template

The OpenShift template provides a quick start to deploy a production-ready instance of Jaeger on OpenShift.
The template includes a `StatefulSet` for Cassandra, a `DaemonSet` for the Agent, as well as regular services
for the individual components such as `Collector`, `Query` and `UI`. Each of those can be scaled individually.

To get started, you'll need to grant your user permissions to create a `DaemonSet`:

```bash
oc login -u system:admin
oc create -f daemonset-admin.yml
oadm policy add-role-to-user daemonset-admin developer
```

Once that is ready, it's only a matter of creating the components from the template:
```bash
oc process -f openshift.yml | oc create -f -
```

The `openshift.yml` file itself is an "OpenShift template", meaning that it can be installed on a regular
OpenShift instance and used later via the UI to create the "instant apps". If the template is installed at
the `openshift` namespace, all users will be able to create a Jaeger instance within their namespaces.

## Getting an OpenShift cluster running

As a developer looking to try this out locally, the easiest is to use the `oc cluster up` command. Getting
this command might be as easy as running `dnf install origin-clients` on a recent Fedora desktop. Refer to
the OpenShift [installation guide or quick start guide](https://install.openshift.com/) for more information. 
Another alternative is to use [`minishift`](https://github.com/minishift/minishift).

## Where to go from here

This template is meant to serve as a starting point for your own installation, or as a quick start tool
for developers interested in having an OpenTracing backend for their microservices deployed on OpenShift.

Even though this template uses a stateful Cassandra, backing storage is set to `emptyDir`. It's more 
appropriate to create a `Volume` and use it instead. 
 
Additionally, the Cassandra image is not any officially supported image. This will be changed soon, once
an official image for OpenShift is released.

And finally, the template assumes that Jaeger is being installed at the `myproject` namespace, which is the
default namespace the user is logged in after running `oc cluster up`

## Uninstalling

If you need to remove the Jaeger components created by this template, run:

```bash
oc delete all,daemonset -l jaeger
```