## Rotas já implementadas em V1
Implemented Routes in V1

*  Aqui você encontrará um resumo das rotas implementadas em V1 incluindo os métodos, endpoints, corpos das requisições e respostas esperadas. / Here you will find a summary of the routes implemented in V1 including methods, endpoints, request bodies, and expected responses.

*  Pesquisas com strings são case-sensitives, ainda irei mudar isso. / String searches are case-sensitive; I will change this in the future.

## Rotas de verificação de disponibilidade:
Routes for Availability Check:

### GET "v1/testreadiness"
- Não é necessário enviar corpo. / No request body required.
- resposta esperada – `200 Ok` / Expected response – `200 OK`

```json
{}
```
### GET "v1/testerror"
- Não é necessário enviar corpo.
- resposta esperada - `401 Unauthorized`

```json
 {
	"error": "Unauthorized: Something went wrong"
}
```

## Rotas aplicadas sobre artigos:
Routes applied to articles:

### GET "v1/articles"
- Não é necessário enviar corpo. / No request body required.
- resposta esperada – `200 Ok` / Expected response – `200 OK`
Caso existam artigos registrados retorna no máximo os últimos 1000, no formato desse exemplo: / If there are registered articles, it returns a maximum of the last 1000, in the format of this example:
- [Retorno](#retorno-comum-para-artigos) comum para artigos. / [Return](#common-response-for-articles) common for articles.

### GET "v1/articles/home"
- Não é necessário enviar corpo. / No request body required.
- resposta esperada – `200 Ok` / Expected response – `200 OK`
Caso existam artigos registrados retorna no máximo os últimos 50, no formato desse exemplo: / If there are registered articles, it returns a maximum of the last 50, in the format of this example:
- [Retorno](#retorno-comum-para-artigos) comum para artigos. / [Return](#common-response-for-articles) common for articles.

### GET "v1/articles/title/{title}"
- Não é necessário enviar corpo. / No request body required.
- resposta esperada – `200 Ok` / Expected response – `200 OK`
Caso exista artigo com esse título o retorna da seguinte maneira: / If an article exists with this title, it returns in the following way:
- [Retorno](#retorno-comum-para-artigos) comum para artigos. / [Return](#common-response-for-articles) common for articles.

### GET "v1/articles/subject/{subject} 
- Não é necessário enviar corpo. / No request body required.
- resposta esperada – `200 Ok` / Expected response – `200 OK`
Caso existam artigos com esses assuntos os retorna da seguinte maneira: / If there are articles with these subjects, it returns them in the following way:
- [Retorno](#retorno-comum-para-artigos) comum para artigos. / [Return](#common-response-for-articles) common for articles.

### GET "v1/articles/author/{author}"
- Não é necessário enviar corpo. / No request body required.
- resposta esperada – `200 Ok` / Expected response – `200 OK`
Caso existam artigos com esse autor: / If there are articles by this author:
- [Retorno](#retorno-comum-para-artigos) comum para artigos. / [Return](#common-response-for-articles) common for articles.

### GET"/articles/field/{field}"
- Não é necessário enviar corpo. / No request body required.
- resposta esperada – `200 Ok` / Expected response – `200 OK`
Caso existam artigos desse campo os retorna da seguinte maneira: / If there are articles from this field, it returns them in the following way:
- [Retorno](#retorno-comum-para-artigos) comum para artigos. / [Return](#common-response-for-articles) common for articles.

### Get "v1/articles/keywords/{keywords}"
- Não é necessário enviar corpo. / No request body required.
- resposta esperada – `200 Ok` / Expected response – `200 OK`
- Caso existam artigos com uma das palavras procuradas separadas por vírgula: / If there are articles with any of the searched words separated by commas:
- [Retorno](#retorno-comum-para-artigos) comum para artigos. / [Return](#common-response-for-articles) common for articles.

### GET"v1/articles/id/{id}"
- Não é necessário enviar corpo. / No request body required.
- resposta esperada – `200 Ok` / Expected response – `200 OK`
Caso existam artigos com esse id os retorna da seguinte maneira: / If there are articles with this ID, it returns them in the following way:
```json
{
	"article": {
		"ID": "fad14955-6c2c-4b81-8371-e581f2aab209",
		"Title": "Sample Article",
		"AuthorID": "5f32d05a-948e-4f7d-879f-2a840988048c",
		"Subject": "Science",
		"Field": "Physics",
                             "File": "Code Base64",
		"Description": "This is a sample article.",
		"Keywords": "sample, article, test",
		"SubmissionDate": "2023-08-09T21:02:48.271936-03:00",
		"AcceptanceDate": "0000-12-31T21:00:00-03:00",
		"IsAccepted": true,
		"LikedBy": null,
		"Citations": null,
		"Shares": 10,
		"CoAuthors": "John Doe, Jane Smith",
		"CoverImage": "abc"
	},
	"relatedTables": {
		"numLikes": 1,
		"likedByNames": [
			"Joe Doe"
		]
	},
	"message": "Article and related data retrieved successfully"
}
```
### GET "v1/articles/likedby/{id}",
- Não é necessário enviar corpo. / No request body required.
- resposta esperada – `200 Ok` / Expected response – `200 OK`
- Caso o artigo em questão tenha sido “gostado” por algum usuário: / If the article in question has been "liked" by any user:
```json
{
	"liked_by_users": [
		{
			"id": "5f32d05a-948e-4f7d-879f-2a840988048c",
			"username": "Joe Doe"
		}
	],
	"message": "Liked users fetched successfully"
}
```

### GET "v1/articles/isciting/{id}"
- Não é necessário enviar corpo. / No request body required.
- resposta esperada – `200 Ok` / Expected response – `200 OK`
- Caso o artigo em questão esteja citando outro: / If the article in question is citing another:
```json
{
	"articles": [
		{
			"id": "7de8316c-eeb5-4e48-a057-3b2b439e4ac7",
			"title": "Sample recomended"
		}
	],
	"message": "Article citing information fetched successfully"
}
```

### GET "v1/articles/citedby/{id}"
- Não é necessário enviar corpo. / No request body required.
- resposta esperada – `200 Ok` / Expected response – `200 OK`
- Caso o artigo em questão seja citando por outro: / If the article in question is being cited by another:
```json
{
	"articles": [
		{
			"id": "fad14955-6c2c-4b81-8371-e581f2aab209",
			"title": "Sample Article"
		}
	],
	"message": "Article cited information fetched successfully"
}
```

### PATCH "/articles/share/{id}"
- Não é necessário enviar corpo. / No request body required.
- resposta esperada – `200 Ok` / Expected response – `200 OK`
- Exemplo caso o artigo exista: / Example if this article exists:
```json
{
	"articleID": "fad14955-6c2c-4b81-8371-e581f2aab209",
	"shares": 5,
	"message": "Article shares incremented successfully"
}
```
## Rotas aplicadas sobre usuários:
Routes applied to users:

### GET "v1/getuser/{userID}"
- Não é necessário enviar corpo. / No request body required.
- resposta esperada – `200 Ok` / Expected response – `200 OK`
- Exemplo caso o usuário exista: / Example if this user exists:
```json
{
	"firstName": "Joe",
	"lastName": "Doe",
	"nickname": "johndoe",
	"dateOfBirth": "1989-12-31T21:00:00-03:00",
	"field": "Software Engineering",
	"interests": null,
	"biography": "A passionate software engineer.",
	"lastEducation": "",
	"lcourse": "",
	"currently": "",
	"ccourse": "",
	"openToColab": true,
	"following": [
		"Joe Doe"
	],
	"profilePhoto": " image/5"
}
```

## Rota para todas as imagens:
Route for all images:

### GET "v1/image/{id}"
- Não é necessário enviar corpo. / No request body required.
- resposta esperada – `200 Ok` / Expected response – `200 OK`
- Exemplo caso imagem exista: / Example if this image exists:
```json
{
	"image_base64": "code base64",
    "message": "Image retrieved successfully"
}
```
## Rotas para registro e login:
Routes for registration and login:

### POST "v1/register"
- Corpo: / Body:
```json
{
  "FirstName": "Joe",
  "LastName": "Doe",
  "Email": "exemple@examplee.com",
  "Password": "123456",
  "DateOfBirth": "1990-01-01",
  "Nickname": "johndoe",
  "Field": "Software Engineering",
  "Biography": "A passionate software engineer.",
  "OpenToColab": true,
  "ProfilePhoto": "code base64"
}
```
- resposta esperada – `201 Created` / Expected response – `201 Created`
- Exemplo em caso de sucesso: / Example case user register works:
```json
{
	"user": {
		"first_name": "Ane",
		"last_name": "Doe",
		"email": "exemple@gexample.com"
	},
	"message": "User registered!"
}
```
### POST "v1/login"
- Corpo: / Body:
```json
{
  "Email": "exemple@example.com",
  "Password": "123456"
}
```
- resposta esperada – `200 Ok` / Expected response – `200 OK`
- Exemplo em caso de sucesso no login: / Example case user login works:
```json
{
	"user_id": "5f32d05a-948e-4f7d-879f-2a840988048c",
	"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1c2VyIiwiZXhwIjoxNjkyOTU3NDczLCJzdWIiOiI1ZjMyZDA1YS05NDhlLTRmN2QtODc5Zi0yYTg0MDk4ODA0OGMifQ.sirVG8Bp1zqPxUm-u3el9kSU_C1Px4_gLF4rg9DqXWQ",
	"public_key": "MIIBCgKCAQEArhN2YisPB7Alr1SQFIBiD72R33kco6PqUmVh0lmXpptpva0RpwJsq3Wq6aYxKJ8vTIkdFtQ4HI7h/TADBY1s6VoARnJ66ucLhnnkXLkR9u57tR0IrhT8wdlCs+j3BWJYWilnvTg4a0Rsm4SAS9XfN3g0fH+2Oj6jS7nRWrbq+gQHztyU8zppSjJaLyltC175YHlMOvsnwSjnPLv0N0ldeet7sUFe+50couKqIT0q9kiILB1d/QSz2wWKs9WUgmHS1y83gz3RHqCWZTwvd7LqSMM/2sWMdukc385dB8S9hL370n6P9fpXw8DqXCBfZ11rBXlW2Kx6seEztV3G6DjCSwIDAQAB",
	"message": "Login successful!"
}
```

## RETORNO COMUM PARA ARTIGOS
COMMON RESPONSE FOR ARTICLES

```json
[
	{
		"id": "7de8316c-eeb5-4e48-a057-3b2b439e4ac7",
		"title": "Sample recomended",
		"author_name": "Joe Doe",
		"subject": "Science",
		"field": "Software Engineering",
		"description": "This is a sample article.",
		"keywords": "sample, article, test",
		"submission_date": "2023-08-11T13:55:12.494859-03:00",
		"liked_by": [],
		"shares": 0,
		"cover_image": "link da imagem"
	},
	{
		"id": "fad14955-6c2c-4b81-8371-e581f2aab209",
		"title": "Sample Article",
		"author_name": "Joe Doe",
		"subject": "Science",
		"field": "Physics",
		"description": "This is a sample article.",
		"keywords": "sample, article, test",
		"submission_date": "2023-08-09T21:02:48.271936-03:00",
		"liked_by": [
			"Joe Doe"
		],
		"shares": 10,
		"cover_image": "link da imagem"
	}
]
```