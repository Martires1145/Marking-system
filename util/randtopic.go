package util

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"sync"
	"time"
)

var startTime = viper.GetString("web.time")

type SnowFlakeIdWorker struct {

	// 开始时间戳
	twepoch int64

	// 机器ID所占的位数
	workerIdBits int64

	// 数据标识ID所占的位数
	dataCenterIdBits int64

	// 支持的最大机器ID
	maxWorkerId int64

	// 支持的最大机房 ID
	maxDataCenterId int64

	// 序列在ID中占的位数
	sequenceBits int64

	// 机器ID向左移位数
	workerIdShift int64

	// 机房ID向左移位数
	dataCenterIdShift int64

	// 时间截向左移位数
	timestampLeftShift int64

	// 生成序列的掩码最大值
	sequenceMask int64

	// 工作机器ID
	workerId int64

	// 机房ID
	dataCenterId int64

	/**
	 * 毫秒内序列
	 */
	sequence int64

	// 上次生成ID的时间戳
	lastTimestamp int64

	// 锁
	lock sync.Mutex
}

func (p *SnowFlakeIdWorker) init(dataCenterId int64, workerId int64) {
	t, _ := time.Parse("2006-01-02", startTime)
	// 开始时间戳；这里是2021-06-01
	p.twepoch = t.Unix()
	// 机器ID所占的位数
	p.workerIdBits = 5
	// 数据标识ID所占的位数
	p.dataCenterIdBits = 5
	// 支持的最大机器ID，最大是31
	p.maxWorkerId = -1 ^ (-1 << p.workerIdBits)
	// 支持的最大机房ID，最大是 31
	p.maxDataCenterId = -1 ^ (-1 << p.dataCenterIdBits)
	// 序列在ID中占的位数
	p.sequenceBits = 12
	// 机器ID向左移12位
	p.workerIdShift = p.sequenceBits
	// 机房ID向左移17位
	p.dataCenterIdShift = p.sequenceBits + p.workerIdBits
	// 时间截向左移22位
	p.timestampLeftShift = p.sequenceBits + p.workerIdBits + p.dataCenterIdBits
	// 生成序列的掩码最大值，最大为4095
	p.sequenceMask = -1 ^ (-1 << p.sequenceBits)

	if workerId > p.maxWorkerId || workerId < 0 {
		panic(errors.New(fmt.Sprintf("Worker ID can't be greater than %d or less than 0", p.maxWorkerId)))
	}
	if dataCenterId > p.maxDataCenterId || dataCenterId < 0 {
		panic(errors.New(fmt.Sprintf("DataCenter ID can't be greater than %d or less than 0", p.maxDataCenterId)))
	}

	p.workerId = workerId
	p.dataCenterId = dataCenterId
	// 毫秒内序列(0~4095)
	p.sequence = 0
	// 上次生成 ID 的时间戳
	p.lastTimestamp = -1
}

// 生成ID，注意此方法已经通过加锁来保证线程安全
func (p *SnowFlakeIdWorker) nextId() int64 {
	p.lock.Lock()
	defer p.lock.Unlock()

	timestamp := p.timeGen()
	// 如果当前时间小于上一次 ID 生成的时间戳，说明发生时钟回拨，为保证ID不重复抛出异常。
	if timestamp < p.lastTimestamp {
		panic(errors.New(fmt.Sprintf("Clock moved backwards. Refusing to generate id for %d milliseconds", p.lastTimestamp-timestamp)))
	}

	if p.lastTimestamp == timestamp {
		// 同一时间生成的，则序号+1
		p.sequence = (p.sequence + 1) & p.sequenceMask
		// 毫秒内序列溢出：超过最大值
		if p.sequence == 0 {
			// 阻塞到下一个毫秒，获得新的时间戳
			timestamp = p.tilNextMillis(p.lastTimestamp)
		}
	} else {
		// 时间戳改变，序列重置
		p.sequence = 0
	}
	// 保存本次的时间戳
	p.lastTimestamp = timestamp

	// 移位并通过或运算拼到一起
	return ((timestamp - p.twepoch) << p.timestampLeftShift) |
		(p.dataCenterId << p.dataCenterIdShift) |
		(p.workerId << p.workerIdShift) | p.sequence
}

func (p *SnowFlakeIdWorker) tilNextMillis(lastTimestamp int64) int64 {
	timestamp := p.timeGen()
	for timestamp <= lastTimestamp {
		timestamp = p.timeGen()
	}
	return timestamp
}

func (p *SnowFlakeIdWorker) timeGen() int64 {
	return time.Now().UnixNano() / 1e6
}

func GetClientID() string {
	idWorker := &SnowFlakeIdWorker{}
	idWorker.init(0, 1)
	return fmt.Sprintf("%d", idWorker.nextId())
}

func GetTopic() (topic1 string, topic2 string) {
	topic1 = GetClientID()
	time.Sleep(time.Millisecond * 10)
	topic2 = GetClientID()
	return
}
