# 부록 — 도구·서비스 사전

> 본 책에 등장하는 클라우드 서비스·오픈소스·외부 SaaS·법규 인증을 한 곳에 모았다.
> 항목당 4줄 — **무엇 / 왜 (역할) / 본 케이스 사용처 / 함정**.
>
> 이 부록은 본문의 약어 풀이를 보완한다 — 깊이 알고 싶을 때 펼치는 사전. 처음부터 끝까지 읽지 않아도 된다.

## 목차

- [1. 클라우드 / 인프라](#1-클라우드--인프라)
- [2. 컨테이너 / 오케스트레이션](#2-컨테이너--오케스트레이션)
- [3. 메시징 / 스트리밍](#3-메시징--스트리밍)
- [4. 데이터베이스](#4-데이터베이스)
- [5. 백엔드 언어 / 프레임워크](#5-백엔드-언어--프레임워크)
- [6. 프론트엔드 / 모바일](#6-프론트엔드--모바일)
- [7. 데이터 / AI / 모델](#7-데이터--ai--모델)
- [8. 관측성 / 운영 도구](#8-관측성--운영-도구)
- [9. CI/CD / IaC](#9-cicd--iac)
- [10. 보안 도구 / 프로토콜](#10-보안-도구--프로토콜)
- [11. 한국 외부 SaaS / 결제 / 인증](#11-한국-외부-saas--결제--인증)
- [12. 한국 공공 인증 / 법규](#12-한국-공공-인증--법규)
- [13. 한국 에너지 / 공공 시스템](#13-한국-에너지--공공-시스템)

---

## 1. 클라우드 / 인프라

### Naver Cloud Platform (NCP)

- **무엇** — 네이버가 운영하는 한국 퍼블릭 클라우드. AWS / GCP / Azure 와 동급 서비스 카탈로그.
- **왜** — 한국 리전·한국어 24/7 지원·공공 사업 가점·세종시 데이터허브 협업 사례·**CSAP 등급 보유** (공공 위탁 사업 필수).
- **본 케이스** — **단일 클라우드** ([케이스 8](../case-study/README.md#8-인프라--클라우드-정책)). NKS / Cloud DB / Object Storage / CDSS / SENS / HyperCLOVA X 등 약 15개 매니지드 사용.
- **함정** — 매니지드 폭이 AWS 대비 좁다. ClickHouse / Vault / Confluent Kafka 매니지드 없음 → NKS 위 self-host 부담. trace UX (Cloud Insight APM) 가 Tempo 대비 신생.

### AWS / GCP / Azure

- **무엇** — 글로벌 3대 퍼블릭 클라우드.
- **왜** — 매니지드 폭·인력 풀·글로벌 리전.
- **본 케이스** — **사용 X**. 1권 (원픽) 은 AWS 메인. 2권은 공공 사업·데이터 주권으로 NCP 단일 결정.
- **함정** — 한국 공공 위탁에서는 CSAP 미보유 시 입찰 제약. AWS Seoul 리전이라도 위탁 계약 검토 필요.

### NCP Cloud DB for PostgreSQL

- **무엇** — NCP 매니지드 PostgreSQL. 백업·복제·패치 자동.
- **왜** — 26명 조직이 DB DBA 풀타임 둘 인력 없음.
- **본 케이스** — 회원 / 매칭 / 정산 / DR 정책 / P2P / 복지 (OLTP 전용 인스턴스).
- **함정** — TimescaleDB extension 활성화는 NCP 의 PostgreSQL 매니지드 지원 여부 확인 필수 (지원 안 하면 NKS 위 self-host).

### NCP Object Storage

- **무엇** — AWS S3 호환 객체 저장소.
- **왜** — DR 백업·정적 자산·한전 dump 보관·로그 cold tier.
- **본 케이스** — 한전 AMI 일 dump / continuous aggregate 익스포트 / DR 백업.
- **함정** — S3 호환이지만 API 미세 차이 — pre-signed URL 만료 / 멀티파트 업로드 사이즈 한도 등 SDK 별 검증 필요.

### NCP CDN+

- **무엇** — NCP 의 콘텐츠 전송 네트워크 (한국 + 일부 글로벌 엣지).
- **왜** — 시민 웹·옥외 디스플레이 정적 자산 캐싱·이미지 변환.
- **본 케이스** — Next.js 정적 빌드·이미지·폰트 (Pretendard subset) 배포.
- **함정** — 글로벌 엣지는 AWS CloudFront 대비 적음. 본 케이스는 한국 단일이라 문제 없음.

### NCP Cloud Functions

- **무엇** — AWS Lambda 류 서버리스 함수 실행.
- **왜** — 야간 배치·이벤트 트리거·간단 ETL.
- **본 케이스** — reconciliation 배치 (한전 dump 대조) / 분기 보고서 자동 생성 / 동의 철회 시 자동 차단 작업.
- **함정** — 콜드 스타트 / 실행 시간 제한 / 동시 실행 수 제한 — Lambda 와 한도가 다름. 폭주 DR 응답에 쓰면 안 됨 (3장 3.3).

### NCP Secret Manager

- **무엇** — 시크릿 (API 키·DB 비밀번호·인증서) 관리 매니지드.
- **왜** — Vault self-host 부담 회피. KCMVP 인증 암호 모듈 연동.
- **본 케이스** — 한전 / KPX / 복지부 API 키 / DB 자격증명 / 인증서.
- **함정** — 시크릿 회전 자동화는 통합 SDK 별 차이. 폐기된 시크릿이 NKS Pod 환경변수에 남으면 위험 — 회전 시 Pod 재시작 강제.

### NCP Sub Account (IAM)

- **무엇** — 사용자 / 역할 / 권한 그룹 관리.
- **왜** — 운영자별 최소 권한 (**PoLP**, Principle of Least Privilege).
- **본 케이스** — 인프라 4 / 백엔드 6 / 데이터 3 등 역할별 권한 그룹.
- **함정** — AWS IAM 만큼 정책 표현이 정교하지 않을 수 있음 — **ABAC** (Attribute-Based Access Control) 류 패턴은 검토 필요.

### NCP Cloud Activity Tracker

- **무엇** — NCP API 호출 / 콘솔 로그인 감사 로그.
- **왜** — 공공 위탁 분기 보고 / 사고 시 추적.
- **본 케이스** — 누가·언제·무엇을 했는지 (시민 데이터 접근·시크릿 조회 등) 5년 보존.
- **함정** — Activity Tracker 가 잡지 않는 영역 (애플리케이션 내부 SQL 등) 은 별도 **SIEM** (Security Information and Event Management) 으로 보완.

---

## 2. 컨테이너 / 오케스트레이션

### Docker / 컨테이너

- **무엇** — OS 격리·이미지 패키징의 표준. Linux cgroup + namespace 위에서 동작.
- **왜** — "내 PC 에서는 됐는데" 해결 / 배포 단위 통일.
- **본 케이스** — 모든 백엔드 / 프론트 서비스·CI 빌드.
- **함정** — 이미지 크기 (alpine vs ubuntu) / 보안 스캔 누락 / 멀티스테이지 빌드 미사용 시 시크릿 leak.

### Kubernetes

- **무엇** — 컨테이너 오케스트레이션. Pod / Deployment / Service / Ingress / ConfigMap / Secret 등 추상으로 클러스터를 선언적으로 관리.
- **왜** — 12개+ 서비스를 사람이 직접 띄우는 게 불가능. 오토스케일·롤링 배포·자기 복구·시크릿 관리 표준.
- **본 케이스** — NKS 위에서 모놀리스 + Go 추출 3개 (IoT/DR/알람) + 프론트 + 콘솔 3개 + Kiosk = 약 12 Pod 군.
- **함정** — 26명 조직에는 무거움. 매니지드 (NKS / EKS) 가 거의 강제. 자체 운영은 SRE 3명 이상이 필수.

### NKS (Naver Kubernetes Service)

- **무엇** — NCP 의 매니지드 Kubernetes (EKS / GKE 와 같은 위치).
- **왜** — control plane 운영 부담 회피 + NCP 통합 (Load Balancer·시크릿·로그).
- **본 케이스** — 단일 클러스터·prod + staging 2 namespace.
- **함정** — Kubernetes 버전 업그레이드 정책 / Add-on 카탈로그가 EKS 만큼 풍부하진 않음.

### Helm

- **무엇** — Kubernetes 의 패키지 매니저. YAML 템플릿화 + 버전·릴리스 관리.
- **왜** — 환경 (dev / staging / prod) 별 같은 차트에 다른 values.
- **본 케이스** — 모든 서비스 Helm 차트 + 사내 차트 저장소.
- **함정** — 차트 복잡도 폭주 / values.yaml 가 길어지면 가독성 ↓. 단순 서비스는 plain manifest 가 나음.

### ArgoCD

- **무엇** — Git 의 상태를 클러스터 실제 상태로 동기화하는 GitOps 도구.
- **왜** — 배포 = Git PR. 감사 가능 / 롤백 = revert.
- **본 케이스** — 모든 배포 ArgoCD 동기화 (NKS self-host). 운영 / staging namespace 별 분리.
- **함정** — ArgoCD 자체 장애 시 신규 배포 중단. NKS 자체가 죽으면 ArgoCD 도 같이 죽음 — DR 시나리오 검토.

### Service Mesh (Istio / Linkerd) — **본 케이스 미사용**

- **무엇** — 서비스 간 통신을 사이드카로 가로채 mTLS / 트래픽 라우팅 / 관측성 자동화.
- **왜 (도입 시)** — 마이크로서비스 50개+ 에서 가치.
- **본 케이스** — **미사용**. 12개 서비스 + 26명 조직에 운영 부담 과잉 (1권과 동일 결정).
- **함정** — "있으면 좋겠다" 로 도입 시 사이드카 메모리·네트워크 hop·디버깅 복잡도 폭발.

---

## 3. 메시징 / 스트리밍

### Apache Kafka

- **무엇** — 분산 로그 기반 스트리밍 플랫폼. Topic (논리 채널) → Partition (병렬 단위) → Consumer Group (병렬 컨슘) 의 3 추상.
- **왜** — 대량 이벤트 (IoT) / 비동기 도메인 이벤트 / 폭주 흡수 (DR 응답) 의 표준.
- **본 케이스** — NCP CDSS (Kafka 호환) 로 IoT ingest + DR 응답 + 도메인 이벤트. **partition 32+** ([케이스 10.1](../case-study/README.md#101-구체-신호-수치-기준)).
- **함정** — Partition 수 = 컨슈머 병렬도 상한. 적게 잡으면 폭주 시 lag. 많이 잡으면 메타데이터 부담. 한 번 정한 partition 수 늘리는 건 가능하지만 키 분포가 깨질 수 있어 신중.

### NCP Cloud Data Streaming Service (CDSS)

- **무엇** — NCP 의 Kafka 호환 매니지드.
- **왜** — Kafka self-host 부담 회피 (Zookeeper·브로커 패치·스토리지 관리).
- **본 케이스** — IoT ingest (일 2,700만) + DR 응답 (월 5~10회 폭주) + 도메인 이벤트.
- **함정** — Confluent Cloud 같은 매니지드 (Schema Registry / ksqlDB / Connect) 만큼 풍부하지 않음. Schema Registry 는 별도 self-host 검토.

### NCP SENS (Simple & Easy Notification Service)

- **무엇** — SMS / 알림톡 (카카오) / Push 통합 발송 매니지드.
- **왜** — 통신사 직접 연동 부담 회피.
- **본 케이스** — 시민 알림 / **취약계층 안전 알람 3채널 중 SMS / 자동 전화 (TTS)** ([7장 2](../07-운영과-조직/README.md#2-slo--도메인별--안전-sla)).
- **함정** — 알림톡은 카카오 검수 필요. SMS 도달률은 통신사·번호 변경에 따라 변동.

### RabbitMQ / NATS — **본 케이스 미사용 / 대안 참고**

- **무엇** — Kafka 와 다른 메시징 — RabbitMQ 는 AMQP 기반 라우팅 강점, NATS 는 가벼움.
- **왜 (대안)** — Kafka 가 무거운 도메인 / 짧은 latency 중심.
- **본 케이스** — Kafka (CDSS) 단일 — 26명 조직이 두 도구 운영 X.
- **함정** — Kafka 만 알면 "전부 Topic" 으로 보임 — 라우팅 / 우선순위 큐가 진짜 필요한 곳을 식별 못 함.

---

## 4. 데이터베이스

### PostgreSQL

- **무엇** — 오픈소스 관계형 DB 의 사실상 표준. JSON·전문검색·extension 풍부.
- **왜** — 트랜잭션 안정 + 확장 가능 (TimescaleDB / PostGIS / pgvector 등).
- **본 케이스** — 회원 / 매칭 / 정산 / DR 정책 / P2P / 복지 (NCP Cloud DB).
- **함정** — VACUUM / autovacuum 튜닝·커넥션 폭주·long transaction 으로 인한 bloat — 단순해 보이지만 운영 깊이 있음.

### TimescaleDB

- **무엇** — PostgreSQL extension. **Hypertable** (시간 + 공간 파티셔닝 자동) + **Continuous Aggregate** (자동 갱신 머티리얼라이즈드 뷰) + **Compression** (column store 형태).
- **왜** — IoT / 메트릭 같은 시계열을 PostgreSQL 스킬로 다룰 수 있음.
- **본 케이스** — IoT 시계열 (스마트미터·태양광·웨어러블) — **전용 인스턴스** (OLTP 와 분리, [케이스 11.2.2](../case-study/README.md#1122-db)).
- **함정** — chunk 크기 / CAGG 갱신 주기 / retention drop 정책이 운영 비용을 좌우. OLTP 와 같은 인스턴스에 두면 WAL / VACUUM 경합으로 사고.

### ClickHouse

- **무엇** — OLAP 컬럼 지향 DB. 분석 쿼리에 수십~수백 배 빠름.
- **왜** — TimescaleDB 도 분석 가능하지만 대량 집계는 ClickHouse 가 압도.
- **본 케이스** — DR 분석 / 시민 보고서 / 동·면 단위 통계 (NKS 위 self-host — NCP 매니지드 없음).
- **함정** — 트랜잭션 X / 업데이트 약함 / 작은 row insert 비효율 — append-only 분석 워크로드 전용. 운영 부담 (Zookeeper / Keeper) 도 있음.

### Redis

- **무엇** — 인메모리 키-값 + 자료구조 (List·Set·Sorted Set·Stream·Pub/Sub).
- **왜** — 캐시 / 세션 / 분산 락 / 실시간 알림.
- **본 케이스** — 시민 앱 첫 화면 캐시·세션·DR 호출 Pub/Sub (NCP Managed Redis).
- **함정** — 영속성 (AOF / RDB) 옵션과 일관성 트레이드오프. 분산 락은 Redlock 의 안전성 논쟁 — critical 한 정산에는 PG advisory lock 이 안전.

### MongoDB / DynamoDB — **본 케이스 미사용**

- **무엇** — NoSQL 문서 / 키-값.
- **왜 (대안)** — 스키마 자유도·수평 확장.
- **본 케이스** — 사용 X. PostgreSQL + JSONB 로 충분.
- **함정** — 1권 (커머스) 은 DynamoDB 등 멀티 DB. 2권은 26명에 다양성을 줄임.

### OpenSearch / Elasticsearch — **본 케이스 미사용**

- **무엇** — 분산 검색 + 분석 엔진.
- **왜 (대안)** — 형태소 분석 (Nori 한국어) + LTR 등 검색 품질.
- **본 케이스** — PostgreSQL FTS 로 진단사·시공사 검색 처리. 검색 비중 낮음.
- **함정** — 검색 비중이 핵심이 되면 (예: 진단사 1만 명) PostgreSQL FTS 한계 도달.

---

## 5. 백엔드 언어 / 프레임워크

### Spring Boot (Java / Kotlin)

- **무엇** — JVM 위 가장 보편적인 백엔드 프레임워크. Spring Framework + 자동 설정 + 임베디드 톰캣.
- **왜** — 도메인 / 트랜잭션 / ORM (JPA) / 보안 / 메시징 생태계가 가장 풍부.
- **본 케이스** — **모놀리스** (회원·매칭·정산·복지·P2P·DR 정책) — Kotlin.
- **함정** — JVM heap 메모리 / 콜드 스타트 — IoT 같은 stateful 폭주 처리에는 무거움.

### Kotlin

- **무엇** — JetBrains 가 만든 JVM 언어. 자바 호환 + null 안전 + 코루틴.
- **왜** — 자바 라이브러리 그대로 + 보일러플레이트 ↓ + 코루틴으로 비동기.
- **본 케이스** — 모놀리스 Spring Boot 의 1차 언어.
- **함정** — 자바 인력 풀로 채용 시 코틀린 학습 곡선 (1~2주). 코루틴은 잘못 쓰면 컨텍스트 leak.

### Go (Golang)

- **무엇** — Google 이 만든 컴파일 언어. 메모리 적음·goroutine 동시성·표준 라이브러리 풍부.
- **왜** — 네트워크 서비스 / 데이터 파이프라인 / CLI 에 강함. 콜드 스타트 빠름.
- **본 케이스** — **IoT ingester / DR 응답기 / 알람 처리기** — Spring 모놀리스가 못 감당하는 부하 격리.
- **함정** — 제네릭이 약함 (1.18+ 도입 후 개선). ORM 이 약함 (수동 SQL 친화).

### Node.js (NestJS) — **본 케이스 미사용 / 1권 참고**

- **무엇** — V8 위 JavaScript 런타임. NestJS 는 Spring 영감의 프레임워크.
- **왜 (대안)** — 프론트팀과 같은 언어 / I/O 다중화.
- **본 케이스** — 사용 X.
- **함정** — CPU 무거운 작업 (시계열 집계·암호화) 에 약함.

### Python (FastAPI) — **본 케이스 미사용 / 1권 참고**

- **무엇** — FastAPI 는 async + 자동 OpenAPI 생성.
- **왜 (대안)** — 데이터 / ML 인접에 강함.
- **본 케이스** — ML 모델 (Prophet·XGBoost) 학습은 Python notebook·운영은 모델 서빙 (사이드카) 또는 미리 계산해 PG 저장.
- **함정** — 시계열 운영 데이터 직접 처리는 Go 우위 (메모리·동시성).

---

## 6. 프론트엔드 / 모바일

### Next.js

- **무엇** — React 기반 풀스택 프레임워크. SSR / SSG / ISR / PPR + Server Action + 파일 기반 라우팅.
- **왜** — 정적 + 동적 혼합 / SEO / 첫 화면 빠름.
- **본 케이스** — 시민 웹 / 진단사·사회복지사 콘솔 / 옥외 디스플레이 Kiosk **모두 NKS 내부 호스팅** (Vercel 미사용 — 데이터 주권, [4장 1](../04-프론트엔드/README.md#1-채널-구성)).
- **함정** — Vercel 외 self-host 시 ISR · Image Optimization · Edge Function 일부 기능 제약. NKS 위 Next.js Standalone 빌드 + 자체 이미지 변환 필요.

### React

- **무엇** — Meta 가 만든 UI 라이브러리. 컴포넌트 + Hooks.
- **왜** — 인력 풀 가장 크고 생태계 (테스트·접근성·디자인 시스템) 풍부.
- **본 케이스** — Next.js 의 기반.
- **함정** — 16+ 의 Concurrent / Server Components 도입으로 모델 학습 곡선 상승.

### Swift / SwiftUI

- **무엇** — Apple 의 iOS / macOS 네이티브 언어 + 선언적 UI.
- **왜** — 시니어 접근성 (VoiceOver·동적 폰트) 가 가장 풍부. **SFSpeechRecognizer** (Native 음성) 가 한국어 시니어 WER 안정.
- **본 케이스** — 시민 iOS 앱. KWCAG 2.2 AA + 단순 모드.
- **함정** — SwiftUI 는 빠르게 발전 중이라 iOS 버전별 동작 차이. iOS 15 미만 지원하면 UIKit 병행.

### Kotlin / Jetpack Compose

- **무엇** — Android 의 선언적 UI (SwiftUI 의 안드로이드판).
- **왜** — TalkBack · 동적 폰트 · 음성 (SpeechRecognizer Native).
- **본 케이스** — 시민 Android 앱.
- **함정** — Compose 는 비교적 새로움 — 일부 서드파티 라이브러리 미지원.

### Recharts / D3 — **그래프 라이브러리**

- **무엇** — Recharts (React 친화 / 가벼움) vs D3 (저수준 / 자유도 높음).
- **왜** — 시민 앱 / 옥외 디스플레이의 전력 시각화.
- **본 케이스** — Recharts (lazy load) + 첫 화면은 SVG 직접 그린 sparkline.
- **함정** — 그래프 라이브러리 풀 번들이 60~80KB — 첫 화면 250KB 예산에 무겁다 → lazy.

---

## 7. 데이터 / AI / 모델

### Apache Airflow / dbt / Spark — **본 케이스 미사용 / 1권 참고**

- **무엇** — Airflow (워크플로우) / dbt (SQL 변환·테스트) / Spark (분산 처리).
- **왜 (대안)** — 대량 ETL / 분석 데이터 마트.
- **본 케이스** — 26명에 과잉. NCP Cloud Functions + ClickHouse 직접 집계로 대체.
- **함정** — "데이터팀 = Airflow" 공식이 1권은 맞지만 작은 조직에서는 README + cron 으로 충분.

### Prophet

- **무엇** — Meta 가 만든 시계열 예측. 추세 + 계절성 + 휴일 효과 + 외생 변수.
- **왜** — 가벼움 + 설명 가능성 + 한국 휴일 / 계절 패턴에 강함.
- **본 케이스** — **DR 응답 예측 동·면 단위** ([5장 4](../05-데이터와-AI/README.md#4-dr-예측-모델)). MAPE 8~12% 목표. 가구별은 보조.
- **함정** — 가구별 단위는 분산이 커서 MAPE 25~40% — 정산에 직접 쓰면 한전 정합 감사에서 폭주.

### XGBoost / LightGBM

- **무엇** — 그래디언트 부스팅 결정 트리. 표 형식 데이터 분류 / 회귀의 사실상 표준.
- **왜** — 빠름 + SHAP 으로 설명 가능 + 작은 데이터에도 강함.
- **본 케이스** — **취약계층 식별** ([5장 5](../05-데이터와-AI/README.md#5-취약계층-식별-모델)). 룰 + 모델 + 인간 리뷰 3 안전망의 모델 단계.
- **함정** — 학습 데이터 편향이 그대로 모델에 박힘 — SHAP 사후 점검 + 임계값 거버넌스 필수.

### SHAP (SHapley Additive exPlanations)

- **무엇** — 게임 이론 Shapley 값 기반의 모델 예측 해석.
- **왜** — "왜 이 가구가 취약 점수가 높은가" 를 변수별 기여도로 분해.
- **본 케이스** — XGBoost 출력에 SHAP 적용 → **사회복지사 콘솔에 자연어 변환** ("난방비 평균 230% / 야간 활동 ↓ / 1인 노인").
- **함정** — SHAP 값 자체는 사회복지사가 못 이해 → 자연어 변환 레이어 필수. raw SHAP plot 만 보여주면 안 씀.

### HyperCLOVA X (HCX)

- **무엇** — 네이버가 만든 한국어 특화 LLM (대규모 언어 모델).
- **왜** — 한국어 정확도 + 공공 사업 가점 + NCP 통합.
- **본 케이스** — 시민 챗봇 (에너지 절약 팁 RAG) — **Q3 도입 + ZRR (Zero Retention) 계약 체결 후** ([케이스 12](../case-study/README.md#12-향후-로드맵-가정)).
- **함정** — ZRR 은 기본 X — 별도 기업 계약 필요. 일반 API 키는 입력이 학습 / 로그 보존될 수 있음. 가격·정책 분기마다 변동.

### RAG (Retrieval Augmented Generation)

- **무엇** — LLM 에 외부 지식 (사내 FAQ / 문서) 을 검색 결과로 주입해 답변 품질을 올림.
- **왜** — LLM 환각 감소 + 도메인 지식 (에너지 절약 팁) 주입.
- **본 케이스** — 사내 FAQ + 동·면 단위 통계 (가구 식별 X) 를 컨텍스트로 챗봇.
- **함정** — 컨텍스트에 PII 가 흘러가면 ZRR 위반. **AI Gateway 단에서 마스킹 강제**.

### AI Gateway

- **무엇** — LLM 호출 앞단의 프록시. 라우팅·캐시·비용 추적·PII 마스킹·fallback.
- **왜** — 여러 LLM (HCX / GPT / Claude) 을 한 곳에서 정책 관리.
- **본 케이스** — PII 마스킹 강제 + ZRR 계약 없는 모델로의 fallback 차단.
- **함정** — fallback 라우팅이 마스킹 우회하는 사례 — fallback 전에도 동일 마스킹 + ZRR 계약 모델로만.

### OpenTelemetry (OTel)

- **무엇** — 메트릭·로그·트레이스 수집의 오픈 표준. SDK + Collector.
- **왜** — 백엔드 별 (Datadog / Cloud Insight / Grafana) 종속을 푸는 추상층.
- **본 케이스** — 모든 서비스 OTel SDK → NCP Cloud Insight 로 전송. 향후 다른 백엔드로 이전해도 SDK 그대로.
- **함정** — Collector 운영 / 샘플링 정책 / 카디널리티 폭주 (gauge 라벨 너무 많음) — OTel 표준이지만 운영 복잡도는 별개.

---

## 8. 관측성 / 운영 도구

### NCP Cloud Insight

- **무엇** — NCP 매니지드 관측성 — 메트릭 + APM (트레이스).
- **왜** — 데이터 주권 (한국 리전) + 공공 위탁 감사 요건.
- **본 케이스** — 모든 NKS Pod / DB / CDSS 메트릭 + 트레이스.
- **함정** — APM 트레이스 검색 UX 가 Grafana Tempo 대비 미흡 — NKS 위 Grafana OSS 보조 (시각화만, 저장은 NCP 유지).

### NCP Cloud Log Analytics

- **무엇** — NCP 매니지드 로그 — 수집 / 검색 / 보존.
- **왜** — 공공 위탁 5년 감사 보존 + 시민 데이터 접근 로그.
- **본 케이스** — 모든 NKS 컨테이너 로그 + 감사 로그.
- **함정** — 로그 폭주 시 비용 — 디버그 로그는 단기 / 감사 로그는 5년 분리 보존.

### Grafana Cloud — **본 케이스 미사용**

- **무엇** — Grafana 가 운영하는 매니지드 관측성 (메트릭+로그+트레이스).
- **왜 (대안)** — 1권 (원픽) 의 1차 선택.
- **본 케이스** — **사용 X** — 2026 시점 한국 리전 부재로 데이터 주권 위반.
- **함정** — "오픈소스 친화 = 안전" 인상이지만 데이터는 미국 / EU / 일본 리전 — 공공 위탁은 못 씀.

### Loki / Tempo / Prometheus — **본 케이스 self-host 회피**

- **무엇** — Grafana 의 로그 / 트레이스 / 메트릭 OSS.
- **왜 (대안)** — 매니지드 비용 절감.
- **본 케이스** — 26명 self-host 부담 — NCP 매니지드 사용. Tempo 만 trace UX 보조로 NKS 위 (저장은 한국 리전 Object Storage).
- **함정** — Prometheus 카디널리티 폭주·Loki TSDB 디스크·Tempo 보존 — 모두 SRE 시간 소요.

### Sentry

- **무엇** — 애플리케이션 에러 추적 SaaS. 스택트레이스 + 사용자 컨텍스트 + 릴리스 추적.
- **왜** — 1권과 동일 — 매니지드 외 self-host (Sentry self-hosted) 도 가능하지만 26명에 부담.
- **본 케이스** — 모든 백엔드 / 프론트 에러. **PII 스크러빙 필수** + 공공 도메인 (복지·취약계층) 은 별도 채널.
- **함정** — Sentry 도 외부 SaaS — PII 스크러빙 누락 시 시민 데이터 유출. SDK 의 beforeSend 훅으로 강제.

---

## 9. CI/CD / IaC

### GitHub Actions

- **무엇** — GitHub 의 CI / 워크플로우 자동화. workflow YAML + 마켓플레이스 액션.
- **왜** — GitHub 사용 중이면 추가 인프라 X. 무료 분량 + secret 통합.
- **본 케이스** — PR 검사 / 빌드 / 테스트 / 이미지 빌드 → NCP Container Registry push.
- **함정** — runner 가 GitHub 인프라 → 한국 ↔ GitHub 네트워크 가끔 느림. 큰 빌드는 self-hosted runner 검토.

### ArgoCD

- (위 [2. 컨테이너](#argocd) 참조)

### Terraform

- **무엇** — HashiCorp 의 IaC. HCL 언어로 클라우드 리소스를 선언적으로.
- **왜** — 콘솔 클릭 = 사라지는 지식. Terraform = 코드 = 리뷰 가능 = 회복 가능.
- **본 케이스** — NCP Provider 로 NKS / VPC / Cloud DB / CDSS / Object Storage / IAM 전부.
- **함정** — state 파일 잠금 / 충돌 / drift — 작은 조직도 remote state 필수 (NCP Object Storage + locking).

### Pulumi — **대안**

- **무엇** — Terraform 과 같은 IaC 인데 일반 언어 (TypeScript / Go / Python).
- **왜 (대안)** — 코드 추상화 강함.
- **본 케이스** — Terraform 채택 (한국 인력 풀 / NCP Provider 지원 우선).

---

## 10. 보안 도구 / 프로토콜

### mTLS (mutual TLS)

- **무엇** — 양방향 TLS — 서버뿐 아니라 클라이언트도 인증서 제시.
- **왜** — IoT 게이트웨이 / 내부 서비스 / 사회복지사 콘솔 — "누구든 https 만 알면 접근" 차단.
- **본 케이스** — 가구 게이트웨이 ↔ NKS / 사회복지사 콘솔 / 외부 공공 API.
- **함정** — 인증서 회전 자동화 (EST/SCEP) 가 없으면 만료로 5만 가구 다운.

### EST (Enrollment over Secure Transport) / SCEP

- **무엇** — IoT 단말이 자동으로 인증서를 발급 / 갱신받는 표준 프로토콜.
- **왜** — 5만 가구 × 인증서 1년 갱신 = 일 137건 자동 처리. 사람이 못 함.
- **본 케이스** — 사내 PKI ↔ 가구 게이트웨이 ([6장 6.1](../06-보안과-컴플라이언스/README.md#61-게이트웨이--미터-인증--인증서-자동-발급--갱신)).
- **함정** — 폐기 (**CRL** Certificate Revocation List / **OCSP** Online Certificate Status Protocol) 도 자동화해야 함. 이사 / 단말 교체 시 즉시.

### OTA (Over-The-Air firmware update) — Dual Partition

- **무엇** — 펌웨어를 원격으로 업데이트. A/B 듀얼 파티션 = 신규 펌웨어 다운로드·검증·전환·실패 시 롤백.
- **왜** — 단일 파티션 OTA 는 실패 시 단말 brick.
- **본 케이스** — 가구 게이트웨이 펌웨어. Canary (10가구) → 100 → 1k → 전체.
- **함정** — 펌웨어 서명 검증 누락 시 악성 펌웨어로 5만 단말 동시 장악.

### SBOM (Software Bill of Materials)

- **무엇** — 소프트웨어 / 펌웨어에 들어간 모든 라이브러리·버전 목록. SPDX / CycloneDX 형식.
- **왜** — CVE 알람 자동화 / 공공 / 미국 정부 (EO 14028) 가 요구.
- **본 케이스** — 게이트웨이 펌웨어마다 SBOM + 사내 CVE 알람.
- **함정** — SBOM 만 만들고 모니터링 안 하면 종이만. CVE 등록 → SBOM 매칭 → 다음 OTA 자동 패치 파이프라인 필요.

### TPM / Secure Boot — **검토 대상**

- **무엇** — 하드웨어 보안 모듈 (TPM) + 부팅 시 펌웨어 서명 검증.
- **왜** — 가구 게이트웨이 변조 / Root 권한 획득 방어.
- **본 케이스** — H/W 비용 vs 위험 트레이드오프로 단가 분석 후 결정 (2026년 검토 중).
- **함정** — H/W 별 TPM 칩 / SOC 통합 보안 영역 (ARM TrustZone) 등 옵션 다양 — 단가 / 운영 부담 검증 필수.

### Vault (HashiCorp) — **본 케이스 NCP Secret Manager 로 대체**

- **무엇** — 시크릿 관리·동적 자격증명·키 회전·정책.
- **왜 (대안)** — 가장 풍부한 시크릿 관리.
- **본 케이스** — NCP Secret Manager 로 대체 (self-host 부담 회피 + KCMVP 연동).

---

## 11. 한국 외부 SaaS / 결제 / 인증

### PASS (이통3사 본인인증)

- **무엇** — SKT / KT / LGU+ 가 공동 운영하는 모바일 본인인증.
- **왜** — 시니어 친화 (앱 + 휴대폰 인증) + 가장 보편적.
- **본 케이스** — 시민 가입 1차 본인인증 ([케이스 11.2.4](../case-study/README.md#1124-인증--토큰)).
- **함정** — PASS 앱 미설치 시니어는 SMS 본인인증 fallback. NICE 보조 권장.

### NICE / KCB (신용평가사)

- **무엇** — 전통 본인인증 (휴대폰 / 공인인증서 / 카드).
- **왜** — PASS 가 안 되는 사용자 대비.
- **본 케이스** — PASS 보조 + 진단사·시공사 신원 검증.
- **함정** — 본인인증 비용 (건당 수십~수백 원) 누적 — 가입 / 매칭 / 인증 횟수 분석 필수.

### 카카오 / 네이버 간편로그인 (OAuth2)

- **무엇** — 카카오·네이버 계정으로 가입 / 로그인.
- **왜** — 가입 컨버전 ↑ / 시민 친숙도.
- **본 케이스** — 시민 가입의 보조 경로 (PASS 가 1차).
- **함정** — 카카오·네이버 약관 변경 시 사용자 재동의 / 이메일 미공개 정책 등 케이스별 처리.

### 토스페이먼츠

- **무엇** — 토스의 PG (Payment Gateway) — 카드·계좌·간편결제.
- **왜** — UX·SDK 품질·정산 투명성.
- **본 케이스** — 매칭 시공사 정산 / 일반 결제.
- **함정** — 1권의 PG 4 사 (토스·NHN·KCP·이니시스) 대비 단일 — 장애 시 fallback 없음. 26명 조직 결정 (위험 등록부에 명시).

### 세종페이 (지역화폐)

- **무엇** — 세종특별자치시가 발행하는 지역화폐. 가맹점 결제 + 시민 환급.
- **왜** — DR 보상 환급의 시민 수령 채널 — 지역경제 + 행정 가산.
- **본 케이스** — DR 보상 환급 / 일부 정산.
- **함정** — 지역화폐 API 스펙은 지자체별 다름 — 별도 협력 + 분기 reconciliation.

### FCM (Firebase Cloud Messaging) / APNs (Apple Push Notification service)

- **무엇** — Google / Apple 의 푸시 인프라.
- **왜** — 모바일 앱 푸시의 사실상 표준.
- **본 케이스** — 시민 앱 알림 (DR 호출 / 환급 / 안전 알람 1순위 채널).
- **함정** — OS 절전 / 알림 거부 / 토큰 만료로 단독 도달률 95~97%. **3채널 OR** 합성 필수 ([7장 2](../07-운영과-조직/README.md#2-slo--도메인별--안전-sla)).

---

## 12. 한국 공공 인증 / 법규

### CSAP (Cloud Security Assurance Program)

- **무엇** — KISA 가 부여하는 공공 위탁 클라우드 보안 인증. 일반 / 표준 / 간편 등급.
- **왜** — 공공 위탁 입찰의 사실상 의무.
- **본 케이스** — **NCP 의 CSAP 등급 보유 사실 확인 후 채택** — 우리는 별도 받지 않지만 위탁 계약서에 등급 명시.
- **함정** — AWS Seoul 리전이라도 CSAP 미보유면 공공 위탁 입찰 어려움.

### KCMVP (Korean Cryptographic Module Validation Program)

- **무엇** — 국정원이 부여하는 행정·공공 정보 암호 모듈 검증.
- **왜** — 시민 식별번호 / 복지 데이터 등 행정 데이터는 KCMVP 인증 라이브러리 의무.
- **본 케이스** — 시크릿 / 행정 데이터 암호화 모듈 (예: 국정원 등재 모듈) 사용. bcrypt / openssl 직접 사용은 위반.
- **함정** — 인증 모듈은 일반 라이브러리 대비 성능 / 호환성 손실 — 도입 전 PoC.

### PIA (Privacy Impact Assessment / 개인정보 영향평가)

- **무엇** — 공공기관 위탁 시 의무. 개보법 제33조 / 시행령.
- **왜** — 시민 데이터 처리의 위험 사전 평가.
- **본 케이스** — **매년 갱신 + 신규 도메인 (예: 웨어러블) 추가 시 재평가**. 외부 평가기관 (NIA 등록) 위탁.
- **함정** — 보안 1명 자체 수행 X — 외부 자문 필수. 분기 / 연간 일정 미리 잡지 않으면 위탁 평가에서 감점.

### ISMS-P (Information Security Management System with Personal Information)

- **무엇** — KISA 가 부여하는 정보보호·개인정보보호 통합 인증.
- **왜** — 매출 100억+ 시 ISMS 의무, ISMS-P 는 상위.
- **본 케이스** — 매출 30~40억 이라 의무 아님 — 단 **공공 위탁 계약서가 사실상 ISMS-P 수준 보안 요구**.
- **함정** — "의무 X" 이라고 안 받으면 공공 위탁에서 떨어짐 — 사실상 의무.

### KWCAG (Korean Web Content Accessibility Guidelines)

- **무엇** — 한국형 웹 접근성 지침. WCAG 2.x 의 한국 적용. 2.2 가 최신.
- **왜** — 공공 / 공익 성격 웹·앱·키오스크 접근성 의무.
- **본 케이스** — 시민 앱 / 웹 / 콘솔 / 옥외 디스플레이 모두 **KWCAG 2.2 AA 강제**.
- **함정** — WCAG 자동 점검 도구 (axe-core 등) 만으론 70% 점검 — 수동 점검 (스크린리더 / 키보드만) 필수.

### 「분산에너지 활성화 특별법」 (2024.6.14 시행, 법률 19470호)

- **무엇** — 시민 분산자원 (태양광 / ESS / EV) 의 시장 참여 합법화 — 특화지역 / 통합발전소 / 소규모전력중개사업 3 경로.
- **왜** — 시민 P2P 거래 / DR 참여의 법적 근거.
- **본 케이스** — **소규모전력중개사업자 등록** 경로 (Q2 목표). 세종 특화지역 지정 여부는 별도 확인.
- **함정** — 3 경로 차이를 한 덩어리로 다루면 사업 라이선스 잘못 등록.

### 「전력시장운영규칙」 (KPX)

- **무엇** — 한국전력거래소 (KPX) 가 운영하는 전력 시장 운영 규칙. DR (수요반응자원) 시장 / 거래·보상 / SLA 근거.
- **왜** — DR 사업의 실질 규칙은 「에너지이용합리화법」 이 아닌 KPX 운영규칙.
- **본 케이스** — DR 호출 응답 SLA / 보고 형식 / 보상 계산.
- **함정** — 「에너지이용합리화법」 을 DR 근거로 인용하면 사실 오류 — 정정 필수 (본 책 6장 1).

### 「개인정보 보호법」 — 핵심 조문

- **무엇** — 개인정보 수집·이용·제공·파기·권리. 2024 개정 시행.
- **왜** — 시민 데이터의 1차 법규.
- **본 케이스** — 제15·17·22·23·32조의2 (PIA)·**제37조의2 자동화 결정 거부권**.
- **함정** — 제37조의2 (자동화 거부권) 누락 — 본 책은 명시 ([6장 1.3](../06-보안과-컴플라이언스/README.md#13-개인정보-보호법-핵심-조문)).

---

## 13. 한국 에너지 / 공공 시스템

### 한전 AMI (Advanced Metering Infrastructure)

- **무엇** — 한국전력공사의 차세대 검침 인프라. 스마트미터 ↔ 한전 중앙. 1차 데이터 소유는 한전.
- **왜** — 가구 전력 사용량 데이터 / DR 호출 / 정산 정합.
- **본 케이스** — 한전 ↔ 우리 — **API + 일 단위 dump 혼합** ([2장 3.1](../02-클라우드와-인프라/README.md#31-한전-ami-연계)). 일 단위 reconciliation 필수.
- **함정** — 한전 API 변경 통보가 짧을 수 있음 (정부 정책 변동) — 캐싱 / 어댑터 레이어 필수.

### KPX (한국전력거래소)

- **무엇** — 도매 전력시장 / DR 시장 운영자. 분산자원 거래·DR 호출·보상 정산.
- **왜** — DR 사업의 호출자 + 시장 관리.
- **본 케이스** — DR 호출 수신 → 시민 응답 집계 → KPX 보고 → 보상 정산.
- **함정** — 호출 빈도 / 응답 SLA 가 시즌 / 정책에 따라 변동 — 분기 1회 운영규칙 갱신 확인.

### 세종시 데이터허브

- **무엇** — 세종특별자치시의 스마트시티 통합 데이터 게이트웨이. 교통 / 환경 / 에너지 등.
- **왜** — 시 차원 공공 데이터 교환의 한 곳.
- **본 케이스** — **양방향** — 우리는 동·면 집계 통계 송신, 시는 전기차 충전소 / 기상 / 인구 통계 제공.
- **함정** — 시 데이터허브 API 스펙 변경 = 우리 API 변경 (위탁 계약). 자율성 일부 양보.

### 복지부 / 사회보장정보시스템 (행복e음)

- **무엇** — 보건복지부의 사회보장 통합 시스템. 취약계층 식별 / 위탁 정산.
- **왜** — 에너지 빈곤 / 취약계층 자동 식별 / 행정 연계.
- **본 케이스** — 별도 위탁 협약 + 데이터 공유 동의 + 처리 위탁 계약 + **시민 별도 동의 후만 접근**.
- **함정** — 행복e음 데이터는 민감정보 (개보법 제23조) — 동의 / 마스킹 / 접근 기록 / 동의 철회 즉시 차단 모두 필요.

### KISA (한국인터넷진흥원)

- **무엇** — 한국 정보보호 진흥기관. ISMS-P / CSAP 인증·사고 통지 / 보안 가이드.
- **왜** — 사고 통지 72h / 보안 가이드 / 인증 평가.
- **본 케이스** — 사고 시 정보주체 5일 / 보호위·KISA 72h 통지 의무.
- **함정** — 통지 양식 / 첨부 / 후속 보고 절차 — 사고 터지고 알면 늦음. 사전 drill 권장.

---

## 더 읽을거리 (도구·서비스 별)

> 책 / 공식 자료 / 법규는 [ADR-0009 룰 4](../../adr/0009-챕터-작성-표준.md) 형식 — 한국어 1차, URL 은 확인된 것만.

### 클라우드 / 쿠버네티스
- *Kubernetes Up & Running*, 3판 — Brendan Burns et al. (O'Reilly, 2022) — 쿠버네티스 기초 / 운영
- *Cloud Native Patterns* — Cornelia Davis (Manning, 2019) — 클라우드 네이티브 사고법
- NCP 공식 문서 — Naver Cloud Platform (URL 확인 필요)

### 메시징 / 스트리밍 / 데이터
- *Kafka: The Definitive Guide*, 2판 — Gwen Shapira et al. (O'Reilly, 2021)
- *Designing Data-Intensive Applications* — Martin Kleppmann (O'Reilly, 2017) — 시계열 / 분산 / 스트림 전반
- TimescaleDB 공식 문서 — Timescale (URL 확인 필요)
- ClickHouse 공식 문서 — ClickHouse Inc. (URL 확인 필요)

### 백엔드 / 프론트
- *Effective Kotlin* — Marcin Moskala (자체 출판, 2019)
- *100 Go Mistakes and How to Avoid Them* — Teiva Harsanyi (Manning, 2022)
- *Patterns of Enterprise Application Architecture* — Martin Fowler (Addison-Wesley, 2002) — 시대를 타지 않는 백엔드 패턴
- Next.js / React 공식 문서 — Vercel / Meta (URL 확인 필요)

### 데이터 / AI
- *Interpretable Machine Learning*, 2판 — Christoph Molnar (자체 출판, 2022) — SHAP / 설명 가능성
- *Forecasting: Principles and Practice*, 3판 — Rob J. Hyndman, George Athanasopoulos (OTexts, 2021) — Prophet 사고 배경

### 보안 / IoT / 공공
- *Practical IoT Hacking* — Fotios Chantzis et al. (No Starch Press, 2021) — 게이트웨이 / 펌웨어
- OWASP IoT Top 10 — OWASP Foundation (URL 확인 필요)
- 「분산에너지 활성화 특별법」 (2024.6.14, 법률 19470호) — 국가법령정보센터 (URL 확인 필요)
- 「개인정보 보호법」 — 국가법령정보센터 (URL 확인 필요)
- CSAP / KCMVP — KISA / 국가정보원 등재 모듈 목록 (URL 확인 필요)
- KWCAG 2.2 — 한국지능정보사회진흥원 (URL 확인 필요)

---

## 부록 갱신 규칙

- 본 부록은 **본 책에 실제로 등장하는** 도구·서비스만 다룬다.
- 신규 도구 추가 시 — 본문에서 첫 등장하는 챕터에서 이 부록으로 링크.
- 본 케이스 사용 여부를 명시 — **사용 / 사용 X / 대안 / 검토** 4 카테고리.
- "왜" 와 "함정" 이 핵심 — 마케팅 카탈로그가 아니라 의사결정 보조.
