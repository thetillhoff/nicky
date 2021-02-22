![Go build](https://github.com/tillhoff/nicky/workflows/Go/badge.svg?event=push)
![Go version](https://img.shields.io/badge/Go--version-1.13.1-informational)
![Go Report Card](https://goreportcard.com/badge/tillhoff/nicky)
![GitHub](https://img.shields.io/github/license/tillhoff/nicky)

# nicky
The YAML-nitpicker of your trust.

## idea
You can use this software to validate YAML-files (f.e. helm's `values.yaml`) against a schema. As (most) YAML can be converted to JSON, the already existing JSON-schema implementation from https://json-schema.org/ is used.

In addition to their features, this tool allows schemata to be written in yaml as well, which will then be automatically get converted to json when validation occurs. As the writing of the schema file contains a lot rendundant information, an (un)folding process is included in the conversion. To get more information on that feature, read the file [unfolding_documentation.md](./unfolding_documentation.md). Alternatively you can disable this feature completely with the flag `--disable-unfolding`.

## running
To see example calls, run `./nicky --help` or `./nicky.exe --help` according to your OS or look at the example at [Taskfile.yml#validate](./Taskfile.yml).

## writing/improving a schema file
To create an initial (folded) schema based on a given YAML-file please look at the example at [Taskfile.yml#generate_schema](./Taskfile.yml).
To completely implement a new schema, one would probably need to look up the possibilities of the used JSON-schema. The fastest way to the overall specification is [here (core)](https://json-schema.org/draft/2019-09/json-schema-core.html), [here (structuring)](https://json-schema.org/understanding-json-schema/structuring.html#structuring) for more information on how to structure complex schemata and [here (combining)](https://json-schema.org/understanding-json-schema/reference/combining.html) for dependency keywords like `any of`, `all of`, `one of` and `not`. These can be understood in a context of "Match _any of_ these subschemata".
## unfolding a schema file
To just unfold a given folded schema-file, please look at the example at [./Taskfile.yml#unfold_schema](./Taskfile.yml).

## developing
Every method of this software is extracted to another file with the same name as the function. Within this file there may be corresponding helper functions only needed for this method. This ensures that each method is modular and selfcontained as much as sensibly possible.
If any method contains multiple steps, these steps are always splitted with at least three lines of comments for better readability of the code.
Global variables (in the meaning of being accessed from multiple methods/files) have the prefix `global_` to their name.

## dependencies
This application requires two dependencies for building.
For more information please read [./Taskfile.yml#install-dependencies](./Taskfile.yml#install-dependencies).

## building
This application was developed in golang version 1.13.1 and can thus be compiled on at least that version.
For example build-commands please look at [./Taskfile.yml#build](./Taskfile.yml).

## additional notes
The process of converting `type: null` from YAML to JSON reads `null` to type `nil` in Golang and converting this to JSON results in `"type":null` instead of `"type":"null"` which would be the desired value for the JSON-Schema. To work around this issue, the YAML-partial should be `type: "null"` as it will then get correctly treated as a string.
