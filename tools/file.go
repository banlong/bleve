package tools
import (
	"os"
	"log"
	"io"
	"mime/multipart"
	"container/list"
	"io/ioutil"
	"time"
	"strings"
	"strconv"
	"errors"
	"bufio"
	"encoding/xml"
	"fmt"
)

func SaveBinFile2Disk(srcFile []byte, destDir string, fileName string) error  {
	//Create new file
	fileAbsName := destDir + fileName

	// write the whole body at once
	err := ioutil.WriteFile(fileAbsName, srcFile, 0644)
	if err != nil {
		log.Println("-- create file failed:", err)
		return err
	}
	//log.Println("-- saved:", fileAbsName)
	return nil
}

func SaveMultipartFile2Disk(srcFile multipart.File, destDir string, fileName string) error  {
	if(srcFile == nil){
		return errors.New("-- cannot save an empty file")
	}

	//Create new file
	fileAbsName := destDir + fileName
	newFile, err := os.Create(fileAbsName)
	if err != nil {
		log.Println("-- create file failed:", err)
		return err
	}

	//Save binary content into new file
	numOfBytes, err := io.Copy(newFile, srcFile)
	if err != nil {
		log.Println("-- save file err:", err)
		return err
	}
	log.Printf("-- saved %s, %d bytes  \n", fileAbsName, numOfBytes)
	defer newFile.Close()
	return nil
}

// GET ALL FILE NAMES IN THE FOLDER
func GetFileNames(dir string)  *list.List{
	names := list.New()
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		names.PushBack(file.Name())
		//log.Println(file.Name())
	}
	return names
}

// DELETE ALL IN A FOLDER
func DeleteDir(dir string)  error{
	tries := 1
	var err error
	//log.Printf("-- cleaning temp dir... \n")
	for ;tries <= 20;tries++{
		err = os.RemoveAll(dir)
		if err != nil {
			//log.Printf("%d try, remove folder error \n", tries)
			time.Sleep(5*time.Second)
		}
	}
	if(err == nil){
		//log.Printf("-- %s removed, cleaning completed \n", dir)
	}
	return err
}

// DELETE A FILE
func DeleteFile(filename string)  error{
	tries := 1
	var err error
	for ;tries <= 20;tries++{
		err = os.Remove(filename)
		if err != nil {
			//log.Printf("%d try, remove file error \n", tries)
			time.Sleep(10 * time.Second)
		}
	}
	if(err == nil){
		log.Println("-- %s removed \n", filename)
	}
	return err
}

//ZERO PADDING, return padding + "i" --> zeroPad(1, 4) = "0001"
func ZeroPad(input int, returnLen int) string{
	data:= strconv.Itoa(input)
	len := len(data)
	gap:= returnLen - len;
	if(gap < 1){
		return data;
	}
	data = strings.Repeat("0", gap) + data;
	return data
}


func ReadFile(tFile *os.File) (int, []byte, error) {
	mBuf, err := ioutil.ReadAll(tFile)

	if err != nil {
		return 0, nil, err
	} else {
		return len(mBuf), mBuf, nil
	}
}

// CREATE DIR FOR DATABASE FILES
func CreateDir(dir string) error{
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0711)
		if err != nil {
			log.Println(" -- error creating " +  dir)
			return err
		}
		//log.Printf("-- created directory %s \n", dir)
	}else{
		//log.Printf("-- %s is ready\n", dir)
	}
	return nil
}

func IsExist(file string) bool{
	//Check existence of the file/dir
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}else{
		return true
	}
}


func GetBytes(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		//log.Println(err)
		return nil, err
	}

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes := make([]byte, size)

	// read file into bytes
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)

	return bytes, nil
}
func ParseList(listFile string) *list.List {
	inFile, _ := os.Open(listFile)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	ret := list.New()
	for scanner.Scan() {
		name := scanner.Text()
		if name != "" {
			keys := strings.Split(name, ".")
			key := keys[0] + ".mp4"
			ret.PushBack(key)
		}
	}
	return ret

}

func DurationToSeconds(duration string) string {
	//input must be hh:mm:ss
	timepot := strings.Split(duration, ":")
	if len(timepot) < 3 {
		log.Println("len:", len(timepot))
		return ""
	}

	hh := timepot[0]
	mm := timepot[1]
	ss := timepot[2]

	//	log.Println("hh:", hh)
	//	log.Println("mm:", mm)
	//	log.Println("ss:", ss)

	hv, _ := strconv.Atoi(hh)
	mv, _ := strconv.Atoi(mm)
	sv, _ := strconv.Atoi(ss)

	val := hv*3600 + mv*60 + sv
	retval := strconv.Itoa(val)
	return retval
}

func ReadXml(filePath string) string {
	type Tag struct {
		Key   string `xml:"mediaPresentationDuration,attr"`
		Value string `xml:",chardata"`
	}
	mpdStruct := new(Tag)
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		log.Println(err)
		return ""
	}
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes := make([]byte, size)

	// read file into bytes
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	xml.Unmarshal(bytes, mpdStruct)
	log.Println(mpdStruct.Key)
	return mpdStruct.Key
}

func GetISODuration(seconds int) string {
	t := time.Date(0, 0, 0, 0, 0, seconds, 0, time.Local)
	log.Println(t.Round(time.Second))
	log.Println(t.Hour())
	log.Println(t.Minute())
	log.Println(t.Second())
	return "PT" + strconv.Itoa(t.Hour()) + "H" + strconv.Itoa(t.Minute()) + "M" + strconv.Itoa(t.Second()) + ".000S"
}


func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}


func CopyFile2(src, dst string) (err error)  {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer w.Close()

	// do the actual work
	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}
	return nil
}