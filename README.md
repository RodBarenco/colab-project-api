If you prefer to read it in [English](#project-colab---backend---english)
Si prefieres leerlo en [Español](#proyecto-colab---backend---español)

[![](https://img.shields.io/badge/FlowChart%20-8AbBE2)](./FlowChart.png) [![](https://img.shields.io/badge/Setup-Guide%20-55ea85)](./setup-guide.md) [![](https://img.shields.io/badge/Tests%20-e86c45)](./test.md) [![](https://img.shields.io/badge/V1-routes%20-ffff1d)](./v1_routes.md) [![](https://img.shields.io/badge/🔒-cryptography%20-aa15ee)](./crypto.md)

# Projeto Colab - Backend - Português

Estou satisfeito com o andamento do projeto, e gostaria de destacar que ele está sendo desenvolvido também com fins educacionais, visando o aprendizado e aprimoramento de minhas habilidades de desenvolvimento.

Nesta etapa inicial, já temos uma estrutura básica do back-end implementada. As rotas estão organizadas em três grupos principais: gerais (v1), de usuários (v2) e de administradores (v3). Isso foi possível através da utilização de um middleware para realizar a separação e garantir a segurança das rotas.

Além disso, existe uma função de autenticação que utiliza JWT (JSON Web Tokens) para garantir a segurança e controle de acesso às rotas. Isso permitirá que apenas usuários autenticados e autorizados acessem determinadas partes da aplicação.

## Sobre o Projeto - Motivação

Muitas vezes, alunos de faculdades, grupos de pesquisa ou mesmo pessoas do ensino médio não têm um local adequado para publicar seus artigos acadêmicos e trabalhos de arte. O Projeto Colab visa preencher essa lacuna, fornecendo uma plataforma moderna para o compartilhamento de conteúdo acadêmico e artístico.

## Objetivos

O objetivo principal do Projeto Colab é criar uma plataforma amigável e acessível, onde estudantes, pesquisadores e artistas possam compartilhar e disseminar seus trabalhos de forma livre e aberta.

### Funcionalidades Planejadas

- Criação e gerenciamento de perfis de usuários
- Envio para aceitação e publicação além de edição de artigos e trabalhos de arte
- Compartilhamento rápido
- Comentários e interação entre os usuários (seguir autores, dar like em artigos, se comunicar com autores que estejam abertos para colaboração)
- Pesquisa e filtragem de conteúdo por assunto, autor e mais

## Contribuições

Contribuições são bem-vindas! Sinta-se à vontade para enviar pull requests com melhorias, correções de bugs e novas funcionalidades.
Pessoas de todas as áreas são bem-vindas.

## Próximo Passo:

Já dicionaremos testes para as funções de procura de artigos no Router V1. Isso garantirá que a busca por artigos funcione conforme o esperado e que futuras alterações não afetem seu funcionamento. Testes sólidos são essenciais para manter a qualidade do código e evitar problemas. 

Também implementamos funções de registro e login(v1) com autenticação e autorização(v2) de áreas específicas e ações individuais do usuário. Dessa forma, garantiremos que cada usuário tenha acesso apenas às áreas permitidas e possa realizar ações apropriadas, aumentando a segurança e a privacidade.

Sendo assim nos proximos passos estão:
- Testes para Funções de Procura de Artigos (Router V1). 
- Testes para Função de Criação de Artigos.
- Criação de rotas e funções para administradores.

## Aviso Legal

Este é um projeto de código aberto e está sendo desenvolvido para fins educacionais. Não há garantia de que este projeto seja implantado em produção ou utilizado comercialmente.

## Licença

O Projeto Colab é licenciado sob a Licença MIT. Consulte o arquivo LICENSE para obter mais informações.

---

# Project Colab - Backend - English

I am pleased with the progress of the project, and I would like to emphasize that it is also being developed for educational purposes, aiming to learn and enhance my development skills.

In this initial stage, we already have a basic structure of the backend implemented. The routes are organized into three main groups: general (v1), user-specific (v2), and administrator-specific (v3). This was achieved by using a middleware to perform this separation and ensure the security of the routes.

Furthermore, there is an authentication function that uses JWT (JSON Web Tokens) to ensure the security and access control of the routes. This will allow only authenticated and authorized users to access certain parts of the application.

## About the Project - Motivation

Often, college students, research groups, or even high school students do not have a suitable place to publish their academic articles and artworks. Project Colab aims to fill this gap by providing a modern platform for sharing academic and artistic content.

## Objectives

The main objective of Project Colab is to create a friendly and accessible platform where students, researchers, and artists can freely and openly share and disseminate their work.

### Planned Features

- User profile creation and management
- Submission for acceptance and publication, as well as editing of articles and artworks
- Quick sharing
- User interactions and comments (follow authors, like articles, communicate with authors open to collaboration)
- Search and filtering of content by subject, author, and more

## Contributions

Contributions are welcome! Feel free to send pull requests with improvements, bug fixes, and new features.
People from all fields are welcome to contribute.

## Next Step

We have already added tests for the article search function in Router V1. This will ensure that the search for articles works as expected and that future changes will not affect how it works. Solid testing is essential to maintain code quality and avoid problems.

We have also implemented registration and login function(v1) with authentication and authorization(v2) of specific areas and individual user actions. In this way, we will ensure that each user has access only to the permitted areas and can take appropriate actions, increasing security and privacy.

So the next steps are:
- Tests for Article Search Functions (Router V1).
- Tests for Article Creation Function.
- Creation of routes and functions for administrators.

## Legal Notice

This is an open-source project and is being developed for educational purposes. There is no guarantee that this project will be deployed in production or used commercially.

## License

Project Colab is licensed under the MIT License. See the LICENSE file for more information.

---

# Proyecto Colab - Backend - Español

Estoy satisfecho con el progreso del proyecto, y me gustaría destacar que también se está desarrollando con fines educativos, con el objetivo de aprender y mejorar mis habilidades de desarrollo.

En esta etapa inicial, ya tenemos una estructura básica del backend implementada. Las rutas están organizadas en tres grupos principales: generales (v1), específicas para usuarios (v2) y específicas para administradores (v3). Esto se logró utilizando un middleware para realizar esta separación y garantizar la seguridad de las rutas.

Además, existe una función de autenticación que utiliza JWT (JSON Web Tokens) para garantizar la seguridad y el control de acceso a las rutas. Esto permitirá que solo los usuarios autenticados y autorizados accedan a ciertas partes de la aplicación.

## Sobre el Proyecto - Motivación

A menudo, los estudiantes universitarios, grupos de investigación e incluso personas de secundaria no tienen un lugar adecuado para publicar sus artículos académicos y trabajos artísticos. El Proyecto Colab tiene como objetivo llenar este vacío, proporcionando una plataforma moderna para compartir contenido académico y artístico.

## Objetivos

El objetivo principal del Proyecto Colab es crear una plataforma amigable y accesible, donde estudiantes, investigadores y artistas puedan compartir y difundir sus trabajos de manera libre y abierta.

### Funcionalidades Planeadas

- Creación y gestión de perfiles de usuarios
- Envío para aceptación y publicación, así como edición de artículos y trabajos artísticos
- Compartir rápidamente
- Comentarios e interacción entre los usuarios (seguir a autores, dar "me gusta" en artículos, comunicarse con autores que estén abiertos a colaborar)
- Búsqueda y filtrado de contenido por tema, autor y más

## Contribuciones

¡Las contribuciones son bienvenidas! Siéntete libre de enviar pull requests con mejoras, correcciones de errores y nuevas funcionalidades.
Personas de todas las áreas son bienvenidas.

## Próximo Paso:

Ya hemos agregado pruebas para la función de búsqueda de artículos en Router V1. Esto asegurará que la búsqueda de artículos funcione según lo esperado y que los cambios futuros no afecten su funcionamiento. Las pruebas sólidas son esenciales para mantener la calidad del código y evitar problemas.

También hemos implementado la función de registro y inicio de sesión (v1) con autenticación y autorización (v2) de áreas específicas y acciones individuales de los usuarios. De esta manera, aseguraremos que cada usuario solo tenga acceso a las áreas permitidas y pueda realizar acciones apropiadas, aumentando la seguridad y la privacidad.

Por lo tanto, los próximos pasos son:

- Pruebas para las Funciones de Búsqueda de Artículos (Router V1).
- Pruebas para la Función de Creación de Artículos.
- Creación de rutas y funciones para administradores.

## Aviso Legal

Este es un proyecto de código abierto y se está desarrollando con fines educativos. No se garantiza que este proyecto sea implementado en producción o utilizado comercialmente.

## Licencia

El Proyecto Colab está licenciado bajo la Licencia MIT. Consulta el archivo LICENSE para obtener más información.

