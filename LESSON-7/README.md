# log-analyzer — Guia de Análise de Logs no Linux

Contexto de referência para o custom agent `log-analyzer` (`.kiro/agents/log-analyzer.json`).
Este agente é **somente leitura**: inspeciona e diagnostica, mas nunca altera, rotaciona,
trunca ou apaga logs, e nunca para/reinicia serviços ou containers.

> Cuidado com segredos/PII em logs (tokens, senhas, e-mails). Referencie por chave/contexto,
> não reproduza valores sensíveis.

---

## 1. journald (systemd)

```bash
# Erros e acima, na última hora
journalctl -p err..alert --since "1 hour ago"

# Logs de um serviço específico
journalctl -u nginx.service --since today

# Desde o último boot, seguindo o kernel
journalctl -k -b

# Janela de tempo precisa
journalctl --since "2026-09-26 08:00" --until "2026-09-26 09:30"

# Por PID ou executável
journalctl _PID=1234
journalctl /usr/sbin/sshd

# Saída em JSON para pós-processamento
journalctl -u myapp -o json --since "30 min ago"
```

## 2. Arquivos em /var/log

```bash
# Autenticação / acesso negado (Debian: auth.log | RHEL: secure)
grep -Ei "failed|invalid|denied|refused" /var/log/auth.log

# Kernel: OOM, segfault, panic
grep -Ei "oom|out of memory|segfault|panic|call trace" /var/log/kern.log

# Logs comprimidos e rotacionados
zgrep -Ei "error|fatal" /var/log/syslog.*.gz
zcat /var/log/nginx/error.log.2.gz | tail -n 100

# Nginx / Apache
tail -n 200 /var/log/nginx/error.log
grep -E " 5[0-9]{2} " /var/log/nginx/access.log | awk '{print $9}' | sort | uniq -c | sort -rn

# Top IPs em access log
awk '{print $1}' /var/log/nginx/access.log | sort | uniq -c | sort -rn | head
```

## 3. Padrões úteis de busca

```bash
# Erros comuns em qualquer log
grep -Ein "error|fatal|panic|timeout|refused|denied|exception|traceback" <arquivo>

# Contagem de ocorrências por tipo
grep -Eio "error|warn|fatal" app.log | sort | uniq -c

# Contexto ao redor de um erro (3 linhas antes/depois)
grep -n -C3 "OutOfMemory" app.log
```

## 4. Containers — Docker / Podman

```bash
# Últimos logs com timestamp e janela de tempo
docker logs --since 30m --timestamps <container>
docker logs --tail 200 <container>
podman logs --since 1h -t <container>

# Compose
docker compose logs --since 15m --tail 100 <serviço>

# Estado do container: restart count, OOMKilled, health
docker inspect -f '{{.RestartCount}} {{.State.OOMKilled}} {{.State.Health.Status}}' <container>

# Eventos do runtime (reinícios, kills, health)
docker events --since 1h --filter event=die --filter event=oom
podman events --since 1h
```

## 5. Kubernetes

```bash
# Logs do container atual e do anterior (após crash)
kubectl logs <pod> -c <container> --since=15m
kubectl logs <pod> -p          # container anterior (previous)

# Eventos e motivos de reinício / OOMKilled
kubectl describe pod <pod>
kubectl get events --sort-by=.lastTimestamp

# Nível de node (containerd/CRI-O)
crictl ps -a
crictl logs <container-id>
```

## 6. lnav — navegação e correlação multi-arquivo

```bash
# Abrir vários logs (detecção automática de formato + timeline unificada)
lnav /var/log/syslog /var/log/nginx/error.log

# Consulta SQL sobre logs estruturados (dentro do lnav, prompt `;`)
#   ;SELECT log_time, log_level, log_body FROM logline WHERE log_level >= 'error';

# Modo batch (não interativo)
lnav -n -c ';SELECT count(*), log_level FROM logline GROUP BY log_level' /var/log/syslog
```

## 7. logcli / Loki (LogQL)

```bash
# Descoberta de labels e streams
logcli labels
logcli series '{job="nginx"}'

# Consulta por janela de tempo
logcli query '{job="myapp"} |= "ERROR"' --since=1h --limit=200

# Taxa de erros por minuto
logcli query 'sum(rate({job="myapp"} |= "ERROR" [1m]))' --since=30m

# Consulta instantânea
logcli instant-query 'count_over_time({job="nginx"} |= "500" [5m])'
```

---

## Metodologia recomendada

1. **Escopo** — identifique serviço/host afetado e a janela de tempo do incidente.
2. **Coleta** — filtre por severidade (`-p err..alert`) e período (`--since/--until`).
3. **Correlação** — cruze eventos entre fontes (app ↔ kernel ↔ runtime ↔ proxy).
4. **Causa raiz** — isole o primeiro evento causal, não apenas o sintoma final.
5. **Remediação** — recomende comandos de correção para o usuário executar (o agente não os executa).
