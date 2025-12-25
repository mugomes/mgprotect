# MGProtect

MGProtect é uma biblioteca em Go para proteção de software via serial e validação de licença. Ela permite criar sistemas de licenciamento seguros que verificam:

- Serial válido
- Produto e versão correta
- Integridade e assinatura digital
- Máquina específica (via hash do hardware)

**Importante:** Para usar MGProtect corretamente, você precisa gerar previamente a chave pública, a chave interna e o serial utilizando a ferramenta [mggenerator](https://www.mugomes.com.br/p/mgprotectedgenerator.html). Esta ferramenta garante que os dados de licença sejam assinados e criptografados de forma compatível com MGProtect.

## Funcionalidades Principais

- Validação de serial usando ed25519
- Validação de licença por máquina
- Assinatura HMAC de dados locais
- Suporte a múltiplas versões de produto

## Instalação

`go get github.com/mugomes/mgprotect`

## Exemplo de Uso

```
package main

import (
    "crypto/ed25519"
    "fmt"
    "mugomes/mgprotect"
)

func main() {
    // Inicializa a proteção
    mgp := mgprotect.New()

    // Configura produto, versão e chaves geradas pelo mggenerator
    mgp.SetProductID(0x01)
    mgp.SetMajorVersion(1)
    
    // pubKey deve vir do mggenerator
    mgp.SetPublicKey(ed25519.PublicKey{valorporvirgula})

    // Chave interna para HMAC
    mgp.SetInternalKey([]byte("SUA_CHAVE_INTERNA_DO_MGGENERATOR"))

    // Chave de criptografia do serial
    mgp.SetK("SUA_CHAVE_K_DO_MGGENERATOR")

    // Validar um serial fornecido pelo usuário
    serial := "digite o serial aqui"
    result := mgp.Validate(serial)

    switch result {
    case 1:
        fmt.Println("Serial válido!")
    case mgprotect.ERRO_SERIAL_INVALIDO:
        fmt.Println("Serial inválido")
    case mgprotect.ERRO_CHECKSUM_INVALIDO:
        fmt.Println("Checksum inválido")
    case mgprotect.ERRO_PRODUTO_INVALIDO:
        fmt.Println("Produto inválido")
    case mgprotect.ERRO_LICENCA_NAO_VALIDA_PARA_VERSAO:
        fmt.Println("Licença não válida para esta versão")
    case mgprotect.ERRO_ASSINATURA_INVALIDA:
        fmt.Println("Assinatura inválida")
    }
}
```

Para inicialização do software você pode usar o validador dessa forma:

```
// Salva a licença no arquivo
err := mgp.SaveLicense(licensePath, serial)
if err != nil {
    fmt.Println("Erro ao salvar licença:", err)
    return
}
fmt.Println("Licença salva com sucesso!")

// Agora podemos carregar e validar a licença
if mgp.LoadAndValidate(licensePath) {
    fmt.Println("Licença válida!")
} else {
    fmt.Println("Licença inválida ou não corresponde a esta máquina")
}
```

A biblioteca pode ser usada em projetos comerciais ou pessoais, mas o uso seguro requer a geração correta das chaves e serial via [mggenerator](https://www.mugomes.com.br/p/mgprotectedgenerator.html).

## Information

 - [Page MGProtect](https://www.mugomes.com.br/p/mgprotect.html)

## Requirement

 - Go 1.25.5

## Support

- GitHub: https://github.com/sponsors/mugomes
- More: https://www.mugomes.com.br/p/apoie.html

## License

Copyright (c) 2025 Murilo Gomes Julio

Licensed under the [MIT](https://github.com/mugomes/mgprotect/blob/main/LICENSE) license.

All contributions to the MGProtect are subject to this license.
