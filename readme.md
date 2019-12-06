# nicky

## idea
You can use this software to validate YAML-files (f.e. helm's values.yaml) against a schema. As YAML is convertible to JSON, the already existing JSON-schema from https://json-schema.org/ can be used. The schema can be provided as yaml file too, which will then be converted to json. As the writing of the schema file contains many rendundant information, an unfolding process is included. To read more on that read the file unfolding_documentation.md. Alternatively you can disable this with the flag --disable-unfolding.

## dependencies
This application requires two dependencies for building.
See ./Taskfile.yml#install-dependencies.

## building
This application was developed in golang version 1.13.1 and can thus be compiled on that version.
For the build-commands see ./Taskfile.yml#build.

## running
To see example calls, run ./nicky --help or ./nicky.exe --help according to your OS.

## writing/improving a schema file
To implement or improve a new schema, one would probably need to look up the possibilities of the used JSON-schema. The fastest way to the overall specification is [here](https://json-schema.org/draft/2019-09/json-schema-core.html), [here](https://json-schema.org/understanding-json-schema/structuring.html#structuring) for more information on how to structure complex schemata and [here](https://json-schema.org/understanding-json-schema/reference/combining.html) for dependency keywords like "any of", "all of", "one of" and "not". These can be understood in a context of 'Match any of these subschemata'.

## developing
Every method of this software is extracted to another file, with the same name as the function. Within this file there may be any required helper functions only needed for this method. By this is ensured, that each method is modular and selfcontaining as much as possible.
If any method contains multiple steps, these are always splitted with at least three lines of comments.
Variables, which are global in the meaning of being accessed from multiple methods/files have the prefix 'global_' to their name.

## additional notes
the process of converting ```type: null``` from YAML to JSON reads ```null``` to type ```nil``` in golang and converting this to JSON results in ```"type":null``` instead of ```"type":"null"```. The latter would be the desired value for JSON-Schema. To work around this issue, the YAML-partial should be ```type: "null"``` as it will then get correctly treated as a string.