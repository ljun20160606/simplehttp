<p align="center">
    <img src="doc/simplehttp.jpg" width="325"/>
</p>
<p align="center">基于 <code>Go</code> + <code>Http</code>😋</p>
<p align="center">
    🔥 <a href="#快速开始">快速开始</a>
</p>

<p align="center">
    <a href="https://golang.org/doc/go1.10"><img src="https://img.shields.io/badge/go-v1.10.0-blue.svg"></a>
    <a href="http://commitizen.github.io/cz-cli"><img src="https://img.shields.io/badge/commitizen-friendly-brightgreen.svg"></a>
</p>

***

## 快速开始

在Github中搜索simplehttp

````go
package main

import "github.com/ljun20160606/simplehttp"

func main() {
     simplehttp.
            Get("https://github.com/search").
            Query("utf8", "✓").
            Query("q", "simplehttp").
            Send().
            String()
}
````

更多 [example/github.go](./example/github.go)

