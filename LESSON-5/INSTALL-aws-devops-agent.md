# Instalação — AWS DevOps Agent Power

Registro da instalação do power [`aws-devops-agent`](https://github.com/kirodotdev/powers/tree/main/aws-devops-agent) neste workspace, via caminho de autenticação **Opção B (SigV4)**.

- **Data:** 2026 (LESSON-3)
- **SO:** Linux (Pop!_OS 24.04 / kernel 7.1.5)
- **Workspace:** `LESSON-3`
- **Caminho de auth escolhido:** Opção B — SigV4 (proxy local com credenciais AWS)
- **Região do AgentSpace:** `us-east-1`

---

## O que é este power

Não é um pacote instalável via npm/pip. É um **Kiro Power**: um conjunto de arquivos de configuração que você copia para o `.kiro/` do workspace:

- `mcp.json` → define os servidores MCP (proxies para a API do AWS DevOps Agent)
- `steering/*.md` → regras de comportamento (roteamento de intenção, workflows, troubleshooting)

Fornece tools de inteligência operacional AWS: investigação de incidentes, otimização de custos, revisão de arquitetura, mapa de topologia, chat com o agente, release testing e release readiness review.

---

## Estrutura do power (origem)

Do repositório `kirodotdev/powers`, diretório `aws-devops-agent/`:

```
aws-devops-agent/
├── mcp.json                  # 3 servidores: remoto (bearer), sigv4, aws-mcp (fallback)
├── POWER.md                  # metadados + documentação do power
└── steering/
    ├── setup.md              # setup, diagnóstico e troubleshooting
    ├── steering.md           # roteamento de tools, fallback, tratamento de erros (alwaysApply)
    ├── release-testing.md    # workflow de release testing (UI/API)
    ├── release-readiness.md  # workflow de release readiness review
    └── ecs-incident-walkthrough.md  # exemplo completo de incidente ECS 503
```

---

## Passos executados

### 1. Inspeção do power (clone temporário)

Feito um sparse clone raso do repositório para ler os arquivos, depois removido:

```bash
cd /tmp && git clone --depth 1 --filter=blob:none --sparse \
  https://github.com/kirodotdev/powers.git powers-check
cd powers-check && git sparse-checkout set aws-devops-agent
# ... leitura dos arquivos ...
rm -rf /tmp/powers-check    # limpeza
```

### 2. Verificação de ferramentas do ambiente

| Ferramenta | Status | Observação |
|------------|--------|------------|
| `git` | ✅ 2.43.0 | |
| `curl` | ✅ | |
| `node` / `npm` | ✅ v24.21.0 | via nvm |
| `python3` / `pip3` | ✅ | via pyenv |
| `uvx` | ⚠️ ausente → **instalado** | necessário p/ SigV4 |
| `aws` CLI | ✅ 2.37.4 | credenciais pendentes |

### 3. Instalação do `uv`/`uvx` (dependência da Opção B)

```bash
curl -LsSf https://astral.sh/uv/install.sh | sh
# instalado em: /home/<USER>/.local/bin/{uv,uvx}
# versão: uvx 0.12.19
```

> ⚠️ `~/.local/bin` pode não estar no PATH que o Kiro passa aos servidores MCP.
> Por isso o `mcp.json` usa o **caminho absoluto** `/home/<USER>/.local/bin/uvx` no campo `command`.

### 4. Cópia dos steering files para o workspace

```bash
mkdir -p .kiro/steering
cp /tmp/powers-check/aws-devops-agent/steering/*.md .kiro/steering/
```

Resultado — arquivos ativos em `.kiro/steering/`:

- `setup.md`
- `steering.md` (alwaysApply)
- `release-testing.md`
- `release-readiness.md`
- `ecs-incident-walkthrough.md`

### 5. `mcp.json` (⚠️ PENDENTE — ação manual do usuário)

A gravação em `.kiro/settings/mcp.json` é **bloqueada por política de permissões** do agente.
**Você precisa criar este arquivo manualmente.** Como a região é `us-east-1`, ela está fixada
diretamente (sem depender da env var `${DEVOPS_AGENT_REGION}`, o que dispensa aprová-la no Kiro):

Inclui `AWS_PROFILE` para o proxy assumir a sessão SSO correta:

```json
{
  "mcpServers": {
    "aws-devops-agent-sigv4": {
      "command": "uvx",
      "timeout": 120000,
      "transport": "stdio",
      "args": [
        "mcp-proxy-for-aws@latest",
        "https://connect.aidevops.us-east-1.api.aws/mcp",
        "--service", "aidevops",
        "--region", "us-east-1",
        "--profile", "<AWS_ACCOUNT_ID>:<SSO_PROFILE_NAME>"
      ],
      "env": {
        "AWS_PROFILE": "<AWS_ACCOUNT_ID>:<SSO_PROFILE_NAME>"
      }
    }
    // ... aws-mcp e demais servidores preservados como já estavam ...
  }
}
```

> O profile é fixado via `--profile` nos **args** (flag do `mcp-proxy-for-aws v1.7.0`) — isso
> passa como argumento de linha de comando e **não** depende da allowlist do Kiro (MCP Approved
> Env Vars). O `env.AWS_PROFILE` é reforço redundante. A sessão SSO precisa estar ativa:
> `aws sso login --sso-session <SSO_SESSION_NAME>`.

> Nota: o `aws-mcp` (fallback) **já existia** no `mcp.json` global com o mesmo profile (via
> `--metadata AWS_PROFILE=...`). Apenas o servidor primário `aws-devops-agent-sigv4` foi
> adicionado; os demais servidores (`enel-api`, `iam-policy-autopilot`, `lean-ctx`, `obsidian-*`)
> foram preservados.

---

## Passos pendentes

Restou **um** item, do lado da AWS:

1. ~~Criar `.kiro/settings/mcp.json`~~ — ✅ **feito**. Como a escrita direta em `~/.kiro/settings/`
   é bloqueada por política do Kiro (`kiro-scope`), o agente gerou `~/mcp.json.novo` (merge completo,
   validado) e o usuário copiou por cima. O `aws-devops-agent-sigv4` já está no arquivo ativo.

2. ~~Autenticar via SSO~~ — ✅ **feito** (`aws sso login --sso-session <SSO_SESSION_NAME>`). Refazer
   quando a sessão expirar. Profile: `<AWS_ACCOUNT_ID>:<SSO_PROFILE_NAME>`.

3. ~~Aprovar env var no Kiro~~ — ✅ **dispensável**: profile fixado via `--profile` nos args.

4. **Criar um AgentSpace** ⚠️ **PENDENTE — único bloqueio restante**. `list-agent-spaces`
   retorna `agentSpaces: []`. Crie um AgentSpace no console do AWS DevOps Agent (Operator Web App)
   na conta `<AWS_ACCOUNT_ID>`, região `us-east-1`. Sem ele, os servidores MCP conectam e autenticam,
   mas `chat`/`investigate` não têm alvo (toda operação exige um `agentSpaceId`).

5. **Reiniciar o Kiro** — ⏳ se ainda não reiniciou após copiar o `mcp.json`.

---

## Verificação pós-setup

Depois de configurar credenciais + região:

```bash
export PATH="$HOME/.local/bin:$PATH"
uvx --version                                  # proxy disponível
aws sts get-caller-identity                    # credenciais válidas
uvx mcp-proxy-for-aws@latest --help            # o proxy sobe
aws devops-agent list-agent-spaces --region us-east-1   # acesso ao DevOps Agent
```

Se `list-agent-spaces` retornar `AccessDeniedException` → falta `AIDevOpsAgentFullAccess` na role.

---

## Estado atual (snapshot desta sessão)

| Item | Estado |
|------|--------|
| `uvx` instalado | ✅ `/home/<USER>/.local/bin/uvx` (0.12.19) |
| AWS CLI | ✅ 2.37.4 |
| `mcp-proxy-for-aws` | ✅ v1.7.0 (tem flag `--profile`) |
| Steering files copiados | ✅ `.kiro/steering/` (5 arquivos) |
| `.kiro/settings/mcp.json` | ✅ aplicado (via `~/mcp.json.novo` → copiado pelo usuário) |
| `aws-devops-agent-sigv4` no arquivo ativo | ✅ confirmado |
| Profile fixado (`--profile` nos args) | ✅ não depende da allowlist do Kiro |
| Credenciais AWS (SSO) | ✅ conta `<AWS_ACCOUNT_ID>`, role `AdministratorAccess` |
| Acesso à API DevOps Agent | ✅ `list-agent-spaces` respondeu (sem AccessDenied) |
| Região | ✅ `us-east-1` (fixada no mcp.json) |
| **AgentSpace provisionado** | ❌ **nenhum** (`agentSpaces: []`) — criar no Operator Web App |
| Reiniciar Kiro | ⏳ após copiar o mcp.json |

**Conclusão:** setup local e de autenticação 100% completo. Único bloqueio para uso efetivo das
tools (`chat`, `investigate`, etc.) é a **inexistência de um AgentSpace** na conta/região — precisa
ser provisionado no Operator Web App do AWS DevOps Agent.

---

## Troubleshooting rápido (SigV4)

| Erro | Causa | Correção |
|------|-------|----------|
| `aws-mcp`/sigv4 sem tools | `uvx` não instalado ou creds ausentes | `uvx --version`, depois `aws sts get-caller-identity` |
| `MCP error -32000: Connection closed` | proxy saiu — creds expiradas ou `uvx` fora do PATH | usar caminho absoluto no `command`; revalidar creds |
| `ExpiredTokenException` | credenciais AWS expiradas | `aws sso login` |
| `AccessDeniedException` | falta permissão IAM | anexar `AIDevOpsAgentFullAccess` |
| tools não aparecem após configurar | MCP inicializa no launch do Kiro | **reiniciar o Kiro** |

Referência completa: `.kiro/steering/setup.md`.
