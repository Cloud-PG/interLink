#!/bin/bash



export INTERLINKCONFIGPATH="$PWD/kustomizations/InterLinkConfig.yaml"

OIDC_ISSUER="${OIDC_ISSUER:-https://dodas-iam.cloud.cnaf.infn.it/}"
AUTHORIZED_GROUPS="${AUTHORIZED_GROUPS:-intw}"
AUTHORIZED_AUD="${AUTHORIZED_AUD:-intertw-vk}"
API_HTTP_PORT="${API_HTTP_PORT:-8080}"
API_HTTPS_PORT="${API_HTTPS_PORT:-443}"
INTERLINKPORT="${INTERLINKPORT:-3000}"
INTERLINKURL="${INTERLINKURL:-http://0.0.0.0}"
INTERLINKPORT="${INTERLINKPORT:-3000}"
INTERLINKURL="${INTERLINKURL:-http://0.0.0.0}"
INTERLINKCONFIGPATH="${INTERLINKCONFIGPATH:-$HOME/.config/interlink/InterLinkConfig.yaml}"

install () {
    mkdir -p $HOME/.local/interlink/logs || exit 1
    mkdir -p $HOME/.local/interlink/bin || exit 1
    mkdir -p $HOME/.config/interlink/ || exit 1
    # download interlinkpath in $HOME/.config/interlink/InterLinkConfig.yaml
    curl -o $HOME/.config/interlink/InterLinkConfig.yaml https://raw.githubusercontent.com/Cloud-PG/interLink/main/kustomizations/InterLinkConfig.yaml

    ## Download binaries to $HOME/.local/interlink/bin
    echo "https://github.com/Cloud-PG/interLink/releases/download/v0.0.1/intertlink-sidecar-slurm_0.0.1_Linux_$(uname -m).tar.gz"
    curl -L -o slurm.tar.gz https://github.com/Cloud-PG/interLink/releases/download/v0.0.1/intertlink-sidecar-slurm_0.0.1_Linux_$(uname -m).tar.gz \
        && tar -xzvf slurm.tar.gz -C $HOME/.local/interlink/bin/
    rm slurm.tar.gz 
    curl -L -o interlink.tar.gz https://github.com/Cloud-PG/interLink/releases/download/v0.0.1/intertlink_0.0.1_Linux_$(uname -m).tar.gz \
        && tar -xzvf interlink.tar.gz -C $HOME/.local/interlink/bin/
    rm interlink.tar.gz

    ## Download oauth2 proxy
    curl -L -o oauth2-proxy-v7.4.0.linux-$(uname -m).tar.gz https://github.com/oauth2-proxy/oauth2-proxy/releases/download/v7.4.0/oauth2-proxy-v7.4.0.linux-$(uname -m).tar.gz
    tar -xzvf oauth2-proxy-v7.4.0.linux-$(uname -m).tar.gz -C $HOME/.local/interlink/bin/
    rm oauth2-proxy-v7.4.0.linux-$(uname -m).tar.gz

}

start () {
    ## Set oauth2 proxy config
    $HOME/.local/interlink/bin/oauth2-proxy \
        --client-id DUMMY \
        --client-secret DUMMY \ 
        --http-address http://0.0.0.0:$API_HTTP_PORT \
        --oidc-issuer-url $OIDC_ISSUER \
        --pass-authorization-header true \
        --provider oidc \
        --redirect-url http://localhost:8081 \
        --oidc-extra-audience intertw-vk \
        --upstream	$INTERLINKURL:$INTERLINKPORT \
        --allowed-group $AUTHORIZED_GROUPS \
        --validate-url ${OIDC_ISSUER}token \
        --oidc-groups-claim groups \
        --email-domain=* \
        --cookie-secret 2ISpxtx19fm7kJlhbgC4qnkuTlkGrshY82L3nfCSKy4= \
        --skip-auth-route="*='*'" \
        --skip-jwt-bearer-tokens true &> $HOME/.local/interlink/logs/oauth2-proxy.log &

    echo $! > $HOME/.local/interlink/oauth2-proxy.pid

    ## start link and sidecar

    $HOME/.local/interlink/bin/interlink &> $HOME/.local/interlink/logs/interlink.log &
    echo $! > $HOME/.local/interlink/interlink.pid

    $HOME/.local/interlink/bin/slurm-sd  &> $HOME/.local/interlink/logs/slurm-sd.log &
    echo $! > $HOME/.local/interlink/slurm-sd.pid
}

stop () {
    kill $(cat $HOME/.local/interlink/oauth2-proxy.pid)
    kill $(cat $HOME/.local/interlink/interlink.pid)
    kill $(cat $HOME/.local/interlink/slurm-sd.pid)
}

case "$1" in
    install)
        install
        ;;
    start) 
        start
        ;;
    stop)
        stop
        ;;
    restart)
        username=${OPTARG}
        ;;
esac
