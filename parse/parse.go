/*
downloader不负责解析工作
该模块不具有实际意义，仅展示如何从redis使用下载后的内容
*/

package parse

// Parser 解析器
type Parser interface{
	Parse(resp *Response)
	Unmarshal(page string) (*Response, error)
}