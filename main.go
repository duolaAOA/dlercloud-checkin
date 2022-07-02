package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	email := flag.String("u", "", "login email")
	password := flag.String("p", "", "login password")
	info := flag.Bool("i", false, "get user information")
	sweepstakes := flag.Bool("s", false, "get free traffic")
	helpPtr := flag.Bool("h", false, "dlercloud checkin help")
	flag.Usage = func() {
		_, _ = fmt.Fprint(os.Stderr,
			`Usage: ./dler  OPTIONS [arg...]
         dler [ -u dler -p passwdxxx | -help ]
Options:
  -u               login email
  -p               login password
  -i               get user information
  -c               get free traffic
  -h               Print usage
`)
	}
	flag.Parse()

	// show help
	if *helpPtr {
		flag.Usage()
		os.Exit(0)
	}

	if *email == "" || *password == "" {
		fmt.Printf("incorrect username or password...\r\n")
		os.Exit(1)
	}

	dler := NewClient(*email, *password)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if *info {
		userInfo, err := dler.GetUserInfo(ctx)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		str := `
    今天已用流量:		%s
    已用流量: 			%s
    剩余流量: 			%s
    套餐每月总计流量:		%s
    到期时间: 			%s
    套餐: 			%s
    账号余额: 			%s
    积分: 			%s`
		inf := fmt.Sprintf(str, userInfo.TodayUsed, userInfo.Used, userInfo.Unused, userInfo.Traffic, userInfo.PlanTime, userInfo.Plan, userInfo.Money, userInfo.Integral)
		fmt.Println(inf)
	}

	if *sweepstakes {
		ck, err := dler.TryToCheckin(ctx)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		str := `
    本次抽中流量:		%s
    今天已用流量:		%s
    已用流量: 			%s
    剩余流量: 			%s
    套餐每月总计流量:		%s`
		inf := fmt.Sprintf(str, ck.Checkin, ck.TodayUsed, ck.Used, ck.Unused, ck.Traffic)
		fmt.Println(inf)
	}
}
