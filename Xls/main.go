package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tealeg/xlsx"
)

var tableWord = "===="
var arrayWord = "++++"
var lineWord = "----"
var noneWord = "-1"
var nTypeString = 1
var nTypeFile = 2
var nTypeOther = 3

var strTypeString = "(string)"
var strTypeFile = "(file)"
var strID = "id"

var outDir = "dataIsland"
var strSheetIsland = "island"
var strSheetBook = "book"

//var outBaseDir = "base"

type oneColumn struct {
	name  string
	nType int
}

func newoneColumn(str string, n int) oneColumn {
	m := oneColumn{}
	m.name = str
	m.nType = n
	return m
}

type oneValue struct {
	valus []string
}

func (m oneValue) toString() string {
	str := ""
	for _, v := range m.valus {
		if str == "" {
			str += fmt.Sprintf("%v", v)
		} else {
			str += fmt.Sprintf("%v%v", arrayWord, v)
		}
	}
	return str
}

func (m oneValue) toCheckFiles(rPath map[string]string) error {
	var err error
	for _, v := range m.valus {
		if v != noneWord {
			_, ok := rPath[v]
			if !ok {
				logPrint("******** error ***** file:%v is not found ********", v)
				err = errors.New("error")
			}
		}
	}
	return err
}

func (m oneValue) toCheckId(rPath map[string]string) error {
	var err error
	for _, v := range m.valus {
		if v != noneWord {

			_, ok := rPath[v]

			if !ok {
				logPrint("******** error ***** id:%v is not found ********", v)
				err = errors.New("error")
			}
		}
	}
	return err
}

type oneRow struct {
	valus []oneValue
}

func newoneRow() oneRow {
	m := oneRow{}
	m.valus = make([]oneValue, 0, 100)
	return m
}

func (m oneRow) toString() string {
	str := ""
	str += lineWord
	var index int
	for _, v := range m.valus {
		if index == 0 {
			str += fmt.Sprintf("%v", v.toString())
		} else {
			str += fmt.Sprintf("%v%v", tableWord, v.toString())
		}
		index++
	}

	return str
}

type oneSheet struct {
	columns []oneColumn
	rows    []oneRow
}

func newSheet() oneSheet {
	m := oneSheet{}
	m.columns = make([]oneColumn, 0, 100)
	m.rows = make([]oneRow, 0, 100)
	return m
}

func (m oneSheet) toString() string {
	str := ""
	for _, v := range m.columns {
		if str == "" {
			str += fmt.Sprintf("%v", v.name)
		} else {
			str += fmt.Sprintf("%v%v", tableWord, v.name)
		}

	}
	for _, v := range m.rows {
		str += v.toString()
	}
	return str
}

func (m oneSheet) toCheckFiles(rPath map[string]string) error {
	var err error
	for i, c := range m.columns {
		if c.nType == nTypeFile {
			for _, v := range m.rows {
				errV := v.valus[i].toCheckFiles(rPath)
				if err == nil && errV != nil {
					err = errV
				}
			}
		}
	}
	return err
}

func (m oneSheet) toCheckId(rPath map[string]string) error {
	var err error

	for i, c := range m.columns {

		if c.nType == nTypeOther {

			for _, v := range m.rows {
				//logPrint("%v,%v", len(v.valus), i)
				errV := v.valus[i].toCheckId(rPath)
				if err == nil && errV != nil {
					err = errV
				}
			}
		}
	}
	return err
}

//所有表
var allSheets = make(map[string]oneSheet)

//所有的ID
var allId = make(map[string]string)

var rFileMap = make(map[string]string) //

func main() {

	if getAllFilelist("data/files") != nil {
		return
	}

	//遍历所有xls
	if loadAllXls() != nil {
		return
	}
	//检查文件
	if checkFileId() != nil {
		return
	}
	//保存
	if SaveFilesV2() != nil {
		return
	}
	logPrint("真棒没错")
	return
}

//log输出
func logPrint(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	fmt.Println(str)
}

//获取文件目录找xls用
func getFilelist(path string, exetend string) (rPath []string) {
	rPath = make([]string, 0, 100)
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if strings.Contains(path, exetend) && !strings.Contains(path, "~") {
			rPath = append(rPath, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
	return rPath
}

//获取文件目录
func getAllFilelist(path string) error {
	logPrint("load all files")
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		rFileMap[f.Name()] = path
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
		return err
	}
	return nil
}

func loadAllXls() error {
	logPrint("load all xls")
	rPathXls := getFilelist("data", ".xlsx")
	for i := 0; i < len(rPathXls); i++ {
		//读取xls文件
		xlFile, err := xlsx.OpenFile(rPathXls[i])
		if err != nil {
			logPrint("load xls :%v ERROR", rPathXls[i])
			return errors.New("err")
		}
		for _, sheet := range xlFile.Sheets {
			tSheet := newSheet()

			for nRow, row := range sheet.Rows {
				tRow := newoneRow()

				for nCell, cell := range row.Cells {
					str, err := cell.String()

					if err != nil {
						logPrint("load xls :%v ,sheet:%v,row:%v,cell:%v ERROR 1", rPathXls[i], sheet.Name, nRow, nCell)
						return errors.New("err")
					}
					if strings.Contains(str, tableWord) || strings.Contains(str, arrayWord) || strings.Contains(str, lineWord) || str == noneWord {
						logPrint("load xls :%v ,sheet:%v,row:%v,cell:%v ERROR,it has %v,%v,%v or %v", rPathXls[i], sheet.Name, nRow, nCell, tableWord, arrayWord, lineWord, noneWord)
						return errors.New("err")
					}
					if nRow == 0 {
						//第一行
						if str == "" {
							logPrint("load xls :%v ,sheet:%v,row:%v,cell:%v ERROR 2", rPathXls[i], sheet.Name, nRow, nCell)
							return errors.New("err")
						}
						if nCell == 0 && str != (strID+strTypeString) {
							//检测是否id开头
							logPrint("load xls :%v ,sheet:%v,row:%v,cell:%v must be id", rPathXls[i], sheet.Name, nRow, nCell)
							return errors.New("err")
						}
						if strings.Contains(str, strTypeString) {
							//是string啥都不做
							str1 := strings.Replace(str, strTypeString, "", -1)
							if str1 == "" {
								logPrint("load xls :%v ,sheet:%v,row:%v,cell:%v ERROR", rPathXls[i], sheet.Name, nRow, nCell)
								return errors.New("err")
							}
							tSheet.columns = append(tSheet.columns, newoneColumn(str1, nTypeString))
						} else if strings.Contains(str, strTypeFile) {
							//是file会找文件
							str1 := strings.Replace(str, strTypeFile, "", -1)
							if str1 == "" {
								logPrint("load xls :%v ,sheet:%v,row:%v,cell:%v ERROR", rPathXls[i], sheet.Name, nRow, nCell)
								return errors.New("err")
							}
							tSheet.columns = append(tSheet.columns, newoneColumn(str1, nTypeFile))
						} else {
							//外联ID，回去对应的表里找ID
							if str == "" {
								logPrint("load xls :%v ,sheet:%v,row:%v,cell:%v ERROR", rPathXls[i], sheet.Name, nRow, nCell)
								return errors.New("err")
							}

							tSheet.columns = append(tSheet.columns, newoneColumn(str, nTypeOther))
						}
					} else {
						if len(tSheet.columns) <= nCell {
							logPrint("load xls :%v ,sheet:%v,row:%v,cell:%v ERROR", rPathXls[i], sheet.Name, nRow, nCell)
							return errors.New("err")
						}
						if str == "" && nCell == 0 {
							logPrint("load xls :%v ,sheet:%v,row:%v,cell:%v can not be null", rPathXls[i], sheet.Name, nRow, nCell)
							return errors.New("err")
						}
						if str == "" {
							str = noneWord
						}
						tRow.valus = append(tRow.valus, oneValue{strings.Split(str, "|")})
						if nCell == 0 {
							//是否已经有了，已经有了就滚蛋
							_, ok := allId[str]
							if ok {
								logPrint("load xls :%v ,sheet:%v,row:%v,cell:%v the id is repeat", rPathXls[i], sheet.Name, nRow, nCell)
								return errors.New("err")
							}
							allId[str] = sheet.Name
						}
						//每一个
						//tSheet.columns[nCell].valus = append(tSheet.columns[nCell].valus, oneColumnValue{strings.Split(str, "|")})
					}

					//logPrint(str)
				}

				if len(tRow.valus) <= 0 && nRow != 0 {
					logPrint("load xls :%v ,sheet:%v,row:%v can not be null", rPathXls[i], sheet.Name, nRow)
					return errors.New("err")
				}

				for len(tRow.valus) != len(tSheet.columns) && nRow != 0 {
					str := noneWord
					tRow.valus = append(tRow.valus, oneValue{strings.Split(str, "|")})

				}

				if nRow != 0 {
					tSheet.rows = append(tSheet.rows, tRow)
				}
			}
			if len(tSheet.columns) <= 0 {
				logPrint("******** waring ***** %v:%v is none ********", rPathXls[i], sheet.Name)
			} else if len(tSheet.rows) <= 0 {
				logPrint("******** waring ***** %v:%v is empty ********", rPathXls[i], sheet.Name)
			} else {
				allSheets[sheet.Name] = tSheet
			}
		}
	}
	return nil
}

func checkFileId() error {
	var err error
	logPrint("check files ")
	for _, v := range allSheets {
		errv := v.toCheckFiles(rFileMap)
		if errv != nil && err == nil {
			err = errv
		}
	}

	//检查ID
	logPrint("check ID ")
	for _, v := range allSheets {
		//logPrint("%v", s)
		errv := v.toCheckId(allId)
		if errv != nil && err == nil {
			err = errv
		}
	}
	return err
}

func SaveFiles() error {
	logPrint("save files ")
	var errR = errors.New("fail")
	//保存文件
	err := os.RemoveAll(outDir)
	if err != nil {
		logPrint("remove dir %v fail", outDir)
		return errR
	}
	err = os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		logPrint("create dir %v fail", outDir)
		return errR
	}
	//循环创建文件夹
	//创建所有的岛的
	//创建所有的book
	sheetIsland, ok := allSheets[strSheetIsland]
	if !ok {
		logPrint("I can not find sheet %v ", strSheetIsland)
		return errR
	}
	for _, v := range sheetIsland.rows {
		//遍历
		srtIsLandDir := outDir + "/" + v.valus[0].valus[0] + "/"
		err = os.MkdirAll(srtIsLandDir, os.ModePerm)
		if err != nil {
			logPrint("create dir %v fail", srtIsLandDir)
			return errR
		}
		srtIsLandBaseDir := srtIsLandDir
		// srtIsLandBaseDir := srtIsLandDir + "_" + outBaseDir + "/"
		// err = os.MkdirAll(srtIsLandBaseDir, os.ModePerm)
		// if err != nil {
		// 	logPrint("create dir %v fail", srtIsLandBaseDir)
		// 	return errR
		// }
		// 循环每一咧，找是ID的，然后拷贝
		for i, vC := range sheetIsland.columns {
			if vC.nType == nTypeFile && i != 0 {
				for _, vFileName := range v.valus[i].valus {
					//拷贝文件
					if vFileName != noneWord {
						fileSrc, ok := rFileMap[vFileName]
						if !ok {
							logPrint("******** error ***** file:%v is not found ******** ", v)
							return errors.New("error")
						}
						//文件是否在
						_, err = CopyFileByName(fileSrc, srtIsLandBaseDir+vFileName)
						if err != nil {
							logPrint("%v copy %v fail ", fileSrc, srtIsLandBaseDir+vFileName)
							return errors.New("error")
						}
					}

				}
			}
			if vC.nType == nTypeOther && i != 0 {
				if vC.name == strSheetBook {
					for _, vBook := range v.valus[i].valus {
						//先创建bookdir
						//再拷贝
						//srtIsLandBookDir := srtIsLandDir + "_" + vBook + "/"
						srtIsLandBookDir := outDir + "/" + vBook + "/"
						err = os.MkdirAll(srtIsLandBookDir, os.ModePerm)
						if err != nil {
							logPrint("create book dir %v to %v fail", vBook, srtIsLandBookDir)
							return errR
						}
						err = copyFile(vBook, srtIsLandBookDir)
						if err != nil {
							logPrint("copy book file fail %v to %v fail", vBook, srtIsLandBookDir)
							return errR
						}
					}
				} else {
					for _, vOther := range v.valus[i].valus {
						//先创建bookdir
						//再拷贝
						if vOther != noneWord {
							err = copyFile(vOther, srtIsLandBaseDir)
							if err != nil {
								logPrint("copy file fail %v to %v fail", vOther, srtIsLandBaseDir)
								return errR
							}
						}

					}
				}
			}
		}
	}
	srtConfigDir := outDir + "/dataConfig/"
	err = os.MkdirAll(srtConfigDir, os.ModePerm)
	if err != nil {
		logPrint("create dir %v fail", srtConfigDir)
		return errR
	}
	//把一个BOOK的文件都考到一个book下面？
	var strAll string
	for k, v := range allSheets {
		//
		if strAll == "" {
			strAll += k
		} else {
			strAll += (tableWord + k)
		}

		if SaveStrToFile(v.toString(), srtConfigDir+k+".txt") != nil {
			logPrint("save config  Fail %v", k)
			return errors.New("err")
		}
	}
	if SaveStrToFile(strAll, srtConfigDir+"dataConfig.txt") != nil {
		logPrint("save config  Fail %v", strAll)
		return errors.New("err")
	}
	return nil
}

func SaveFilesV2() error {
	logPrint("save files ")
	var errR = errors.New("fail")
	//保存文件
	err := os.RemoveAll(outDir)
	if err != nil {
		logPrint("remove dir %v fail", outDir)
		return errR
	}
	err = os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		logPrint("create dir %v fail", outDir)
		return errR
	}

	for k, v := range allSheets {
		//所有的岛
		strSheetDir := outDir + "/" + k
		// err = os.MkdirAll(strSheetDir, os.ModePerm)
		// if err != nil {
		// 	logPrint("create dir %v fail", strSheetDir)
		// 	return errR
		// }
		for _, vLine := range v.rows {
			srtLineDir := strSheetDir + vLine.valus[0].valus[0] + "/"
			bCreatDir := false
			for i, vC := range v.columns {
				if vC.nType == nTypeFile && i != 0 {
					for _, vFileName := range vLine.valus[i].valus {
						//拷贝文件
						if vFileName != noneWord {
							fileSrc, ok := rFileMap[vFileName]
							if !ok {
								logPrint("******** error ***** file:%v is not found ******** ", v)
								return errors.New("error")
							}
							if !bCreatDir {
								err = os.MkdirAll(srtLineDir, os.ModePerm)
								if err != nil {
									logPrint("create dir %v fail", srtLineDir)
									return errR
								}
								bCreatDir = true
							}

							//文件是否在
							_, err = CopyFileByName(fileSrc, srtLineDir+vFileName)
							if err != nil {
								logPrint("%v copy %v fail ", fileSrc, srtLineDir+vFileName)
								return errors.New("error")
							}
						}
					}
				}
			}
		}
	}

	srtConfigDir := outDir + "/dataConfig/"
	err = os.MkdirAll(srtConfigDir, os.ModePerm)
	if err != nil {
		logPrint("create dir %v fail", srtConfigDir)
		return errR
	}
	//把一个BOOK的文件都考到一个book下面？
	var strAll string
	for k, v := range allSheets {
		//
		if strAll == "" {
			strAll += k
		} else {
			strAll += (tableWord + k)
		}

		if SaveStrToFile(v.toString(), srtConfigDir+k+".txt") != nil {
			logPrint("save config  Fail %v", k)
			return errors.New("err")
		}
	}
	if SaveStrToFile(strAll, srtConfigDir+"dataConfig.txt") != nil {
		logPrint("save config  Fail %v", strAll)
		return errors.New("err")
	}
	return nil
}

func copyFile(id string, path string) error {
	var copyList = make(map[string]bool)
	return copyFileList(id, path, copyList)
}

func copyFileList(id string, path string, copyList map[string]bool) error {
	//var copyList = make(map[string]bool)
	//在所有的里面找到这货
	if id == noneWord {
		return nil
	}
	_, ok := copyList[id]
	if ok {
		return nil
	}
	//拷贝

	sheetName, ok := allId[id]
	if !ok {
		logPrint("id: %v is not define but used , i have check it before see this tell me ", id)
		return errors.New("err")
	}
	//v
	vSheet, ok := allSheets[sheetName]
	if !ok {
		logPrint("sheet: %v is not define but used , i have check it before see this tell me ", sheetName)
		return errors.New("err")
	}
	//
	for _, v := range vSheet.rows {
		//遍历
		if id == v.valus[0].valus[0] {
			//需要拷贝
			// 循环每一咧，找是ID（nTypeOther） 递归，file的拷贝,
			for i, vC := range vSheet.columns {
				if vC.nType == nTypeOther && i != 0 {
					for _, vOther := range v.valus[i].valus {
						err := copyFileList(vOther, path, copyList)
						if err != nil {
							return err
						}
					}
				}
				if vC.nType == nTypeFile && i != 0 {
					for _, vFileName := range v.valus[i].valus {
						//拷贝文件
						if vFileName != noneWord {
							fileSrc, ok := rFileMap[vFileName]
							if !ok {
								logPrint("******** error ***** file:%v is not found ******** ", v)
								return errors.New("error")
							}
							//文件是否在
							_, err := CopyFileByName(fileSrc, path+vFileName)
							if err != nil {
								logPrint("%v copy %v fail ", fileSrc, path+vFileName)
								return errors.New("error")
							}
						}

					}
				}
			}
			return nil
		}
	}
	logPrint("ID: %v is not found in sheet ? why ! tell me", id)
	return errors.New("err")
}

func CopyFileByName(src, des string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err)
	}
	defer srcFile.Close()

	desFile, err := os.Create(des)
	if err != nil {
		fmt.Println(err)
	}
	defer desFile.Close()

	return io.Copy(desFile, srcFile)
}

func SaveStrToFile(src, des string) error {
	dstFile, err := os.Create(des)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer dstFile.Close()
	dstFile.WriteString(src)
	return nil
}
