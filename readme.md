# nicky

## idea
You can use this software to validate YAML-files (f.e. helm's ```values.yaml```) against a schema. As YAML is convertible to JSON, the already existing JSON-schema implementation from https://json-schema.org/ can be used. But in addition to their features, the schema for this tool can be provided as yaml too, which will then auotmatically get converted to json. As the writing of the schema file contains many rendundant information, an unfolding process is included to the conversion. To get more information on that unfolding, read the file ```./unfolding_documentation.md```. Alternatively you can disable this feature with the flag ```--disable-unfolding```.

## dependencies
This application requires two dependencies for building.
For more information please read ```./Taskfile.yml#install-dependencies```.

## building
This application was developed in golang version 1.13.1 and can thus be compiled on at least that version.
For example build-commands please look at ```./Taskfile.yml#build```.

## running
To see example calls, run ```./nicky --help``` or ```./nicky.exe --help``` according to your OS or look at the example at ```./Taskfile_{{GOOS}}.yml#run```.

## writing/improving a schema file
To create an initial folded schema based on a given yaml-file please look at the example at ```./Taskfile_{{GOOS}}.yml#generate_schema```.
To implement or improve a new schema, one would probably need to look up the possibilities of the used JSON-schema. The fastest way to the overall specification is [here](https://json-schema.org/draft/2019-09/json-schema-core.html), [here](https://json-schema.org/understanding-json-schema/structuring.html#structuring) for more information on how to structure complex schemata and [here](https://json-schema.org/understanding-json-schema/reference/combining.html) for dependency keywords like "any of", "all of", "one of" and "not". These can be understood in a context of 'Match any of these subschemata'.

## unfolding a schema file
To just unfold a given folded schema-file, please look at the example at ```./Taskfile_{{GOOS}}.yml#unfold_schema```.

## developing
Every method of this software is extracted to another file with the same name as the function. Within this file there may be required helper functions only needed for this method. This ensures that each method is modular and selfcontained as much as sensibly possible.
If any method contains multiple steps, these steps are always splitted with at least three lines of comments for better readability of the code.
Global variables (in the meaning of being accessed from multiple methods/files) have the prefix ```global_``` to their name.

## additional notes
The process of converting ```type: null``` from YAML to JSON reads ```null``` to type ```nil``` in Golang and converting this to JSON results in ```"type":null``` instead of ```"type":"null"```. The latter would be the desired value for JSON-Schema. To work around this issue, the YAML-partial should be ```type: "null"``` as it will then get correctly treated as a string.
