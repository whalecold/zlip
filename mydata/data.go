package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"mydata/utils"
	"net"
	"os"
	"sort"
	"strings"
	"sync"
)

//http://challenge.yuansuan.cn/

//全局数据
var RecvBuffer []byte
var BufferChan chan []byte
var LastCheckSum = 0x0
var FillNum = byte(0xAB)

type RecvData struct {
	Seq      int32
	CheckSum uint32
	Data     []byte
}

type RecvDataSlice []RecvData

func (r RecvDataSlice) Len() int {
	return len(r)
}

func (r RecvDataSlice) Swap(i, j int) {
	//r[i].Seq, r[j].Seq = r[j].Seq, r[i].Seq
	//r[i].Data, r[j].Data = r[j].Data, r[i].Data
	r[i], r[j] = r[j], r[i]
}

func (r RecvDataSlice) Less(i, j int) bool {
	return r[i].Seq < r[j].Seq
}

var SeqSlice []int
var ErrorSlice []int
var DataInfoSlice RecvDataSlice
var SumErrorNum int

func PrepareLogin(conn net.Conn) error {
	temp := make([]byte, 1024)
	readLen, err := conn.Read(temp)
	if err != nil {
		return err
	}

	handshake := string(temp[:readLen])
	handshake = strings.Replace(handshake, "\n", "", -1)

	//fmt.Printf("recv handshake : %v\n", handshake)

	strSlice := strings.Split(handshake, ":")
	if len(strSlice) < 2 {
		return errors.New("error Handshake!")
	}

	verifyInfo := fmt.Sprintf("IAM:%s:zls1129@gmail.com\n", strSlice[1])
	//fmt.Printf("%s\n", verifyInfo)
	_, err = conn.Write([]byte(verifyInfo))
	if err != nil {
		return err
	}

	readLen, err = conn.Read(temp)
	if err != nil {
		return err
	}

	fmt.Printf("recv result : %v \n", string(temp[:readLen]))
	return nil
}

var DataNum int32
var dataNumer uint32
var RightDataNum int

func UnpackDataChunk(chunk []byte) bool {
	if len(chunk) == 0 {
		return true
	}
	sequence := binary.BigEndian.Uint32(chunk[:4])
	checksum := binary.BigEndian.Uint32(chunk[4:8])
	//serinum := chunk[4:8]
	length := binary.BigEndian.Uint32(chunk[8:12])
	data := chunk[12 : length+12]
	//fmt.Printf("sequence : %v | checksum : %v | length : %v | data : %v \n", sequence, checksum, length, data)
	//fmt.Printf("checksum : %v  \n", serinum)

	lastSum := sequence
	dataNumer = dataNumer + length
	for index := 0; index*4 < len(data); index++ {
		temp := make([]byte, 4)
		dataLen := len(data)
		if index*4+4 > dataLen {
			temp1 := data[index*4 : dataLen]
			utils.DeepCopy(&temp, &temp1)
			for i := len(temp); i < 4; i++ {
				temp = append(temp, FillNum)
			}
		} else {
			temp1 := data[index*4 : index*4+4]
			utils.DeepCopy(&temp, &temp1)
		}
		lastSum ^= binary.BigEndian.Uint32(temp)
	}
	DataNum++

	result := make([]byte, 4)
	binary.BigEndian.PutUint32(result, lastSum)
	//fmt.Printf("data : %v\n", result)
	if lastSum != checksum {
		//fmt.Printf("index %v not same sum" +
		//	"%v and %v\n", sequence, serinum, result)
		ErrorSlice = append(ErrorSlice, int(sequence))
		SumErrorNum += len(data)
	} else {
		SeqSlice = append(SeqSlice, int(sequence))
		recvDataTemp := make([]byte, len(data))
		utils.DeepCopy(&recvDataTemp, &data)
		DataInfoSlice = append(DataInfoSlice, RecvData{int32(sequence),
			checksum, recvDataTemp})
		RightDataNum += len(recvDataTemp)
	}
	//fmt.Printf("seq : %v\n", sequence)
	return UnpackDataChunk(chunk[length+12:])
}

func main() {

	RecvBuffer = make([]byte, 1024*128)
	BufferChan = make(chan []byte)
	DataInfoSlice = make(RecvDataSlice, 0, 10240)
	SeqSlice = make([]int, 0, 1024)
	ErrorSlice = make([]int, 0, 1024)

	readLen := ReadFile("calcaData", RecvBuffer)
	RecvBuffer = RecvBuffer[:readLen]
	fmt.Printf("read len %v\n", len(RecvBuffer))
	UnpackDataChunk(RecvBuffer)
	WritePicFile()

	return

	GetDataFromTcp()
	WritePicFile()
	return
}

func GetDataFromTcp() {
	var wg sync.WaitGroup

	conn, err := net.Dial("tcp", "challenge.yuansuan.cn:7042")
	if err != nil {
		fmt.Printf("connect failed !\n")
		return
	}

	if err = PrepareLogin(conn); err != nil {
		fmt.Printf("PrepareLogin failed %v!\n", err)
		return
	}

	wg.Add(1)
	go func(con net.Conn) {
		for {
			length, error := con.Read(RecvBuffer)
			fmt.Printf("recv data length %v\n", length)
			if error != nil || length == 0 {
				fmt.Printf("Read failed %v!\n", error)
				close(BufferChan)
				wg.Done()
				break
			}
			newTemp := RecvBuffer[:length]
			TempBuffer := make([]byte, length)
			utils.DeepCopy(&TempBuffer, &newTemp)
			//RecvBuffer = RecvBuffer[:0]
			WriteFile("calcaData", TempBuffer)
			BufferChan <- TempBuffer
		}
	}(conn)

	wg.Add(1)
	go func() {
		for chunk := range BufferChan {
			UnpackDataChunk(chunk)
			fmt.Printf("data right : %v \n", DataNum)
		}
		wg.Done()
	}()
	wg.Wait()
}

func WritePicFile() {
	sort.Ints(SeqSlice)
	sort.Ints(ErrorSlice)
	fmt.Printf("sort : %v\n", SeqSlice)
	fmt.Printf("sort err : %v\n", ErrorSlice)
	fmt.Printf("get num : %v\n", dataNumer)

	sort.Sort(DataInfoSlice)
	lastData := make([]byte, 0, 1024)

	//lastSum := int32(-1)
	//var checkSum uint32
	//var tempVar []byte
	var lostLen int
	seqMap := make(map[int32]interface{})
	var MaxSumNumber int
	for _, value := range DataInfoSlice {

		MaxSumNumber += len(value.Data)
		if _, ok := seqMap[value.Seq]; ok {
			lostLen += len(value.Data)
			continue
		} else {
			seqMap[value.Seq] = 1
		}
		//if lastSum == value.Seq {
		//	//fmt.Printf("befor : %v \n now : %v\n", tempVar, value.Data)
		//	//fmt.Printf("befor : %v \n now : %v\n", checkSum, value.CheckSum)
		//	//tempVar = value.Data
		//	//checkSum = value.CheckSum
		//	//flag := CheckSameSlice(tempVar, value.Data)
		//	//fmt.Printf("is Same slice : %v\n", flag)
		//	lostLen += len(value.Data)
		//	continue
		//}
		//lastSum = value.Seq
		fmt.Printf("se-- %v\n", value.Seq)
		fmt.Printf("se--=== %v\n", len(value.Data))
		lastData = append(lastData, value.Data...)
		//checkSum = value.CheckSum
		//tempVar = value.Data
	}

	//fmt.Printf("---- %s\n", string(lastData))
	fmt.Printf("len ;; %v\n", len(lastData))
	fmt.Printf("len ;; %v\n", lostLen)
	fmt.Printf("len ;; %v\n", SumErrorNum)
	fmt.Printf("len ;111; %v\n", RightDataNum)
	fmt.Printf("len ;222; %v\n", MaxSumNumber)
	WriteFile("pic.png", lastData)
}

func CheckSameSlice(a, b []byte) bool {
	if a == nil || b == nil || len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func WriteFile(file string, data []byte) {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Printf("open file error")
		return
	}
	defer f.Close()
	writeLen, err := f.Write(data)
	if err != nil {
		fmt.Printf("write error")
		return
	}
	fmt.Printf("write len %v\n", writeLen)
}

func ReadFile(file string, data []byte) int {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("open file error")
		return 0
	}
	defer f.Close()
	writeLen, err := f.Read(data)
	if err != nil {
		fmt.Printf("write error")
		return 0
	}
	data = data[:writeLen]
	fmt.Printf("read len %v\n", writeLen)
	return writeLen
}
