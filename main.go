package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"os"
)

type Estatistica map[rune]uint

func analisaEstatistica(r io.Reader) *Estatistica {
	buffer := bufio.NewReader(r)

	estatistica := make(Estatistica)

	for {
		r, _, err := buffer.ReadRune()
		if err != nil {
			break
		}

		value, exists := estatistica[r]

		if !exists {
			estatistica[r] = 1
		} else {
			estatistica[r] = value + 1
		}
	}

	return &estatistica
}

func montarAlfabeto(estatistica *Estatistica) *Alfabeto {
	a := make(Alfabeto, 0)

	for r, peso := range *estatistica {
		a = append(a, NewNoVazio(r, peso))
	}

	return &a
}

func constroiAlfabeto(a *Alfabeto) *Alfabeto {
	no1, no2 := a.RetiraMenor(), a.RetiraMenor()

	novoNo := NewNo(no1, no2)

	*a = append(*a, novoNo)

	return a
}

func montarTabelaDeTroca(no *No, tabela *Tabela, codigo []bool) {

	if no.IsFolha() {
		codigoFinal := make([]bool, len(codigo))
		copy(codigoFinal[:], codigo[:])
		(*tabela)[no.valor] = codigoFinal
	} else {
		if no.esquerdo != nil {
			montarTabelaDeTroca(no.esquerdo, tabela, append(codigo, false))
		}
		if no.direito != nil {
			montarTabelaDeTroca(no.direito, tabela, append(codigo, true))
		}
	}
}

func populaBits(r io.Reader, tabela *Tabela) *SeqBits {
	seqBits := NewSeqBits()

	buffer := bufio.NewReader(r)

	for {
		r, _, err := buffer.ReadRune()
		if err != nil {
			break
		}

		seqBits.AdicionaBits(tabela.Bits(r))
	}

	seqBits.CarriageReturn()

	return seqBits
}

func compacta(src, dst string) {
	// src := "./biblia.txt"
	// dst := "./biblia.fzip"

	input, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	// stringExemplo := "aaa bb c"
	// estatistica := analisaEstatistica(strings.NewReader(stringExemplo))
	estatistica := analisaEstatistica(input)
	alfabeto := montarAlfabeto(estatistica)

	for len(*alfabeto) > 1 {
		alfabeto = constroiAlfabeto(alfabeto)
	}

	raiz := alfabeto.Raiz()

	tabela := NewTabela()

	montarTabelaDeTroca(raiz, tabela, []bool{})

	tabelaJSON, err := json.Marshal((*tabela).ToJson())
	if err != nil {
		panic(err)
	}

	// seqBits := populaBits(strings.NewReader(stringExemplo), tabela)
	input, _ = os.Open(src)
	seqBits := populaBits(input, tabela)

	f, err := os.Create(dst)
	if err != nil {
		panic(err)
	}

	f.Write(tabelaJSON)
	f.Write([]byte("\n"))
	f.Write(seqBits.Bytes())

}

func descompacta(src, dst string) {
	// src := "./biblia.fzip"
	// dst := "./biblia_descompactada.txt"
	f, err := os.Open(src)
	if err != nil {
		panic(err)
	}

	tabelaJSON := []byte{}
	bits := []byte{}

	buffer := bufio.NewReader(f)

	for {
		b, _ := buffer.ReadByte()
		if b == byte('\n') {
			break
		}
		tabelaJSON = append(tabelaJSON, b)
	}

	for {
		b, err := buffer.ReadByte()
		if err != nil {
			break
		}
		bits = append(bits, b)
	}

	tabelaReversa := NewTabelaReversa(bits, tabelaJSON)

	var texto bytes.Buffer

	for {
		r, err := tabelaReversa.NextRune()
		if err != nil {
			break
		}
		texto.WriteRune(r)
	}

	newFile, err := os.Create(dst)
	if err != nil {
		panic(err)
	}

	newFile.Write(texto.Bytes())

}

func main() {
	var compactaFlag bool
	var descompactaFlag bool
	var source string
	var destination string

	flag.BoolVar(&compactaFlag, "c", false, "Usado para compactar o arquivo")
	flag.BoolVar(&descompactaFlag, "x", false, "Usado para descompactar o arquivo")
	flag.StringVar(&source, "i", "", "Arquivo de origem para ser compactado/descompactado")
	flag.StringVar(&destination, "o", "", "Arquivo de destino para ser compactado/descompactado")

	flag.Parse()

	if !compactaFlag && !descompactaFlag {
		panic("Deve selecionar pelo menos uma operação!")
	}
	if compactaFlag && descompactaFlag {
		panic("Deve selecionar apenas uma operação!")
	}
	if source == "" {
		panic("Você deve definir um arquivo de origem")
	}
	if destination == "" {
		panic("Você deve definir um arquivo de destino")
	}

	if compactaFlag {
		compacta(source, destination)
	}

	if descompactaFlag {
		descompacta(source, destination)
	}

}
