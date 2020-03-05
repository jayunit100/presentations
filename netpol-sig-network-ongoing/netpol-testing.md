# NetworkPolicy testing framework with truth tables

Summary goes here.

## Motivation and Goals

The Y goes here?

## Design

What the framework does.

## Current state

- An implementation can be found [here](https://github.com/vmware-tanzu/antrea/tree/master/hack/netpol).
- Runs about 14 test cases in about 10mins.
- In comparision, e2e tests focused on NetworkPolicy takes close to an hour.
- Recently integrated with Antrea CI.

## Future work

- More test cases for a full coverage of NetworkPolicy spec.
- Cleanup flag to determine if resources created by NetPol must be deleted.
- Node specific test cases.
- Extend framework to run scale tests for NetworkPolicy.
- Extend framework to test other K8s resources with truth tables.
