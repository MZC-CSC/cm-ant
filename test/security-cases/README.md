# 보안 게이트 검증용 의도적 결함 케이스 (테스트 전용 · 미머지)

이 폴더의 파일들은 **보안 게이트(security-gate.yml)가 실제로 결함을 잡는지** 확인하기 위해 *일부러 심은 더미 결함*이다. 모두 가짜/더미 값이며 실제 시크릿·실제 취약 코드가 아니다.

| 파일 | 의도한 결함 | 기대 발화 게이트 |
|------|-----------|----------------|
| `fake_key.txt` | 더미 PRIVATE KEY 블록 (`.pem`은 `.gitignore`로 추적 제외라 `.txt`에 넣음) | gitleaks(private-key) |
| `fake_aws.txt` | 더미 AWS 키 | gitleaks(aws-access-token·generic-api-key) |
| `insecure_sample.go` | 하드코딩 비밀번호(G101)·약한 난수(G404) | gosec |
| `Dockerfile` | `latest` 태그·`USER root` 미스컨피그 | Trivy(misconfig) |

> 이 브랜치(`test-security-gate-cases`)는 develop·upstream에 **절대 머지하지 않는다.** 게이트 검증이 끝나면 폐기한다. (cmig-workflow POLICY-SECURITY §9.1)
