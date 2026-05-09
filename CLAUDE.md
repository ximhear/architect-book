# Architect Book

대한민국 서비스를 설계·구축할 수 있는 **소프트웨어 아키텍트**가 되기 위한 학습서를 집필하는 저장소.
독자는 클라우드 인프라부터 백엔드·프론트엔드까지 서비스 전체를 조망하고, 스스로 의사결정할 수 있는 역량을 갖추는 것을 목표로 한다.

## 집필 목적과 독자

- **목적**: "한 명의 아키텍트가 대한민국 환경에서 동작하는 실서비스를 처음부터 끝까지 설계할 수 있게 한다"
- **주독자**:
  - 시니어 백엔드/프론트엔드 개발자 → 아키텍트 전환을 준비
  - 스타트업 CTO/리드 개발자 → 서비스 전체를 책임져야 하는 위치
  - 클라우드/인프라 엔지니어 → 애플리케이션 영역까지 확장하려는 사람
- **전제 지식**: 한 가지 이상 언어로 3년 이상 실무 경험. 기초 문법은 다루지 않는다.

## 책의 범위 (Scope)

전체 서비스를 구성하는 데 필요한 모든 레이어를 포괄한다.

### 1. 아키텍처 기초
- 아키텍처 스타일: Monolith, Modular Monolith, Microservices, Event-Driven, Serverless
- 도메인 주도 설계(DDD), 헥사고날/클린 아키텍처
- 비기능 요구사항(성능, 가용성, 보안, 확장성, 운영성) 트레이드오프
- C4 모델, ADR(Architecture Decision Records)

### 2. 클라우드 / 인프라
- 국내 우선: **NCP(Naver Cloud), KT Cloud, NHN Cloud** + AWS/GCP/Azure
- 컨테이너/오케스트레이션: Docker, Kubernetes, Helm
- IaC: Terraform, Pulumi
- 네트워킹: VPC, 로드밸런서, CDN(국내 망 특수성 포함)
- 관측성: Prometheus, Grafana, OpenTelemetry, Loki, Datadog
- CI/CD: GitHub Actions, ArgoCD, GitLab CI

### 3. 백엔드
- 언어/프레임워크: **Java(Spring Boot), Kotlin, Go, Node.js(NestJS), Python(FastAPI)**
- API 설계: REST, GraphQL, gRPC
- 인증/인가: OAuth2, OIDC, JWT, **PASS·KISA 본인확인, 간편로그인(카카오/네이버)**
- 데이터: RDB(PostgreSQL, MySQL), NoSQL(MongoDB, Redis, DynamoDB), 검색(Elasticsearch/OpenSearch)
- 메시징/스트리밍: Kafka, RabbitMQ, NATS
- 캐시 / 동시성 / 분산 트랜잭션 / Saga / Outbox 패턴

### 4. 프론트엔드
- 웹: **React, Next.js, Vue/Nuxt, Svelte/SvelteKit**
- 모바일: React Native, Flutter, 네이티브 연동
- 상태관리, SSR/SSG/ISR/PPR, 캐시 전략
- BFF, 마이크로프론트엔드, 디자인 시스템
- 접근성, i18n(한/영), 웹뷰 통합

### 5. 데이터/AI
- 데이터 파이프라인(Airflow, Spark, dbt)
- 데이터 웨어하우스(BigQuery, Snowflake, Redshift)
- AI/LLM 통합: RAG, 벡터DB, AI Gateway, MCP

### 6. 보안 / 컴플라이언스 (대한민국 특화)
- **개인정보보호법 / 정보통신망법 / 전자금융감독규정**
- **ISMS-P, KISA 가이드라인, 망분리, 가명·익명처리**
- 결제: PG/VAN, 카카오페이/네이버페이/토스, 본인인증
- 암호화: KCMVP 검증 모듈, HSM, KMS

### 7. 운영 / 조직
- SRE, SLO/SLI, 장애 대응(Incident Management), 포스트모템
- 팀 토폴로지, 컴퍼넌트 오너십, 플랫폼 엔지니어링
- 비용 최적화(FinOps)

### 8. 케이스 스터디
실제 가상의 대한민국 서비스를 설계하며 책 전체를 관통한다.
- 예: "월 1,000만 MAU 커머스", "B2B SaaS", "핀테크/간편결제", "공공/금융 시스템"

## 집필 원칙

1. **한국어 1차**, 기술 용어는 영어 병기 (예: "관측성(Observability)")
2. **국내 환경 우선**: AWS-only가 아닌, NCP·KT Cloud 등 국내 클라우드와 망분리·금융망 특수성을 함께 다룬다
3. **트레이드오프 중심**: "정답"이 아닌 "선택지와 그 결과"를 제시한다
4. **다이어그램 우선**: 글보다 그림. C4, 시퀀스, 데이터 플로우를 적극 활용
5. **ADR 형식**: 주요 의사결정은 "Context → Decision → Consequences" 구조로 기록
6. **최신 기술 스택 반영**: 매년 한 번 이상 갱신을 전제로 한 구조
7. **케이스 단일 출처 (Single Source of Truth)**: 책의 모든 원픽 사례는 `manuscript/case-study/README.md` 를 단일 출처로 한다. 챕터에서 새 케이스 수치·상태를 만들지 않는다 — 필요하면 케이스 스터디에 먼저 추가한 뒤 인용한다 ([ADR-0008](adr/0008-케이스-단일-출처-원칙.md))

## 리포지토리 구조 (예정)

```
architect-book/
├── README.md
├── CLAUDE.md                  # 이 파일 — 집필 가이드
├── manuscript/                # 원고 (마크다운)
│   ├── 00-intro/
│   ├── 01-architecture-basics/
│   ├── 02-cloud-infra/
│   ├── 03-backend/
│   ├── 04-frontend/
│   ├── 05-data-ai/
│   ├── 06-security-compliance/
│   ├── 07-operations/
│   └── 08-case-studies/
├── diagrams/                  # 다이어그램 소스 (mermaid, drawio, excalidraw)
├── code-samples/              # 장별 코드 예제
│   └── <chapter>/<example>/
├── adr/                       # 책 자체에 대한 의사결정 기록
└── assets/                    # 이미지, 그림 파일
```

## Claude에게 요청할 때의 가이드

이 저장소에서 작업을 부탁할 때는 다음을 따른다.

### 원고를 쓸 때
- **모든 원고는 한국어로** 작성한다. 코드 주석도 한국어 우선.
- 새 챕터를 만들 때는 `manuscript/<번호>-<주제>/` 디렉토리에 `README.md`로 시작.
- 글의 도입부에는 항상 **"이 장에서 답하는 질문"** 3~5개를 bullet으로 명시.
- 각 절의 끝에는 **"트레이드오프"** 또는 **"체크리스트"** 박스를 둔다.
- 외부 자료를 인용할 때는 출처 URL을 각주처럼 남긴다 (단, URL을 임의로 만들지 말 것 — 확인된 것만).
- **케이스 스터디 인용 규칙** ([원칙 7](#집필-원칙)): 원픽의 조직·트래픽·DB 부하·문제 상황 등 사실을 챕터 본문에 적을 때는 `manuscript/case-study/README.md` 에서 가져온다. 본문에 새 수치/사실을 만들지 말 것. 새 수치가 필요하면 케이스 스터디 문서를 먼저 갱신한 뒤 챕터에서 인용. 가설·사고실험은 본문에 명시 ("이 절의 논의를 위해 X 를 가정한다").
- **챕터 작성 표준** ([ADR-0009](adr/0009-챕터-작성-표준.md)): 4가지 룰을 따른다.
  - **룰 1**: 모든 1차 절 끝에 트레이드오프 / 체크리스트 / 함정 박스 중 하나
  - **룰 2**: 모든 약어·외래어는 챕터 내 첫 등장 시 풀이 (도구 고유명 면제)
  - **룰 3**: 챕터 당 Mermaid 다이어그램 2개 이상 (도입 챕터 00 면제)
  - **룰 4**: 더 읽을거리 인용 형식 통일 — 책 / 공식문서 / 법규 / 온라인글 별도 형식

### 다이어그램을 만들 때
- 우선순위: **Mermaid → Excalidraw → drawio**
- C4 모델 표기를 따른다 (Context → Container → Component → Code)
- 다이어그램 원본은 `diagrams/`에, 렌더된 이미지는 `assets/`에.

### 코드 예제를 작성할 때
- 동작하는 최소 예제를 `code-samples/<chapter>/<example>/` 에 넣는다.
- README.md에 실행 방법, 의존성, 한국 환경에서의 주의점을 명시.
- 예제 언어가 여러 개일 때는 **Java(Spring Boot)와 Go를 기본**으로 하고 필요 시 추가.

### 의사결정을 기록할 때
- 책의 구성·범위·표기법에 대한 결정은 `adr/NNNN-제목.md` 형식으로 추가.
- 형식: Context / Decision / Consequences / Date / Status.

### 하지 말 것
- "AWS만이 정답" 같은 단정적 표현 금지 — 항상 대안과 트레이드오프를 함께.
- 실존 회사의 내부 정보를 추측해 사실인 것처럼 쓰지 말 것.
- 검증되지 않은 통계·벤치마크 수치를 그럴듯하게 만들어내지 말 것 (출처 없으면 명시).
- 한국 법규는 시점에 따라 달라지므로, 인용 시 **확인 필요** 표시를 남길 것.

## 빌드 / 출판 (추후)

원고 형식이 안정되면 다음 중 하나로 빌드한다 (미정):
- mdBook
- Honkit / GitBook
- Quarto
- Vercel + Next.js 기반 정적 사이트

선택은 `adr/`에 기록한다.
