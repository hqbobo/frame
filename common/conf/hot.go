package conf

import (
	"io/ioutil"

	"github.com/fsnotify/fsnotify"
	"github.com/hqbobo/frame/common/log"
)

//HotConfig 支持热更的配置
type HotConfig struct {
	path    string
	watcher *fsnotify.Watcher
	load    func(config string)
	close   chan int
}

//Exit 退出
func (hc *HotConfig) Exit() {
	close(hc.close)
}

func (hc *HotConfig) watch() {
	defer func() {
		log.Debug("退出监控")
		hc.watcher.Close()
	}()

	for {
		select {
		case event, ok := <-hc.watcher.Events:
			if !ok {
				log.Info("return")
				return
			}
			if event.Op == fsnotify.Remove || event.Op == fsnotify.Write {
				data, err := ioutil.ReadFile(hc.path)
				if err != nil {
					log.Warn(err)
				}
				hc.load(string(data))
				hc.watcher.Add(hc.path)
			}

		case err, ok := <-hc.watcher.Errors:
			if !ok {
				log.Info("return")
				return
			}
			log.Error("error:", err)
		case <-hc.close:
			log.Info("退出配置")
		}
	}
}

//NewHotconfig 支持热更的配置
func NewHotconfig(path string, load func(config string)) (*HotConfig, error) {
	var err error
	conf := new(HotConfig)
	conf.path = path
	conf.load = load
	conf.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Warn(err)
	}
	load(string(data))
	conf.watcher.Add(path)
	go conf.watch()
	return conf, nil
}
