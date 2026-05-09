// GraphQL 스키마.
// 모바일 화면이 필요한 만큼만 가져갈 수 있도록 객체 그래프로 표현.
export const typeDefs = `#graphql
  type User {
    id: ID!
    name: String!
  }

  type Product {
    id: ID!
    name: String!
    price: Int!
  }

  type Order {
    id: ID!
    amount: Int!
    status: String!
    user: User!         # 같은 사용자 N개 주문 → DataLoader 로 1 batch
    products: [Product!]!
  }

  type Query {
    orders(userId: ID!): [Order!]!
    order(id: ID!): Order
  }
`;
