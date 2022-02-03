# Equivalent Unified Ideographs Converter

Replace KANGXI and CJK radicals (and strokes) with their equivalent CJK Unified Ideographs.

### Purpose
The purpose of this library is mainly to deal with inconsistent character encodings for the same visually similar ideographs. Depending on where your data has come from, a character may look identical to your eye but have a different underlying unicode value. Attempting to programmatically match or index on these values will fail (as they are visually the same but the value is different). To illustrate the problem, see the following table of characters which are visually similar but have different unicode values. (the CJK radical _grass_):


| Character 	| Unicode Hex 	| Name                       	|
|-----------	|-------------	|----------------------------	|
| ⺾        	| 2EBE        	| CJK RADICAL GRASS ONE      	|
| ⺿        	| 2EBF        	| CJK RADICAL GRASS TWO      	|
| ⻀        	| 2EC0        	| CJK RADICAL GRASS THREE    	|
| 艹        	| 8279        	| CJK UNIFIED IDEOGRAPH 8279 	|


In the above case, you may have one dataset which uses `CJK RADICAL GRASS ONE` characters, and another which uses `CJK RADICAL GRASS TWO`, and these will not match. We can solve this problem by converting the "variant" ideographs to their CJK Unified Ideograph equivalents across all datasets.

### Usage
#### File Transformation
The `cmd/replace/main.go` file can be used to replace all the listed variants in a file with their equivalent unified ideographs. This should be run from the root of the project:

```
go run cmd/replace/main.go -input=example/mismatched.txt
```

Optionally, the replaced data can be output to a new file, leaving the original intact:

```
go run cmd/replace/main.go -input=example/mismatched.txt -output=example/out.txt
replaced 11 characters in example/out.txt with data from example/mismatched.txt
⺁ : 1
⺄ : 1
⺆ : 1
⺇ : 1
⺼ : 1
⺅ : 1
⺩ : 2
⺫ : 1
⺭ : 1
⺶ : 1
```

#### Importable Library
The underlying code is designed to be imported so you can use the functionality in your projects to normalise the radicals and strokes of input strings so they are consistent in your data. 

> The two main functions are `func Replace(s string) []byte` (for simple strings) and `func BufferedReplace(r io.Reader) (out *bytes.Buffer, err error)` (for reading files, etc). 

```
package main

import (
    "fmt"
    eqi "github.com/mochi-co/equivalent-unified-ideograph"
)

func main() {
    o := eqi.Replace(`a⺁b⺄c⺅d⺆e⺇f⺩g⺫h⺭i⺶j⺼`)
    fmt.Println(o) // a厂b乙c亻d冂e𠘨f王g目h礻i羊j肉
}
```

See the `cmd` folder for more examples of use.

### Data
The translation tables for making this conversion are based upon the [unicode.org](unicode.org) [EquivalentUnifiedIdeograph.txt](https://www.unicode.org/Public/UNIDATA/EquivalentUnifiedIdeograph.txt) file originally created by Ken Lunde. A copy of the EquivalentUnifiedIdeograph.txt file is provided as `data/EquivalentUnifiedIdeograph.txt` in order to facilitate the generation of translation tables, and will be updated periodically if and when upstream changes become available. 

#### Regenerating the translation tables (static.go)
The `static.go` file defines two variables containing variants and their unified equivalents - Pairs, an array containing each variant character, unified character, and variant name; and MappedPairs, a map of unified characters keyed on variant character (as strings). These "translation tables" are generated from the data listed in EquivalentUnifiedIdeograph.txt. 

The utility for generating the `static.go` file is provided in this repository as `cmd/updatestatic/main.go`, which should be run from the root of the project:

```
go run cmd/updatestatic/main.go
```

Update Static also takes a couple of optional parameters if you wish to change the input and output (eg. you have a custom unicode remapping file):

```
go run cmd/updatestatic/main.go -index=data/EquivalentUnifiedIdeograph.txt -output=static.go
```

### Contributions

Contributions, Pull Requests and feedback are welcomed and encouraged! Open an issue to report a bug, ask a question, or make a feature request.