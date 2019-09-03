# Threat modelling in kubernetes clusters ++

## Reach out + thanks to the ASF

-  https://kubernetes.io/community/
-  jay@apache.org / @jayunit100 

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

- When BigData was overfunded... there was money to waste.
- Remeber Cascalog, Cascading, DSL's for Mapreduce, Hbase on HDFS, ...
- Pig, Hive, MR 1, MR2, ... growth of tools with integration problem.
- Mahout on Yarn.
- Hive interop w/ Hadoop.
- HBase on Yarn.
- Spark HDFS/S3/... connector.

## What about now? 

The BigData trough of Disillusionment is here.
- Reduced commit frequency to bigtop and other projects.
- Spark commoditized Batch SQL and simple streaming.
- HDFS isnt the only data source, and IoT means it never will be.
- Integration with infrastructure,clouds, and a cheap alt. to EMR/BigQuery is what
people need.
- Dynamic data lakes and tools for building real time analytics pipelines.
- Ultimately Reasoning about dataflow with nifi is needed.

Proprietary Big Data clouds : Can we compete with them ?

# What will we do today

Look at a *very* prototypey sketch of what BigTop could be: A batteries included
alternative to PBDCs.

## No more RPM/Debs

- Use *public docker image* for packaging. 

## No more old tools

- Drill, Hue, Puppet, Oozie, Pig....  Who cares?  People using these tools can 
maintain them on their own - and the BigTop Community isnt big enough to continue
integration testing them.,, and vendors simple arent helping us.
- Build, Deployment, RPM, and Debian == 50% of all bigtop issues.  Lets nix them.
- Focus on a small stack that provides immediate value and build a new community.


# What do Data scientists do ?

- In memory compute
- Ad hoc querying
- Cheap object stores
- Batteries included deployments.

They DONT do : Terraform, Puppet, Ansible, Maven, ... 

.. Lets give them something they can use.

# K8s doesnt have an open source, credible  bigdata or datalake story.  Lets give it one.

- BigDataSig was recently downleveleded to a working group.
- While spark eats the world - vendors have failed to build an open source blueprint that
makes data governance and K8s analytics easy for anyone.
- Rook.io was widely successfull doing this for storage, and is now becoming a defacto
standard for open source PV provisioning.
- Bigtop should do the same.


# Ok ! So lets start going through the code...

- Kubernetes reference installer: Raw Kubeadm "master" and "slave" scripts That anyone can just run
however they want to.


1. Dynamic Storage, NFS, and Host Path docs and testing
2. Helm charts that we regularly test, and run real workloads against.
3. Yamls where we cant - with *manual* deployment instructions.
  - Automation isnt worth it here.  What the community needs is integrated recipes that are hackable,
    not a proscriptive distribution that is overspecified.
4. Intro to how we'll use Kustomize + ConfigMaps for injecting into spark, nifi, presto, etc...

## This is what I propose as the future of bigtop.

- Warehousing: HBase, and Presto
- Streaming: Kafka, Spark
- Workflow: Kafka, HBase, Nifi, Zepplin

## Warehousing: Minio with Presto

Deployment archiecture w/ Presto
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

## HBase, Kafka, Nifi, Zookeeper

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


