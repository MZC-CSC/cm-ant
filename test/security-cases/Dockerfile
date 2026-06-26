# 보안 게이트(Trivy misconfig) 발화 검증용 더미 Dockerfile (테스트 전용 · 미머지)
# 의도적 미스컨피그: 태그 미고정(latest) + 마지막 USER 가 root.
FROM ubuntu:latest

RUN apt-get update && apt-get install -y curl

USER root

CMD ["bash"]
