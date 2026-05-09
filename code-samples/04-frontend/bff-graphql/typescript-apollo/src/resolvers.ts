import DataLoader from 'dataloader';
import { backends, Order, Product, User } from './backends.js';

/**
 * 한 요청마다 새 DataLoader 인스턴스 생성 — 요청 간 캐시 누수 방지.
 * Apollo Server 의 context 함수에서 매 요청마다 createLoaders() 호출.
 */
export function createLoaders() {
  return {
    user: new DataLoader<string, User | null>(async (ids) => {
      const users = await backends.users.batchGet(ids as string[]);
      return ids.map(id => users.find(u => u.id === id) ?? null);
    }),
    product: new DataLoader<string, Product | null>(async (ids) => {
      const products = await backends.products.batchGet(ids as string[]);
      return ids.map(id => products.find(p => p.id === id) ?? null);
    }),
  };
}

export type Loaders = ReturnType<typeof createLoaders>;

interface Ctx {
  loaders: Loaders;
}

export const resolvers = {
  Query: {
    orders: (_: unknown, args: { userId: string }) =>
      backends.orders.listByUser(args.userId),
    order: (_: unknown, args: { id: string }) =>
      backends.orders.byId(args.id),
  },

  Order: {
    // user: order.userId 로 DataLoader 호출 — 같은 batch 안의 동일 userId 는 1번만 호출.
    user: (order: Order, _: unknown, { loaders }: Ctx) =>
      loaders.user.load(order.userId),

    products: (order: Order, _: unknown, { loaders }: Ctx) =>
      Promise.all(order.productIds.map(id => loaders.product.load(id))),
  },
};
