# Cloudmeta
A cloud meta data utility.
![](assets/cloudmeta_logo.png?width=200)

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