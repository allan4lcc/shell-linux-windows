package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

//Servidor que delvolve um Shell para cliente
func err(erro error) {
	if erro != nil {
		fmt.Println("Error -> ", erro)
		os.Exit(3)
		//validandos os erros
	}
}
func tcpBasico() {
	//Abrindo conexao basica tcp
	fmt.Println("Iniciando Servidor...")
	connect, erro1 := net.Listen("tcp", ":8081")
	err(erro1)

	//aceitando conexao tcp
	receberConexao, erro2 := connect.Accept()
	err(erro2)
	fmt.Println("Conexao recibida...")
	defer connect.Close() //a conexao sera ecerrada apos todo processo ser executao
	//Organizando os byts da conexao
	mensagem, erro3 := bufio.NewReader(receberConexao).ReadString('\n')
	err(erro3)
	fmt.Println("Msg recibida pelo cliente -> ", mensagem)

	//retorna mensagen ou um comando para o cliente conectado
	comando := executeCmd(mensagem)
	receberConexao.Write([]byte(comando))

}

//executa comandos do cmd/bash
func executeCmd(comando string) []byte {
	retornaComandoParaCliete, _ := exec.Command(strings.TrimSuffix(comando, "\n")).Output()
	return retornaComandoParaCliete
}

//func principal
func main() {
	tcpBasico()
}
