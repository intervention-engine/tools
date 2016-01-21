Tools
=====

Simple tools used by the Intervention Engine project to work with data.

Installation
------------

This project currently uses Go 1.5.3 and is built using the Go toolchain.

To install Go, follow the instructions found at the [Go website](http://golang.org/doc/install).

Following standard Go practices, you should clone this project to:

```
$GOPATH/src/github.com/intervention-engine/tools
```

For example:

```
cd $GOPATH/src/github.com/intervention-engine
git clone https://github.com/intervention-engine/tools.git
cd $GOPATH/src/github.com/intervention-engine/tools
```

If you plan to install the tools to `$GOPATH/bin`, you may want to set the *GOBIN* environmental variable and add it to your path (perhaps in your `~/.profile` or `~.bashrc`\)

```
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN
```

Go only allows for one `main` function per package, therefore the *tools* repo puts each tool into a separate package. This makes building and installing a little more difficult, but allows us to have all the tools in a single repo.

generate
--------

The *generate* tool is used to generate synthetic patient records as FHIR DSTU2 resources and post them to a FHIR DSTU2 server. These synthetic records include patient info, vitals, conditions, medications, and encounters. Due to Intervention Engine's prominent use case, all synthetic records are tuned to a geriatric population.

To install *generate* to `$GOPATH/bin`:

```
$ cd $GOPATH/src/github.com/intervention-engine/tools/cmd/generate
$ go get
$ go install
```

To build a local *generate* executable:

```
$ cd $GOPATH/src/github.com/intervention-engine/tools/cmd/generate
$ go get
$ go build
```

To use a local *generate* executable:

```
$ ./generate --help
Usage of ./generate:
  -fhirURL string
    	URL for the FHIR server
  -n int
    	Number of patients to generate (default 100)
```

ieuser
------

The *ieuser* tool is used to add, remove, and edit Intervention Engine users in the database. In order to use this tool, the MongoDB database must be running.

To install *ieuser* to `$GOPATH/bin`:

```
$ cd $GOPATH/src/github.com/intervention-engine/tools/cmd/ieuser
$ go get
$ go install
```

To build a local *ieuser* executable:

```
$ cd $GOPATH/src/github.com/intervention-engine/tools/cmd/ieuser
$ go get
$ go build
```

To use a local *ieuser* executable:

```
$ ./ie-user
Usage: command <arguments> (function)
    ------
    add <username> <password> (add single user)
    add_file <filepath> (add users from comma separated file)
    remove <username> (remove single user)
    remove_all (remove all users)
    change_pass <username> <password> (change user's password)
```

uploadfhir
----------

The *uploadfhir* tool can be used to upload a FHIR DSTU2 bundle to a FHIR DSTU2 server.

To install *uploadfhir* to `$GOPATH/bin`:

```
$ cd $GOPATH/src/github.com/intervention-engine/tools/cmd/uploadfhir
$ go get
$ go install
```

To build a local *uploadfhir* executable:

```
$ cd $GOPATH/src/github.com/intervention-engine/tools/cmd/uploadfhir
$ go get
$ go build
```

To use a local *uploadfhir* executable:

```
$ ./uploadfhir --help
Usage of ./uploadfhir:
  -bundle string
    	Path to the a single JSON file containing a bundle to upload
  -fhirURL string
    	URL for the FHIR server
```

uploadhds
---------

The *uploadhds* tool can be used to convert Health Data Standards (HDS) records to FHIR DSTU2 and post them to a FHIR DSTU2 server.

To install *uploadhds* to `$GOPATH/bin`:

```
$ cd $GOPATH/src/github.com/intervention-engine/tools/cmd/uploadhds
$ go get
$ go install
```

To build a local *uploadhds* executable:

```
$ cd $GOPATH/src/github.com/intervention-engine/tools/cmd/uploadhds
$ go get
$ go build
```

To use a local *uploadhds* executable:

```
$ ./uploadhds --help
NAME:
 upload - Convert health-data-standards JSON to FHIR JSON and upload it to a FHIR Server

USAGE:
 uploadhds [global options] command [command options] [arguments...]

VERSION:
 0.0.0

COMMANDS:
 help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
 --offset, -o "0"   How many years to offset dates by
 --fhir, -f         URL for the FHIR server
 --json, -j         Path to the directory of JSON files
 --single, -s   Path to the a single JSON file
 --help, -h     show help
 --version, -v  print the version
```
