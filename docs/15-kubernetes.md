DaemonSet
A DaemonSet defines Pods that provide node-local facilities. 
A DaemonSet ensures that all (or some) Nodes run a copy of a Pod.
As nodes are added to the cluster, Pods are added to them. As nodes are removed from the cluster, those Pods are garbage collected. 
Some typical uses of a DaemonSet are:
running a cluster storage daemon on every node
running a logs collection daemon on every node