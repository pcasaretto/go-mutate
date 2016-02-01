1 - walk the package files
2 - generate mutations
3 - run tests for each mutation
4 - check for killed mutations

https://github.com/goadesign/goa/blob/master/goagen/codegen/workspace.go#L74


1 - type check tests
2 - get all usages of current package
3 - for each usage
4 - walk the definition
5 - generate mutations
6 - run tests for each mutation
7 - check for killed mutations
