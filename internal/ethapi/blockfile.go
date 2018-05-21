package ethapi

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"
)

// Block represents an block info.
type FileBlock struct {
	Timestamp    big.Int     `json:"timestamp"`    // 交易时间
	Number       big.Int     `json:"number"`       // 区块号
	ParentHash   string      `json:"parent_hash"`  // 区块父哈希
	Hash         string      `json:"hash"`         // 区块哈希
	TxCnt        int         `json:"tx_cnt"`       //交易个数
	Transactions interface{} `json:"transactions"` //交易列表
}

type storage struct {
	path     string
	file     *os.File
	writer   *bufio.Writer
	fileName string
	cnt      int64
}

func NewStorage(path string, cnt, blocknumber int64) *storage {
	s := new(storage)
	s.cnt = cnt
	s.fileName = s.getFileName(blocknumber)
	f, err := os.OpenFile(s.fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		panic(err)
	}
	s.file = f
	s.writer = bufio.NewWriter(f)

	return s
}

//InsertBlock insert block info into file
func (s *storage) InsertBlock(block *FileBlock) (*storage, error) {

	fileName := s.getFileName(block.Number.Int64())
	if s.fileName != fileName {
		s.writer.Flush()
		time.Sleep(5 * time.Second)
		s.file.Close()
		tmps := NewStorage(s.path, s.cnt, block.Number.Int64())
		s = tmps
	}

	blockByte, err := json.Marshal(block)
	if err != nil {
		return nil, err
	}

	s.writer.Write(blockByte)
	s.writer.WriteByte('\n')
	return s, nil
}

func (s *storage) ReadFile() error {
	// fileName := s.getFileName(4010000)
	// f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	// if err != nil {
	// 	return err
	// }
	// defer f.Close()

	// w := bufio.NewReader(f)
	// for {
	// 	line, err := w.ReadBytes('\n') //以'\n'为结束符读入一行
	// 	if err != nil || io.EOF == err {
	// 		log.Info("err ", err)
	// 		break
	// 	}
	// 	block := &idb.Block{}
	// 	json.Unmarshal(line, block)
	// }
	return nil
}

func (s *storage) getFileName(blockNumber int64) string {
	return s.path + fmt.Sprintf("eth_block_%04d_%04d.txt", blockNumber%20, blockNumber/s.cnt)
}
