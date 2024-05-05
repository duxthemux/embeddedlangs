  docker run -d \
    --name gitea \
    -e USER_UID=1000 \
    -e USER_GID=1000 \
    -p 3000:3000 \
    -p 222:22 \
    --restart=unless-stopped \
    -v ./stage/gitea:/data \
    gitea/gitea:latest

#c2b2f7053ffe62f05da862700f6cb1fc19b840e5