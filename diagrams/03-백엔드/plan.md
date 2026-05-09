# 3장 — 다이어그램 보강 plan

## 본문에 이미 있는 Mermaid

1. 1.1 헥사고날 구조 (graph LR — adapter/application/domain/port)
2. 4.2 Saga 시퀀스 (sequenceDiagram)
3. 8.1 도메인 분할 (graph TB)
4. 8.4 주문 결제 시퀀스 (sequenceDiagram — PG 라우터 fallback)

## 추가 필요 (Excalidraw)

### [01-hexagonal-metaphor.excalidraw] 헥사고날 메타포 — 6각형 성

- **위치**: 본문 1.1 헥사고날 다이어그램 옆 (Mermaid 보완)
- **목적**: 도메인이 외부 의존성으로부터 어떻게 보호되는지 메타포로
- **표현**:
  - 6각형 성 중앙 — `domain` (순수)
  - 6각 변 — port (in/out)
  - 6각 변 바깥 — adapter (HTTP / JPA / Kafka / PG / 배송사 / 이벤트 컨슈머)
  - 적군 (외부 변경) 화살표가 어댑터에서 멈추는 그림
- **메타포**: "도메인은 성, 어댑터는 성벽"
- **우선순위**: 보통

### [02-saga-compensation.excalidraw] Saga 보상의 보상 — 3계층 안전망

- **위치**: 본문 4.2 보상의 보상 표 옆
- **목적**: DLQ + 수동 reconciliation + 일 단위 reconciler 가 어떻게 떨어진 메시지를 잡는지 그물로
- **표현**:
  - 위에서 떨어지는 메시지 (정상 / 보상 / 보상 실패)
  - 1차 그물 = DLQ (자동, 알람만)
  - 2차 그물 = 운영팀 수동 reconciliation (사람 손)
  - 3차 그물 = 일 단위 reconciler (배치, 어긋난 거 잡음)
  - "RPO 1분" 라벨 — 결제 도메인은 3중 모두 필수
- **메타포**: "낚시 / 곡예 안전망"
- **우선순위**: 높음 (책에서 가장 자주 언급되는 보상 패턴의 시각적 보강)
