package copy_dir

import (
	"os"
	"regexp"
	"strings"
)

// CopyDir 函数接受两个字符串参数，分别表示源文件夹和目标文件夹的路径
// 返回值是拷贝的文件数量和可能发生的错误
func CopyDir(src, dst string) (count int, err error) {
	// 使用正则表达式将路径按照 / 或 \ 分割成切片
	regexp1, err := regexp.Compile(`(/|\\)`)
	if err != nil {
		return 0, err
	}
	srcSplits := regexp1.Split(src, 10000)
	dstSplits := regexp1.Split(dst, 10000)

	// 调用 CopyDirInner 函数，传入源文件夹和目标文件夹的前缀和最后一级名称
	return CopyDirInner(strings.Join(srcSplits[:len(srcSplits)-1], "/"), srcSplits[len(srcSplits)-1], strings.Join(dstSplits[:len(dstSplits)-1], "/"), dstSplits[len(dstSplits)-1])
}

// CopyDirInner 函数接受四个字符串参数，分别表示源文件夹和目标文件夹的前缀和最后一级名称
// 返回值是拷贝的文件数量和可能发生的错误
func CopyDirInner(srcPrefix, src string, dstPrefix, dst string) (count int, err error) {
	// 如果前缀为空，则设置为当前目录
	if srcPrefix == "" {
		srcPrefix = "."
	}
	if dstPrefix == "" {
		dstPrefix = "."
	}
	// 读取源文件夹下的所有文件和子文件夹
	dirs, err := os.ReadDir(srcPrefix + "/" + src)
	if err != nil {
		return 0, err
	}

	// 在目标文件夹下创建同名的子文件夹
	pathCursor := dstPrefix + "/" + dst + "/" + src
	err = os.MkdirAll(pathCursor, 0600)
	if err != nil {
		return 0, err
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			// 如果是子文件夹，则递归调用 CopyDirInner 函数，传入相应的参数
			countSub, err := CopyDirInner(srcPrefix+"/"+src, dir.Name(), dstPrefix+"/"+dst, src)
			if err != nil {
				return 0, err
			}
			count += countSub
		} else {
			// 如果是文件，则读取其内容，并写入到目标文件夹下同名的文件中
			bytesFile, err := os.ReadFile(srcPrefix + "/" + src + "/" + dir.Name())
			if err != nil {
				return 0, err
			}

			err = os.WriteFile(pathCursor+"/"+dir.Name(), bytesFile, 0600)
			if err != nil {
				return 0, err
			}
			count++
		}
	}
	return count, nil
}
