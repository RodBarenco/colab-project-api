If you prefer to read it in [English](#encryption-process-in-the-colab-project)

## Processo de Criptografia no Projeto Colab

1. **Geração de Chaves:**
   - O sistema gera um par de chaves assimétricas para cada automaticamente para administradores (usuários comuns devem registrar sua chave pública mais tarde). Cada par de chaves consiste em uma chave pública e uma chave privada. A chave pública é usada para criptografar os dados, enquanto a chave privada é mantida em segredo pelo administrador/usuário comum e usada para descriptografar.

2. **Envio de Mensagens:**
   - Quando um usuário ou administrador envia uma mensagem para o sistema em uma rota que respondde de forma encriptada, essa criptografia usa a chave pública do destinatário. Isso garante que apenas o destinatário possa descriptografar a mensagem.

3. **Pequenas Respostas Criptografadas:**
   - Se a resposta gerada for pequena o suficiente (por exemplo, menos de 512 bytes), ela é criptografada com a chave pública do destinatário, que pode descriptografar a mensagem diretamente usando sua chave privada.

4. **Grandes Respostas Criptografadas:**
   - Se a resposta for grande (acima do limite definido, como 512 bytes), o sistema gera uma chave AES (Advanced Encryption Standard) aleatória, que é uma chave simétrica. A resposta é então criptografada usando essa chave AES, resultando em um payload criptografado.

5. **Chave AES Criptografada:**
   - A chave AES gerada é criptografada usando a chave pública do destinatário e anexada à resposta criptografada. Agora, a resposta contém tanto o payload criptografado quanto a chave AES criptografada.

6. **Descriptografia de Grandes Respostas:**
   - O destinatário recebe a resposta criptografada e a chave AES criptografada. Primeiro, ele usa sua chave privada para descriptografar a chave AES. Agora, ele possui a chave AES descriptografada.

7. **Descriptografia do Payload:**
   - Usando a chave AES descriptografada, o destinatário pode descriptografar o payload da resposta, obtendo assim a mensagem original.

Esse processo garante a segurança das comunicações, permitindo que os destinatários descriptografem as mensagens recebidas, independentemente do tamanho da resposta, usando suas chaves privadas para pequenas respostas ou a chave AES descriptografada para respostas maiores. É importante manter as chaves privadas em segredo, pois elas são essenciais para descriptografar as mensagens.

---

**Processo de Criptografia pelo Cliente no Projeto Colab:**

1. **Registro da Chave Pública:**
   - Os usuários comuns devem registrar sua chave pública (RSA) no sistema caso queriam enviar mensagens de maneira criptografada, os administradores já recebem suas chaves na hora do registro. 

2. **Envio de Mensagens Criptografadas:**
   - Antes de enviar uma mensagem ao servidor, o cliente verifica se a mensagem é pequena o suficiente (menor que 512 bytes) para ser criptografada apenas com a chave pública disponibilizada pelo servidor no momento do login, sendo grande será necessário usar uma chave AES.

3. **Mensagens Grandes (acima de 512 bytes):** Para mensagens grandes, o cliente deve gerar uma chave AES aleatória, criptografar a mensagem com essa chave AES e, em seguida, criptografa a chave AES com a chave pública disponibilisada pelo servidor.

4. **Formato de Envio para Mensagens Grandes:** O cliente envia o corpo da mensagem da seguinte forma:
   ```JSON
   {
       "aes_key": "hexadecimal da chave AES",
       "encrypted_data": "hexadecimal dos dados criptografados"
   }
   ```

Este processo permite que os clientes criptografem suas mensagens de acordo com o tamanho, protegendo a comunicação. As chaves privadas devem ser mantidas em segredo, pois são essenciais para a descriptografia de mensagens recebidas. A chave pública do servidor também pode ser obtida através da rota `GET /v1/show-pkey`.


-------------------

## Encryption Process in the Colab Project

1. **Key Generation:**
   - The system automatically generates a pair of asymmetric keys for each administrator (regular users must register their public keys later). Each key pair consists of a public key and a private key. The public key is used to encrypt data, while the private key is kept secret by the administrator/regular user and used for decryption.

2. **Message Submission:**
   - When a user or administrator sends a message to the system on a route that responds encrypted, this encryption uses the recipient's public key. This ensures that only the recipient can decrypt the message.

3. **Small Encrypted Responses:**
   - If the generated response is small enough (e.g., less than 512 bytes), it is encrypted with the recipient's public key, allowing the recipient to decrypt the message directly using their private key.

4. **Large Encrypted Responses:**
   - If the response is large (above the defined limit, e.g., 512 bytes), the system generates a random Advanced Encryption Standard (AES) key, which is a symmetric key. The response is then encrypted using this AES key, resulting in an encrypted payload.

5. **Encrypted AES Key:**
   - The generated AES key is encrypted using the recipient's public key and appended to the encrypted response. Now, the response contains both the encrypted payload and the encrypted AES key.

6. **Decryption of Large Responses:**
   - The recipient receives the encrypted response and the encrypted AES key. First, they use their private key to decrypt the AES key. Now, they have the decrypted AES key.

7. **Payload Decryption:**
   - Using the decrypted AES key, the recipient can decrypt the payload of the response, thereby obtaining the original message.

This process ensures communication security, allowing recipients to decrypt received messages, regardless of the response size, using their private keys for small responses or the decrypted AES key for larger responses. It is crucial to keep private keys secret as they are essential for message decryption.

---

**Client-Side Encryption Process in the Colab Project:**

1. **Public Key Registration:**
   - Regular users must register their public key (RSA) in the system if they want to send encrypted messages, while administrators receive their keys during registration.

2. **Sending Encrypted Messages:**
   - Before sending a message to the server, the client checks if the message is small enough (less than 512 bytes) to be encrypted with the public key provided by the server at login. If it's large, an AES key is required.

3. **Large Messages (Above 512 Bytes):** For large messages, the client generates a random AES key, encrypts the message with this AES key, and then encrypts the AES key with the public key provided by the server.

4. **Format for Sending Large Messages:** The client sends the message body in the following format:
   ```JSON
   {
       "aes_key": "AES key hexadecimal",
       "encrypted_data": "hexadecimal of encrypted data"
   }
   ```

This process allows clients to encrypt their messages based on size, ensuring communication security. Private keys must be kept secret, as they are essential for decrypting received messages. The server's public key can also be obtained through the `GET /v1/show-pkey` route.