<script setup lang="ts">
import { useTickerSocket } from './composables/useTickerSocket'

const { isConnected, ticks } = useTickerSocket('ws://localhost:8080/ws')

const formatTime = (ts: number) => {
  if (!ts) return '-'
  // Otomatis deteksi apakah timestamp dari Go dalam detik atau milidetik
  const ms = ts < 1e12 ? ts * 1000 : ts
  return new Date(ms).toLocaleTimeString('id-ID')
}
</script>

<template>
  <main class="container">
    <header class="header">
      <h1>Crypto Ticker (Go WS + Protobuf)</h1>
      <span :class="['badge', isConnected ? 'connected' : 'disconnected']">
        {{ isConnected ? '● Connected (Binary)' : '○ Disconnected' }}
      </span>
    </header>

    <div v-if="Object.keys(ticks).length === 0" class="empty">
      Menunggu stream binary Protobuf dari ws://localhost:8080/ws ...
    </div>

    <div v-else class="grid">
      <div v-for="item in ticks" :key="item.symbol" class="card">
        <div class="card-top">
          <span class="symbol">{{ item.symbol }}</span>
          <span :class="['change', item.changePct >= 0 ? 'up' : 'down']">
            {{ item.changePct >= 0 ? '+' : '' }}{{ item.changePct.toFixed(2) }}%
          </span>
        </div>

        <div class="price">
          ${{ item.price.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}
        </div>

        <div class="timestamp">Update: {{ formatTime(item.timestamp) }}</div>
      </div>
    </div>
  </main>
</template>

<style scoped>
.container {
  max-width: 760px;
  margin: 2.5rem auto;
  padding: 0 1rem;
  font-family: system-ui, -apple-system, sans-serif;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}
.header h1 {
  font-size: 1.4rem;
  margin: 0;
}
.badge {
  padding: 0.3rem 0.75rem;
  border-radius: 999px;
  font-size: 0.85rem;
  font-weight: 600;
}
.connected {
  background: #dcfce7;
  color: #166534;
}
.disconnected {
  background: #fee2e2;
  color: #991b1b;
}
.empty {
  padding: 2.5rem;
  text-align: center;
  color: #6b7280;
  border: 1px dashed #d1d5db;
  border-radius: 0.75rem;
}
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(210px, 1fr));
  gap: 1rem;
}
.card {
  padding: 1.25rem;
  border: 1px solid #e5e7eb;
  border-radius: 0.75rem;
  background: #fff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}
.card-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.symbol {
  font-size: 1.2rem;
  font-weight: 700;
}
.change {
  font-weight: 600;
  font-size: 0.9rem;
}
.up { color: #16a34a; }
.down { color: #dc2626; }
.price {
  font-size: 1.75rem;
  font-weight: 800;
  margin: 0.75rem 0;
  font-variant-numeric: tabular-nums;
}
.timestamp {
  font-size: 0.8rem;
  color: #6b7280;
}
</style>
