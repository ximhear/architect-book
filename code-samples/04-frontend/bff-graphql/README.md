# BFF GraphQL — Apollo Server (TypeScript) (4장 4.3)

> 모바일 BFF — 백엔드 N개를 한 응답으로 합성.
> 화면별 쿼리 자유 + DataLoader 로 N+1 방지.

## 무엇을 보여주는가

- **GraphQL 스키마** — Order / Product / User
- **Resolver** — 백엔드 REST / gRPC 호출 시뮬
- **DataLoader** — 같은 batch 안의 같은 ID 요청 합치기 (N+1 방지)
- **요청 1번 → 응답 N개 합성** — 모바일 데이터 양 절약

## 디렉토리

```
bff-graphql/
└── typescript-apollo/
    ├── package.json
    ├── tsconfig.json
    └── src/
        ├── index.ts          # Apollo Server 부트
        ├── typeDefs.ts       # GraphQL 스키마
        ├── resolvers.ts      # Resolver + DataLoader
        └── backends.ts       # 백엔드 호출 시뮬
```

## 핵심 — DataLoader 로 N+1 방지

```typescript
// "주문 N개 + 각 주문의 사용자 정보" 쿼리
// 일반 resolver: N+1 (주문 N개 → 사용자 N번 호출)
// DataLoader: 1 batch (주문 N개 → 사용자 1번 batch 호출)

const userLoader = new DataLoader<string, User>(async (userIds) => {
  const users = await backends.users.batchGet(userIds);  // 한 번 호출
  return userIds.map(id => users.find(u => u.id === id)!);
});
```

## 책 인용

- 4장 4.3 — BFF GraphQL (모바일) / Next.js Server Action (웹)
- 4장 4.1 — BFF 시퀀스 다이어그램 (병렬 호출 패턴)

## 실행

```bash
cd typescript-apollo
npm install
npm run dev
# 브라우저: http://localhost:4000
```

GraphQL Sandbox 에서 쿼리 시도:

```graphql
query MainScreen {
  orders(userId: "user-1") {
    id
    amount
    user { id name }
    products { id name price }
  }
}
```

## 한국 환경 메모

- **모바일 클라이언트** — 첫 화면 < 200KB ([케이스 5.3](../../../manuscript/case-study/README.md#53-기술-제약)) 부담을 BFF 가 합성으로 줄임
- **GraphQL Federation** — 도메인별 subgraph 가 늘어나면 Apollo Federation 검토
- **Spring 백엔드 연동** — 한국 백엔드 시장은 Spring + DGS Framework (Kotlin) 도 인기 — 향후 추가 예제

## 라이선스

[MIT](../../LICENSE)
