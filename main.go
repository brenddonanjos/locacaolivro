package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/brenddonanjos/locacaolivro/cliente"
	"github.com/brenddonanjos/locacaolivro/livro"
	"github.com/brenddonanjos/locacaolivro/locacao"
)

func main() {
	serverInit()
}

//serverInit monta as rotas e inicia o server
func serverInit() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/cliente/", cliente.ClienteHandler)
	http.HandleFunc("/livro/", livro.LivroHandler)
	http.HandleFunc("/locacao/", locacao.LocacaoHandler)
	log.Println("executando ..")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func Home(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "{\"pagina\":\"Home\", \"mensagem\": \"Bem vindo!\"}")
}
