package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
)

func main() {
	flag.PrintDefaults()
	fileExt := flag.String("e", "jpeg", "Расширение файлов для обработки")
	dirFrom := flag.String("f", "from", "Папка с исходными файлам")
	dirTo := flag.String("t", "to", "Папка с результирующими файлами")
	firstByte := flag.Int("fb", 4, "Первый байт для вырезки")
	dataLength := flag.Int("l", 5, "Длина данных для вырезки")
	combineFiles := flag.Int("c", 0, "Размер объединямых файлов в кБ, если 0 - не объединять ")
	flag.Parse()

	fmt.Printf("* Расширение файлов для обработки\t[*.%s]\n", *fileExt)
	fmt.Printf("* Папка с исходными файлам\t\t[%s]\n", *dirFrom)
	fmt.Printf("* Папка с результирующими файлами\t[%s]\n", *dirTo)
	fmt.Printf("* Первый байт для вырезки\t\t[%d]\n", *firstByte)
	fmt.Printf("* Длина данных для вырезки\t\t[%d]\n", *dataLength)
	fmt.Printf("* Размер объединямых файлов в кБ\t[%d]\n", *combineFiles)

	fileList := fileListFn(dirFrom, fileExt)
	if len(fileList) == 0 {
		log.Fatal("\n\nERROR: Не найдено ни одного файла для обработки")
		return
	}
	fmt.Printf("\n\n********* Для обработки найдено %d файлов *********\n", len(fileList))

	processFileList(fileList, *dirFrom, *dirTo, *firstByte, *dataLength)

	if *combineFiles > 0 {
		combineFilesProcess()
	}

}

func combineFilesProcess() {

}

func processFileList(list []string, from string, to string, firstByte int, length int) {
	makeOutputDir(to)
	var wg sync.WaitGroup
	for _, s := range list {
		wg.Add(1)
		go processOneFile(from, s, firstByte, length, to, &wg)
	}
	wg.Wait()
	fmt.Printf("\n\n********* Обработка закончена *********\n")

}

func processOneFile(from string, s string, firstByte int, length int, to string, wg *sync.WaitGroup) {
	defer wg.Done()
	bytes, err := os.ReadFile(fmt.Sprintf("./%s/%s", from, s))
	if err != nil {
		fmt.Printf("Ошибка %s при обработке файла %s", err, s)
		return
	}
	if len(bytes) < firstByte+length {
		fmt.Printf("Ошибка при обработке файла %s, не хватает данных для вырезки\n", s)
		return
	}
	data := bytes[firstByte : firstByte+length]
	writeBytes(data, s, to)
	fmt.Printf("\nfile %s ready", s)
}

func writeBytes(data []byte, s string, to string) {
	_, err := os.Create(fmt.Sprintf("./%s/%s", to, s))
	if err != nil {
		log.Fatalf("Ошибка создания файлf! %v", err)
		return
	}
	err = os.WriteFile(fmt.Sprintf("./%s/%s", to, s), data, 0777)
	if err != nil {
		log.Fatalf("Ошибка записи в файл! %v", err)
		return
	}
}

func makeOutputDir(to string) bool {
	if _, err := os.Stat(to); os.IsNotExist(err) {
		err := os.Mkdir(to, 0777)
		if err != nil {
			log.Fatalf("Ошибка создания выходной директории! %v", err)
			return true
		}
	}
	return false
}

func fileListFn(from *string, ext *string) []string {
	f, err := os.Open(*from)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var fileList []string
	for _, v := range files {
		if !v.IsDir() && strings.Contains(v.Name(), fmt.Sprintf(".%s", *ext)) {
			fileList = append(fileList, v.Name())
		}
	}
	sort.Strings(fileList)
	fmt.Println(fileList)
	return fileList
}
