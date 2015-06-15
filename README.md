# tools

Simple tools used by the Intervention Engine project to work with data.

## upload

Usage instructions:

	$ ./upload --help
	NAME:
	   upload - Convert health-data-standards JSON to FHIR JSON and upload it to a FHIR Server

	USAGE:
	   upload [global options] command [command options] [arguments...]

	VERSION:
	   0.0.0

	COMMANDS:
	   help, h	Shows a list of commands or help for one command
	   
	GLOBAL OPTIONS:
	   --offset, -o "0"	How many years to offset dates by
	   --fhir, -f 		URL for the FHIR server
	   --json, -j 		Path to the directory of JSON files
	   --single, -s 	Path to the a single JSON file
	   --help, -h		show help
	   --version, -v	print the version
