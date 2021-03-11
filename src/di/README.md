> OpenMix 出品：[https://openmix.org](https://openmix.org/mix-go)

## Mix DI

DI、IoC 容器

DI, IoC container

> 该库还有 php 版本：https://github.com/mix-php/bean

## Overview

一个创建对象以及处理对象依赖关系的库，该库可以实现统一管理依赖，全局对象管理，动态配置刷新等。

## Installation

- 安装

```
go get -u github.com/mix-go/di
```

## Quick start

通过依赖配置实例化一个单例

```go
package main

import (
    "github.com/mix-go/di"
)

type Foo struct {
}

func init() {
    obj := &di.Object{
        Name: "foo",
        New: func() (interface{}, error) {
            i := &Foo{}
            return i, nil
        },
        Singleton: true,
    }
    if err := di.Provide(obj); err != nil {
        panic(err)
    }
}

func main() {
    var foo *Foo
    if err := di.Find("foo", &foo); err != nil {
        panic(err)
    }
    // use foo
}
```

## Reference

依赖配置中引用另一个依赖配置的实例

```go
package main

import (
    "github.com/mix-go/di"
)

type Foo struct {
    Bar *Bar
}

type Bar struct {
}

func init() {
    objs := []*di.Object{
        {
            Name: "foo",
            New: func() (interface{}, error) {
                // reference bar
                var bar *Bar
                if err := di.Find("bar", &bar); err != nil {
                    return nil, err
                }

                i := &Foo{
                    Bar: bar,
                }
                return i, nil
            },
            Singleton: true,
        },
        {
            Name: "bar",
            New: func() (interface{}, error) {
                i := &Bar{}
                return i, nil
            },
        },
    }
    if err := di.Provide(objs...); err != nil {
        panic(err)
    }
}

func main() {
    var foo *Foo
    if err := di.Find("foo", &foo); err != nil {
        panic(err)
    }
    // use foo
}
```

## Refresh singleton

程序执行中配置信息发生变化时，可以刷新单例的实例

```go
obj, err := di.Container().Object("foo")
if err != nil {
    panic(err)
}
if err := obj.Refresh(); err != nil {
    panic(err)
}
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/