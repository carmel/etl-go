# 功能介绍
- 这是一个用golang编写的数据导入和导出的小工具。它以多线程运行，可将excel表格数据(*.xslx)导入到数据库(mysql)，也可通过在配置文件中写sql语句而将查询结果导出到Excel。

- 导入还是导出根据命令行选项参数m来指定，i为导入，e为导出。

- 导入时需通过选项参数p来指定excel的路径

- 注意Excel的格式：页签即表名，表格第一行对应了各个字段的名字。

- 需要自动生成的字段使用参数g指定（多个字段以逗号分隔，生成库为`github.com/rs/xid`）

- 例如将`name,gender,org`三个字段导入到`student`表中: > ./etl -p "./demo.xlsx" -m i -g id

# 下载
下载请前往[Releases页面](https: //github.com/carmel/etl-go/releases)  

# 配置示例
* mysql  
  ```yml
  DiverName: "root: root@tcp(localhost: 3306)/recovery?charset=utf8mb4"
  DB: mysql
  LimitChan: 15
  SQL:
    - SELECT rq 日期, bjmc 班级, if(dmlbm='04', '缺勤', '正常') 考勤, xh 学号,xm 姓名 FROM performance order by bjmc,rq
  ```
* sqlite3  
  ```yml
  DiverName: "file:D:/Spiritual/recovery.db?cache=shared"
  DB: sqlite3
  LimitChan: 15
  SQL:
    - SELECT rq 日期, bjmc 班级, if(dmlbm='04', '缺勤', '正常') 考勤, xh 学号,xm 姓名 FROM performance order by bjmc,rq
  ```
* postgres
  ```yml
  DiverName: "postgres://postgres:tygspg2017@47.104.106.121:5432/sun_dev?sslmode=disable"
  DB: postgres
  LimitChan: 15
  SQL:
    - SELECT rq 日期, bjmc 班级, if(dmlbm='04', '缺勤', '正常') 考勤, xh 学号,xm 姓名 FROM performance order by bjmc,rq
  ```