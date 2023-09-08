If you prefer to read it in [English](#project-colab---backend---english)
Si prefieres leerlo en [Español](#proyecto-colab---backend---español)

[![](https://img.shields.io/badge/FlowChart%20-8AbBE2)](./FlowChart.png) [![](https://img.shields.io/badge/Setup-Guide%20-55ea85)](./setup-guide.md) [![](https://img.shields.io/badge/Tests%20-e86c45)](./test.md) [![](https://img.shields.io/badge/V1-routes%20-ffff1d)](./v1_routes.md) [![](https://img.shields.io/badge/🔒-cryptography%20-aa15ee)](./crypto.md)

# Projeto Colab - Backend - Português

Estou satisfeito com o andamento do projeto, e gostaria de destacar que ele está sendo desenvolvido também com fins educacionais, visando o aprendizado e aprimoramento de minhas habilidades de desenvolvimento.

Nesta etapa inicial, já temos uma estrutura básica do back-end implementada. As rotas estão organizadas em três grupos principais: gerais (v1), de usuários (v2) e de administradores (v3). Isso foi possível através da utilização de middlewares para realizar a separação e garantir a segurança das rotas. Também implementei a possibilidade de entrada de e saída de dados com criptografia em (v2) e obrigatoriedade em (v3).

## Sobre o Projeto - Motivação

Muitas vezes, alunos de faculdades, grupos de pesquisa ou mesmo pessoas do ensino médio não têm um local adequado para publicar seus artigos acadêmicos e trabalhos de arte. O Projeto Colab visa preencher essa lacuna, fornecendo uma plataforma moderna para o compartilhamento de conteúdo acadêmico e artístico.

## Objetivos

O objetivo principal do Projeto Colab é criar uma plataforma amigável e acessível, onde estudantes, pesquisadores e artistas possam compartilhar e disseminar seus trabalhos de forma livre e aberta.

### Rotas para Funcionalidades Implementadas

### Versão 1 (V1):

#### Rotas Gerais (v1):

- `GET /testreadiness`: Rota de teste de prontidão.
- `GET /testerror`: Rota de teste de erro.
- `GET /articles`: Obtém os últimos 1000 artigos.
- `GET /articles/home`: Obtém os últimos 50 artigos para exibição na página inicial.
- `GET /articles/title/{title}`: Obtém artigos por título.
- `GET /articles/subject/{subject}`: Obtém artigos por assunto.
- `GET /articles/author/{author}`: Obtém artigos por autor.
- `GET /articles/field/{field}`: Obtém artigos por campo.
- `GET /articles/keywords/{keywords}`: Obtém artigos por palavras-chave.
- `GET /articles/id/{id}`: Obtém um artigo por ID.
- `GET /articles/likedby/{id}`: Obtém os usuários que gostaram de um artigo.
- `GET /articles/isciting/{id}`: Obtém os artigos que estão sendo citados por um artigo.
- `GET /articles/citedby/{id}`: Obtém os artigos que citam um artigo.
- `PATCH /articles/share/{id}`: Incrementa o número de compartilhamentos de um artigo.
- `GET /getuser/{userID}`: Obtém informações do usuário por ID.
- `PATCH /recoverynow/`: Atualiza a senha do usuário.
- `POST /send-recovery-mail`: Envia um e-mail de recuperação de senha.
- `GET /image/{id}`: Obtém uma imagem em base64 por ID.
- `POST /register`: Rota para registrar um usuário.
- `POST /login`: Rota para fazer login como usuário.
- `POST /Adminlogin`: Rota para fazer login como administrador. (deve ser encriptado)
- `GET /show-pkey`: Obtém a chave pública para criptografar o envio de menssagens.

### Versão 2 (V2):

- As rotas da V2 incluem autenticação e autorização.

#### Rotas Gerais (v2):

- `GET /testerror`: Rota de teste de erro.
- `GET /testreadiness/{userID}`: Rota de teste de prontidão com autenticação.
- `GET /home-articles/{userID}`: Obtém artigos recomendados para o usuário. (deve ter a inclusão do número de meses a serem pesquisados - exemplo: a inclusão de `?monthsAgo=12` mostrará apenas os artigo aceitos nos últimos 12 meses)
- `GET /is-article-liked/{userID}`: Verifica se o usuário gostou do artigo em questão.

#### Rotas de Artigos (v2):

- `POST /create-article/{userID}`: Cria um novo artigo.
- `PATCH /like-article/{userID}`: Adiciona um usuário à lista de gostaram de um artigo.
- `PATCH /unlike-article/{userID}`: Remove um usuário da lista de gostaram de um artigo.
- `PATCH /add-cited-article/{userID}`: Adiciona uma citação a um artigo.
- `PATCH /remove-cited-article/{userID}`: Remove uma citação de um artigo.
- `PATCH /add-key/{userID}`: Adiciona uma chave pública ao usuário. ( deve ser uma chave RSA de 2048 bits)
- `PATCH /follow/{userID}`: Faz um usuário seguir outro.
- `PATCH /unfollow/{userID}`: Faz um usuário deixar de seguir outro.
- `PATCH /update-user/{userID}`: Atualiza campos do usuário. (com ela podemos modificar qualquer dado do usuário)

#### Rotas Encriptadas (v2):

- Rotas com a terminação `-ecpt` indicam que a resposta é encriptada e que a entrada também o deve ser.
- exemplo qualquer rota presente em `(v2)` pode ser chamada da seguinte maneira `PATCH /like-article-ecpt/{userID}`

### Versão 3 (V3):

- Todas as rotas na V3 são completamente encriptadas, tanto na entrada quanto na saída.

#### Rotas Gerais (v3):

- `GET /testreadiness`: Rota de teste de prontidão.
- `GET /testerror`: Rota de teste de erro.
- `GET /unaccepted-article-id/{adminID}`: Obtém IDs de artigos não aceitos.
- `GET /unaccepted-article-by-field/{adminID}`: Obtém artigos não aceitos por campo.

#### Rotas de Administração (v3):

- `PATCH /approve-article/{adminID}`: Aprova um artigo, requer nível de permissão mínimo igual 3 para aprovação direta, ou 3 aprovações de Administradores com menor grau de permissão.
- `PATCH /approve-admin/{adminID}`: Aprova um administrador, requer nível de permissão mínimo igual 3.
- `PATCH /disapprove-admin/{adminID}`: Desaprova um administrador, requer nível de permissão mínimo igual 3.
- `PATCH /mod-permission-admin/{adminID}`: Modifica as permissões de um administrador (PATCH).

#### Rotas de Exclusão (v3):

- `DELETE /delete-article/{adminID}`: Exclui um artigo, requer nível de permissão mínimo igual 3.
- `DELETE /delete-user/{adminID}`: Exclui um usuário, requer nível de permissão mínimo igual 3.
- `DELETE /delete-admin/{adminID}`: Exclui um administrador,requer nível de permissão mínimo igual 3 e que o outro administrador tenha um nível de permissão menor e esteja aprovado.
- `DELETE /clean-articles/{adminID}`: Limpa todos os artigos não aceitos antigos.
- `DELETE /clean-articles-by-date/{adminID}`: Limpa artigos não aceitos antigos até determinada data, apenas o root (permissão 4 pode usar).
- `DELETE /clean-articles-by-date-and-field/{adminID}`: Limpa artigos não aceitos antigos até determinada data podendo ser usado por administradores com permissão 3 do mesmo campo.

Essas são as rotas do projeto "Projeto Colab" com seus tipos de ação correspondentes.

### Funcionalidades a serem Implementadas

- Criação de administradores (o root é criado automaticamente caso não exista).
- Compartilhamento rápido - deve ser implementado pelo cliente.
- Talvez comentários e interação entre os usuários (adição de comentários sobre artigos e possibilidade de se comunicar com autores que estejam abertos para colaboração)

## Contribuições

Contribuições são bem-vindas! Sinta-se à vontade para enviar pull requests com melhorias, correções de bugs e novas funcionalidades.
Pessoas de todas as áreas são bem-vindas.

## Próximo Passo:
Já adicionamos testes E2E para  as funções de login e registro de usuários de v1.

Agora adicionaremos testes para as funções de procura de artigos no Router V1. Isso garantirá que a busca por artigos funcione conforme o esperado e que futuras alterações não afetem seu funcionamento. Testes sólidos são essenciais para manter a qualidade do código e evitar problemas. 

Sendo assim nos proximos passos estão:
- Testes para Funções de Procura de Artigos (Router v1). 
- Testes para Função de Criação de Artigos. (v2)

## Aviso Legal

Este é um projeto de código aberto e está sendo desenvolvido para fins educacionais. Não há garantia de que este projeto seja implantado em produção ou utilizado comercialmente.

## Licença

O Projeto Colab é licenciado sob a Licença MIT. Consulte o arquivo LICENSE para obter mais informações.

---

# Project Colab - Backend - English

I'm pleased with the progress of the project, and I'd like to highlight that it's also being developed for educational purposes, aiming to enhance my development skills and learning.

In this initial stage, we already have a basic backend structure implemented. The routes are organized into three main groups: general (v1), user-specific (v2), and administrator-specific (v3). This segregation and route security have been achieved through the use of middlewares. Additionally, I've implemented data encryption in (v2) and mandatory encryption in (v3) for input and output.

## About the Project - Motivation

Often, college students, research groups, or even high school individuals lack a suitable platform to publish their academic articles and artistic works. The Project Colab aims to fill this gap by providing a modern platform for sharing academic and artistic content.

## Objectives

The main goal of Project Colab is to create a user-friendly and accessible platform where students, researchers, and artists can freely and openly share and disseminate their work.

### Routes for Implemented Features

### Version 1 (V1):

#### General Routes (v1):

- `GET /testreadiness`: Readiness test route.
- `GET /testerror`: Error test route.
- `GET /articles`: Get the latest 1000 articles.
- `GET /articles/home`: Get the latest 50 articles for display on the homepage.
- `GET /articles/title/{title}`: Get articles by title.
- `GET /articles/subject/{subject}`: Get articles by subject.
- `GET /articles/author/{author}`: Get articles by author.
- `GET /articles/field/{field}`: Get articles by field.
- `GET /articles/keywords/{keywords}`: Get articles by keywords.
- `GET /articles/id/{id}`: Get an article by ID.
- `GET /articles/likedby/{id}`: Get users who liked an article.
- `GET /articles/isciting/{id}`: Get articles being cited by an article.
- `GET /articles/citedby/{id}`: Get articles citing an article.
- `PATCH /articles/share/{id}`: Increment the number of shares for an article.
- `GET /getuser/{userID}`: Get user information by ID.
- `PATCH /recoverynow/`: Update the user's password.
- `POST /send-recovery-mail`: Send a password recovery email.
- `GET /image/{id}`: Get an image in base64 by ID.
- `POST /register`: Route for user registration.
- `POST /login`: Route for user login.
- `POST /Adminlogin`: Route for administrator login (must be encrypted).
- `GET /show-pkey`: Get the public key for encrypting message transmission.

### Version 2 (V2):

- V2 routes include authentication and authorization.

#### General Routes (v2):

- `GET /testerror`: Error test route.
- `GET /testreadiness/{userID}`: Readiness test route with authentication.
- `GET /home-articles/{userID}`: Get recommended articles for the user (must include the number of months to search - e.g., including `?monthsAgo=12` will show only articles accepted in the last 12 months).
- `GET /is-article-liked/{userID}`: Check if the user liked the specified article.

#### Article Routes (v2):

- `POST /create-article/{userID}`: Create a new article.
- `PATCH /like-article/{userID}`: Add a user to the list of those who liked an article.
- `PATCH /unlike-article/{userID}`: Remove a user from the list of those who liked an article.
- `PATCH /add-cited-article/{userID}`: Add a citation to an article.
- `PATCH /remove-cited-article/{userID}`: Remove a citation from an article.
- `PATCH /add-key/{userID}`: Add a public key to the user (must be a 2048-bit RSA key).
- `PATCH /follow/{userID}`: Make a user follow another user.
- `PATCH /unfollow/{userID}`: Make a user stop following another user.
- `PATCH /update-user/{userID}`: Update user fields (this can modify any user data).

#### Encrypted Routes (v2):

- Routes with the `-ecpt` suffix indicate that the response is encrypted, and the input should be encrypted as well. For example, any route in `(v2)` can be called as follows: `PATCH /like-article-ecpt/{userID}`.

### Version 3 (V3):

- All routes in V3 are completely encrypted, both for input and output.

#### General Routes (v3):

- `GET /testreadiness`: Readiness test route.
- `GET /testerror`: Error test route.
- `GET /unaccepted-article-id/{adminID}`: Get IDs of unaccepted articles.
- `GET /unaccepted-article-by-field/{adminID}`: Get unaccepted articles by field.

#### Administration Routes (v3):

- `PATCH /approve-article/{adminID}`: Approve an article, requires a minimum permission level of 3 for direct approval or 3 approvals from Administrators with lower permission levels.
- `PATCH /approve-admin/{adminID}`: Approve an administrator, requires a minimum permission level of 3.
- `PATCH /disapprove-admin/{adminID}`: Disapprove an administrator, requires a minimum permission level of 3.
- `PATCH /mod-permission-admin/{adminID}`: Modify an administrator's permissions (PATCH).

#### Deletion Routes (v3):

- `DELETE /delete-article/{adminID}`: Delete an article, requires a minimum permission level of 3.
- `DELETE /delete-user/{adminID}`: Delete a user, requires a minimum permission level of 3.
- `DELETE /delete-admin/{adminID}`: Delete an administrator, requires a minimum permission level of 3, and the other administrator must have a lower permission level and be approved.
- `DELETE /clean-articles/{adminID}`: Clean all old unaccepted articles.
- `DELETE /clean-articles-by-date/{adminID}`: Clean old unaccepted articles up to a certain date, only root (permission 4) can use this.
- `DELETE /clean-articles-by-date-and-field/{adminID}`: Clean old unaccepted articles up to a certain date, can be used by administrators with permission 3 from the same field.

These are the routes for the "Project Colab" project with their corresponding actions.

### Future Features to Implement

- Creation of administrators (the root is created automatically if it doesn't exist).
- Fast sharing - to be implemented by the client.
- Possibly comments and interaction among users (adding comments on articles and the possibility to communicate with authors open to collaboration).

## Contributions

Contributions are welcome! Feel free to submit pull requests with improvements, bug fixes, and new features.
People from all backgrounds are welcome.

## Next Step:
We have already added E2E tests for the login and user registration functions of v1.

Next, we will add tests for the article search functions in Router V1. This will ensure that

 article searches work as expected and that future changes do not affect their functionality. Robust tests are essential to maintain code quality and prevent issues.

Therefore, the next steps include:
- Tests for Article Search Functions (Router v1).
- Tests for Article Creation Function (v2).

## Legal Disclaimer

This is an open-source project developed for educational purposes. There is no guarantee that this project will be deployed in production or used commercially.

## License

The Project Colab is licensed under the MIT License. Please refer to the LICENSE file for more information.

---

# Proyecto Colab - Backend - Español

Estoy satisfecho con el progreso del proyecto y me gustaría destacar que también se está desarrollando con fines educativos, con el objetivo de mejorar mis habilidades de desarrollo y aprendizaje.

En esta etapa inicial, ya tenemos implementada una estructura básica del backend. Las rutas están organizadas en tres grupos principales: generales (v1), específicas para usuarios (v2) y específicas para administradores (v3). Esta segmentación y seguridad de rutas se ha logrado mediante el uso de middlewares. Además, he implementado la encriptación de datos tanto en (v2) como en (v3), donde es obligatoria, tanto en la entrada como en la salida.

## Acerca del Proyecto - Motivación

A menudo, los estudiantes universitarios, los grupos de investigación o incluso las personas de la escuela secundaria no tienen un lugar adecuado para publicar sus artículos académicos y obras artísticas. El Proyecto Colab tiene como objetivo llenar este vacío proporcionando una plataforma moderna para compartir contenido académico y artístico.

## Objetivos

El objetivo principal del Proyecto Colab es crear una plataforma amigable y accesible donde los estudiantes, investigadores y artistas puedan compartir y difundir libremente su trabajo.

### Rutas para Funcionalidades Implementadas

### Versión 1 (V1):

#### Rutas Generales (v1):

- `GET /testreadiness`: Ruta de prueba de preparación.
- `GET /testerror`: Ruta de prueba de error.
- `GET /articles`: Obtiene los últimos 1000 artículos.
- `GET /articles/home`: Obtiene los últimos 50 artículos para mostrar en la página de inicio.
- `GET /articles/title/{title}`: Obtiene artículos por título.
- `GET /articles/subject/{subject}`: Obtiene artículos por tema.
- `GET /articles/author/{author}`: Obtiene artículos por autor.
- `GET /articles/field/{field}`: Obtiene artículos por campo.
- `GET /articles/keywords/{keywords}`: Obtiene artículos por palabras clave.
- `GET /articles/id/{id}`: Obtiene un artículo por ID.
- `GET /articles/likedby/{id}`: Obtiene usuarios a quienes les gustó un artículo.
- `GET /articles/isciting/{id}`: Obtiene los artículos que están siendo citados por un artículo.
- `GET /articles/citedby/{id}`: Obtiene los artículos que citan un artículo.
- `PATCH /articles/share/{id}`: Incrementa el número de veces que se ha compartido un artículo.
- `GET /getuser/{userID}`: Obtiene información de usuario por ID.
- `PATCH /recoverynow/`: Actualiza la contraseña del usuario.
- `POST /send-recovery-mail`: Envía un correo electrónico de recuperación de contraseña.
- `GET /image/{id}`: Obtiene una imagen en formato base64 por ID.
- `POST /register`: Ruta para registrar a un usuario.
- `POST /login`: Ruta para iniciar sesión como usuario.
- `POST /Adminlogin`: Ruta para iniciar sesión como administrador (debe estar encriptada).
- `GET /show-pkey`: Obtiene la clave pública para encriptar el envío de mensajes.

### Versión 2 (V2):

- Las rutas de la V2 incluyen autenticación y autorización.

#### Rutas Generales (v2):

- `GET /testerror`: Ruta de prueba de error.
- `GET /testreadiness/{userID}`: Ruta de prueba de preparación con autenticación.
- `GET /home-articles/{userID}`: Obtiene artículos recomendados para el usuario (debe incluir el número de meses a buscar, por ejemplo, incluir `?monthsAgo=12` mostrará solo los artículos aceptados en los últimos 12 meses).
- `GET /is-article-liked/{userID}`: Verifica si al usuario le gustó el artículo en cuestión.

#### Rutas de Artículos (v2):

- `POST /create-article/{userID}`: Crea un nuevo artículo.
- `PATCH /like-article/{userID}`: Agrega un usuario a la lista de quienes les gustó un artículo.
- `PATCH /unlike-article/{userID}`: Elimina un usuario de la lista de quienes les gustó un artículo.
- `PATCH /add-cited-article/{userID}`: Agrega una cita a un artículo.
- `PATCH /remove-cited-article/{userID}`: Elimina una cita de un artículo.
- `PATCH /add-key/{userID}`: Agrega una clave pública al usuario (debe ser una clave RSA de 2048 bits).
- `PATCH /follow/{userID}`: Hace que un usuario siga a otro.
- `PATCH /unfollow/{userID}`: Hace que un usuario deje de seguir a otro.
- `PATCH /update-user/{userID}`: Actualiza los campos del usuario (esto puede modificar cualquier dato del usuario).

#### Rutas Encriptadas (v2):

- Las rutas con el sufijo `-ecpt` indican que la respuesta está encriptada y que la entrada también debe estar encriptada. Por ejemplo, cualquier ruta en `(v2)` se puede llamar de la siguiente manera: `PATCH /like-article-ecpt/{userID}`.

### Versión 3 (V3):

- Todas las rutas en V3 están completamente encriptadas, tanto para la entrada como para la salida.

#### Rutas Generales (v3):

- `GET /testreadiness`: Ruta de prueba de preparación.
- `GET /testerror`: Ruta de prueba de error.
- `GET /unaccepted-article-id/{adminID}`: Obtiene los IDs de los artículos no aceptados.
- `GET /unaccepted-article-by-field/{adminID}`: Obtiene artículos no aceptados por campo.

#### Rutas de Administración (v3):

- `PATCH /approve-article/{adminID}`: Aprueba un artículo, requiere un nivel mínimo de permisos de 3 para la aprobación

 directa o 3 aprobaciones de administradores con un nivel de permisos menor.
- `PATCH /approve-admin/{adminID}`: Aprueba a un administrador, requiere un nivel mínimo de permisos de 3.
- `PATCH /disapprove-admin/{adminID}`: Desaprueba a un administrador, requiere un nivel mínimo de permisos de 3.
- `PATCH /mod-permission-admin/{adminID}`: Modifica los permisos de un administrador (PATCH).

#### Rutas de Eliminación (v3):

- `DELETE /delete-article/{adminID}`: Elimina un artículo, requiere un nivel mínimo de permisos de 3.
- `DELETE /delete-user/{adminID}`: Elimina a un usuario, requiere un nivel mínimo de permisos de 3.
- `DELETE /delete-admin/{adminID}`: Elimina a un administrador, requiere un nivel mínimo de permisos de 3 y que el otro administrador tenga un nivel de permisos menor y esté aprobado.
- `DELETE /clean-articles/{adminID}`: Limpia todos los artículos no aceptados antiguos.
- `DELETE /clean-articles-by-date/{adminID}`: Limpia artículos no aceptados antiguos hasta una fecha específica, solo el administrador root (permiso 4) puede usarlo.
- `DELETE /clean-articles-by-date-and-field/{adminID}`: Limpia artículos no aceptados antiguos hasta una fecha específica y puede ser utilizado por administradores con permisos 3 en el mismo campo.

Estas son las rutas del proyecto "Proyecto Colab" con sus tipos de acciones correspondientes.

### Funcionalidades a Implementar

- Creación de administradores (el administrador root se crea automáticamente si no existe).
- Compartir rápidamente: debe ser implementado por el cliente.
- Posiblemente comentarios e interacción entre los usuarios (agregar comentarios sobre artículos y la posibilidad de comunicarse con autores que estén abiertos a colaborar).

## Contribuciones

¡Las contribuciones son bienvenidas! Siéntete libre de enviar solicitudes de extracción con mejoras, correcciones de errores y nuevas funcionalidades. Personas de todas las áreas son bienvenidas.

## Próximo Paso

Ya hemos agregado pruebas E2E para las funciones de inicio de sesión y registro de usuarios en v1.

A continuación, agregaremos pruebas para las funciones de búsqueda de artículos en Router V1. Esto asegurará que las búsquedas de artículos funcionen según lo esperado y que los cambios futuros no afecten su funcionamiento. Las pruebas sólidas son esenciales para mantener la calidad del código y evitar problemas.

Por lo tanto, los próximos pasos incluyen:
- Pruebas para las Funciones de Búsqueda de Artículos (Router v1).
- Pruebas para la Función de Creación de Artículos (v2).

## Aviso Legal

Este es un proyecto de código abierto y se está desarrollando con fines educativos. No hay garantía de que este proyecto se implemente en producción o se utilice comercialmente.

## Licencia

El Proyecto Colab está bajo la Licencia MIT. Consulta el archivo LICENSE para obtener más información.

---