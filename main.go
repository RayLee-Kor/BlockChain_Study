package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})

	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := newProofOfWork(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

type BlockChain struct {
	blocks []*Block
}

func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]

	newBlock := NewBlock(data, prevBlock.Hash)

	bc.blocks = append(bc.blocks, newBlock)
}

func NewGenesisBlock() *Block {
	return NewBlock("BlockChain is God", []byte{})
}

func NewBlockChain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}

const targetBits = 24 // 추출 난이도 선택

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func newProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	// 1
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}
	return pow
}

func intToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	// byte 슬라이스
	err := binary.Write(buff, binary.BigEndian, num)

	if err != nil {
		fmt.Println(err)
	}

	return buff.Bytes()
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.block.PrevBlockHash,
		pow.block.Data,
		intToHex(pow.block.Timestamp),
		intToHex(int64(targetBits)),
		intToHex(int64(nonce)),
	}, []byte{})
	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte

	max_number := ^uint(0)
	nonce := 0

	fmt.Printf("블록 마이닝 시작!!! %s\n", pow.block.Data)
	for uint(nonce) < max_number {
		data := pow.prepareData((nonce))
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Println("마이닝 성공!!!")
	return nonce, hash[:]
}

func main() {

	bc := NewBlockChain()

	bc.AddBlock("send 1 bitcoin")
	bc.AddBlock("send 1 klaytn")

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}

}
