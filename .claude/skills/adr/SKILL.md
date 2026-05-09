---
name: adr
description: Architect-book 의 Architecture Decision Record 를 표준 형식 (Context / Decision / Consequences) 으로 새로 만든다. 다음 ADR 번호를 자동으로 계산해 adr/NNNN-<제목>.md 파일을 생성한다. 인자로 결정의 제목을 받는다 (예 - "/adr 빌드툴로 mdBook 채택").
---

# adr

책 자체에 대한 의사결정 (집필 형식, 빌드 도구, 표기법, 챕터 구조 등) 을 ADR 파일로 기록한다.

## 입력

결정의 제목을 자유 텍스트로 받는다. 예:
- `빌드툴로 mdBook 채택`
- `코드 예제 기본 언어를 Spring Boot 와 Go 로 정함`
- `한국 법규 인용은 시점 명시 필수`

## 절차

1. **다음 번호 계산**:
   - `adr/` 디렉토리에서 `NNNN-*.md` 패턴을 찾아 가장 큰 번호 + 1.
   - 없으면 `0001` 부터 시작. 4자리 zero-pad.
2. **파일명 생성**: `adr/NNNN-<슬러그>.md`
   - 슬러그: 한글은 그대로 두고 공백을 `-` 로 치환. 영문 약어는 소문자로.
3. **템플릿 작성**: 아래 템플릿을 사용. Date 는 오늘 날짜 (YYYY-MM-DD), Status 는 `Proposed` 로 초기화.
4. **결과 보고**: 생성된 파일 경로를 알려주고, 사용자에게 Context / Decision / Consequences 세 섹션을 채우도록 안내.

## 템플릿

```markdown
# NNNN. <제목>

- **Date**: YYYY-MM-DD
- **Status**: Proposed
- **Deciders**: (의사결정자 — 보통 책의 저자)

## Context

이 결정이 필요한 배경. 어떤 문제를 풀고 있는지, 어떤 제약이 있는지.
(3~10줄)

## Decision

내린 결정과 그 핵심 근거.
(간결하게. "우리는 X 를 한다" 형태로.)

## Considered Alternatives

대안과 그 대안을 채택하지 않은 이유.

- **A**: (설명) — (탈락 이유)
- **B**: (설명) — (탈락 이유)

## Consequences

### 긍정적

- (얻는 것)

### 부정적

- (감수하는 것)

### 중립

- (그저 달라지는 것)

## Follow-ups

이 결정으로 인해 추가로 해야 할 일 (있다면).

- [ ] (할 일)

---

> **Status 변경 시**: 결정이 확정되면 `Accepted`, 폐기되면 `Deprecated`,
> 새 ADR 로 대체되면 `Superseded by NNNN` 으로 갱신.
```

## 주의

- **기존 ADR 의 Status 를 함부로 바꾸지 말 것**. ADR 은 시간순 의사결정 기록이므로 보존이 원칙. 변경이 필요하면 새 ADR 을 만들어 이전 것을 `Superseded by` 로 표시.
- 번호는 빈 자리(gap) 없이 순차 증가. 삭제된 번호를 재사용하지 말 것.
- 책의 내용 (예: "Spring 보다 Go 가 낫다") 에 대한 결정은 ADR 이 아니라 본문 또는 트레이드오프 박스에 들어간다. ADR 은 **책을 어떻게 만들지** 에 대한 메타 결정만.
