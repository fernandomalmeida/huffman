package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
)

type Alfabeto []*No

func (a *Alfabeto) RetiraMenor() *No {
	var menorIndex int
	menorPeso := uint(math.Inf(+1))
	var menorNo *No

	for i, no := range *a {
		if no.peso < menorPeso {
			menorIndex = i
			menorPeso = no.peso
			menorNo = no
		}
	}

	left := (*a)[:menorIndex]
	right := (*a)[menorIndex+1:]

	*a = append(left, right...)

	return menorNo
}

func (a *Alfabeto) Raiz() *No {
	return a.AsSlice()[0]
}

func (a *Alfabeto) AsSlice() []*No {
	return *a
}

type No struct {
	valor    rune
	peso     uint
	esquerdo *No
	direito  *No
}

func NewNoVazio(r rune, peso uint) *No {
	return &No{
		valor: r,
		peso:  peso,
	}
}

func NewNo(e *No, d *No) *No {
	return &No{
		peso:     e.peso + d.peso,
		esquerdo: e,
		direito:  d,
	}
}

func (n *No) IsFolha() bool {
	return n.esquerdo == nil && n.direito == nil
}

type Tabela map[rune][]bool

func NewTabela() *Tabela {
	return &Tabela{}
}

func (t *Tabela) Bits(r rune) []bool {
	return (*t)[r]
}

func (t *Tabela) ToJson() map[rune]string {
	tabela := make(map[rune]string)

	for key, value := range *t {
		tabela[key] = byteArrayToString(value)
	}

	return tabela
}

func byteArrayToString(bits []bool) string {
	var buffer bytes.Buffer

	for _, bit := range bits {
		buffer.WriteByte(boolToByte(bit))
	}

	return buffer.String()
}

func boolToByte(b bool) byte {
	if b {
		return byte('1')
	} else {
		return byte('0')
	}
}

type SeqBits struct {
	bytes    []byte
	currByte byte
	currPos  uint8
}

const initialCurrPos = 8

func NewSeqBits() *SeqBits {
	return &SeqBits{
		bytes:   make([]byte, 0),
		currPos: initialCurrPos,
	}
}

func (s *SeqBits) AdicionaBits(bits []bool) {
	for _, bit := range bits {
		s.currPos--
		if bit {
			s.currByte = s.currByte | (0x01 << s.currPos)
		}
		if s.currPos == 0 {
			s.CarriageReturn()
		}
	}
}

func (s *SeqBits) CarriageReturn() {
	s.bytes = append(s.bytes, s.currByte)
	s.currByte = 0x00
	s.currPos = initialCurrPos
}

func (s *SeqBits) Bytes() []byte {
	return s.bytes
}

type TabelaReversa struct {
	tabela   map[string]rune
	bytes    []byte
	currByte byte
	currPos  uint8
}

func NewTabelaReversa(bytes []byte, tabelaJson []byte) *TabelaReversa {
	tabela := make(map[rune]string)
	err := json.Unmarshal(tabelaJson, &tabela)
	if err != nil {
		panic(err)
	}

	tabelaReversa := &TabelaReversa{
		currByte: bytes[0],
		bytes:    bytes[1:],
		currPos:  initialCurrPos,
		tabela:   make(map[string]rune),
	}

	for key, value := range tabela {
		tabelaReversa.tabela[value] = key
	}

	return tabelaReversa
}

func (t *TabelaReversa) NextRune() (rune, error) {
	var buffer bytes.Buffer
	buffer.Reset()
	for {

		t.currPos--
		bit := t.currByte & (0x01 << t.currPos)
		if bit > 0x00 {
			buffer.WriteByte(byte('1'))
		} else {
			buffer.WriteByte(byte('0'))
		}

		if t.currPos == 0 {
			if len(t.bytes) == 0 {
				return rune(0), errors.New("No more Runes")
			}
			t.currByte = t.bytes[0]
			t.bytes = t.bytes[1:]
			t.currPos = initialCurrPos
		}

		if val, ok := t.tabela[buffer.String()]; ok {
			return val, nil
		}

	}
}

// func (s *RevSeqBits) AdicionaBits(bits []bool) {
// 	for _, bit := range bits {
// 		s.currPos--
// 		if bit {
// 			s.currByte = s.currByte | (0x01 << s.currPos)
// 		}
// 		if s.currPos == 0 {
// 			s.CarriageReturn()
// 		}
// 	}
// }

// func (s *RevSeqBits) CarriageReturn() {
// 	s.bytes = append(s.bytes, s.currByte)
// 	s.currByte = 0x00
// 	s.currPos = initialCurrPos
// }

// func (s *SeqBits) Bytes() []byte {
// 	return s.bytes
// }
