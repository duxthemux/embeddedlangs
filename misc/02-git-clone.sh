  docker run -d \
    --name git-clone \
    -e USER_UID=1000 \
    -e USER_GID=1000 \
    -e GITHUB_REPOSITORY=http://c2b2f7053ffe62f05da862700f6cb1fc19b840e5@gitea.orb.local:3000/gophercon-2024/motor_calculo.git \
    -e GIT_PUSH_FREQUENCY=5s \
    -e GIT_CLONE_DIR=/data \
    --restart=unless-stopped \
    -v ./stage/gitroot:/data \
    gopherconbr-2024-emb-languages/git-clone:latest
