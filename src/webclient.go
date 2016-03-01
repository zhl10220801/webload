package main

import (
        "fmt"
        "net"
        "net/http"
        "time"
        "io/ioutil"
//       "math/rand"
       "flag"
)
func substr(str string, start, length int) string {
    rs := []rune(str)
    rl := len(rs)
    end := 0
    if start < 0 {
        start = rl - 1 + start
    }
    end = start + length
    if start > end {
        start, end = end, start
    }
    if start < 0 {
        start = 0
    }
    if start > rl {
        start = rl
    }
    if end < 0 {
        end = 0
    }
    if end > rl {
        end = rl
    }
    return string(rs[start:end])
}


func web(i int,url string,timeout int, connect_sucess_number_channel chan int ,connect_timeout_number_channel chan int ,connect_number_channel chan int,connect_command_channel  int ,connect_filesize_channel chan int) {
            connect_process :=i
//           connect_number :=0
          connect_sucess_number :=0
          connect_timeout_number :=0
          connect_filesize_number :=0
           end_time :=int(time.Now().Unix())+connect_command_channel
           fmt.Println("协程",connect_process,end_time)
l :=0
           for  { l++

           v :=  end_time-int(time.Now().Unix())

        if v>0 {
                /* 等待时间 可以 随机产生 0.1-1秒
                r := rand.New(rand.NewSource(time.Now().UnixNano()))
                n :=(r.Intn(10)+1)*1000000
                time.Sleep(time.Duration(n))
                */
                client := &http.Client{
                        Transport: &http.Transport{
                                Dial: func(netw, addr string) (net.Conn, error) {
                                        deadline := time.Now().Add(time.Duration(timeout)*1000000000)
                                        c, err := net.DialTimeout(netw, addr, time.Duration(timeout)*1000000000)
                                        if err != nil {
                                                return nil, err
                                        }
                                        c.SetDeadline(deadline)
                                        return c, nil
                                },
                        },
                }
                reqest, _ := http.NewRequest("GET", url, nil)
                       now := time.Now()
                       log_id := time.Now().UnixNano()
                        t :=time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
                        log_time :=substr(t, 0,19)

                        response, err := client.Do(reqest)
                        if err != nil {

                        end_time := time.Now()
                          connect_timeout_number++
                        var dur_time time.Duration = end_time.Sub(now)
                                fmt.Println("协程",connect_process,"|",connect_timeout_number,"|",log_id,"|",url,"|timeout|",log_time,"|",dur_time, "|timeout")
                                continue
                        }
                        end_time := time.Now()
                        var dur_time time.Duration = end_time.Sub(now)
        body, _ := ioutil.ReadAll(response.Body)
        len :=len(body)
         connect_sucess_number++
         connect_filesize_number +=len
        fmt.Println("协程",connect_process,"|",connect_sucess_number, "|", log_id,"|", url,"|",response.StatusCode,"|", log_time,"|",dur_time, "|",len)
      response.Body.Close()



               } else {

break         
}
}
        defer func() {
                connect_timeout_number_channel <- connect_timeout_number
                connect_sucess_number_channel  <- connect_sucess_number
                connect_number_channel <- connect_sucess_number+connect_timeout_number
                connect_filesize_channel  <- connect_filesize_number 

        }()
}


func main() {
          connect_sucess_number_channel := make(chan int,1000)
          connect_timeout_number_channel := make(chan int,1000)
          connect_number_channel := make(chan int,1000)
          connect_filesize_channel := make(chan int,1000)
        url := flag.String("url", "http://www.sohu.com", "Input your username")
        connect_command_channel := flag.Int("runtime", 2, "run time")

        c := flag.Int("client", 3, "client number")
        timeout := flag.Int("timeout", 2, "net time out")
       flag.Parse()
        for i := 0; i <  *c; i++ {
                go web(i,*url,*timeout,connect_sucess_number_channel,connect_timeout_number_channel,connect_number_channel,*connect_command_channel,connect_filesize_channel)
        }
        time.Sleep(time.Duration(*connect_command_channel)*1000000000)   
time.Sleep(1e9)
fmt.Println(" stop all agent ")
close(connect_sucess_number_channel)
close(connect_timeout_number_channel)
close(connect_number_channel)
close(connect_filesize_channel)

connect_sucess_number_total :=0
for i := range  connect_sucess_number_channel{
        connect_sucess_number_total += i   

                     }
fmt.Println("成功连接数 ", connect_sucess_number_total)

connect_timeout_number_total :=0
for i := range  connect_timeout_number_channel{
        connect_timeout_number_total += i
                     }
fmt.Println("超时数", connect_timeout_number_total)

connect_number_total :=0
for i := range  connect_number_channel{
        connect_number_total += i
                     }
fmt.Println(" 连接总数", connect_number_total)
fmt.Println("测试时间", *connect_command_channel)

connect_filesize_total :=0
for i := range  connect_filesize_channel{
        connect_filesize_total += i
                     }
fmt.Println(" 流量总数 ", connect_filesize_total)

}
