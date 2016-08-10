Tools [![Build Status](https://travis-ci.org/intervention-engine/tools.svg?branch=master)](https://travis-ci.org/intervention-engine/tools)
===========================================================================================================================================

The *tools* repository contains command-line tools that are helpful (or sometimes necessary) for the [Intervention Engine](https://github.com/intervention-engine/ie) project. These include tools for generating and uploading synthetic patient data, uploading FHIR bundles, and converting and uploading Health Data Standards (HDS) records.

Building and Running tools Locally
----------------------------------

Intervention Engine is a stack of tools and technologies. For information on installing and running the full stack, please see [Building and Running the Intervention Engine Stack in a Development Environment](https://github.com/intervention-engine/ie/blob/master/docs/dev_install.md).

For information related specifically to building and running the tools in this repository (*tools*), please refer to the following sections in the above guide. Note that most of these tools require the ie server (or a fhir server) to be running.

-	(Prerequisite) [Install Git](https://github.com/intervention-engine/ie/blob/master/docs/dev_install.md#install-git)
-	(Prerequisite) [Install Go](https://github.com/intervention-engine/ie/blob/master/docs/dev_install.md#install-go)
-	(Prerequisite) [Clone tools Repository](https://github.com/intervention-engine/ie/blob/master/docs/dev_install.md#clone-tools-repository)

Go only allows for one `main` function per package, therefore the *tools* repo puts each tool into a separate package. This makes building and installing a little more difficult, but allows us to have all the tools in a single repo.

generate
--------

The *generate* tool is used to generate synthetic patient records as FHIR DSTU2 resources and post them to a FHIR DSTU2 server. These synthetic records include patient info, vitals, conditions, medications, and encounters. Due to Intervention Engine's prominent use case, all synthetic records are tuned to a geriatric population. Before you can run the *generate* tool, you must install its dependencies via `go get` and build the `generate` executable:

```
$ cd $GOPATH/src/github.com/intervention-engine/tools/cmd/generate
$ go get
$ go build
```

The *generate* tool takes a `-fhirURL` flag to indicate the FHIR server to upload the patients to, as well as a `-n` flag to indicate the number of patients to generate (with the default being 100).

```
$ ./generate -fhirURL http://localhost:3001 -n 20
```

To get usage information, run `generate` with the `-help` flag.

uploadfhir
----------

The *uploadfhir* tool can be used to upload a FHIR DSTU2 bundle to a FHIR DSTU2 server. Before you can run the *uploadfhir* tool, you must install its dependencies via `go get` and build the `uploadfhir` executable:

```
$ cd $GOPATH/src/github.com/intervention-engine/tools/cmd/uploadfhir
$ go get
$ go build
```

The *uploadfhir* tool takes a `-fhirURL` flag to indicate the FHIR server to upload the patients to, as well as a `-bundle` flag to indicate the path to a FHIR DSTU2 bundle to upload.

```
$ ./uploadfhir -fhirURL http://localhost:3001 -bundle /path/to/some/bundle.json
```

To get usage information, run `uploadfhir` with the `-help` flag.

uploadhds
---------

The *uploadhds* tool can be used to convert Health Data Standards (HDS) records to FHIR DSTU2 and post them to a FHIR DSTU2 server. Before you can run the *uploadhds* tool, you must install its dependencies via `go get` and build the `uploadhds` executable:

```
$ cd $GOPATH/src/github.com/intervention-engine/tools/cmd/uploadhds
$ go get
$ go build
```

The *uploadhds* tool takes a `-fhir` flag to indicate the FHIR server to upload the patients to. The *uploadhds* tool can upload a whole directory of HDS patients, indicated by the `-json` flag, or it can upload a single HDS patient, indicated by the `-single` flag. In addition, the *uploadhds* tool can offset dates by a specified number of years, indicated by the `offset` flag (potentially useful for testing).

```
$ ./uploadhds -fhir http://localhost:3001 -single /path/to/some/hds/patient.json
```

To get usage information, run `uploadhds` with the `-help` flag.

License
-------

Copyright 2016 The MITRE Corporation

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
