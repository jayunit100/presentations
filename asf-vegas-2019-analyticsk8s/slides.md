#  A Community for testing the Cloud Native Analytics Stack. ++

## Reach out + thanks to the ASF

-  https://kubernetes.io/community/
-  jay@apache.org / @jayunit100 
-  sid@minio.io / @wlan0
- ... 


### SID Feel free to add change edit whatever, this is v0 of the slides this far.
---

# A qoute about Packaging

What commercialism has brought into Linux has been the incentive to make a good distribution that is easy to use and that has all the packaging issues worked out. 

- Linus Torvalds

---

# What *was* Bigtop ?

From our website... https://bigtop.apache.org/

- Bigtop packages Hadoop RPMs and DEBs, so that you can manage and maintain your Hadoop cluster.
- Bigtop provides an integrated smoke testing framework, alongside a suite of over 50 test files.
- Bigtop provides vagrant recipes, raw images, and (work-in-progress) docker recipes for deploying Hadoop from zero.

# Why did we need to do this? 

Cambrian explosion in bigdata tools... 

- When BigData was overfunded... there was money to waste.. people made their own tools.
- Remeber Cascalog, Cascading, DSL's for Mapreduce, Hbase on HDFS ?
- Pig, Hive, MR 1, MR2, ... growth of tools with integration problem.
- Mahout on Yarn.
- Hive interop w/ Hadoop.
- HBase on Yarn.
- Spark HDFS/S3/... connector.

## What about now? 

The BigData trough of Disillusionment is here.
- Reduced relevance of hadoop, hive, pig, mahout, hbase, other projects.
- Spark commoditized Batch SQL and simple streaming.
- HDFS isnt the only data source, and IoT means it never will be.
- Integration with infrastructure,clouds, and a cheap alt. to EMR/BigQuery is what
people need.
- Data Science is getting more demanding and more integral to everyday apps.
- Dynamic data lakes and tools for building real time analytics pipelines.
- Ultimately Reasoning about dataflow with nifi is needed.

Proprietary Big Data clouds : Can we compete with them ?

# What will we do today

Look at a *very* prototypey sketch of what BigTop could be: A batteries included
alternative to PBDCs.

## We'll start saying *NO* to old tools

Its time to let old things die. 

- Drill, Hue, Puppet, Oozie, Pig....  Who cares?  People using these tools can 
maintain them on their own - and the BigTop Community isnt big enough to continue
integration testing them.,, and vendors simple arent helping us.
- Build, Deployment, RPM, and Debian == 50% of all bigtop issues.  Lets nix them.
- Focus on a small stack that provides immediate value and build a new community.

## Bigtop Packages Hadoop .. Who cares :) ? 

We don't need RPMs and DEBs.  We need tarballs and docker images.
The age of the mutable linux server is ending.

- Use *public docker image* for packaging.
- Support any *NIX distro, alike via the Docker and the CRI.
- Multiple implementations available: Docker not a requirement.
- Tarballs - the easiest way *ever* to run hadoop, spark, ... ARE BACK!

# What do Data scientists do ?

- In memory compute
- Ad hoc querying
- Cheap object stores
- Batteries included deployments
- Infrastructure taken for granted

They DONT do : Terraform, Puppet, Ansible, Maven, ... 

.. Lets give them something they can use... *Kubernetes*

# What is Kubernetes? 

You probably already know.  If not, I'll wing it and tell you.

- An API for cloud functionality that just works, anywhere, and has push button deployments in
every public cloud.
- An autoscaler, service discovery mechanism, storage provisioner, release and upgrade manager
- Oh yeah and its also a container orchestrator.

Its essentially the entire operations aspect for any app - commoditized into an API.

No *actual* intro to K8s will happen here, since there are 1000s online.

# K8s doesnt have an open source, credible  bigdata or datalake story.  Lets give it one.

- BigDataSig was recently downleveleded to a working group.
- While spark eats the world - vendors have failed to build an open source blueprint that
makes data governance and K8s analytics easy for anyone.
- Rook.io was widely successfull doing this for storage, and is now becoming a defacto
standard for open source PV provisioning.
- Bigtop should do the same.

# Spark: How we should do ConfigMaps

static:
- ConfigMaps injected for *all* files that users can change.
- helm/stable == 1.5, we should be on 2.x at least (microsoft/helm has 2x w/ zepplin)
dynamic:
- cloud native driver + native zepplin : Spark operator from google.
- Someone needs to own the native tranlation... https://github.com/GoogleCloudPlatform/spark-on-k8s-operator/issues/531 / https://issues.apache.org/jira/browse/SPARK-24432.

# Statefull services: Nifi, HBase, Kafka -> Zk

- stable/helm charts are usable as is.
- 
How we should handle state

- Reuse Zookeeper or other bookeeping stuff, minimize resource usage and have
thoughtful DR story for CP data stores.

For production clusters, the availability of ZK as a single service
thats heavily resourced against durable storage is important.  Otherwise,
multiple ZK clusters and PVs might need to be looked at in an outage,
and you may have an exorbitantly high cost for a new ZK cluster for each 
service.

... So,...
 

```                                                                 
+-----------------------------------+                               
|                                   |                               
|  Hbase  -----> ZK                 |                               
|  Kafka ------> ZK                 |                               
|  Nifi -------> ZK                 |                               
|                                   |                               
|                                   |                               
|    Unify the zookeeper cluster,   |                               
|    inject it via configmap        |                               
|    to Hbase, Kafka, Nifi.         |                               
|                                   |                               
|    ..Finally, persistent volumes..|                               
|                                   |                               
+-----------------------------------+                               
```                                 

# Minio + Presto : Give users a warehouse

People need to do ad hoc querying.  Package Minio and PResto together for people to use
as a one-stop warehouse for querying data at any scale.
```
+------------------------+                      
 |                        |                      
 |                        |                      
 | VM                     |                      
 |                        |                      
 |  minio                 |                      
 |  /data01/ (50G)        |                      
 |  ./minio server /data01/.                     
 |                        |                      
 |                        |                      
 +------------------------+                      
                                                 
 +-----------------+                             
 |                 |                             
 |                 |                             
 | VM              |                             
 |                 |                             
 | presto          |                             
 | /               | <-- also, hive.metastore.uri
 |   minio.        |                             
 |     properties  |                             
 |       hive.s3.* |                             
 +-----------------+                             
```

# Result: 

An analytics distro that is native to kubernetes, which can be use to cluster, query, and store
pedascalable data, which also has a single, opinionated model for deploying and running spark.

## Possibly use the spark operator here.

# Ok ! So lets start going through the code... and do some demos.

- Kubernetes reference installer: Raw Kubeadm "master" and "slave" scripts That anyone can just run
however they want to.



## Notebooks: Spark with Zepplin

Normally... 

```                                             
+------------------++------------------+         
|  VM              || VM               |         
|                  ||                  |         
|  Spark           || Spark            |         
|     master url   ||   slaves         |         
|     <- workers   ||                  |         
|                  ||                  |         
|                  ||                  |         
|                  ||                  |         
+------------------++------------------+         
```

ALT: Spark operator with autoscaling ?  See SPARK-24432.


### with these files modified by puppet/shell/manually...

```
├── README.md
├── core-site.xml (s3a.filesystem minio urls)
├── log4j.properties
├── spark-defaults.conf
├── spark-deployment.yaml
└── spark-env.sh
```

For bigtop, configmapify these files, with seeded attributes for accessing the default
object store.

### Putting it all together

Ultimately, the future of BigTop's architecture
might look something like this.

```
    +-------------------------+                                       
    |                         |                                       
    |                         |                                       
    |                         |                                       
    |                         |                                       
    | minio                   |                                       
    | /data01/ (50G)          |                                       
    | ./minio server /data01/.                +-------------------+   
    |                         |+------------> |   Ingress(nginx)  |   
    |                         |               |                   |   
    +-------------------------+               |                   |   
                                              |                   |   
   +-----------------++---------------------> +-------------------+   
   |                 |                                     ^    ^     
   |                 |                                     |    |     
   |                 |                                     |    |     
   |                 |                                     |    |     
   | presto          |                                     |    |     
   | /               | <-- also, hive.metastore.uri        |    |     
   |   minio.        |                                     |    |     
   |     properties  |                                     |    |     
   |       hive.s3.* |                                     |    |     
   +-----------------+                                     |    |     
                                                           |    |     
  +------------------++------------------+                 |    |     
  |  VM              || VM               | ----------------+    |     
  |                  ||                  |                      |     
  |  Spark           || Spark            | ├── README.md        |     
  |     master url   ||   slaves         | ├── core-site.xml    |     
  |     <- workers                       | ├── log4j.properties |     
  |                  ||                  | ├── spark-defaults.conf    
  |                  ||                  | ├── spark-deployment.yaml  
  |                  ||                  | └── spark-env.sh     |     
  +------------------++------------------+                      |     
                                                                |     
 +--------------------------------------+                        |     
 |                                      |                        |     
 |  Hbase  -----> ZK                    |                        |     
 |  Kafka ------> ZK                    | -----------------------+     
 |  Nifi -------> ZK                    |                              
 |                                      |                              
 |                                      |                              
 |    Unify the zookeeper cluster,      |                             +
 |    inject it via configmap           |                              
 |    to Hbase, Kafka, Nifi.            |                              
 |                                      |                              
 |    ..Finally, persistent volumes..   |                          +   
 |                                      |                              
 +--------------------------------------+                              
```



