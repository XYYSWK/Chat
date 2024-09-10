package gtype

import (
	"bytes"
	"encoding/hex"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"io"
	"mime/multipart"
	"strconv"
)

var fileTypeMap = make(map[string]string)

func init() {
	fileTypeMap["ffd8ffe000104a464946"] = "jpg"  // JPEG [jpg]
	fileTypeMap["89504e470d0a1a0a0000"] = "png"  // PNG [png]
	fileTypeMap["47494638396126026f01"] = "gif"  // GIF [gif]
	fileTypeMap["49492a00227105008037"] = "tif"  // TIFF [tif]
	fileTypeMap["424d228c010000000000"] = "bmp"  // 16色位图[bmp]
	fileTypeMap["424d8240090000000000"] = "bmp"  // 24位位图[bmp]
	fileTypeMap["424d8e1b030000000000"] = "bmp"  // 256色位图[bmp]
	fileTypeMap["41433130313500000000"] = "dwg"  // CAD [dwg]
	fileTypeMap["3c21444f435459504520"] = "html" // HTML [html]   3c68746d6c3e0  3c68746d6c3e0
	fileTypeMap["3c68746d6c3e0"] = "html"        // HTML [html]   3c68746d6c3e0  3c68746d6c3e0
	fileTypeMap["3c21646f637479706520"] = "htm"  // HTM [htm]
	fileTypeMap["48544d4c207b0d0a0942"] = "css"  // css
	fileTypeMap["696b2e71623d696b2e71"] = "js"   // js
	fileTypeMap["7b5c727466315c616e73"] = "rtf"  // Rich Text Format [rtf]
	fileTypeMap["38425053000100000000"] = "psd"  // Photoshop [psd]
	fileTypeMap["46726f6d3a203d3f6762"] = "eml"  // Email [Outlook Express 6] [eml]
	fileTypeMap["d0cf11e0a1b11ae10000"] = "doc"  // MS Excel 注意：word、msi 和 excel的文件头一样
	fileTypeMap["d0cf11e0a1b11ae10000"] = "vsd"  // Visio 绘图
	fileTypeMap["5374616E64617264204A"] = "mdb"  // MS Access [mdb]
	fileTypeMap["252150532D41646F6265"] = "ps"
	fileTypeMap["255044462d312e350d0a"] = "pdf"  // Adobe Acrobat [pdf]
	fileTypeMap["2e524d46000000120001"] = "rmvb" // rmvb/rm相同
	fileTypeMap["464c5601050000000900"] = "flv"  // flv与f4v相同
	fileTypeMap["00000020667479706d70"] = "mp4"
	fileTypeMap["49443303000000002176"] = "mp3"
	fileTypeMap["000001ba210001000180"] = "mpg" //
	fileTypeMap["3026b2758e66cf11a6d9"] = "wmv" // wmv与asf相同
	fileTypeMap["52494646e27807005741"] = "wav" // Wave [wav]
	fileTypeMap["52494646d07d60074156"] = "avi"
	fileTypeMap["4d546864000000060001"] = "mid" // MIDI [mid]
	fileTypeMap["504b0304140000000800"] = "zip"
	fileTypeMap["526172211a0700cf9073"] = "rar"
	fileTypeMap["235468697320636f6e66"] = "ini"
	fileTypeMap["504b03040a0000000000"] = "jar"
	fileTypeMap["4d5a9000030000000400"] = "exe"        // 可执行文件
	fileTypeMap["3c25402070616765206c"] = "jsp"        // jsp文件
	fileTypeMap["4d616e69666573742d56"] = "mf"         // MF文件
	fileTypeMap["3c3f786d6c2076657273"] = "xml"        // xml文件
	fileTypeMap["494e5345525420494e54"] = "sql"        // xml文件
	fileTypeMap["7061636b616765207765"] = "java"       // java文件
	fileTypeMap["406563686f206f66660d"] = "bat"        // bat文件
	fileTypeMap["1f8b0800000000000000"] = "gz"         // gz文件
	fileTypeMap["6c6f67346a2e726f6f74"] = "properties" // bat文件
	fileTypeMap["cafebabe0000002e0041"] = "class"      // bat文件
	fileTypeMap["49545346030000006000"] = "chm"        // bat文件
	fileTypeMap["04000000010000001300"] = "mxp"        // bat文件
	fileTypeMap["504b0304140006000800"] = "docx"       // docx文件
	fileTypeMap["d0cf11e0a1b11ae10000"] = "wps"        // WPS文字wps、表格et、演示dps都是一样的
	fileTypeMap["6431303a637265617465"] = "torrent"
	fileTypeMap["6D6F6F76"] = "mov"         // Quicktime [mov]
	fileTypeMap["FF575043"] = "wpd"         // WordPerfect [wpd]
	fileTypeMap["CFAD12FEC5FD746F"] = "dbx" // Outlook Express [dbx]
	fileTypeMap["2142444E"] = "pst"         // Outlook [pst]
	fileTypeMap["AC9EBD8F"] = "qdf"         // Quicken [qdf]
	fileTypeMap["E3828596"] = "pwl"         // Windows Password [pwl]
	fileTypeMap["2E7261FD"] = "ram"         // Real Audio [ram]
}

// 获取前面结果字节的二进制
func bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF                             // 0xFF 表示无符号整数，二进制为 1111 1111 占一个字节，和其进行 & 操作的数，会保留最低 8 位，即最低字节（最后一个字节）
		hv := hex.EncodeToString(append(temp, sub)) // 将 temp 和 sub 拼接在一起，然后将字节切片转换为对应的十六进制字符串表示形式
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10)) // 如果转换后的字符串长度小于 2(即 sub 没有值)，表示只有一个字符，需要在前面补 0
		}
		res.WriteString(hv)
	}
	return res.String()
}

// GetFileType 获取文件类型
// 文件前面几个字节来判断
// fSrc：文件字节流（就用前面几个字节）
func GetFileType(file *multipart.FileHeader) (string, errcode.Err) {
	f, err := file.Open()
	if err != nil {
		return "", errcode.ErrServer
	}
	fSrc := make([]byte, 10)
	_, err = io.ReadFull(f, fSrc) // 从文件中读取指定长度的数据
	if err != nil {
		return "", errcode.ErrServer
	}
	fileCode := bytesToHexString(fSrc)

	return fileTypeMap[fileCode], nil
}
