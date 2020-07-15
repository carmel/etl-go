package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
)

var (
	db *sqlx.DB
	// lock sync.Mutex
	conf struct {
		DiverName string   `yaml:"DiverName"`
		DB        string   `yaml:"DB"`
		LimitChan int      `yaml:"LimitChan"`
		SQL       []string `yaml:"SQL"`
	}
	EXCEL_COL = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
)

func init() {
	c, err := ioutil.ReadFile("conf.yml")
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(c, &conf)
	if err != nil {
		log.Fatalln(err)
	}
	//	log.SetPrefix("[Info]")
	//	log.SetFlags(log.LstdFlags | log.LUTC)
	// conf["LimitChan"].(json.Number).Int64()
	db, err = sqlx.Connect(conf.DB, conf.DiverName)
	if err != nil {
		panic(err)
	}
}

func main() {
	defer db.Close()
	runtime.GOMAXPROCS(runtime.NumCPU())

	path := flag.String("p", "", "excel文件路径")
	mode := flag.String("m", "i", "i/e, 导入/导出")
	// title := flag.String("t", "未命名.xlsx", "导出的excel名称")
	flag.Parse()
	if *mode == "i" { // 导入excel数据
		if *path == "" {
			log.Fatalln("未指定要导入的excel路径")
		}
		xlsx, err := excelize.OpenFile(*path)
		if err != nil {
			log.Println(`open excel`, err)
		}
		// sn := xlsx.GetSheetName(1)
		// rows, _ := xlsx.GetRows(sn)

		// xlFile, err := xlsx.OpenFile(conf.Path)
		if err != nil {
			panic(err)
		}
		for _, sheet := range xlsx.WorkBook.Sheets.Sheet {
			if rows, err := xlsx.GetRows(sheet.Name); err == nil && len(rows) != 0 {
				var buffer bytes.Buffer
				buffer.WriteString(`INSERT INTO `)
				buffer.WriteString(sheet.Name)
				buffer.WriteString(`(`)

				buffer.WriteString(strings.Join(rows[0], `,`))
				buffer.WriteString(`)VALUES(`)
				l := len(rows[0])
				buffer.WriteString(strings.TrimSuffix(strings.Repeat(`?,`, l), `,`))
				buffer.WriteString(`)`)

				query := buffer.String()
				limitChan := make(chan bool, conf.LimitChan)
				wg := sync.WaitGroup{}
				for i, row := range rows[1:] {
					fmt.Println(`------正在处理第`, i, `行`)
					limitChan <- true
					wg.Add(1)
					go func(i int, r []string) {
						defer func() {
							wg.Done()
							<-limitChan
							if err := recover(); err != nil {
								// logger.Printf("第%d行: %+v, 错误: %v", i+1, r, err)
								fmt.Printf("第%d行: %+v, 错误: %v\n", i+1, r, err)
							}
						}()
						var args []interface{}
						for _, v := range r {
							args = append(args, v)
						}
						db.MustExec(query, args...)
					}(i+1, row)
				}
				wg.Wait()
				// ticker := time.NewTicker(2 * time.Second) //定时器,每2秒钟执行一次
				// for c := range ticker.C {

				// 	if int32(l) == t {
				// 		close(limitChan)
				// 		fmt.Printf("------%v: 成功导入%d行，5秒后自动关闭", c, t)
				// 		time.Sleep(time.Second * 5)
				// 		ticker.Stop()
				// 		break
				// 	}
				// }
			}
		}
	} else { // 导出指定sql数据

		var rows *sqlx.Rows
		var err error
		var index int
		var title []string
		xlsx := excelize.NewFile()
		for i, sql := range conf.SQL {
			index = 0
			rows, err = db.Queryx(sql)
			if err != nil {
				panic(err)
			}

			if i != 0 {
				xlsx.NewSheet(fmt.Sprintf("%s%d", "Sheet", i+1))
			}

			for rows.Next() {
				index++
				if index == 1 {
					title, _ = rows.Columns()
					for n, v := range title {
						if err = xlsx.SetCellValue(fmt.Sprintf("%s%d", "Sheet", i+1), fmt.Sprintf("%s%d", EXCEL_COL[n], 1), v); err != nil {
							log.Println(`SetCellValue`, err)
						}
					}
				}
				rs, _ := rows.SliceScan()
				for n, v := range rs {
					// fmt.Println(fmt.Sprintf("%s%d", "Sheet", i+1), fmt.Sprintf("%s%d", EXCEL_COL[n+1], index))
					if err = xlsx.SetCellValue(fmt.Sprintf("%s%d", "Sheet", i+1), fmt.Sprintf("%s%d", EXCEL_COL[n], index+1), v); err != nil {
						log.Println(`SetCellValue`, err)
					}
				}
			}
		}

		err = xlsx.SaveAs(fmt.Sprintf("%s.xlsx", strconv.FormatInt(time.Now().Unix(), 10)))
		if err != nil {
			panic(err)
		}

	}
	/*日志写入文件
	f, err := os.Create("imple.log")
	if err != nil {
		log.Fatalf("file open error : %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	*/

	// file, err := os.Create("imple.log")
	// defer file.Close()
	// if nil != err {
	// 	panic(err)
	// }

	// logger := log.New(file, "err_", log.Ldate|log.Ltime|log.Lshortfile)
	//Flags返回Logger的输出选项
	// logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

}
