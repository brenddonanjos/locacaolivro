CREATE DATABASE IF NOT EXISTS golocacao;

USE golocacao;

CREATE TABLE IF NOT EXISTS clientes (
    id integer auto_increment,
    nome varchar(80),
    email varchar(80),
    telefone varchar(80),
    cpf varchar(80),
    ativo TINYINT(1),
    endereco_id integer,
    PRIMARY KEY (id)
) ENGINE = innodb;

CREATE TABLE IF NOT EXISTS enderecos (
    id integer auto_increment,
    cep varchar(80),
    logradouro varchar(190),
    numero varchar(80),
    bairro varchar(80),
    cidade varchar(190),
    uf varchar(4),
    complemento varchar(255),
    PRIMARY KEY (id)
) ENGINE = innodb;

CREATE TABLE IF NOT EXISTS livros (
    id integer auto_increment,
    titulo varchar(255),
    autor varchar(190),
    ano integer,
    edicao varchar(80),
    editora varchar(190),
    PRIMARY KEY (id)
) ENGINE = innodb;

CREATE TABLE IF NOT EXISTS locacoes (
    id integer auto_increment,
    status varchar(80),
    data_locacao datetime,
    prazo_dias integer,
    cliente_id integer,
    livro_id integer,
    PRIMARY KEY (id),
    CONSTRAINT fk_cliente FOREIGN KEY (cliente_id) REFERENCES clientes (id),
    CONSTRAINT fk_livro FOREIGN KEY (livro_id) REFERENCES livros (id)
) ENGINE = innodb