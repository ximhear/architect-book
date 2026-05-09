# 5장 — 다이어그램 보강 plan

## 본문에 이미 있는 Mermaid

1. 1.1 데이터 플랫폼 큰 그림 (graph LR)
2. 3.3 Lambda vs Kappa (graph TB — 비교)
3. 4.3 색인 파이프라인 (graph LR)
4. 6.2 LLM 패턴 (graph TB — 단순/RAG/Agent 비교)

## 추가 필요 (Excalidraw)

### [01-korean-search-flow.excalidraw] 한국어 검색의 형태소 흐름 — 단계별 분해

- **위치**: 본문 4.2 한국어 검색 함정 절
- **목적**: "삼성 갤럭시" 한 단어가 색인·검색 단계를 거치며 어떻게 분해·정규화·매칭되는지 한 그림으로
- **표현**:
  - 입력 "삼성 갤럭시 S24" → 형태소 분석 (Nori) → 토큰 리스트
  - 토큰 리스트 → 시노님 적용 ("Samsung", "SAMSUNG" 동등) → 정규화
  - 정규화된 토큰 → 색인 (OpenSearch)
  - 검색 시: "샴성 갤럭시" → 자모 분해 → edit distance 매칭 → 정규 토큰으로 정정 → 검색
  - 함정 — Nori 사전 hot-reload 미지원, AWS Managed Nori 의 한계, 자모 false positive (ㄱ → 모든 ㄱ 단어)
- **메타포**: "한국어는 영어보다 한 단계 더 깊다"
- **우선순위**: 높음 (본 책의 차별화 영역)
