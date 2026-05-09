// 백엔드 호출 시뮬. 실제는 Spring REST / gRPC 호출.

export interface User { id: string; name: string }
export interface Product { id: string; name: string; price: number }
export interface Order {
  id: string;
  amount: number;
  status: string;
  userId: string;
  productIds: string[];
}

// 학습용 in-memory 데이터
const usersDb: User[] = [
  { id: 'user-1', name: '김원픽' },
  { id: 'user-2', name: '이수민' },
];

const productsDb: Product[] = [
  { id: 'p-100', name: '갤럭시 S24', price: 1_200_000 },
  { id: 'p-101', name: '아이폰 15', price: 1_500_000 },
  { id: 'p-200', name: 'Pretendard 굿즈 티셔츠', price: 25_000 },
];

const ordersDb: Order[] = [
  { id: 'o-1', amount: 1_225_000, status: 'PAID', userId: 'user-1', productIds: ['p-100', 'p-200'] },
  { id: 'o-2', amount: 1_500_000, status: 'PAID', userId: 'user-1', productIds: ['p-101'] },
  { id: 'o-3', amount: 25_000,    status: 'CREATED', userId: 'user-2', productIds: ['p-200'] },
];

// API 호출 시뮬 — 실제는 fetch / gRPC 클라이언트.
export const backends = {
  orders: {
    async listByUser(userId: string): Promise<Order[]> {
      return ordersDb.filter(o => o.userId === userId);
    },
    async byId(id: string): Promise<Order | null> {
      return ordersDb.find(o => o.id === id) ?? null;
    },
  },
  users: {
    async batchGet(ids: readonly string[]): Promise<User[]> {
      console.log(`[backends.users.batchGet] ids=${ids.join(',')}  ← 1번만 호출되어야 함`);
      return usersDb.filter(u => ids.includes(u.id));
    },
  },
  products: {
    async batchGet(ids: readonly string[]): Promise<Product[]> {
      console.log(`[backends.products.batchGet] ids=${ids.join(',')}`);
      return productsDb.filter(p => ids.includes(p.id));
    },
  },
};
