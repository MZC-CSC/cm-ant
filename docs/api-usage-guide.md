# Simple & Sample API Usage Guide

> 이 문서는 프로젝트 GitHub wiki에서 저장소 `docs/`로 이관되었습니다. 앞으로 문서는 wiki가 아닌 저장소 `docs/`에서만 관리합니다.

## 비용 추정 api 호출 흐름

1. 추천 모델 조회
2. 비용 추정 데이터 업데이트 혹은 조회
3. 마이그레이션 프로세스
4. 추정 사용 비용 데이터 업데이트
5. 추정 사용 비용 데이터 조회

---

## 비용 추정 api 명세

> 아래 비용 추정 API의 요청·응답 예시는 dev cm-ant 라이브 환경에서 실제로 호출해 캡처한 값입니다 (2026-07-21 기준). cb-Spider 업데이트로 응답 구조가 바뀌어 그 결과를 반영했습니다.

### 1. 비용 추정 데이터 업데이트 혹은 조회 API

#### 엔드포인트

`POST /api/v1/cost/estimate`

#### 설명

providerName, regionName, instanceType, image(optional) 정보에 대한 비용 추정 데이터를 제공합니다. 만일 Ant 의 데이터베이스에 존재하지 않는 스펙이거나 존재하지만 7 일 이내의 데이터가 아닌 경우 서브 시스템 호출을 통해 비용 정보를 가져와 저장 후 응답을 제공합니다.

#### 요청 본문

```json
{
  "specs": [
    {
      "providerName": "aws",
      "regionName": "ap-northeast-2",
      "instanceType": "t3.small",
      "image": "ubuntu22.04"
    }
  ],
  "specsWithFormat": [
    {
      "commonSpec": "aws+ap-northeast-2+t3.small",
      "commonImage": "aws+ap-northeast-2+ubuntu22.04"
    }
  ]
}
```

- specsWithFormat 과 specs 는 둘 중 하나만 넘어와도 괜찮습니다. (최소 한 쪽은 필요합니다.)
- 모든 요청의 스펙은 취합되어 결과를 제공합니다.
- os, lisence 에 따라 제공되는 가격이 다를 수 있습니다.

`specs`

- `providerName (required):` 요청 프로바이더명 (aws | azure | alibaba | tencent | gcp | ibm)

- `regionName (required)`: CSP 와 매칭되는 region 명

- `instanceType (required)`: CSP 와 매칭되는 instance type

- `image(optional)`: CSP 에 사용되는 os 이미지 정보입니다.

`specsWithFormat`

- `commonSpec (required):` `+` 구분자로 연결된 `providerName+region+instanceType` 정보입니다.

- `commonImage (Optional)`: `+` 구분자로 연결된 `providerName+region+machineImage` 정보입니다.

#### 응답 결과

```json
{
  "successMessage": "Successfully update and get estimate cost info",
  "code": 200,
  "result": {
    "esimateCostSpecResults": [
      {
        "providerName": "aws",
        "regionName": "ap-northeast-2",
        "instanceType": "t3.small",
        "totalMinMonthlyPrice": 18.72,
        "totalMaxMonthlyPrice": 18.72,
        "estimateForecastCostSpecDetailResults": [
          {
            "id": 1,
            "vCpu": "2",
            "memory": "2048 GiB",
            "productDescription": "productFamily= Compute Instance, version= 20260721012550",
            "originalPricePolicy": "OnDemand",
            "pricePolicy": "OnDemand",
            "unit": "PerHour",
            "currency": "USD",
            "price": "0.0260000000",
            "calculatedMonthlyPrice": 18.72,
            "priceDescription": "$0.026 per On Demand Linux t3.small Instance Hour",
            "lastUpdatedAt": "2026-07-21T11:48:09.319204Z"
          }
        ]
      }
    ]
  }
}
```

#### 비용 추정 데이터 조회 (읽기 전용)

이미 저장된 비용 추정 데이터는 별도의 조회 API 로 확인할 수 있습니다. Ant 데이터베이스에서만 조회하므로 서브 시스템 호출 비용이 발생하지 않습니다.

`GET /api/v1/cost/estimate`

```text
Query Param:
providerName(required), regionName(required), instanceType, vCpu, memory, osType, page, size
```

- `providerName (required)`: 조회할 프로바이더명
- `regionName (required)`: 조회할 region 명
- `instanceType`, `vCpu`, `memory`, `osType`: 필터 조건 (optional)
- `page`, `size`: 페이지네이션 (optional)

예: `GET /api/v1/cost/estimate?providerName=aws&regionName=ap-northeast-2&instanceType=t3.small`

```json
{
  "successMessage": "Successfully get price info.",
  "code": 200,
  "result": {
    "estimateCostInfoResult": [
      {
        "id": 1,
        "providerName": "aws",
        "regionName": "ap-northeast-2",
        "instanceType": "t3.small",
        "vCpu": "2",
        "memory": "2048 GiB",
        "productDescription": "productFamily= Compute Instance, version= 20260721012550",
        "originalPricePolicy": "OnDemand",
        "pricePolicy": "OnDemand",
        "unit": "PerHour",
        "currency": "USD",
        "price": "0.0260000000",
        "calculatedMonthlyPrice": 18.72,
        "priceDescription": "$0.026 per On Demand Linux t3.small Instance Hour",
        "lastUpdatedAt": "2026-07-21T11:48:09.319204Z",
        "imageName": ""
      }
    ],
    "resultCount": 1
  }
}
```

---

### 2. 추정 사용 비용 데이터 업데이트 API

#### 엔드포인트

`POST /api/v1/cost/estimate/forecast`

#### 

#### 설명

마이그레이션 된 워크로드에 대한 운영 비용 데이터를 ant 로 불러와 저장합니다.
저장 가능한 기간은 현재 부터 **최대 14일간의 비용데이터만 저장 가능하며, 호출 1건당 Cost Explorer Api 사용 비용 0.01$ 가 연결된 aws 계정에 청구**됩니다.

Tumblebug 에서 사용하는 nsId 와 infraId 에 해당하는 리소스 (현재는 instance id 만 취득 가능) 의 정보를 가지고 와 해당하는 가격 정보를 업데이트합니다.

현재 제공되는 provider 는 aws 이며 추후 확장 예정입니다.

#### 요청 파라미터

```json
{
  "nsId": "mig01",
  "infraId": "infra101"
}
```

- `nsId`: Tumblebug 네임스페이스 ID
- `infraId`: Tumblebug 인프라(구 mci) ID

#### 응답 결과

```json
{
  "successMessage": "Successfully update estimate forecast cost info.",
  "code": 200,
  "result": {
    "fetchedDataCount": 0,
    "updatedDataCount": 0,
    "insertedDataCount": 0
  }
}
```

- `fetchedDataCount`: 서브 시스템에서 조회된 비용 데이터 건수
- `updatedDataCount`: 기존 데이터 중 갱신된 건수
- `insertedDataCount`: 새로 저장된 건수

---

### 3. 추정 사용 비용 데이터 조회 API

#### 엔드포인트

`GET /api/v1/cost/estimate/forecast`

#### 설명

사용자의 필터 조건에 맞는 추정 사용 비용 데이터를 ant 의 데이터베이스에서 조회합니다. 추정 사용 비용  데이터 업데이트를 진행한 후 조회를 진행해야 합니다.
 ant 의 데이터베이스에서 조회하기 때문에 별도의 api 호출 비용이 발생하지 않습니다.

```text
Query Param: startDate(format 2026-07-14), endDate(format 2026-07-21), costAggregationType, provider, nsIds, infraIds, resourceTypes, resourceIds, dateOrder, resourceTypeOrder
```

#### 

#### 요청 파라미터

- startDate: 2026-07-14 형식의 시작 날짜(포함)
- endDate: 2026-07-21 형식의 종료 날짜(포함)
- costAggregationType: 비용 데이터 조회 묶음 간격 (*daily|weekly|monthly)
- provider: provider 로 필터링 하기 위한 값 (현재는 aws 만 제공)
- nsIds: 조회 원하는 nsId 목록
- infraIds: 조회 원하는 infraId 목록
- resourceType: 조회를 원하는 리소스 타입 (VM|VNet|DataDisk|Etc|)
- resourceIds: 조회 원하는 리소스 id 목록
- dateOrder: 날짜 기준 오름차순 내림차순 (asc|desc)
- resourceTypeOrder: resourceType 오름차순 내림차순 (asc|desc)

예: `GET /api/v1/cost/estimate/forecast?startDate=2026-07-14&endDate=2026-07-21&costAggregationType=daily&nsIds=mig01&infraIds=infra101`

#### 

#### 응답 결과

아래는 해당 인프라에 아직 청구된 비용 데이터가 없어 결과가 비어 있는 실제 응답입니다. 추정 사용 비용 데이터 업데이트로 비용 데이터가 쌓이면 `getEstimateForecastCostInfoResults` 목록이 채워지고 `resultCount` 도 함께 증가합니다.

```json
{
  "successMessage": "Successfully get estimate forecast cost",
  "code": 200,
  "result": {
    "getEstimateForecastCostInfoResults": [],
    "resultCount": 0
  }
}
```

목록에 데이터가 있을 때 각 항목은 `category` · `date` · `provider` · `resourceId` · `resourceType` · `totalCost` · `unit` 필드를 가집니다.

---

### 5. 사용 비용 추정 RAW 데이터 업데이트 API

#### 설명

마이그레이션 된 워크로드에 대한 운영 비용 데이터를 ant 로 불러와 저장합니다.
저장 가능한 기간은 현재 부터 최대 14일간의 비용데이터만 저장 가능하며, 호출 1건당 Cost Explorer Api 사용 비용 0.01$ 가 연결된 aws 계정에 청구됩니다.

이는 ns와 infra 와 상관없이 원하는 리소스의 id 를 통해서 비용 데이터를 ant 로 저장할 수 있습니다.

```json
{
  "costResources": [
    {
      "resourceType": "VM",
      "resourceIds": [
        "i-02dxxxxxxdec4"
      ]
    },
    {
      "resourceType": "VNet",
      "resourceIds": [
        "eni-06cxxxxxx02e34"
      ]
    },
    {
      "resourceType": "DataDisk",
      "resourceIds": [
        "vol-0917xxxxxxe2f6"
      ]
    }
  ],
  "awsAdditionalInfo": {
    "ownerId": "xxxxxxxxxxxx",
    "regions": [
      "ap-northeast-2"
    ]
  }
}
```

#### 요청 파라미터

- costResources: 비용 정보를 조회하기 위해 필요한 내용
  - resourceType: 어떤 리소스의 비용 정보를 조회할 것인지 정의 (VM|VNet|DataDisk)
  - resourceIds: 해당 타입에 해당하는 리소스 id 리스트
- awsAdditionalInfo: aws 비용 조회에 추가적으로 필요한 데이터
  - ownerId: 12자리 숫자로 구성된 aws 값
  - regions: eni 조합 시 필요한 region 이름 (VNet 비용 조회 용도)

---

## 성능 평가 api 호출 흐름

1. migration 된 워크로드 선택
   - nsId, mciId, vmId
2. metrics agent 설치 요청
3. 성능 평가를 위한 시나리오 사용자 입력 후 성능 평가 수행
4. 워크로드를 대상으로 수행한 성능 평가 상태 조회
   - 워크로드에 상태 표시가 가능하도록
5. 워크로드를 대상으로 수행한 성능 평가 결과 조회

### 참고 사항

현재 마이그레이션 된 워크로드와 성능 평가 실행 사이에 연결 관계가 프레임워크에 정의되어 있지 않습니다. <br>
성능 평가 흐름상 워크로드와 성능 평가는 링크가 필요해 보이며 지난번 말씀주셨던 것 처럼 연결 고리를 추가 예정입니다.

---

## 성능 평가 api 명세

### 1. 성능평가 지표 수집 에이전트 설치 (optional)

#### 엔드포인트

`POST /api/v1/load/monitoring/agent/install`

#### 설명

성능 평가시 추가적으로 지표 정보를 수집하기 위해 사용할 메트릭 에이전트를 설치하는 단계입니다. 메트릭 에이전트를 통해 부하 발생 시 인프라의 CPU, 메모리, 네트워크 I/O, 디스크 I/O 등의 지표 정보를 수집합니다.

별도로 설치하지 않고 성능 평가를 진행하는 경우 latency 와 request per second 등 http 요청에 대한 지표만 수집합니다.

성능 평가를 수행하면 agent 로 metrics 를 수집하는 옵션을 선택할 수 있습니다.

```json
{
  "nsId": "nsId",
  "mcisId": "mcisId",
  "vmId": "vmId"  # optional 
}
```

#### 요청 프로퍼티

- nsId: 네임스페이스 ID
- mciId: MCI ID
- vmId: VM ID

-----

### 2. 성능 평가용 Load Generator 설치 (optional)

#### 엔드포인트

`POST /api/v1/load/generators`

#### 설명

성능 평가 실행을 위해 부하를 발생하기 위한 Load Generator를 설치하는 단계입니다. 
Local 또는 Remote 환경에 설치할 수 있습니다.

별도로 설치하지 않고 성능 평가를 실행하는 경우 부하 발생 실행 요청의 프로퍼티인 `installLocation` 을 확인하여 부하 발생기 설치를 진행합니다.

```json
{
  "installLocation": "local"  // local | remote
}
```

#### 프로퍼티

- installLocation: 설치 위치 (local 또는 remote)
  
  - local: 현재 ANT가 작동하는 환경(Ubuntu 기반)에 설치
  
  - remote: Tumblebug의 recommendVM을 통해 VM을 프로비저닝한 후 프로비저닝 된 VM 에 설치
    
    | t2.medium | 2   | 4.0 |
    | --------- | --- | --- |

-----

### 3. 성능 평가를 위한 부하 발생 시작

#### 엔드포인트

`POST /api/v1/load/tests/run`

#### 설명

성능 평가를 위한 부하 테스트를 실행하는 단계입니다.

앞선 단계인 성능평가 지표 수집 에이전트 설치 및 성능 평가 부하 발생기 설치하는 과정을 포함하고 있는 API 입니다.

```json
{
        "agentHostname": "",
        "collectAdditionalSystemMetrics": true,
        "httpReqs": [
            {
                "bodyData": "",
                "hostname": "xx.xxx.xx.xxx",
                "method": "get",
                "path": "",
                "port": "80",
                "protocol": "http"
            }
        ],
        "installLoadGenerator": {
            "installLocation": "remote"
        },
        "nsId": "mig01",
        "mciId": "mmci01",
        "vmId": "vm01-1",
        "testName": "first test gogo",
        "virtualUsers": "10",
        "duration": "20",
        "rampUpTime": "20",
        "rampUpSteps": "3"
    }
```

부하 발생 실행 요청은 즉시 반환되며, 응답으로 `loadTestKey` 를 받습니다. 실제 부하 테스트(사전 점검 → 부하 발생기 설치 → 부하 발생 → 결과 수집)는 비동기로 진행되므로, 진행 상태는 아래 **4. 부하 발생 상태 확인** 을 통해 조회합니다.

#### 프로퍼티

- installLoadGenerator: Load Generator 설치 위치 (앞서 설치한 경우 설치 스킵)
  - installLocation: 설치 위치 (remote 또는 local)
- testName: 테스트 이름
- virtualUsers: 가상 사용자 수
- duration: 테스트 실행 시간 (초 단위)
- rampUpTime: 부하 증가 시간 (초 단위)
- rampUpSteps: 부하 증가 단계 수
- hostname: 테스트 대상 호스트명
- port: 포트 번호
- agentHostname: 에이전트가 설치된 호스트명 (Optional)
- collectAdditionalSystemMetrics: 에이전트를 통해 메트릭 수집 여부 확인
- httpReqs: HTTP 요청 리스트
  - method: HTTP 메소드 (GET, POST)
  - protocol: 통신 프로토콜 (http, https)
  - hostname: 요청할 호스트명
  - port: 포트 번호
  - path: 요청 경로
  - bodyData: 요청 본문 데이터 (Optional)

-----

### 4. 부하 발생 상태 확인

#### 엔드포인트

부하 발생 상태는 두 가지 방법으로 조회할 수 있습니다.

- `GET /api/v1/load/tests/state/{loadTestKey}` — 특정 테스트 키의 실행 상태 조회
- `GET /api/v1/load/tests/state/last` (operationId `GetLastLoadTestExecutionState`) — 특정 워크로드(node)에 대해 **마지막으로 수행된** 실행 상태 조회. node 는 쿼리 파라미터 `nsId` / `mciId` / `vmId` 로 지정합니다.

> node id 는 이름이라 재사용될 수 있으므로, "이 워크로드의 마지막 실행"을 계속 보여주는 화면은 응답의 `nodeUid` 로 그 실행이 이후 교체된 VM의 것인지 구분해야 합니다.

#### 설명

부하 발생 상태를 확인합니다. 부하 발생이 진행되고 있는 상태와 정보를 반환합니다.

```text
Path Param
loadTestKey: 성능 평가 실행 시 반환된 테스트 키
```

응답에는 실행의 진행을 세부적으로 표현하기 위한 필드가 포함됩니다.

- `executionStatus`: 실행 전체 상태. 값은 `on_processing`(사전 점검·설치·부하 발생 등 진행 중) · `on_fetching`(부하는 끝났고 결과 수집 중) · `successed`(완료, 결과 조회 가능) · `test_failed`(완료 전 중단, 원인은 `failureMessage`)
- `startAt` / `finishAt`: 실행 시작·종료 시각 (진행 중이면 `finishAt` 은 null)
- `expectedFinishAt`, `totalExpectedExecutionSecond`: duration + ramp-up 으로 계산한 예상 종료. 부하 이전의 점검·설치는 포함하지 않으므로 진행 바의 힌트로만 사용합니다.
- `failureMessage`: 실패 시 한 줄 사유 (성공 시 비어 있음)
- `nodeUid`: 대상 node 의 uid. node id 는 재사용되므로 실행이 어느 VM에 속하는지 식별하는 데 사용합니다.
- `steps`: 실행의 단계별 진행을 **트리**로 표현합니다. 최상위는 실행이 거치는 *phase* 이고, 각 phase 의 하위 단계는 `children` 에 담깁니다. phase 만 필요한 경우 `children` 을 무시하고 phase 행만 읽으면 됩니다.

`steps[]` 각 노드(phase 또는 sub-step)의 필드:

- `seq`: 실행 내 순서
- `name`: 단계 식별자. phase 는 단일 이름(`precheck`), sub-step 은 `phase.sub` 형태(`precheck.target_reachable`)
- `status`: `pending` · `running` · `ok` · `failed` · `skipped`
- `message`: 현재 상태를 나타내는 짧은 한 줄 (지금 무엇을 하는지 표시할 때 사용)
- `detail`: 실패 시 보여줄 상세 진단·오류 원인 (여러 줄일 수 있음)
- `elapsedSec`: 해당 단계 소요 시간. 완료된 단계는 전체 구간, 진행 중이면 지금까지의 시간. phase 값은 children 에서 합산됩니다.
- `children`: phase 의 하위 단계 목록

#### 응답 예시

```json
{
        "code": 200,
        "result": {
            "compileDuration": "436.008µs",
            "createdAt": "2024-11-05T06:22:08.832848Z",
            "executionDuration": "1m4.032562237s",
            "executionStatus": "on_processing",
            "id": 1,
            "loadGeneratorInstallInfo": {
                "createdAt": "2024-11-05T06:19:27.85666Z",
                "id": 1,
                "installLocation": "local",
                "installPath": "/opt/ant/jmeter",
                "installType": "jmeter",
                "installVersion": "5.6",
                "status": "installed",
                "updatedAt": "2024-11-05T06:22:08.824086Z"
            },
            "loadTestKey": "1730787567844-loadtest-key-example",
            "nodeUid": "vm-9f2c1a7b-3e40-4d21-9c88-0a1b2c3d4e5f",
            "startAt": "2024-11-05T06:22:08.828219Z",
            "finishAt": null,
            "expectedFinishAt": "2024-11-05T06:23:08.828219Z",
            "totalExpectedExecutionSecond": 60,
            "failureMessage": "",
            "updatedAt": "2024-11-05T06:23:12.869756Z",
            "steps": [
                {
                    "seq": 1,
                    "name": "precheck",
                    "status": "ok",
                    "message": "Checking the environment",
                    "elapsedSec": 5,
                    "children": [
                        {
                            "seq": 1,
                            "name": "precheck.target_exists",
                            "status": "ok",
                            "elapsedSec": 1
                        },
                        {
                            "seq": 3,
                            "name": "precheck.target_reachable",
                            "status": "ok",
                            "message": "Target port 80 reachable",
                            "elapsedSec": 3
                        },
                        {
                            "seq": 4,
                            "name": "precheck.metric_port_open",
                            "status": "skipped",
                            "message": "Metrics not requested"
                        }
                    ]
                },
                {
                    "seq": 2,
                    "name": "generator_install",
                    "status": "ok",
                    "elapsedSec": 12
                },
                {
                    "seq": 3,
                    "name": "agent_install",
                    "status": "skipped",
                    "message": "Metrics not requested"
                },
                {
                    "seq": 4,
                    "name": "jmx_prepare",
                    "status": "ok",
                    "elapsedSec": 1
                },
                {
                    "seq": 5,
                    "name": "jmeter_run",
                    "status": "running",
                    "message": "Generating load",
                    "elapsedSec": 22
                },
                {
                    "seq": 6,
                    "name": "result_fetch",
                    "status": "pending",
                    "children": [
                        {
                            "seq": 1,
                            "name": "result_fetch.file_result",
                            "status": "pending"
                        },
                        {
                            "seq": 2,
                            "name": "result_fetch.file_cpu",
                            "status": "pending"
                        }
                    ]
                }
            ]
        },
        "successMessage": "Successfully retrieved load test execution state information"
    }
```

> `steps` 트리(phase ↔ sub-step)의 전체 구조·각 phase 의 하위 단계·`result_fetch` 결과 파일 수집 방식 등 상세는 [docs/load-test-status-api.md](load-test-status-api.md) 를 참고하세요. 실패 시 각 단계 메시지의 의미는 [docs/load-test-troubleshooting.md](load-test-troubleshooting.md) 를 참고하세요.

------

### 5. 성능 평가 결과 확인

#### 엔드포인트

`GET /api/v1/load/tests/result`

`GET /api/v1/load/tests/result/metrics`

#### 설명

부하 발생이 완료된 후, 성능 평가 결과를 조회하는 단계입니다.

```text
Query Param

loadTestKey: 성능 평가 실행 시 반환된 테스트 키
format: 결과 형식 (normal | aggregate)
```

만약 모니터링 에이전트가 설치된 경우, 메트릭 에이전트를 통해 수집된 cpu, memory, disk, network 등 데이터를 metrics 결과 조회 API 호출을 통해 확인해볼 수 있습니다.

```text
Query Param:
loadTestKey: 성능 평가 실행 시 반환된 테스트 키
```

#### 성능 평가 결과 조회 응답

aggregate 의 경우 결과에 대한 집계된 값입니다. 데이터 샘플링 없이 모든 데이터에 대한 집계된 값을 제공합니다.

normal 은 시계열 데이터로 표현되는 응답값입니다. 결과데이터는 100ms 단위로 평균값을 계산하여 가장 근접한 값을 찾아내 샘플링한 결과값을 추출하여 제공합니다.

```json
 "[aggregate]": {
    "code": 0,
    "errorMessage": "string",
    "result": [
      {
        "average": 0,  # 평균 latency (ms)
        "errorPercent": 0,  # error 비율
        "label": "string", # 테스트 이름
        "maxTime": 0,  # 최대 latency  (ms)
        "median": 0,  # 중간 값 latency (ms)
        "minTime": 0,   # 최소 latency (ms)
        "ninetyFive": 0,  # latency 95% 값 (ms)
        "ninetyNine": 0,  # latency 99% 값 (ms)
        "ninetyPercent": 0,  # latencyt 90% 값 (ms)
        "receivedKB": 0,  # 받은 KB / sec
        "requestCount": 0,  # 발생한 총 요청 count
        "sentKB": 0,  # 보낸 KB / sec
        "throughput": 0  # 요청에 대한 처리량 / sec
      }
    ],
    "successMessage": "string"
  },
  "[normal]": {
    "code": 0,
    "errorMessage": "string",
    "result": [
      {
        "label": "string",
        "results": [  # 단일 요청에 대한 단일 결과 값들의 목록
          {
            "bytes": 0,
            "connection": 0,
            "elapsed": 0,  # 개별 요청 총 소요 시간 (ms)
            "idleTime": 0,
            "isError": true,  # 에러 여부
            "latency": 0,  # 개별 요청 지연 시간 (ms)
            "no": 0,
            "sentBytes": 0,
            "timestamp": "string",  # 개별 요청 발생 시간
            "url": "string"
          }
        ]
      }
    ],
    "successMessage": "string"
  }
}
```

#### 성능 평가 지표 조회 응답

```json
{
  "code": 0,
  "errorMessage": "string",
  "result": [
    {
      "label": "string",  # 아래 수집 지표에 대한 목록에서 확인 가능
      "metrics": [  
        {
          "isError": true,  # 해당 요청에 대한 에러 여부
          "timestamp": "string",  # 개별 요청 발생 시간
          "unit": "string",  # 아래 수집 지표에 대한 목록에서 확인 가능
          "value": "string"  # 수집 지표에 대한 값
        }
      ]
    }
  ],
  "successMessage": "string"
}
```

##### 성능 평가 수집 지표에 대한 목록

```json
"cpu_all_combined": {  # 전체 cpu 사용률
        Unit:     "%",
    },
    "cpu_all_idle": {  # 유휴 cpu 비율
        Unit:     "%",
    },
    "memory_all_used": {  # 메모리 사용률
        Unit:     "%",
    },
    "memory_all_free": {  # 유휴 메모리 비율
        Unit:     "%",
    },
    "memory_all_used_kb": {  # 메모리 사용 kb
        Unit:     "mb",
    },
    "memory_all_free_kb": {  # 유휴 메모리 kb
        Unit:     "mb",
    },
    "disk_read_kb": {  # 디스크 읽기 kb
        Unit:     "kb",
    },
    "disk_write_kb": {   # 디스크 쓰기 kb
        Unit:     "kb",
    },
    "disk_use": {  # 디스크 사용율
        Unit:     "%",
    },
    "disk_total": {  # 총 디스크 용량 mb
        Unit:     "mb",
    },
    "network_recv_kb": {  # 네트워크 수신 kb
        Unit:     "kb",
    },
    "network_sent_kb": {  # 네트워크 발신 kb
        Unit:     "kb",
    },
```
