# Wording
In the following documentation the word _origin_ refers to a fully featured Helm validation schema, described in YAML. To validate against the origin, one has only to convert it to JSON.

# Beginning already included
The beginning of the origin always contains the key-values `id`, `$schema`, `title`and `type`.
They can be ditched and are automatically added when the field `name` is present on the root object.

```
---
name: helm
```
get's unfolded to
```
---
$id: helm.schema.json
$schema: 'http://json-schema.org/draft-07/schema#'
title: helm
type: object
```

# Object notation
An object in the origin requires the fields `type` and `properties` and the optional field `additionalProperties`.
The field `type` is automatically added if the object contains the field `properties`. The field `additionalProperties: true` is automatically added, if the `properties` field contains `...:` at the end, else it is added with false.
```
someobject:
  properties:
    prop1: [...]
    prop2: [...]
    prop3: [...]
```
get's unfolded to
```
someobject:
  type: object
  additionalProperties: false
  properties:
    prop1: [...]
    prop2: [...]
    prop3: [...]
```
and
```
someobject:
  properties:
    prop1: [...]
    prop2: [...]
    prop3: [...]
    ...:
```
get's unfolded to
```
someobject:
  type: object
  additionalProperties: true
  properties:
    prop1: [...]
    prop2: [...]
    prop3: [...]
```

# Array notation
A list in the origin requires the field `type` and `items`.
The field `type` is automatically added if the list contains the field `items`. The field `additionalItems: true` is automatically added if the `items` field contains `...:` at the end, else it is added with false.
```
somelist:
  items:
    elem1: [...]
    elem2: [...]
    elem3: [...]
```
get's converted to
```
somelist:
  type: array
  additionalItems: false
  items:
    type: object
    properties:
      elem1: [...]
      elem2: [...]
      elem3: [...]
```
and
```
somelist:
  items:
    elem1: [...]
    elem2: [...]
    elem3: [...]
    ...:
```
get's converted to
```
somelist:
  type: array
  additionalItems: true
  items:
    type: object
    properties:
      elem1: [...]
      elem2: [...]
      elem3: [...]
```

# Pattern notation
A string in the origin requires the field `type` and most often this string should match again a regex pattern stored in the field `pattern`.
The field `type` is automatically added if the string contains the field `pattern`.
```
somestring:
  pattern: (\w+)
```
get's converted to
```
somestring:
  type: string
  pattern: (\w+)
```

# Simplified primitive type declaration
If one wants to declare a variable which has to be of type integer or string, by default there are two layers required. This feature removes this rendundancy and thus simplifies primitive declaration. Only if there are no additional subkeys like `minimum` added and the parentkey is either `properties` or `items`, this simplified declaration is in effect.
```
port: integer
name: string
nodes: array
```
get's converted to
```
port:
  - type: integer
name:
  - type: string
nodes:
  - type: array
```

# Switch-case instead of multiple if-then clauses
## not implemented yet
The default notation on switch case is quite write-intensive and contains many redundant informations. Thus, to avoid all this fuzz, it makes sense to introduce a switch-case functionality on properties.
Alternative: anchors
```
someobject:
  properties:
    aprop: [...]
    anotherprop: [...]
  switch:
    {what:} // is this necessary?
      properties:
        aprop
    cases:
    - equals: a
      properties:
        thirdprop: [a...]
    - equals: b
      properties:
        thirdprop: [b...]
    - #default
      properties:
        thirdprop: [c...]
      
```
get's converted to
```
someobject:
  type: object
  additionalProperties: false
  properties:
    aprop: [...]
    anotherprop: [...]
  allOf:
  - if:
      properties:
        aprop: a
    then:
      properties:
        thirdprop: [a...]
  - if:
      properties:
        aprop: b
    then:
      properties:
        thirdprop: [b...]
  - if:
      not:
        anyOf:
          - properties:
              aprop: a
          - properties:
              aprop: b
    then:
      properties:
        thirdprop: [c...]
```