![](https://img.shields.io/badge/LICENSE-AGPL-blue.svg)
[![Build Status](https://travis-ci.org/andyxning/eventarbiter.svg?branch=master)](https://travis-ci.org/andyxning/eventarbiter)

### eventarbiter
----
Kubernetes emits events when some important things happend internally.

For example, when the CPU or Memory pool Kubernetes cluster provides can not satisfy the request application made, an `FailedScheduling` event will be emitted and the message contained in the event will explain what is the reason for the `FailedScheduling` with event message like `pod (busybox-controller-jdaww) failed to fit in any node\nfit failure on node (192.168.0.2): Insufficient cpu\n` or `pod (busybox-controller-jdaww) failed to fit in any node\nfit failure on node (192.168.0.2): Insufficient memory\n`.

Also, if the application malloc a lot of memory which exceeds the `limit` watermark, kernel OOM Killer will arise and kill processes randomly. Under this circumstance, Kubernetes will emits an `SystemOOM` event with event message like `System OOM encountered.`

Note that we may use various monitor stack for Kubernetes and we can send an alarm if the average `usage` of memory exceeds the 80 percent of `limit` in the past two minutes. However, if the memory malloc operation is done in a short duration, the monitor may not work properly to send an alarm on it for that the memory `usage` will rise up highly in a short duration and after that it will be killed and restarted with memory `usage` being normal. Resource fragment exists in Kubernetes cluster. We may encounter a situation that the total remaining memory and cpu pool can satisfy the request of application but the scheduler can not schedule the application instances. This is caused that the remaining cpu and memory resource is split across all the `minion` nodes and any single `minion` can not make cpu or memory resource for the application.

Something that can not be handled by monitor can be handled by events. `eventarbiter` can watch for events, filter out events indicating bad status in Kubernetes cluster.

**eventarbiter supports callback when one of the listening events happends. eventarbiter DO NOT send event alarms for you and you should do this using yourself using callback.**


### Comparison
----
There are already some projects to do somthing about Kubernetes events.
* [Heapster](https://github.com/kubernetes/heapster) has a component `eventer`. `eventer` can watch for events for a Kubernetes cluster and supports `ElasticSearch`, `InfluxDB` or `log` sink to store them. It is really useful for collecting and storing Kubernetes events. We can monitor what happends in the cluster without logging into each `minion`. `eventarbiter` also import the logic of watching Kubernetes from `eventer`.
* [kubewatch](https://github.com/skippbox/kubewatch) can only watch for Kubernetes events about the creation, update and delete for Kubernetes `object`, such as `Pod` and `ReplicationController`. `kubewatch` can also send an alarm through `slack`. However, `kubewatch` is limited in the events can be watched and the limited alarm tunnel. With `eventarbiter`'s `callback` sink, you can `POST` the event alarm to a `transfer station`. And after that you can do anything with the event alarm, such as sending it with email or sending it with `PagerDuty`. It is on your control. :)

### Event Alarm Reason
----
|Event|Description|
|-----|-----------|
|node_notready|occures when a `minion`(`kubelet`) node changed to `NotReady`| 
|node_notschedulable|occures when a `minion`(`kubelet`) node changed status to `SchedulableDisabled`|
|node_systemoom|occures when a an application is OOM killed on a 'minion'(`kubelet`) node|
|node_rebooted|occures when a `minion`(`kubelet`) node is restrated|
|pod_backoff| occures when an container in a `pod` can not be started normally. In our situation, this may be caused by the image can not be pulled or the image specified do not exist|
|pod_failed| occures when an container in the `pod` can not be started normally. In our situation, this may be caused by the image can not be pulled or the image specified do not exist|
|pod_failedsync|occures when an container in the `pod` can not be started normally. In our situation, this may be caused by the image can not be pulled or the image specified do not exist|
|pod_insufficientcpu| occures when an application can not be scheduled du to insufficient cpu in the cluster|
|pod_insufficientmemory|occures when an application can not be scheduled du to insufficient memory in the cluster|
|pod_unhealthy|occures when the pod health check failed|

### Usage
----
Just like `eventer` in `Heapster` project. `eventarbiter` supports the `source` and `sink` command line arguments.

For `source` argument, the usage is **the same as [what it does in `eventer`](https://github.com/kubernetes/heapster/blob/master/docs/source-configuration.md)**.

For `sink` argument, the usage is like [`eventer` sink](https://github.com/kubernetes/heapster/blob/master/docs/sink-configuration.md). `eventerarbiter` supports `stdout` and `callback`.
* `stdout` can log the event alarm to `stdout` with `json` format.
* `callback` is a HTTP API with `POST` method enabled. The event alarm will be `POST`ed to the `callback` URL.
  * `--sink=callback:CALLBACK_URL`
  * `CALLBACK_URL` should return HTTP `200` or `201` for success. All other HTTP return status
  code will be considered failure.
For `environment` argument, you can set a comma separated key-value pairs as an `Environment` map
 field in event alert object. This can be used as a `context` to pass whatever you want.

Additionally, `eventerarbiter` also supports an `event_filter` argument. Event alarm reasons specified in `event_filter` will be filtered out from `eventarbiter`.

The normal commands to start an instance of `eventerarbiter` will be
* dev
  * `eventarbiter -source='kubernetes:http://127.0.0.1:8080?inClusterConfig=false' -logtostderr=true -event_filter=pod_unhealthy -max_procs=3 -sink=stdout`
* production
  * `eventarbiter -source='kubernetes:http://127.0.0.1:8080?inClusterConfig=false' -logtostderr=true -event_filter=pod_unhealthy -max_procs=3 -sink=callback:http://127.0.0.1:3086`
  * There is also a faked http service in `script/dev` listening in `3086` with `/` endpoints. 

### Build
* make build
  * Note: `eventarbiter` requires Go1.7
