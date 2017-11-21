# Image Mock

Crie imagens mockadas rapidamente para seus protótipos web

## Como funciona?

Image Mock é um servidor web para gerar imagens de mock com um simples acesso a uma url

Por exemplo:

```html
<img src="http://localhost:8080/300" alt="300x300" />
```
![300x300](https://raw.githubusercontent.com/augustopimenta/imagemock/master/exemplos/300.png)

Alterando a cor de fundo:
```html
<img src="http://localhost:8080/300?bg=ff0000" alt="300x300" />
```
![300x300](https://raw.githubusercontent.com/augustopimenta/imagemock/master/exemplos/300_color.png)

Alterando a cor do texto:
```html
<img src="http://localhost:8080/300?fg=000" alt="300x300" />
```
![300x300](https://raw.githubusercontent.com/augustopimenta/imagemock/master/exemplos/300_tcolor.png)

Alterando o texto:
```html
<img src="http://localhost:8080/300?t=Exemplo" alt="300x300" />
```
![300x300](https://raw.githubusercontent.com/augustopimenta/imagemock/master/exemplos/300_text.png)

Alterando a forma:
```html
<img src="http://localhost:8080/300?r=150" alt="300x300" />
```
![300x300](https://raw.githubusercontent.com/augustopimenta/imagemock/master/exemplos/300_r.png)

Você pode ainda alterar o seu tamanho e mesclar todas essas opções

```
http://localhost:8080/300x450?bg=f00&fg=000&t=Hi!
```

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
