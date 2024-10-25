package poker

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func Req() {
    fmt.Println("开始请求数据")
    // doLogin()
    // aa()

    // fetchData()
    // fetchData2()
    // decode()
    // body := []byte(fetchData())
    // clockInV2()
//  
    // P,err := formatJson(body)
    // P,err := formatJson(body)
// 
    // if err != nil {
        // log.Fatalf("格式化 JSON 数据失败: %v", err)
    // }
    // fmt.Println(P)

    

}


func decryptAES(ciphertext, key string) (string, error) {
    decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", fmt.Errorf("Base64 解码失败: %v", err)
    }

    block, err := aes.NewCipher([]byte(key))
    if err != nil {
        return "", fmt.Errorf("创建 AES 密钥失败: %v", err)
    }

    if len(decodedCiphertext) < aes.BlockSize {
        return "", fmt.Errorf("密文太短")
    }

    iv := decodedCiphertext[:aes.BlockSize]
    decodedCiphertext = decodedCiphertext[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(decodedCiphertext, decodedCiphertext)

    return string(decodedCiphertext), nil
}

func aa() {
	token1 := "4A57A25E4577D0A554CEF24703B2D544EE6335F8051302E54F7E287D07C99C01531F0106AE91ECE316F5DA142217F77D20A62C5574993377C7224F18E44EDE178C510C85AC5496ADD20770FF27E15B5297F1B56C3BDB6167DC3C25741AEADD8708E49E2EAFAD0C800465472F13E57E8D915CD239C42301C2E18F8A9066D170A3BBF0C4AF771BA5BC2BBC22BAFD430DFDE18897C6C3C7CB543449C29FC9FCEA9559806D7555BD4322EC0EFBA630801476"
    key := "your-32-byte-key-here" // 需要替换为实际的密钥

    decodedToken1, err := decryptAES(token1, key)
    if err != nil {
        log.Fatalf("解密 token1 失败: %v", err)
    }
    fmt.Println("token1 解密后的内容:", decodedToken1)
}
func encode() {
    // 原始字符串
    originalString := "fa0c07c6-476c-4a70-b84b-957bb4efee8a"

    // Base64 编码
    encodedString := base64.StdEncoding.EncodeToString([]byte(originalString))
    fmt.Println("Base64 编码后的字符串:", encodedString)
}
func decode(){

    // Base64 编码后的字符串
    encodedString := "019c0f5b2400fbd2fc3dfece19bbfd92fdfacb7c70e02e367460872887de588b41000d4bd06735701f9be6750fac01d76c6eebd9438f873e9c741346dc12c9145f3d8d8fd8"

    // Base64 解码
    decodedBytes, err := base64.StdEncoding.DecodeString(encodedString)
    if err != nil {
        log.Fatalf("解码失败: %v", err)
    }
    fmt.Println("Base64 解码后的字符串:", string(decodedBytes))
}
func formatJson(data []byte) (string, error) {
    log.Println("开始格式化数据",data)
    var jsonData map[string]interface{}

    if err := json.Unmarshal(data, &jsonData); err != nil {
        return "", err
    }
    var todayData []map[string]interface{}

    // 遍历 JSON 数据并转换时间戳
    if dataArray, ok := jsonData["data"].([]interface{}); ok {
        for i, item := range dataArray {
            if itemMap, ok := item.(map[string]interface{}); ok {
                if timestamp, ok := itemMap["time"].(float64); ok {
                    itemMap["timeStr"] = formatTimestamp(timestamp)
                    if(isToday(timestamp)){
                        todayData = append(todayData, itemMap)
                    }
                   
                    dataArray[i] = itemMap 
                }
            }
        }
    }

    // 遍历todayData数据,如果有，且在打卡时间内，说明上午或下午打卡了，否则未打卡
    for i, item := range todayData {
        fmt.Println("第",i+1,"次打卡")
        if timestamp, ok := item["time"].(float64); ok {
            fmt.Println("打卡时间：",formatTimestamp(timestamp))
            if(isDuringlockInTime(timestamp)){
                fmt.Println("上午打卡了")
            }else if(!isDuringlockInTime(timestamp)){
                fmt.Println("上午未打卡")
            }else if(isDuringlockOutTime(timestamp)){
                fmt.Println("下午打卡了")
            }else if(!isDuringlockOutTime(timestamp)){
                fmt.Println("下午未打卡")
            }else{
                fmt.Println("超出打卡时间 打卡")
            }
        }
    }

    jsonData["data"] = todayData

// 将修改后的 jsonData 重新编码为 JSON 字节数组
    modifiedData, err := json.Marshal(jsonData)
    if err != nil {
        return "", err
    }

    // 格式化 JSON 数据
    var prettyJSON bytes.Buffer
    if err := json.Indent(&prettyJSON, modifiedData, "", "    "); err != nil {
        return "", err
    }

    return prettyJSON.String(), nil
}

// 当天的时间戳
func isToday(timestamp float64) bool {
    t := time.Unix(0, int64(timestamp)*int64(time.Millisecond))
    now := time.Now()

    return t.Year() == now.Year() && t.YearDay() == now.YearDay()
}
func isDuringlockInTime(timestamp float64) bool {
    // 6点~8点半
    t := time.Unix(0, int64(timestamp)*int64(time.Millisecond))
    if t.Hour() >= 6 && t.Hour() < 8 {
        return true
    }
    return false
  
}
func isDuringlockOutTime(timestamp float64) bool {
    // 17点半~18点
    t := time.Unix(0, int64(timestamp)*int64(time.Millisecond))
    if t.Hour() >= 17 && t.Hour() < 18 {
        return true
    }
    return false
}



func formatTimestamp(timestamp float64) string {
    t := time.Unix(0, int64(timestamp)*int64(time.Millisecond))

    // fmt.Println(t)
    return t.Format("2006-01-02 15:04:05")
}

func mockJsonString() []byte{
//    读取json文件
    //   1
    path := "D:\\project\\my_go_project\\tdd_go\\poker\\data.json"


    data, err := os.ReadFile(path)

    if err != nil {
        log.Fatalf("读取文件失败: %v", err)
    }

    return data

 
}
func clockInV2() {
    url := "https://oa.cmzq-office.com/standard-daily-attendance-web/attendance/clockInV2"
    method := "POST"

    payload := []byte(`{"ruleId":"66e254cac147b805bd30a0ab","type":2,"attendanceAddr":"广东省广州市天河区广州信息港","outsideOffice":1,"clockInOrder":1,"orgId":"100001377309","updateDriver":0,"driver":"77c42311e0537676a7dea0772a42b678"}`)

    client := &http.Client{}
    req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
    if err != nil {
        log.Fatalf("创建请求失败: %v", err)
    }

    req.Header.Add("host", "oa.cmzq-office.com")
    req.Header.Add("sec-ch-ua-mobile", "?1")
    req.Header.Add("sec-ch-ua-platform", `"Android"`)
    req.Header.Add("user-agent", "Mozilla/5.0 (Linux; Android 13; M2012K11AC Build/TKQ1.221114.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/129.0.6668.100 Mobile Safari/537.36  WorkbenchType_1/nativeMoa/appVersion_3.5.5.1011")
    req.Header.Add("sec-ch-ua", `"Android WebView";v="129", "Not=A?Brand";v="8", "Chromium";v="129"`)
    req.Header.Add("content-type", "application/json; charset=utf-8")
    req.Header.Add("token", "4A57A25E4577D0A554CEF24703B2D544EE6335F8051302E54F7E287D07C99C01531F0106AE91ECE316F5DA142217F77D20A62C5574993377C7224F18E44EDE178C510C85AC5496ADD20770FF27E15B5297F1B56C3BDB6167DC3C25741AEADD874E6E1079ADD1D0E8C6199D57D632F47306D408D31E296D6C889CA46CE8C32721DF28F6B45B42B16CB2D5FC7672C81E000CB468C5D3F8A4920833DAA4BAFC805559806D7555BD4322EC0EFBA630801476")
    req.Header.Add("client-type", "2")
    req.Header.Add("accept", "*/*")
    req.Header.Add("origin", "https://oa.cmzq-office.com")
    req.Header.Add("x-requested-with", "com.cmri.ercs.yqx")
    req.Header.Add("sec-fetch-site", "same-origin")
    req.Header.Add("sec-fetch-mode", "cors")
    req.Header.Add("sec-fetch-dest", "empty")
    req.Header.Add("referer", "https://oa.cmzq-office.com/h5/ydbg/lightapp/standard-dailyattendanceh5/index.html?enterId=100001377309&andfxUserid=3cd64348801e3cfad9086fd74bfd8615")
    req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
    req.Header.Add("cookie", "YK01d99ff9=019c0f5b2400fbd2fc3dfece19bbfd92fdfacb7c70e02e367460872887de588b41000d4bd06735701f9be6750fac01d76c6eebd9438f873e9c741346dc12c9145f3d8d8fd8; SESSION=NTRhNTNjYTktYjcwNC00ZTk0LTgyODAtNDU1NTJhMjVjM2Fi; YK01c92017=019c0f5b2471ff825f95006745f703fb402fc4434008ecc816dbb16c7b40cb9df03de4f08a997a191b86f54bef961abe2a0184645d")
    req.Header.Add("priority", "u=1, i")
    req.Header.Add("Connection", "close")

    res, err := client.Do(req)
    if err != nil {
        log.Fatalf("请求失败: %v", err)
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Fatalf("读取响应失败: %v", err)
    }

    fmt.Println("响应:", string(body))
}

func updateRecord() {
    url := "https://oa.cmzq-office.com/standard-daily-attendance-web/attendance/updateRecord"
    method := "POST"

    payload := []byte(`{"ruleId":"66e254cac147b805bd30a0ab","type":2,"attendanceAddr":"广东省广州市天河区广州信息港","outsideOffice":1,"orgId":"100001377309","clockInOrder":1,"clockInId":"671a1415617be70a78e496d8"}`)

    client := &http.Client{}
    req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
    if err != nil {
        log.Fatalf("创建请求失败: %v", err)
    }

    req.Header.Add("host", "oa.cmzq-office.com")
    req.Header.Add("connection", "keep-alive")
    req.Header.Add("sec-ch-ua-mobile", "?1")
    req.Header.Add("sec-ch-ua-platform", "Android")
    req.Header.Add("user-agent", "Mozilla/5.0 (Linux; Android 13; M2012K11AC Build/TKQ1.221114.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/129.0.6668.100 Mobile Safari/537.36  WorkbenchType_1/nativeMoa/appVersion_3.5.5.1011")
    req.Header.Add("sec-ch-ua", `"Android WebView";v="129", "Not=A?Brand";v="8", "Chromium";v="129"`)
    req.Header.Add("content-type", "application/json; charset=utf-8")
    req.Header.Add("token", "4A57A25E4577D0A554CEF24703B2D544EE6335F8051302E54F7E287D07C99C01531F0106AE91ECE316F5DA142217F77D20A62C5574993377C7224F18E44EDE178C510C85AC5496ADD20770FF27E15B5297F1B56C3BDB6167DC3C25741AEADD87D17A7AE8BDA1F3D302EE49F525E8E9FC0674A79B602523E119E7723F6B9DC8CFDF28F6B45B42B16CB2D5FC7672C81E0055B26AC8B899A7CD8BF19E517B2550BBDEAB78E7C8DE3CB08DE0E9C969404D6C")
    req.Header.Add("client-type", "2")
    req.Header.Add("accept", "*/*")
    req.Header.Add("origin", "https://oa.cmzq-office.com")
    req.Header.Add("x-requested-with", "com.cmri.ercs.yqx")
    req.Header.Add("sec-fetch-site", "same-origin")
    req.Header.Add("sec-fetch-mode", "cors")
    req.Header.Add("sec-fetch-dest", "empty")
    req.Header.Add("referer", "https://oa.cmzq-office.com/h5/ydbg/lightapp/standard-dailyattendanceh5/index.html?enterId=100001377309&andfxUserid=3cd64348801e3cfad9086fd74bfd8615")
    req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
    req.Header.Add("cookie", "YK01d99ff9=019c0f5b2400fbd2fc3dfece19bbfd92fdfacb7c70e02e367460872887de588b41000d4bd06735701f9be6750fac01d76c6eebd9438f873e9c741346dc12c9145f3d8d8fd8; SESSION=NTRhNTNjYTktYjcwNC00ZTk0LTgyODAtNDU1NTJhMjVjM2Fi; YK01c92017=019c0f5b2471ff825f95006745f703fb402fc4434008ecc816dbb16c7b40cb9df03de4f08a997a191b86f54bef961abe2a0184645d")

    res, err := client.Do(req)
    if err != nil {
        log.Fatalf("请求失败: %v", err)
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Fatalf("读取响应失败: %v", err)
    }

    fmt.Println("响应:", string(body))
}

func doLogin2() {
    url := "https://oa.cmzq-office.com/standard-daily-attendance-web/login/doLogin"
    method := "POST"

    payload := []byte(`{}`)

    client := &http.Client{}
    req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
    if err != nil {
        log.Fatalf("创建请求失败: %v", err)
    }

    req.Header.Add("host", "oa.cmzq-office.com")
    req.Header.Add("sec-ch-ua-platform", `"Android"`)
    req.Header.Add("sec-ch-ua", `"Android WebView";v="129", "Not=A?Brand";v="8", "Chromium";v="129"`)
    req.Header.Add("sec-ch-ua-mobile", "?1")
    req.Header.Add("client-type", "2")
    req.Header.Add("user-agent", "Mozilla/5.0 (Linux; Android 13; M2012K11AC Build/TKQ1.221114.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/129.0.6668.100 Mobile Safari/537.36  WorkbenchType_1/nativeMoa/appVersion_3.5.5.1011")
    req.Header.Add("orgid", "100001377309")
    req.Header.Add("content-type", "application/json; charset=utf-8")
    req.Header.Add("token", "4A57A25E4577D0A554CEF24703B2D544EE6335F8051302E54F7E287D07C99C01531F0106AE91ECE316F5DA142217F77D20A62C5574993377C7224F18E44EDE178C510C85AC5496ADD20770FF27E15B5297F1B56C3BDB6167DC3C25741AEADD8756FDD22F31945AAD7575B428CEC6FD2EF462E3D5F3621E115AEEBC34DFBBB1D28022E521E95129637529A3C97D785C4FA7D1950C992BAF149613E40EEAB6F9186F80626BF31AC6F812B60D152001319D")
    req.Header.Add("accept", "*/*")
    req.Header.Add("origin", "https://oa.cmzq-office.com")
    req.Header.Add("x-requested-with", "com.cmri.ercs.yqx")
    req.Header.Add("sec-fetch-site", "same-origin")
    req.Header.Add("sec-fetch-mode", "cors")
    req.Header.Add("sec-fetch-dest", "empty")
    req.Header.Add("referer", "https://oa.cmzq-office.com/h5/ydbg/lightapp/standard-dailyattendanceh5/index.html?enterId=100001377309&andfxUserid=3cd64348801e3cfad9086fd74bfd8615")
    req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
    req.Header.Add("cookie", "YK01d99ff9=019c0f5b2400fbd2fc3dfece19bbfd92fdfacb7c70e02e367460872887de588b41000d4bd06735701f9be6750fac01d76c6eebd9438f873e9c741346dc12c9145f3d8d8fd8; SESSION=N2RjYzk5ZmQtMzlhNy00ZDI2LTk0YzQtOTY4YjI3NGJiNTE1; YK01c92017=019c0f5b2471ff825f95006745f703fb402fc4434008ecc816dbb16c7b40cb9df03de4f08a997a191b86f54bef961abe2a0184645d")

    req.Header.Add("priority", "u=1, i")
    req.Header.Add("Connection", "close")

    res, err := client.Do(req)
    if err != nil {
        log.Fatalf("请求失败: %v", err)
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    log.Println("响应:11", res)
    if err != nil {
        log.Fatalf("读取响应失败: %v", err)
    }

    fmt.Println("响应:", string(body))
}


func fetchData2() {
    url := "https://oa.cmzq-office.com/standard-daily-attendance-web/statistics/queryRegionRecord"
    method := "POST"

    payload := []byte(`{"orgId":"100001377309","startDate":1729440000000,"endDate":1730044799999}`)

    client := &http.Client{}
    req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
    if err != nil {
        log.Fatalf("创建请求失败: %v", err)
    }

    req.Header.Add("host", "oa.cmzq-office.com")
    req.Header.Add("sec-ch-ua-mobile", "?1")
    req.Header.Add("sec-ch-ua-platform", `"Android"`)
    req.Header.Add("user-agent", "Mozilla/5.0 (Linux; Android 13; M2012K11AC Build/TKQ1.221114.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/129.0.6668.100 Mobile Safari/537.36  WorkbenchType_1/nativeMoa/appVersion_3.5.5.1011")
    req.Header.Add("sec-ch-ua", `"Android WebView";v="129", "Not=A?Brand";v="8", "Chromium";v="129"`)
    req.Header.Add("content-type", "application/json; charset=utf-8")
    req.Header.Add("token", "4A57A25E4577D0A554CEF24703B2D544EE6335F8051302E54F7E287D07C99C01531F0106AE91ECE316F5DA142217F77D20A62C5574993377C7224F18E44EDE178C510C85AC5496ADD20770FF27E15B5297F1B56C3BDB6167DC3C25741AEADD87D8E1FD96E73273DD9EA5F4E15B802CB72A757771225B15EF7C8ECE5F249806CBC9607739D59E0B81B07EC2AAF858267FF65998F73422C5708C475CC1E42D0E641D944D5A67083537B483CD6BE13B355E")
    req.Header.Add("client-type", "2")
    req.Header.Add("accept", "*/*")
    req.Header.Add("origin", "https://oa.cmzq-office.com")
    req.Header.Add("x-requested-with", "com.cmri.ercs.yqx")
    req.Header.Add("sec-fetch-site", "same-origin")
    req.Header.Add("sec-fetch-mode", "cors")
    req.Header.Add("sec-fetch-dest", "empty")
    req.Header.Add("referer", "https://oa.cmzq-office.com/h5/ydbg/lightapp/standard-dailyattendanceh5/index.html?enterId=100001377309&andfxUserid=3cd64348801e3cfad9086fd74bfd8615")
    req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
    req.Header.Add("cookie", "YK01d99ff9=019c0f5b2400fbd2fc3dfece19bbfd92fdfacb7c70e02e367460872887de588b41000d4bd06735701f9be6750fac01d76c6eebd9438f873e9c741346dc12c9145f3d8d8fd8; SESSION=ODZkZjg3ZmUtMzI1OS00Y2UyLWFkOWItZDlhYTZjNzEyNmUz; YK01c92017=019c0f5b2471ff825f95006745f703fb402fc4434008ecc816dbb16c7b40cb9df03de4f08a997a191b86f54bef961abe2a0184645d")
    req.Header.Add("priority", "u=1, i")
    req.Header.Add("Connection", "close")

    res, err := client.Do(req)
    if err != nil {
        log.Fatalf("请求失败: %v", err)
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Fatalf("读取响应失败: %v", err)
    }

    fmt.Println("响应:", string(body)[0:100])
}

func fetchData() []byte {
    log.Println("我是fetchData111")
    url := "https://oa.cmzq-office.com/standard-daily-attendance-web/statistics/queryRegionRecord"
    method := "POST"

    payload := []byte(`{"orgId":"100001377309","startDate":1729440000000,"endDate":1730044799999}`)

    client := &http.Client{}
    req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
    if err != nil {
        log.Fatalf("创建请求失败: %v", err)
    }

    req.Header.Add("host", "oa.cmzq-office.com")
    req.Header.Add("connection", "keep-alive")
    req.Header.Add("sec-ch-ua-mobile", "?1")
    req.Header.Add("sec-ch-ua-platform", `"Android"`)
    req.Header.Add("user-agent", "Mozilla/5.0 (Linux; Android 13; M2012K11AC Build/TKQ1.221114.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/129.0.6668.100 Mobile Safari/537.36  WorkbenchType_1/nativeMoa/appVersion_3.5.5.1011")
    req.Header.Add("sec-ch-ua", `"Android WebView";v="129", "Not=A?Brand";v="8", "Chromium";v="129"`)
    req.Header.Add("content-type", "application/json; charset=utf-8")
    req.Header.Add("cookie", "YK01d99ff9=019c0f5b2400fbd2fc3dfece19bbfd92fdfacb7c70e02e367460872887de588b41000d4bd06735701f9be6750fac01d76c6eebd9438f873e9c741346dc12c9145f3d8d8fd8; SESSION=ZmEwYzA3YzYtNDc2Yy00YTcwLWI4NGItOTU3YmI0ZWZlZThh; YK01c92017=019c0f5b2471ff825f95006745f703fb402fc4434008ecc816dbb16c7b40cb9df03de4f08a997a191b86f54bef961abe2a0184645d")
    req.Header.Add("token", "4A57A25E4577D0A554CEF24703B2D544EE6335F8051302E54F7E287D07C99C01531F0106AE91ECE316F5DA142217F77D20A62C5574993377C7224F18E44EDE178C510C85AC5496ADD20770FF27E15B5297F1B56C3BDB6167DC3C25741AEADD87D8E1FD96E73273DD9EA5F4E15B802CB72A757771225B15EF7C8ECE5F249806CBC9607739D59E0B81B07EC2AAF858267FF65998F73422C5708C475CC1E42D0E641D944D5A67083537B483CD6BE13B355E")
    req.Header.Add("client-type", "2")
    req.Header.Add("accept", "*/*")
    req.Header.Add("origin", "https://oa.cmzq-office.com")
    req.Header.Add("x-requested-with", "com.cmri.ercs.yqx")
    req.Header.Add("sec-fetch-site", "same-origin")
    req.Header.Add("sec-fetch-mode", "cors")
    req.Header.Add("sec-fetch-dest", "empty")
    req.Header.Add("referer", "https://oa.cmzq-office.com/h5/ydbg/lightapp/standard-dailyattendanceh5/index.html?enterId=100001377309&andfxUserid=3cd64348801e3cfad9086fd74bfd8615")
    req.Header.Add("accept-language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")

    res, err := client.Do(req)
    if err != nil {
        log.Fatalf("请求失败: %v", err)
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Fatalf("读取响应失败: %v", err)
    }
    fmt.Println("响应:1111", string(body)[0:100])
    return body
}