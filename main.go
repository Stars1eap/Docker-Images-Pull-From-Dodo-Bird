package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	apiBaseURL = "https://docker.aityp.com/api/v1/image?search="
)

type ImageResponse struct {
	Count   int     `json:"count"`
	Error   bool    `json:"error"`
	Results []Image `json:"results"`
}

type Image struct {
	Source    string `json:"source"`
	Mirror    string `json:"mirror"`
	Platform  string `json:"platform"`
	Size      string `json:"size"`
	CreatedAt string `json:"createdAt"`
}

func main() {
	pullCmd := flag.NewFlagSet("pull", flag.ExitOnError)
	platform := pullCmd.String("platform", "", "指定平台 (如: linux/amd64)")
	tag := pullCmd.String("tag", "", "指定标签")
	interactive := pullCmd.Bool("i", false, "交互式选择")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "pull":
		pullCmd.Parse(os.Args[2:])
		if pullCmd.NArg() == 0 {
			fmt.Println("错误: 需要指定镜像名")
			pullCmd.Usage()
			os.Exit(1)
		}
		imageName := pullCmd.Arg(0)
		handlePull(imageName, *platform, *tag, *interactive)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("未知命令: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`
Docker镜像拉取工具 FROM 渡渡鸟镜像同步站（https://docker.aityp.com/）
作者: Stars1eap
工具开源地址：https://github.com/Stars1eap/Docker-Images-Pull-From-Dodo-Bird
镜像源: 渡渡鸟（https://docker.aityp.com/）
版本: 1.0

使用方法:
  ./dimages pull [选项] <镜像名>

选项:
  -platform string   指定平台 (如: linux/amd64, linux/arm64)
  -tag string       指定标签
  -i                交互式选择镜像

示例:
  ./dimages pull nginx	获取 nginx 镜像
  ./dimages pull -platform linux/arm64 nginx	获取 linux/arm64 平台 nginx 镜像
  ./dimages pull -tag latest nginx	获取 nginx:latest 镜像
  ./dimages pull -i nginx	多个镜像情况下查询50个镜像`)
}

func handlePull(imageName, platform, tag string, interactive bool) {
	images, err := searchImages(imageName)
	if err != nil {
		fmt.Printf("查询镜像失败: %v\n", err)
		os.Exit(1)
	}

	if len(images) == 0 {
		fmt.Printf("未找到镜像: %s\n", imageName)
		os.Exit(1)
	}

	// 过滤镜像
	filteredImages := filterImages(images, platform, tag)

	if len(filteredImages) == 0 {
		fmt.Println("没有找到匹配的镜像")
		os.Exit(1)
	}

	var selectedImage *Image
	if interactive || len(filteredImages) > 1 {
		selectedImage, err = selectImageInteractively(filteredImages)
		if err != nil {
			fmt.Printf("选择镜像失败: %v\n", err)
			os.Exit(1)
		}
	} else {
		selectedImage = &filteredImages[0]
	}

	fmt.Printf("拉取镜像: %s\n", selectedImage.Mirror)
	if err := pullImage(selectedImage.Mirror); err != nil {
		fmt.Printf("拉取镜像失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("成功拉取镜像: %s\n", selectedImage.Mirror)
}

func searchImages(imageName string) ([]Image, error) {
	url := apiBaseURL + imageName

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回错误状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var imageResponse ImageResponse
	if err := json.Unmarshal(body, &imageResponse); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	if imageResponse.Error {
		return nil, fmt.Errorf("API返回错误")
	}

	return imageResponse.Results, nil
}

func filterImages(images []Image, platform, tag string) []Image {
	var filtered []Image

	for _, img := range images {
		match := true

		if platform != "" && img.Platform != platform {
			match = false
		}

		if tag != "" {
			sourceParts := strings.Split(img.Source, ":")
			if len(sourceParts) < 2 || sourceParts[1] != tag {
				match = false
			}
		}

		if match {
			filtered = append(filtered, img)
		}
	}

	return filtered
}

func selectImageInteractively(images []Image) (*Image, error) {
	fmt.Printf("\n找到 %d 个镜像:\n", len(images))
	for i, img := range images {
		fmt.Printf("[%d] %s\n", i+1, img.Source)
		fmt.Printf("    镜像: %s\n", img.Mirror)
		fmt.Printf("    平台: %s, 大小: %s, 创建时间: %s\n\n",
			img.Platform, img.Size, formatDate(img.CreatedAt))
	}

	fmt.Print("请选择要拉取的镜像 (输入编号): ")
	var input string
	fmt.Scanln(&input)

	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(images) {
		return nil, fmt.Errorf("无效的选择")
	}

	return &images[choice-1], nil
}

func formatDate(dateStr string) string {
	// 简化日期显示
	if len(dateStr) > 10 {
		return dateStr[:10]
	}
	return dateStr
}

func pullImage(mirrorURL string) error {
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("未找到docker命令: %v", err)
	}

	cmd := exec.Command("docker", "pull", mirrorURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
