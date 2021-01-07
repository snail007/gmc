
# GMC Web & API Framework

<img align="right" src="/doc/images/logo2.png" width="200" height="auto"/>  

[![Build Status](https://travis-ci.com/snail007/gmc.svg?branch=master)](https://travis-ci.com/snail007/gmc)
[![codecov](https://codecov.io/gh/snail007/gmc/branch/master/graph/badge.svg)](https://codecov.io/gh/snail007/gmc)
[![Go Report](https://goreportcard.com/badge/github.com/snail007/gmc)](https://goreportcard.com/report/github.com/snail007/gmc)
[![API Reference](https://img.shields.io/badge/go.dev-reference-blue)](https://pkg.go.dev/github.com/snail007/gmc)
[![LICENSE](https://img.shields.io/github/license/snail007/gmc)](#)

GMC is a smart and flexible golang web and api development framework. GMC goal is high performance, good productivity and write less code to do more things.

# Contents

[USER GUIDE](https://snail007.github.io/gmc/) ｜ [使用指南](https://snail.gitee.io/gmc/zh)

## Attention
this project is undergoing development, will be changed frequently.

## Pull Request
PR is welcomed, but you should keep well code specification, such as : comment, testing, benchmark, example.

A package must be include the fllowing files:   

`xxx` is package name.  

| File | Description |
| ---- | ---- |
| xxx.go | main file |
| xxx_test.go | testing file, code coverage must be than 90% |
| example_test.go  | each public function example code |
| benchmark_test.go | benchmark testing file |
| doc.go | description of package |
| README.md | testing and benchmarkresult must be include |

You can reference the package sync/gpool to get more information about code specification.
