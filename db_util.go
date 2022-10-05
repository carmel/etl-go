package main

// import (
// 	"log"
// 	"math"
// 	"reflect"
// 	"strings"
// 	"trainee/util"

// 	"github.com/jmoiron/sqlx"
// )

// func ListMap(rows *sqlx.Rows, call func(map[string]interface{}) (string, string)) (result []map[string]interface{}) {
// 	for rows.Next() {
// 		tmp := make(map[string]interface{})
// 		rows.MapScan(tmp)
// 		for k, encoded := range tmp {
// 			switch encoded.(type) {
// 			case []byte:
// 				tmp[k] = string(encoded.([]byte))
// 			}
// 		}
// 		if call != nil {
// 			key, res := call(tmp)
// 			tmp[key] = res
// 		}
// 		result = append(result, tmp)
// 	}
// 	return
// }

// func SaveModel(o interface{}, s bool) bool {
// 	e := reflect.TypeOf(o)
// 	v := reflect.ValueOf(o)
// 	tab := util.AntiCamelCase(e.Name())
// 	var f reflect.StructField
// 	var add []string
// 	var edit []string
// 	var fd string
// 	for i := 0; i < e.NumField(); i++ {
// 		f = e.Field(i)
// 		fd = util.AntiCamelCase(f.Name)
// 		if v.Field(i).Interface() != reflect.Zero(f.Type).Interface() {
// 			add = append(add, fd)
// 			edit = append(edit, fd+"=:"+fd)
// 		}
// 	}
// 	var err error
// 	if s { //新增操作
// 		_, err = DB.NamedExec(`INSERT INTO `+tab+`(`+strings.Join(add, ",")+`)VALUES(`+`:`+strings.Join(add, ",:")+`)`, o)

// 	} else {
// 		_, err = DB.NamedExec(`UPDATE `+tab+` SET `+strings.Join(edit, ",")+` WHERE id = :id`, o)
// 	}
// 	if err != nil {
// 		util.Log(err)
// 		return false
// 	}
// 	return true
// }

// func SaveMapV1(t string, m map[string]interface{}) bool {
// 	var err error
// 	var add []string
// 	var edit []string
// 	var f string
// 	for k, _ := range m {
// 		f = util.AntiCamelCase(k)
// 		add = append(add, f)
// 		edit = append(edit, f+"=:"+f)
// 	}

// 	if v, ok := m["Id"]; ok && v != "" {
// 		log.Println(`UPDATE ` + t + ` SET ` + strings.Join(edit, ",") + ` WHERE id = :id`)
// 		_, err = DB.NamedExec(`UPDATE `+t+` SET `+strings.Join(edit, ",")+` WHERE id = :id`, m)
// 	} else {
// 		m["Id"] = util.UUID()
// 		add = append(add, "Id")
// 		log.Println(`INSERT INTO ` + t + `(` + strings.Join(add, ",") + `)VALUES(` + `:` + strings.Join(add, ",:") + `)`)
// 		_, err = DB.NamedExec(`INSERT INTO `+t+`(`+strings.Join(add, ",")+`)VALUES(`+`:`+strings.Join(add, ",:")+`)`, m)
// 	}

// 	if err != nil {
// 		log.Println(err)
// 		return false
// 	}
// 	return true
// }

// func SaveMap(t string, m map[string]interface{}) string {
// 	var err error
// 	var add []string
// 	var edit []string
// 	for f, v := range m {
// 		var flag bool
// 		switch z := v.(type) { // 零值过滤
// 		case float32: // 注意float32与float64不可写在一起，因在case路由中，如果不能精准到单路线，v还是一个interface{}
// 			flag = math.Abs(float64(z)-0) < 0.0000001
// 		case float64:
// 			flag = math.Abs(z-0) < 0.0000001
// 		case int, int32, int64:
// 			flag = v == 0
// 		case string:
// 			flag = v == ""
// 		case nil:
// 			flag = true
// 			// util.Log(`type is`, v)
// 		}
// 		if flag {
// 			continue
// 		}
// 		add = append(add, f)
// 		if f != util.AppConfig.Id {
// 			edit = append(edit, f+"=:"+f)
// 		}
// 	}
// 	util.Log(`UPDATE `+t+` SET `+strings.Join(edit, ",")+` WHERE `+util.AppConfig.Id+` = :`+util.AppConfig.Id+``, m)
// 	if v, ok := m[util.AppConfig.Id]; ok && v != "" {
// 		// m["modify_time"] = time.Now().Unix()
// 		util.Log(`UPDATE `+t+` SET `+strings.Join(edit, ",")+` WHERE `+util.AppConfig.Id+` = :`+util.AppConfig.Id+``, m)
// 		_, err = DB.NamedExec(`UPDATE `+t+` SET `+strings.Join(edit, ",")+` WHERE `+util.AppConfig.Id+` = :`+util.AppConfig.Id+``, m)
// 		// DB.NamedExec(`INSERT INTO l_alteration(id, tab, o_val, n_val, modify_time, modifier)VALUES(:id, :tab, :o_val, :n_val, :modify_time, :modifier)`, m)
// 	} else {
// 		m[util.AppConfig.Id] = util.UUID()
// 		add = append(add, util.AppConfig.Id)
// 		util.Log(`INSERT INTO `+t+`(`+strings.Join(add, ",")+`)VALUES(`+`:`+strings.Join(add, ",:")+`)`, m)
// 		_, err = DB.NamedExec(`INSERT INTO `+t+`(`+strings.Join(add, ",")+`)VALUES(`+`:`+strings.Join(add, ",:")+`)`, m)
// 	}

// 	if err != nil {
// 		util.Log(t, `save:`, err)
// 		return ""
// 	}
// 	return m[util.AppConfig.Id].(string)
// }
