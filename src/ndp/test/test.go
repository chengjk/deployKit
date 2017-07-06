package test

import (
	"fmt"
	"strings"
)

func main() {
	var str="e:/web{version}{version}.zip"

	str=strings.Replace(str,"{version}","v1.0",1)

	fmt.Println(str)


	//url := "http://172.30.10.171/FacebookPMD/EC/snapshots/v0.8.5_006/RELEASE-NOTE.txt"
	//resp, err := http.Get(url)
	//if err != nil {
	//	panic(err)
	//}
	//filePath := path.Join("./upload/" + path.Base(url))
	//os.MkdirAll("./upload", 0777)
	//f, err := os.Create(filePath)
	//if err != nil {
	//	panic(err)
	//}
	//io.Copy(f, resp.Body)
}
