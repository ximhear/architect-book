# 다이어그램 (보강 작업)

> 본문 inline Mermaid 로는 표현이 부자연스러운 자유 도형 / 메타포 / 풍경화를 위한 디렉토리.
> 도구 우선순위는 [ADR-0006](../adr/0006-다이어그램-도구-우선순위.md) 참조.

## 도구 우선순위

| 순위 | 도구 | 사용처 | 본 디렉토리 처리 |
|------|------|--------|----------------|
| 1 | **Mermaid** | 시퀀스 / 플로우 / C4 / ER | 본문 inline (`manuscript/<chapter>/`) — 본 디렉토리 외 |
| 2 | **Excalidraw** | 자유 도형 / 손그림 / 메타포 | `.excalidraw` JSON + PNG export |
| 3 | **drawio** | 위 둘로 안 되는 복잡 시스템도 | `.drawio.xml` + PNG export |

## 디렉토리 구조

```
diagrams/
├── README.md                  # 이 파일
└── <챕터별>/
    ├── plan.md                # 어떤 다이어그램이 필요한지 명세
    ├── *.excalidraw           # Excalidraw 원본 (편집 가능)
    ├── *.drawio.xml           # drawio 원본 (선택)
    └── exports/               # PNG / SVG export (assets/ 와 짝)
        └── *.png
```

## 작업 흐름

1. `<챕터>/plan.md` 에 어떤 다이어그램이 필요한지 명세
2. Excalidraw (https://excalidraw.com) 또는 drawio 에서 그림
3. 원본 (`.excalidraw` JSON 또는 `.drawio.xml`) 을 `<챕터>/` 에 저장
4. PNG / SVG export 를 `<챕터>/exports/` 또는 `assets/<챕터>/` 로
5. 본문에서 `assets/<챕터>/<이름>.png` 로 참조

## 현재 상태

각 챕터 `plan.md` 에 **필요한 다이어그램의 명세** 만 정리되어 있음. 실제 작도는 다음 단계.

| 챕터 | 본문 Mermaid 수 | 추가 필요 (Excalidraw) |
|------|---------------|----------------------|
| 1 | 3 | 2건 — NFR 충돌 거미줄, 7.2 진화 풍경 |
| 2 | 3 | 1건 — 멀티 클라우드 메타포 (한국 지도) |
| 3 | 4 | 2건 — 헥사고날 메타포, Saga 보상 흐름 풍경 |
| 4 | 2 | 2건 — 채널 구성 메타포, 라이브 화면 합성 |
| 5 | 4 | 1건 — 검색 한국어 형태소 흐름 |
| 6 | 3 | 2건 — 한국 법규 지형도, 망분리 3종 비교 |
| 7 | 3 | 2건 — 팀 토폴로지 사진, 인시던트 사람 흐름 |
| 8 | 4 | 1건 — 원픽 전체 풍경화 (한 장 마스터) |

총 13건. 우선순위는 각 `plan.md` 참조.

## Mermaid vs Excalidraw 결정 기준

- **Mermaid 로 충분** — 박스·화살표·시퀀스·플로우 / 텍스트 라벨 / 그리드 배치
- **Excalidraw 필요** — 자유 곡선 / 손그림 느낌 / 메타포 / 사진 같은 풍경 / 색·아이콘 강조
- **drawio 필요** — 위 둘로도 안 되는 매우 복잡한 시스템도 (드물게)

## 파일명 규칙

- `<영역>-<짧은이름>.excalidraw` — 예: `nfr-conflicts.excalidraw`, `live-screen.excalidraw`
- 챕터 안에서 본문 등장 순서대로 번호 부여 권장 — `01-nfr-conflicts.excalidraw`
