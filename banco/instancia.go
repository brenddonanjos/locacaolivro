package banco

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DB_NAME = "golocacao"
	DB_USER = "root"
	DB_PASS = "root"
)

func Start() *sql.DB {
	checkDb() //verifica se o banco existe

	//instancia o banco específico
	db, err := sql.Open("mysql", DB_USER+":"+DB_PASS+"@/"+DB_NAME)
	if err != nil {
		panic(err)
	}

	//monta a estrutura do banco (simulando uma migration)
	exec(db, "use "+DB_NAME)

	exec(db, `create table if not exists clientes (
		id integer auto_increment,
		nome varchar(80),
		email varchar(80),
		telefone varchar(80),
		cpf varchar(80),
		ativo TINYINT(1),
		endereco_id integer,
		PRIMARY KEY (id)
	) ENGINE = innodb`)

	exec(db, `create table if not exists enderecos (
		id integer auto_increment,
		cep varchar(80),
		logradouro varchar(190),
		numero varchar(80),
		bairro varchar(80),
		cidade varchar(190),
		uf varchar(4),
		complemento varchar(255),
		PRIMARY KEY (id)
	) ENGINE = innodb`)

	exec(db, `create table if not exists livros (
		id integer auto_increment,
		titulo varchar(255),
		autor varchar(190),
		ano integer,
		edicao varchar(80),
		editora varchar(190),
		PRIMARY KEY (id)
	) ENGINE = innodb`)

	exec(db, `create table if not exists locacoes (
		id integer auto_increment,
		status varchar(80),
		data_locacao datetime,
		prazo_dias integer,
		cliente_id integer,
		livro_id integer,
		PRIMARY KEY (id),
		CONSTRAINT fk_cliente FOREIGN KEY (cliente_id) REFERENCES clientes (id),
		CONSTRAINT fk_livro FOREIGN KEY (livro_id) REFERENCES livros (id)
	) ENGINE = innodb`)

	return db
}

func checkDb() {
	db, err := sql.Open("mysql", DB_USER+":"+DB_PASS+"@/")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	exec(db, "CREATE DATABASE IF NOT EXISTS "+DB_NAME)
}

func exec(db *sql.DB, sql string) sql.Result {
	result, err := db.Exec(sql)

	if err != nil {
		panic(err)
	}

	return result
}
