If you prefer to read it in [English](#introduction)
# Projeto Colab - Testes

## Introdução

Esse projeto ainda está no começo e por isso o número de testes unitários para as funções implementadas no pacote "handlers" ainda é limitado. O objetivo dos testes unitários é garantir que as funções estejam funcionando corretamente e produzindo os resultados esperados.

## Executando os Testes

Os testes estão localizados na pasta "backend-src/handlers" e escritos usando a biblioteca de testes padrão do Go. Cada função do pacote "handlers" terá uma série de testes separados de modo a isolar suas diferentes funcionalidades e respostas esperadas para cada situação e assim facilitar a identificação de problemas específicos.

Para executar os testes, certifique-se de:

1. O arquivo .test.env esteja presente na pasta "backend-src" e contenha as configurações adequadas para conexão com o banco de dados. Esse arquivo é essencial para os testes realizarem as conexões corretamente e garantirem que os dados sejam manipulados adequadamente durante os testes.

2. O programa esteja rodando e com o argumento "2" no início ao ser executado. Isso pode ser feito ao chamar o programa e passar "2" como argumento na linha de comando. Por exemplo:

```bash
./seu_programa 2
```
3. O terminal a ser usado para o teste deve estar no diretório "backend-src/handlers" do projeto.

4.  Utilize o comando `go test -v` para iniciar a execução dos testes.

```bash
go test -v
```

## Função `RegisterHandler`

A função `RegisterHandler` é a principal função de registro de usuário do pacote "handlers". Até o momento, foram escritos apenas testes para entradas de `FirstName`, `LastName`, `Email`, `PassWord`, `DateOfBirth`, `Nickname`, `Field`, e `Biography` dessa função, além dos testes:

1. `TestRegisterHandler_Standard`: Testa o registro de um usuário padrão e verifica se a resposta está correta.
2. `TestRegisterHandler_DuplicateUser`: Testa o registro de um usuário com email duplicado e verifica se a resposta de erro é retornada corretamente.
---
Este é apenas o início dos testes unitários para o pacote "handlers". Mais testes serão adicionados no futuro para abranger diferentes cenários e garantir a robustez e confiabilidade do código.

Fique à vontade para contribuir com mais testes ou fazer melhorias nos testes existentes para garantir a qualidade contínua do código.

---

# Project Colab - Unit Tests

## Introduction

This project is still in its early stages, and as such, the number of unit tests for the functions implemented in the "handlers" package is still limited. The purpose of unit tests is to ensure that the functions are working correctly and producing the expected results.

## Running the Tests

The tests are located in the "backend-src/handlers" folder and written using the Go standard testing library. Each function in the "handlers" package will have a series of separate tests to isolate its different functionalities and expected responses for each situation, facilitating the identification of specific problems.

To run the tests, make sure to:

1. Have the ".test.env" file present in the "backend-src" folder, containing the appropriate configurations for connecting to the database. This file is essential for the tests to establish correct connections and ensure that data is manipulated appropriately during testing.

2. Ensure that the program is running with the argument "2" at the beginning when executed. This can be done by calling the program and passing "2" as an argument in the command line. For example:

```bash
./your_program 2
```
3. The terminal to be used for testing must be in the "backend-src/handlers" directory of the project.

4. Use the command ``go test -v` to start the test execution.

```bash
go test -v
```

## Function `RegisterHandler`

The `RegisterHandler` function is the main user registration function in the "handlers" package. So far, tests have been written only for the inputs of `FirstName`, `LastName`, `Email`, `PassWord`, `DateOfBirth`, `Nickname`, `Field`, and `Biography` of this function, in addition to the following tests:

1. `TestRegisterHandler_Standard`: Tests the registration of a standard user and verifies if the response is correct.
2. `TestRegisterHandler_DuplicateUser`: Tests the registration of a user with a duplicate email and verifies if the error response is returned correctly.
---
This is just the beginning of the unit tests for the "handlers" package. More tests will be added in the future to cover different scenarios and ensure the code's robustness and reliability.

Feel free to contribute more tests or make improvements to the existing tests to ensure the continuous quality of the code.

---