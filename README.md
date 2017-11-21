# Image Mock

Crie imagens mockadas rapidamente para seus protótipos web

## Compilando e Executando

Você pode compilá-lo diretamente se estiver com o go instalado em sua maquina: (Acesse o diretório onde se encontra os arquivos do projeto)

```bash
# go get ./...
# go build
```

Você pode então executar o binário:

```bash
# ./imagemock
```

## Porta padrão

Utilizando o parametro -p é possível especificar a porta que o servidor web utilizará (Por padrão é utilizada a porta 8080):

```bash
# ./imagemock -p 8888
```

## Docker

Você pode utilizar o docker para rodar a aplicação criando uma imagem(Acesse o diretório onde se encontra o arquivo de Dockerfile):

```bash
# docker build -t imagemock .
```

E então:

```bash
# docker run -d -p 8080:8080 imagemock
```
