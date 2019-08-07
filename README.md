# Cloudmeta
A cloud meta data utility.
<p align="center">
<img width="400" src="https://raw.githubusercontent.com/spotmaxtech/cloudmeta/master/assets/cloudmeta_logo.png">
</p>

# Why cloudmeta?
With it, you can get popular cloud meta data, such as aws instance info,
spot price, od-demand price, spot interruption info, regions info and so on.

*For now, SpotMax team are providing open database for you*

# Who may use it?
If you are writing some automation code with aws/aliyun, I think you can get some idea in these project!

# Supporting cloud platform
1. aws
2. aliyun

# Usage 

```bash
go get github.com/spotmaxtech/cloudmeta
```

```go
package main

import (
	"fmt"
	"github.com/spotmaxtech/cloudmeta"
	)

func main() {
	meta := cloudmeta.DefaultAWSMetaDb()
	region := meta.Region().GetRegionInfo("us-east-1")
	fmt.Println(region)
}
```