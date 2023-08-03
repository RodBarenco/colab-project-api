If you prefer to read it in [English](#instructions-for-setting-up-and-running-the-project)


## Instruções para Montar e Rodar o Projeto

A seguir estão as etapas para montar e rodar o projeto. Certifique-se de ter o Go Language instalado em sua máquina ou MV. O projeto está sendo desenvolvido na versão Go 1.20.3. Caso necessário verifique a documentação oficial de Go [aqui](https://go.dev/)

1. **Clonar o Repositório**
   - Clone este repositório para o diretório desejado em sua máquina.

2. **Instalar Dependências**
   - Navegue para a pasta `./backend-src` do projeto.
   - Verifique o arquivo `go.mod` que lista todas as bibliotecas necessárias.

3. **Configuração do Banco de Dados**
   - O projeto está utilizando o PostgreSQL como banco de dados padrão.
   - Caso deseje utilizar outro banco de dados, modificações deverão ser feitas no código.
   - Estamos usando o GORM como ORM para lidar com o banco de dados. Sugerimos que você verifique a documentação do GORM [aqui](https://gorm.io/docs/index.html) para mais informações.

4. **Build do Projeto**
   - Após configurar as dependências, você pode buildar o projeto.
   - Certifique-se de estar na pasta `./backend-src`.
   - Execute o comando de build do Go para gerar o executável.

5. **Configuração do Arquivo .env**
   - Antes de executar o arquivo gerado, é necessário criar um arquivo `.env` no mesmo diretório.
   - O arquivo `.env` deve conter as seguintes variáveis:

```dotenv
PORT=numero_da_porta
DATABASE_URL="string_de_conexao_com_banco_de_dados"
SECRET=chave_secreta_do_JWT
```
## Execução do Projeto

- Após criar o arquivo .env, você pode executar o arquivo gerado com o seguinte comando:

```bash
./nome_do_executavel 1
```

### Testando em outra porta e banco de dados

- Caso seja necessário testar o projeto em outra porta e/ou com outro banco de dados, você pode criar um arquivo chamado .test.env.
O arquivo .test.env deve conter as configurações desejadas:

```dotenv
PORT=numero_da_outra_porta
DATABASE_URL="string_de_conexao_com_outro_banco_de_dados"
SECRET=outra_chave_secreta_do_JWT
```

- Para executar o arquivo gerado usando o .test.env, você deve passar o número 2 como argumento:

```bash
./nome_do_executavel 2
```
---
Com essas etapas concluídas, você deve estar pronto para montar e rodar o projeto com sucesso. Caso encontre algum problema, consulte as instruções acima novamente ou verifique a documentação das bibliotecas utilizadas para obter mais informações.

Lembre-se de substituir `nome_do_executavel` pelo nome real do executável gerado, assim como os números de portas, strings de conexão com de banco de dados e chaves secretas.

---
## Instructions for Setting Up and Running the Project

Below are the steps to set up and run the project. Make sure you have Go Language installed on your machine or VM. The project is being developed using Go version 1.20.3. If needed, check the official Go documentation [here](https://go.dev/).

1. **Clone the Repository**
   - Clone this repository to your desired directory on your machine.

2. **Install Dependencies**
   - Navigate to the `./backend-src` folder of the project.
   - Check the `go.mod` file to see the list of all required libraries.

3. **Database Configuration**
   - The project uses PostgreSQL as the default database.
   - If you wish to use a different database, modifications will need to be made in the code.
   - We are using GORM as the ORM to handle the database. We suggest checking the GORM documentation [here](https://gorm.io/docs/index.html) for more information.

4. **Project Build**
   - After configuring the dependencies, you can build the project.
   - Make sure you are in the `./backend-src` folder.
   - Execute the Go build command to generate the executable.

5. **Configuration of .env File**
   - Before running the generated file, create a `.env` file in the same directory.
   - The `.env` file should contain the following variables:

```dotenv
PORT=port_number
DATABASE_URL="database_connection_string"
SECRET=JWT_secret_key
```
## Running the Project

- After creating the .env file, you can run the generated file with the following command:

```bash
./executable_name 1
```

### Testing on Another Port and Database

- If you need to test the project on another port and/or with a different database, you can create a file named .test.env.
The .test.env file should contain the desired configurations:

```dotenv
PORT=another_port_number
DATABASE_URL="another_database_connection_string"
SECRET=another_JWT_secret_key
```
- To execute the generated file using .test.env, you should pass the number 2 as an argument:

```bash
./executable_name 2
```

---
With these steps completed, you should be ready to set up and run the project successfully. If you encounter any issues, please refer back to the instructions above or consult the documentation of the used libraries for further information.

Remember to replace executable_name with the actual name of the generated executable, as well as the port numbers, database connection strings, and secret keys.