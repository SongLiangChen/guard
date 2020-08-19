package main

import (
	"encoding/json"
	"fmt"
	"guard/cmd"
	"guard/config"
	"guard/model"
	"io/ioutil"
	"os"
	"time"

	"github.com/SongLiangChen/common"
	"github.com/judwhite/go-svc/svc"
)

var (
	alterTempleProcessShutDown = `进程异常
服务器[%v]
进程[%v]
多次尝试重启失败，请立刻排查`

	alterTempleProcessRestart = `进程异常
服务器[%v]
进程[%v]
重启成功`

	alterTempleCPU = `进程异常
服务器[%v]
进程[%v]
当前CPU使用率[%v]大于阈值百分之[%v]`

	alterTempleMEM = `进程异常
服务器[%v]
进程[%v]
当前MEM使用率[%v]大于阈值百分之[%v]`
)

type Node struct {
	Ps     []*model.Process
	exit   chan bool
	exited chan bool
}

func (n *Node) Init(env svc.Environment) error {
	if err := config.InitConfig(); err != nil {
		return err
	}

	println(config.GetConfig().Section("system").Key("alert").MustInt(0))

	n.Ps = make([]*model.Process, 0)
	buf, err := ioutil.ReadFile("./data/json/process.json")
	if err != nil {
		return err
	}

	if err = json.Unmarshal(buf, &n.Ps); err != nil {
		return err
	}

	n.exit = make(chan bool)
	n.exited = make(chan bool)

	return nil
}

func (n *Node) Start() error {
	go goOnPatrol(n.Ps, n.exit, n.exited)
	return nil
}

func (n *Node) Stop() error {
	close(n.exit)
	<-n.exited
	return nil
}

func main() {
	if err := svc.Run(&Node{}); err != nil {
		println(err.Error())
	}
}

func goOnPatrol(ps []*model.Process, exit, exited chan bool) {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	cnt := make(map[string]int)
	hostname := os.Getenv("RUNMODE")

	for {
		select {
		case <-exit:
			goto FINISH

		case <-ticker.C:
			for _, p := range ps {
				s := cmd.GetProcessStat(p.Name)
				if s == nil {
					println(p.Name, "获取进程信息失败，请检查配置")
					continue
				}
				if s.Pid == "" {
					cnt[p.Name]++
					if cnt[p.Name] >= 3 {
						DingDingAlert(fmt.Sprintf(alterTempleProcessShutDown, hostname, p.Name))
					}

					if p.RestartSh != "" {
						println(p.Name, "进程未启动，尝试启动...")
						_ = cmd.RestartProcess(p.RestartSh)

						time.Sleep(time.Second)
						if ts := cmd.GetProcessStat(p.Name); ts != nil && ts.Pid != "" {
							println(p.Name, "启动成功")
						} else {
							cnt[p.Name]++
							println(p.Name, "启动失败")
						}
						continue
					}

					continue
				}
				if cnt[p.Name] >= 3 {
					DingDingAlert(fmt.Sprintf(alterTempleProcessRestart, hostname, p.Name))
				}
				cnt[p.Name] = 0

				if p.MaxCPU < s.CPU {
					DingDingAlert(fmt.Sprintf(alterTempleCPU, hostname, p.Name, s.CPU, p.MaxCPU))
				}
				if p.MaxMem < s.MEM {
					DingDingAlert(fmt.Sprintf(alterTempleMEM, hostname, p.Name, s.MEM, p.MaxMem))
				}
			}
		}
	}

FINISH:
	exited <- true
}

func DingDingAlert(content string) {
	alert := config.GetConfig().Section("system").Key("alert").MustInt(0)
	if alert == 0 {
		return
	}

	var urlStr = `https://oapi.dingtalk.com/robot/send?access_token=5226f0309f4431b35ffaa0939e6d11665243b21de8650aefcfdcf3f5932196be`

	data := make(map[string]interface{})
	data["msgtype"] = "text"
	data["text"] = map[string]string{
		"content": content,
	}
	buf, _ := json.Marshal(data)

	agent := common.NewHttpAgent()
	_, buf, err := agent.Post(urlStr).ContentType(common.TypeJSON).Timeout(3 * time.Second).SendData(buf).End()
	if err != nil {
		println(err.Error())
	} else {
		println(string(buf))
	}
}

func FeishuAlert(content string) {
	alert := config.GetConfig().Section("system").Key("alert").MustInt(0)
	if alert == 0 {
		return
	}

	var urlStr = ``

	data := make(map[string]string)
	data["title"] = "服务进程异常"
	data["text"] = content
	buf, _ := json.Marshal(data)

	agent := common.NewHttpAgent()
	_, _, _ = agent.Post(urlStr).ContentType(common.TypeJSON).Timeout(3 * time.Second).SendData(buf).End()
}
