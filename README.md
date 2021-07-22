# ibm-hpcs-tke-sdk

This repository implements a set of utility functions used to configure the crypto units assigned to a Hyper Protect Crypto Services (HPCS) service instance in the IBM Cloud.

A crypto unit is a portion of an IBM cryptographic coprocessor plugged into an IBM System Z mainframe.  A crypto unit is reserved for the exclusive use of a single customer.  Each crypto unit has a current master key register holding a 32-byte AES key.  The key encrypts the contents of an associated key store.  Operational work loads send requests to the crypto unit with secret data encrypted using the master key.  The secret data is unencrypted and used only inside the secure hardware of the crypto unit.

Before an HPCS service instance can be used for operational work loads, the master key registers in the crypto units assigned to the service instance must be set.  Terraform commands can be used to perform this initialization.  IBM's Terraform plug-in calls the utility functions in this repository to issue the low-level administrative commands needed to create crypto unit administrators, set the signature thresholds, and load the current master key registers.

Four functions are provided by the TKE SDK:

* Query -- Queries the current state of the crypto units assigned to an HPCS service instance.  Returned information includes the number, type, and location of crypto units, what administrators are installed, the signature thresholds, and the master key register status and verification patterns.

* CheckTransition -- Validates inputs provided in the Terraform resource block for an HPCS service instance, reads the current state of the crypto units, and determines whether the desired final state described by the resource block can be reached from the initial state of the crypto units.

* Update -- Issues low-level administrative commands to the crypto units assigned to an HPCS service instance to establish the desired final state described in the Terraform resource block for the service instance.  The Update function internally calls the Query and CheckTransition functions.

* Zeroize -- Clears the current master key registers, removes administrators, and sets signature thresholds to zero to prepare the crypto units of an HPCS service instance for deleting the service instance.

## Organization of the TKE SDK

The TKE SDK is organized as four packages:

1. github.com/IBM/ibm-hpcs-tke-sdk/common -- basic infrastructure for submitting commands to a crypto unit
2. github.com/IBM/ibm-hpcs-tke-sdk/ep11cmds -- set of parts each handling a single administrative command type
3. github.com/IBM/ibm-hpcs-tke-sdk/rest -- extends Go language capabilities for processing HTTP requests
4. github.com/IBM/ibm-hpcs-tke-sdk/tkesdk -- the four TKE SDK utility functions and common high-level functions they use

