# Experimental parser Angular template

> This repository only shows what a parser on the `Go` might look like

### Benchmark

100k line of template (not uses expressions)

| Parser                            | ms   |
| --------------------------------- | ---- |
| @angular/compiler (parseTemplate) | 1700 |
| ng-template-parser                | 65   |
