package main

import (
	"crypto/tls"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"os"
	"send-fiule-to-emails/cache"
	"send-fiule-to-emails/util"
	"strconv"
	"strings"
	"time"
)
import "github.com/360EntSecGroup-Skylar/excelize/v2"

var mail *gomail.Dialer
var envFilePath = ".env"
var cacheFilePath = "cache.json"
var excel = "邮箱.xlsx"

func init() {
	if err := godotenv.Load(envFilePath); err != nil {
		panic("读取环境变量失败")
	}
	fmt.Println("加载环境变量成功")
	host := os.Getenv("EMAIL_HOST")
	user := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASSWORD")

	port, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	mail = gomail.NewDialer(host, port, user, password)
	mail.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	fmt.Println("初始化邮箱成功")
	cache.Init(cacheFilePath)
	fmt.Println("初始化缓存成功")
}

func main() {
	var err error
	keyEmails, err := ReadNameEmails(excel)
	if err != nil {
		panic(fmt.Sprintf("读取邮箱失败: %s", err.Error()))
	}
	fmt.Printf("共有%d个邮箱待接收文件\n\n", len(keyEmails))
	var failedNames []string
	for name, email := range keyEmails {
		fmt.Println("------------------------")
		fmt.Printf("开始处理: %s\n", name)
		if cache.Get(name) != nil {
			fmt.Printf("%s 邮件已发送, 跳过处理\n", name)
			continue
		}
		file, err := util.FindFileName1(name, "files")
		if err != nil {
			fmt.Printf("%s 未找到匹配文件, 跳过处理\n", name)
			failedNames = append(failedNames, name)
			continue
		}
		if !util.IsFileExists(file) {
			fmt.Printf("%s 待发送文件不存在, 跳过处理\n", name)
			failedNames = append(failedNames, name)
			continue
		}
		// 发送邮件
		time.Sleep(5 * time.Second)
		if err := SendFile(name, email, file); err != nil {
			fmt.Printf("%s 发送邮件失败, 跳过处理: %s\n", name, err.Error())
			failedNames = append(failedNames, name)
			continue
		}
		fmt.Println("发送邮件成功")
		// 设置为已发送
		cache.Set(name, email)
		fmt.Printf("处理成功: %s\n", name)
	}
	fmt.Println("处理完成")
}

type KeyEmailMap map[string]string

func ReadNameEmails(excelPath string) (KeyEmailMap, error) {
	if !util.IsFileExists(excelPath) {
		return nil, fmt.Errorf("%s 不存在", excelPath)
	}
	keyEmails := KeyEmailMap{}
	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		return nil, err
	}
	rows, err := f.GetRows(f.GetSheetList()[0])
	for _, row := range rows {
		if len(row) < 2 {
			continue
		}
		if row[0] == "" || row[1] == "" {
			continue
		}
		keyEmails[row[0]] = row[1]
	}
	return keyEmails, nil
}

func SendFile(name string, email string, file string) error {
	t := strings.Split(file, "/")
	fileName := t[len(t)-1]
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_FROM"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", fmt.Sprintf("%s成绩单", name))
	h := make(map[string][]string, 0)
	h["Content-Type"] = []string{fmt.Sprintf(`application/octet-stream; charset=utf-8; name="%s"`, fileName)} //要设置这个，否则中文会乱码
	fileSetting := gomail.SetHeader(h)
	m.Attach(file, fileSetting, gomail.Rename(fileName))
	return mail.DialAndSend(m)
}
