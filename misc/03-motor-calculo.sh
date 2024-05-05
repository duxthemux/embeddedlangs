  docker run -d \
    --name moto-calculo-service \
    --network=sys \
    -p 8080:8080 \
    --restart=unless-stopped \
    -v ./stage/gitroot/proposta:/data \
    gopherconbr-2024-emb-languages/motor-calculo:latest
