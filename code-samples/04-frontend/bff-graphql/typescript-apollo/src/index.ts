import { ApolloServer } from '@apollo/server';
import { startStandaloneServer } from '@apollo/server/standalone';
import { typeDefs } from './typeDefs.js';
import { resolvers, createLoaders } from './resolvers.js';

const server = new ApolloServer({
  typeDefs,
  resolvers,
});

const { url } = await startStandaloneServer(server, {
  listen: { port: 4000 },
  // 매 요청마다 새 DataLoader 생성 — 요청 간 캐시 격리
  context: async () => ({
    loaders: createLoaders(),
  }),
});

console.log(`🚀 BFF GraphQL ready at: ${url}`);
console.log(`Try: query { orders(userId: "user-1") { id user { name } products { name price } } }`);
