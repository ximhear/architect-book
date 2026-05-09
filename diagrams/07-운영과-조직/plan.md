# 7장 — 다이어그램 보강 plan

## 본문에 이미 있는 Mermaid

1. 1.1 DevOps / SRE / Platform 관계 (graph TB — 점선 공유)
2. 3.1 인시던트 단계 (flowchart LR)
3. 5.1 팀 토폴로지 4유형 + 인터랙션 모드 (graph TB)

## 추가 필요 (Excalidraw)

### [01-team-topology-photo.excalidraw] 팀 토폴로지 4가지 — 사진 / 메타포

- **위치**: 본문 5.1 팀 토폴로지 표 옆
- **목적**: 추상적인 4가지 팀 유형 (Stream-aligned / Platform / Enabling / Complicated Subsystem) 을 사진/메타포로
- **표현**:
  - Stream-aligned — 강물 (단일 흐름)
  - Platform — 다리 / 발판 (다른 팀이 그 위에 서서 일함)
  - Enabling — 코치 (한시적 강화, 떠남)
  - Complicated Subsystem — 시계 내부 톱니 (전문성 깊음)
  - 4개를 한 풍경에 배치 + 인터랙션 모드 화살표
- **메타포**: "팀은 박스가 아니라 풍경"
- **우선순위**: 보통

### [02-incident-people-flow.excalidraw] 인시던트 — 사람의 흐름

- **위치**: 본문 3절 인시던트 관리
- **목적**: 알람부터 회고까지 IC / Operations Lead / Communications Lead / Scribe 가 각각 어떤 시점에 등장하는지
- **표현**:
  - 가로 시간축 (탐지 → 격리 → 복구 → 회고)
  - 위쪽에 4 역할의 막대 — 각자 활동 구간 표시
  - 사이사이 핵심 결정 — IC 지명 / 외부 통보 / 복구 완료 / 회고
  - 우측 — "비난 없는 회고가 정직한 보고로 이어짐" 콜아웃
- **메타포**: "장애는 시스템이 아니라 사람의 흐름"
- **우선순위**: 보통
