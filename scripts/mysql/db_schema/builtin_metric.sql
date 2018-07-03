use falcon_portal;
insert into metric(name, alias, value_type, unit) values
("cpu.idle", "CPU空闲", "float64", "%"),
("cpu.busy", "CPU繁忙", "float64", "%"),
("cpu.user", "用户CPU时间", "float64", "%"),
("cpu.nice", "低优先级用户模式所占用的CPU时间", "float64", "%"),
("cpu.system", "系统CPU时间", "float64", "%"),
("cpu.iowait", "IO等待时间", "float64", "%"),
("cpu.irq", "CPU硬中断时间", "float64", "%"),
("cpu.softirq", "CPU软中断时间", "float64", "%"),
("cpu.steal", "虚拟进程的等待时间", "float64", "%"),
("cpu.guest", "运行虚拟进程的CPU时间", "float64", "%"),
("cpu.switches", "CPU切换次数", "unit64", "%"),

("mem.memtotal", "内存总量", "uint64", "B"),
("mem.memused", "内存使用量", "uint64", "B"),
("mem.memfree", "内存剩余", "uint64", "B"),
("mem.swaptotal", "交换内存总量", "uint64", "B"),
("mem.swapused", "交换内存用量", "uint64", "B"),
("mem.swapfree", "交换内存剩余", "uint64", "B"),
("mem.memfree.percent", "内存剩余百分比", "float64", "%"),
("mem.memused.percent", "内存用量百分比", "float64", "%"),
("mem.swapfree.percent", "交换内存剩余百分比", "float64", "%"),
("mem.swapused.percent", "交换内存用量百分比", "float64", "%");
